package app_agent_chat_runner

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_registry"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/domain_service"
	biz_agent "github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/biz_entity"
	po_agent "github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/po_entity"
	biz_entity_app_config "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
)

type RunnerRuntimeParameters struct {
	message               *po_entity.Message
	isFirstChunk          bool
	functionCallState     bool
	interactionStep       int
	maxInteractionSteps   int
	toolCalls             []*ToolCall
	toolCallNames         string
	toolCallInputs        map[string]map[string]any
	historyPromptMessages []*po_entity.PromptMessage
	assistantThoughts     []po_entity.IPromptMessage
	toolResponse          []*ToolResponseItem
	fullAnswer            string
}

type ToolInvokeMeta struct {
	TimeCost   float64        `json:"time_cost"`
	Error      error          `json:"error"`
	ToolConfig map[string]any `json:"tool_config"`
}

type ToolResponseItem struct {
	ToolCallID   string          `json:"tool_call_id"`
	ToolCallName string          `json:"tool_call_name"`
	ToolResponse string          `json:"tool_response"`
	Meta         *ToolInvokeMeta `json:"meta"`
}

type FunctionCallAgentRunner struct {
	agentThoughtCount         int
	applicationGenerateEntity *biz_entity_app_generate.AgentChatAppGenerateEntity
	appConfig                 *biz_entity_app_config.AgentChatAppConfig
	conversation              *po_entity.Conversation
	agentDomain               *domain_service.AgentDomain
	queueManager              *biz_entity.StreamGenerateQueue
	agentFlusher              biz_agent.AgentFlusher
	promptToolMessages        []*biz_entity.PromptMessageTool
	modelCaller               model_registry.IModelRegistryCall
	promptMessages            []po_entity.IPromptMessage
	toolRuntimeMap            map[string]*biz_agent.ToolRuntimeConfiguration
	*RunnerRuntimeParameters
	providerType string
}

func NewFunctionCallAgentRunner(tenantID string, applicationGenerateEntity *biz_entity_app_generate.AgentChatAppGenerateEntity, conversation *po_entity.Conversation, agentDomain *domain_service.AgentDomain, queueManager *biz_entity.StreamGenerateQueue, agentFlusher biz_agent.AgentFlusher,
	promptToolMessage []*biz_entity.PromptMessageTool, promptMessage []po_entity.IPromptMessage, toolRuntimeMap map[string]*biz_agent.ToolRuntimeConfiguration, modelCaller model_registry.IModelRegistryCall, providerType string) *FunctionCallAgentRunner {

	return &FunctionCallAgentRunner{
		applicationGenerateEntity: applicationGenerateEntity,
		conversation:              conversation,
		toolRuntimeMap:            toolRuntimeMap,
		agentDomain:               agentDomain,
		queueManager:              queueManager,
		providerType:              providerType,
		agentFlusher:              agentFlusher,
		promptToolMessages:        promptToolMessage,
		promptMessages:            promptMessage,
		modelCaller:               modelCaller,
		RunnerRuntimeParameters: &RunnerRuntimeParameters{
			interactionStep:       1,
			isFirstChunk:          true,
			historyPromptMessages: make([]*po_entity.PromptMessage, 0),
			functionCallState:     true,
			maxInteractionSteps:   int(math.Min(float64(applicationGenerateEntity.MaxIteration), 5)),
			toolResponse:          make([]*ToolResponseItem, 0),
		},
	}
}

type ToolCall struct {
	ToolCallID   string
	ToolCallName string
	TollCallArgs map[string]any
}

func (fca *FunctionCallAgentRunner) Run(ctx context.Context, message *po_entity.Message, query string, streamQueue chan *biz_entity.MessageQueueMessage) error {

	fca.message = message

	for fca.functionCallState && fca.interactionStep <= int(fca.maxInteractionSteps) {
		var (
			response string
		)

		fca.functionCallState = false

		if fca.interactionStep == fca.maxInteractionSteps {
			fca.promptToolMessages = make([]*biz_entity.PromptMessageTool, 0)
		}

		fca.organizePromptMessage()

		if fca.interactionStep > 1 {
			go fca.interactionInvokeLLM(ctx)
		}

		agentThought, err := fca.handleStreamAgentMessageQueue(ctx, streamQueue)
		if err != nil {
			return err
		}

		agentThought.Tool = fca.toolCallNames
		agentThought.ToolInput = fca.toolCallInputs

		if err := fca.agentDomain.UpdateAgentThought(ctx, agentThought); err != nil {
			return err
		}

		if err := fca.agentFlusher.AgentThoughtToStreamResponse(ctx, agentThought.ID); err != nil {
			return err
		}

		assistantMessage := biz_entity.NewAssistantPromptMessage(response)

		if len(fca.toolCalls) > 0 {
			for _, toolCall := range fca.toolCalls {
				assistantMessage.ToolCalls = append(assistantMessage.ToolCalls, &biz_entity.ToolCall{
					ID:   toolCall.ToolCallID,
					Type: "function",
					Function: &biz_entity.ToolCallFunction{
						Name:      toolCall.ToolCallName,
						Arguments: toolCall.TollCallArgs,
					},
				})
			}
		}

		fca.assistantThoughts = append(fca.assistantThoughts, assistantMessage)

		agentThought.ToolInput = fca.toolCallInputs
		agentThought.Tool = fca.toolCallNames

		if err := fca.agentDomain.UpdateAgentThought(ctx, agentThought); err != nil {
			return err
		}

		if err := fca.agentFlusher.AgentThoughtToStreamResponse(ctx, agentThought.ID); err != nil {
			return err
		}

		var toolArtifacts []*biz_agent.ToolArtifact

		for _, toolCall := range fca.toolCalls {

			toolRuntimeIns, ok := fca.toolRuntimeMap[toolCall.ToolCallName]

			if !ok {
				toolArtifacts = append(toolArtifacts, &biz_agent.ToolArtifact{
					ToolCallID:   toolCall.ToolCallID,
					ToolCallName: toolCall.ToolCallName,
					ToolResponse: fmt.Sprintf("there is not a tool named %s", toolCall.ToolCallName),
					Meta:         biz_agent.ErrorInvokeMetaIns(fmt.Sprintf("there is not a tool named %s", toolCall.ToolCallName)),
				})
			}

			toolEngine := domain_service.NewToolEngine(toolRuntimeIns, message, fca.providerType)

			toolInvokeResponse := toolEngine.AgentInvoke(ctx, toolCall.TollCallArgs, fca.applicationGenerateEntity.UserID, fca.appConfig.TenantID, biz_agent.InvokeFrom(fca.applicationGenerateEntity.InvokeFrom))

			toolArtifacts = append(toolArtifacts, &biz_agent.ToolArtifact{
				ToolCallID:   toolCall.ToolCallID,
				ToolCallName: toolCall.ToolCallName,
				ToolResponse: toolInvokeResponse.InvokeToolPrompt,
				Meta:         toolInvokeResponse.ToolInvokeMeta,
			})

			if toolInvokeResponse.InvokeToolPrompt != "" {
				fca.assistantThoughts = append(fca.assistantThoughts, &po_entity.ToolPromptMessage{
					PromptMessage: &po_entity.PromptMessage{
						Content: toolInvokeResponse.InvokeToolPrompt,
						Role:    po_entity.TOOL,
						Name:    toolCall.ToolCallName,
					},
					ToolCallID: toolCall.ToolCallID,
				})
			}
		}

		if len(toolArtifacts) > 0 {
			observation, meta := fca.getObservationAndMeta(toolArtifacts)
			agentThought.ToolMetaStr = meta
			agentThought.Observation = observation

			if err := fca.agentDomain.UpdateAgentThought(ctx, agentThought); err != nil {
				return err
			}
		}

		if err := fca.agentFlusher.AgentThoughtToStreamResponse(ctx, agentThought.ID); err != nil {
			return err
		}

		fca.interactionStep += 1
	}

	llmResult := &biz_entity.LLMResult{
		Model:         fca.applicationGenerateEntity.Model,
		PromptMessage: fca.promptMessages,
		Message: &biz_entity.AssistantPromptMessage{
			PromptMessage: &po_entity.PromptMessage{
				Content: fca.fullAnswer,
			},
		},
		Usage: biz_entity.NewEmptyLLMUsage(),
	}

	event := biz_entity.NewAppQueueEvent(biz_entity.MessageEnd)
	fca.queueManager.FinalManual(&biz_entity.QueueMessageEndEvent{
		AppQueueEvent: event,
		LLMResult:     llmResult,
	})

	return nil
}

func (fca *FunctionCallAgentRunner) interactionInvokeLLM(ctx context.Context) {

	fca.modelCaller.InvokeLLM(ctx, fca.promptMessages, fca.queueManager, fca.applicationGenerateEntity.ModelConf.Parameters, fca.promptToolMessages, make([]string, 0), fca.applicationGenerateEntity.UserID, nil)
}

func (fca *FunctionCallAgentRunner) handleStreamAgentMessageQueue(ctx context.Context, streamQueue chan *biz_entity.MessageQueueMessage) (*po_agent.MessageAgentThought, error) {

	agentThought, err := fca.CreateAgentThought(ctx, fca.message.ID, "", "", make(map[string]map[string]any, 0), []string{})

	if err != nil {
		return nil, err
	}

	for resultChunk := range streamQueue {
		if fca.isFirstChunk {
			err := fca.agentFlusher.AgentThoughtToStreamResponse(ctx, agentThought.ID)
			if err != nil {
				return nil, err
			}
			fca.isFirstChunk = false
		}

		if chunkEvent, ok := resultChunk.Event.(*biz_entity.QueueAgentMessageEvent); ok {

			if chunkEvent.Chunk.Delta.FinishReason == biz_entity.AGENT_END {
				fca.fullAnswer = chunkEvent.Chunk.Delta.Message.GetContent()
				break
			}

			if fca.checkTools(chunkEvent.Chunk) {
				fca.functionCallState = true
				fca.toolCalls = append(fca.toolCalls, fca.extractToolCalls(chunkEvent.Chunk)...)
				fca.toolCallNames = fca.getToolNames(fca.toolCalls)
				fca.toolCallInputs = fca.getToolInputs(fca.toolCalls)
			}
			deltaText := chunkEvent.Chunk.Delta.Message.Content
			if err := fca.agentFlusher.AgentMessageToStreamResponse(deltaText.(string)); err != nil {
				return nil, err
			}
		}
	}

	return agentThought, nil
}

// func (fca *FunctionCallAgentRunner) initSystemMessage(prePrompt string, promptMessages []*po_entity.PromptMessage) []*po_entity.PromptMessage {

// 	if len(promptMessages) == 0 && prePrompt != "" {
// 		return append(fca.historyPromptMessages, po_entity.NewSystemMessage(prePrompt))
// 	}

// 	var appendedPromptMessage []*po_entity.PromptMessage

// 	if promptMessages[0].Role != po_entity.SYSTEM && prePrompt != "" {
// 		appendedPromptMessage = append(appendedPromptMessage, po_entity.NewSystemMessage(prePrompt))
// 		appendedPromptMessage = append(appendedPromptMessage, promptMessages...)
// 	}

// 	return appendedPromptMessage
// }

func (fca *FunctionCallAgentRunner) checkTools(llmChunk *biz_entity.LLMResultChunk) bool {
	return len(llmChunk.Delta.Message.ToolCalls) > 0
}

func (fca *FunctionCallAgentRunner) extractToolCalls(llmChunk *biz_entity.LLMResultChunk) []*ToolCall {
	var toolCalls []*ToolCall
	for _, prompt := range llmChunk.Delta.Message.ToolCalls {
		toolCalls = append(toolCalls, &ToolCall{
			ToolCallID:   prompt.ID,
			ToolCallName: prompt.Function.Name,
			TollCallArgs: prompt.Function.Arguments,
		})
	}

	return toolCalls
}

func (fca *FunctionCallAgentRunner) getToolNames(tools []*ToolCall) string {
	var toolNames []string
	for _, tool := range tools {
		toolNames = append(toolNames, tool.ToolCallName)
	}
	return strings.Join(toolNames, ";")
}

func (faa *FunctionCallAgentRunner) getObservationAndMeta(toolArtifacts []*biz_agent.ToolArtifact) (map[string]string, map[string]*po_agent.ToolEngineInvokeMeta) {

	var (
		observation = make(map[string]string)
		meta        = make(map[string]*po_agent.ToolEngineInvokeMeta)
	)

	for _, artifact := range toolArtifacts {
		observation[artifact.ToolCallName] = artifact.ToolResponse
		meta[artifact.ToolCallName] = biz_agent.ConvertToPoMeta(artifact.Meta)
	}

	return observation, meta

}
func (fca *FunctionCallAgentRunner) getToolInputs(tools []*ToolCall) map[string]map[string]any {
	var toolInput = make(map[string]map[string]any)

	for _, tool := range tools {
		toolInput[tool.ToolCallName] = tool.TollCallArgs
	}
	return toolInput
}

func (fca *FunctionCallAgentRunner) organizePromptMessage() {
	if len(fca.assistantThoughts) > 0 {
		fca.promptMessages = append(fca.promptMessages, fca.assistantThoughts...)
	}
}

func (fca *FunctionCallAgentRunner) CreateAgentThought(ctx context.Context, messageID, message, toolName string, toolInput map[string]map[string]any, messageFileIDs []string) (*po_agent.MessageAgentThought, error) {

	thoughtObject := &po_agent.MessageAgentThought{
		MessageID:     messageID,
		Tool:          toolName,
		ToolInput:     toolInput,
		Message:       message,
		Position:      fca.agentThoughtCount + 1,
		Currency:      "USD",
		CreatedByRole: "account",
		MessageFiles:  messageFileIDs,
		CreatedBy:     fca.applicationGenerateEntity.EasyUIBasedAppGenerateEntity.UserID,
	}

	thought, err := fca.agentDomain.AgentRepo.CreateAgentThought(ctx, thoughtObject)

	if err != nil {
		return nil, err
	}

	fca.agentThoughtCount += 1
	return thought, nil
}
