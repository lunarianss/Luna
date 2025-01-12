package app_agent_chat_runner

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/lunarianss/Luna/infrastructure/log"
	"github.com/lunarianss/Luna/internal/api-server/core/model_runtime/model_registry"
	"github.com/lunarianss/Luna/internal/api-server/domain/agent/domain_service"
	biz_agent "github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/biz_entity"
	po_agent "github.com/lunarianss/Luna/internal/api-server/domain/agent/entity/po_entity"
	biz_entity_app_config "github.com/lunarianss/Luna/internal/api-server/domain/app/entity/biz_entity/provider_app_config"
	biz_entity_chat_prompt_message "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/chat_prompt_message"
	biz_entity_openai_standard_response "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/openai_standard_response"
	biz_entity_base_stream_generator "github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/biz_entity/stream_base_generator"
	"github.com/lunarianss/Luna/internal/api-server/domain/chat/entity/po_entity"
	biz_entity_app_generate "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/provider_app_generate"
)

type RunnerRuntimeParameters struct {
	message               *po_entity.Message
	functionCallState     bool
	interactionStep       int
	maxInteractionSteps   int
	toolCalls             []*ToolCall
	toolCallNames         string
	toolCallInputs        map[string]string
	historyPromptMessages []*biz_entity_chat_prompt_message.PromptMessage
	assistantThoughts     []biz_entity_chat_prompt_message.IPromptMessage
	toolResponse          []*ToolResponseItem
	taskState             *biz_entity_base_stream_generator.ChatAppTaskState
	fullAssistant         string
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
	biz_entity_base_stream_generator.IStreamGenerateQueue
	agentFlusher       biz_agent.AgentFlusher
	promptToolMessages []*biz_entity_chat_prompt_message.PromptMessageTool
	modelCaller        model_registry.IModelRegistryCall
	promptMessages     []biz_entity_chat_prompt_message.IPromptMessage
	toolRuntimeMap     map[string]*biz_agent.ToolRuntimeConfiguration
	*RunnerRuntimeParameters
	providerType string
	secretKet    string
	fileBaseUrl  string
	bucket       string
}

func NewFunctionCallAgentRunner(tenantID string, applicationGenerateEntity *biz_entity_app_generate.AgentChatAppGenerateEntity, conversation *po_entity.Conversation, agentDomain *domain_service.AgentDomain, queueManager biz_entity_base_stream_generator.IStreamGenerateQueue, agentFlusher biz_agent.AgentFlusher,
	promptToolMessage []*biz_entity_chat_prompt_message.PromptMessageTool, promptMessage []biz_entity_chat_prompt_message.IPromptMessage, toolRuntimeMap map[string]*biz_agent.ToolRuntimeConfiguration, modelCaller model_registry.IModelRegistryCall, appConfig *biz_entity_app_config.AgentChatAppConfig, providerType string, secretKet, fileBaseUrl string, bucket string) *FunctionCallAgentRunner {

	return &FunctionCallAgentRunner{
		applicationGenerateEntity: applicationGenerateEntity,
		conversation:              conversation,
		secretKet:                 secretKet,
		fileBaseUrl:               fileBaseUrl,
		toolRuntimeMap:            toolRuntimeMap,
		agentDomain:               agentDomain,
		IStreamGenerateQueue:      queueManager,
		providerType:              providerType,
		bucket:                    bucket,
		agentFlusher:              agentFlusher,
		promptToolMessages:        promptToolMessage,
		promptMessages:            promptMessage,
		appConfig:                 appConfig,
		modelCaller:               modelCaller,
		RunnerRuntimeParameters: &RunnerRuntimeParameters{
			interactionStep:       1,
			historyPromptMessages: make([]*biz_entity_chat_prompt_message.PromptMessage, 0),
			functionCallState:     true,
			maxInteractionSteps:   int(math.Min(float64(applicationGenerateEntity.MaxIteration), 5)),
			toolResponse:          make([]*ToolResponseItem, 0),
			taskState: &biz_entity_base_stream_generator.ChatAppTaskState{
				LLMResult: biz_entity_base_stream_generator.NewEmptyLLMResult(),
			},
		},
	}
}

type ToolCall struct {
	ToolCallID   string
	ToolCallName string
	TollCallArgs string
}

func (fca *FunctionCallAgentRunner) Run(ctx context.Context, message *po_entity.Message, query string) (*biz_entity_base_stream_generator.ChatAppTaskState, error) {

	fca.message = message

	for fca.functionCallState && fca.interactionStep <= int(fca.maxInteractionSteps) {
		var (
			messageFileIDs []string
		)

		fca.functionCallState = false

		if fca.interactionStep == fca.maxInteractionSteps {
			fca.promptToolMessages = make([]*biz_entity_chat_prompt_message.PromptMessageTool, 0)
		}

		fca.organizePromptMessage()

		if fca.interactionStep > 1 {
			go fca.interactionInvokeLLM(ctx)
		}

		agentThought, err := fca.handleStreamAgentMessageQueue(ctx)
		if err != nil {
			return nil, err
		}

		assistantMessage := biz_entity_chat_prompt_message.NewAssistantToolPromptMessage(fca.fullAssistant)

		if len(fca.toolCalls) > 0 {
			for _, toolCall := range fca.toolCalls {
				assistantMessage.ToolCalls = append(assistantMessage.ToolCalls, &biz_entity_openai_standard_response.ToolCall{
					ID:   toolCall.ToolCallID,
					Type: "function",
					Function: &biz_entity_openai_standard_response.ToolCallFunction{
						Name:      toolCall.ToolCallName,
						Arguments: toolCall.TollCallArgs,
					},
				})
			}
		}

		// before tool call update agent thought
		fca.assistantThoughts = append(fca.assistantThoughts, assistantMessage)

		agentThought, err = fca.SaveAgentThought(ctx, agentThought, fca.toolCallNames, fca.toolCallInputs, fca.fullAssistant, nil, nil, fca.fullAssistant, []string{}, nil)

		if err != nil {
			return nil, err
		}

		if err := fca.agentFlusher.AgentThoughtToStreamResponse(ctx, agentThought.ID); err != nil {
			return nil, err
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

			toolEngine := domain_service.NewToolEngine(toolRuntimeIns, message, fca.providerType, fca.agentDomain, fca.bucket)

			toolInvokeResponse := toolEngine.AgentInvoke(ctx, toolCall.TollCallArgs, fca.applicationGenerateEntity.UserID, fca.appConfig.TenantID, biz_agent.InvokeFrom(fca.applicationGenerateEntity.InvokeFrom))

			for _, messageFile := range toolInvokeResponse.MessageFiles {
				err := fca.agentFlusher.AgentMessageFileToStreamResponse(ctx, messageFile.MessageFile, fca.secretKet, fca.fileBaseUrl)

				if err != nil {
					return nil, err
				}

				messageFileIDs = append(messageFileIDs, messageFile.MessageFile)
			}

			toolArtifacts = append(toolArtifacts, &biz_agent.ToolArtifact{
				ToolCallID:   toolCall.ToolCallID,
				ToolCallName: toolCall.ToolCallName,
				ToolResponse: toolInvokeResponse.InvokeToolPrompt,
				Meta:         toolInvokeResponse.ToolInvokeMeta,
			})

			if toolInvokeResponse.InvokeToolPrompt != "" {
				fca.assistantThoughts = append(fca.assistantThoughts, &biz_entity_chat_prompt_message.ToolPromptMessage{
					PromptMessage: &biz_entity_chat_prompt_message.PromptMessage{
						Content: toolInvokeResponse.InvokeToolPrompt,
						Role:    biz_entity_chat_prompt_message.TOOL,
						Name:    toolCall.ToolCallName,
					},
					ToolCallID: toolCall.ToolCallID,
				})
			}
		}

		if len(toolArtifacts) > 0 {
			observation, meta := fca.getObservationAndMeta(toolArtifacts)

			toolInvokeMeta := make(map[string]*biz_agent.ToolEngineInvokeMeta)

			for _, toolArtifact := range toolArtifacts {
				toolInvokeMeta[toolArtifact.ToolCallName] = toolArtifact.Meta
			}

			// after tool call update agent thought
			agentThought, err = fca.SaveAgentThought(ctx, agentThought, "", nil, "", observation, meta, "", messageFileIDs, nil)

			if err != nil {
				return nil, err
			}
		}

		if err := fca.agentFlusher.AgentThoughtToStreamResponse(ctx, agentThought.ID); err != nil {
			return nil, err
		}

		fca.interactionStep += 1
		fca.toolCalls = make([]*ToolCall, 0)
		fca.toolCallNames = ""
		fca.toolCallInputs = make(map[string]string)
		fca.fullAssistant = ""
	}

	return fca.taskState, nil
}

func (fca *FunctionCallAgentRunner) interactionInvokeLLM(ctx context.Context) {
	fca.IStreamGenerateQueue.Debug()
	fca.IStreamGenerateQueue = fca.Fork()

	go fca.IStreamGenerateQueue.Listen()

	// log.Info(color.GreenString("========== 下一次 Agent 迭代 ==========="))

	// util.LogCompleteInfo(fca.promptMessages)

	fca.modelCaller.InvokeLLM(ctx, fca.promptMessages, fca.IStreamGenerateQueue, fca.applicationGenerateEntity.ModelConf.Parameters, fca.promptToolMessages, make([]string, 0), fca.applicationGenerateEntity.UserID, nil)
}

func (fca *FunctionCallAgentRunner) handleStreamAgentMessageQueue(ctx context.Context) (*po_agent.MessageAgentThought, error) {
	// isFirstChunk := true
	isOccurredErr := false
	isNormalQuit := false

	agentThought, err := fca.CreateAgentThought(ctx, fca.message.ID, "", "", make(map[string]string, 0), []string{})

	if err != nil {
		return nil, err
	}

	err = fca.agentFlusher.AgentThoughtToStreamResponse(ctx, agentThought.ID)

	if err != nil {
		return nil, err
	}

	resultQueue, finalQueue, errorQueue := fca.GetQueues()

QuitLoop:
	for {
		select {
		case resultMessage := <-resultQueue:
			fca.handleResultChunk(resultMessage)
		case finalMessage, ok := <-finalQueue:
			isNormalQuit = true
			if !ok {
				break QuitLoop
			}
			fca.handleFinalChunk(finalMessage)
		case errorMessage, ok := <-errorQueue:
			isOccurredErr = true
			if !ok {
				break QuitLoop
			}
			if mc, ok := errorMessage.Event.(*biz_entity_base_stream_generator.QueueErrorEvent); ok {
				log.Errorf("found queue error event: %#+v", mc.Err)
				return nil, mc.Err
			}
		}
	}

	if isNormalQuit {
		for len(resultQueue) > 0 {
			resultMessage := <-resultQueue
			fca.handleResultChunk(resultMessage)
		}

		fca.CloseOutNormalExit()
	}

	if isOccurredErr {
		for len(resultQueue) > 0 {
			resultMessage := <-resultQueue
			fca.handleResultChunk(resultMessage)
		}
		fca.CloseOutErr()
	}

	return agentThought, nil
}

func (fca *FunctionCallAgentRunner) handleFinalChunk(message *biz_entity_base_stream_generator.MessageQueueMessage) {
	if mc, ok := message.Event.(*biz_entity_base_stream_generator.QueueMessageEndEvent); ok {
		if fca.checkTools(mc.LLMResult) {
			fca.functionCallState = true
			fca.toolCalls = append(fca.toolCalls, fca.extractToolCalls(mc.LLMResult)...)
			fca.toolCallNames = fca.getToolNames(fca.toolCalls)
			fca.toolCallInputs = fca.getToolInputs(fca.toolCalls)
		}
		fca.fullAssistant = mc.LLMResult.Message.GetContent()
		fca.taskState.LLMResult = mc.LLMResult
	}
}

func (fca *FunctionCallAgentRunner) handleResultChunk(message *biz_entity_base_stream_generator.MessageQueueMessage) error {
	if chunkEvent, ok := message.Event.(*biz_entity_base_stream_generator.QueueAgentMessageEvent); ok {
		deltaText := chunkEvent.Chunk.Delta.Message.Content
		if err := fca.agentFlusher.AgentMessageToStreamResponse(deltaText.(string)); err != nil {
			return err
		}
	}
	return nil
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

func (fca *FunctionCallAgentRunner) checkTools(llmChunk *biz_entity_base_stream_generator.LLMResult) bool {
	return len(llmChunk.Message.ToolCalls) > 0
}

func (fca *FunctionCallAgentRunner) extractToolCalls(llmChunk *biz_entity_base_stream_generator.LLMResult) []*ToolCall {
	var toolCalls []*ToolCall
	for _, prompt := range llmChunk.Message.ToolCalls {
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
func (fca *FunctionCallAgentRunner) getToolInputs(tools []*ToolCall) map[string]string {
	var toolInput = make(map[string]string)

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

func (fca *FunctionCallAgentRunner) CreateAgentThought(ctx context.Context, messageID, message, toolName string, toolInput map[string]string, messageFileIDs []string) (*po_agent.MessageAgentThought, error) {

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

func (fca *FunctionCallAgentRunner) SaveAgentThought(ctx context.Context, agentThought *po_agent.MessageAgentThought, toolName string, toolInput map[string]string, thought string, observation map[string]string, toolInvokeMeta map[string]*po_agent.ToolEngineInvokeMeta, answer string, messageIDs []string, llmUsage *biz_entity_base_stream_generator.LLMUsage) (*po_agent.MessageAgentThought, error) {

	if thought != "" {
		agentThought.Thought = thought
	}

	if toolName != "" {
		agentThought.Tool = toolName
	}

	if len(toolInput) > 0 {
		agentThought.ToolInput = toolInput
	}

	if len(observation) > 0 {
		agentThought.Observation = observation
	}

	if answer != "" {
		agentThought.Answer = answer
	}

	if len(messageIDs) > 0 {
		agentThought.MessageFiles = messageIDs
	}

	if agentThought.ToolMetaStr == nil {
		agentThought.ToolMetaStr = make(map[string]*po_agent.ToolEngineInvokeMeta, 0)
	}

	if len(toolInvokeMeta) > 0 {
		agentThought.ToolMetaStr = toolInvokeMeta
	}

	if agentThought.Tool != "" {
		tools := strings.Split(agentThought.Tool, ";")

		labels := agentThought.ToolLabelsStr

		if labels == nil {
			labels = make(map[string]map[string]any, 0)
		}

		for _, tool := range tools {
			if tool == "" {
				continue
			}

			_, ok := labels[tool]

			if ok {
				continue
			}

			toolRuntime := fca.agentDomain.ToolManager.GetToolByIdentity(tool)
			if toolRuntime != nil {
				labels[tool] = map[string]any{"en_US": toolRuntime.Identity.Label.En_US, "zh_Hans": toolRuntime.Identity.Label.Zh_Hans}
			} else {
				labels[tool] = map[string]any{"en_US": tool, "zh_Hans": tool}
			}
		}
	}

	if err := fca.agentDomain.UpdateAgentThought(ctx, agentThought); err != nil {
		return nil, err
	}

	return agentThought, nil
}
