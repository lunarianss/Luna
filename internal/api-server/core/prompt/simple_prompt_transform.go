package prompt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/lunarianss/Luna/internal/api-server/core/app"
	"github.com/lunarianss/Luna/internal/api-server/core/app/app_config"
	"github.com/lunarianss/Luna/internal/api-server/core/prompt/utils"
	"github.com/lunarianss/Luna/internal/api-server/entities/message"
	"github.com/lunarianss/Luna/internal/api-server/model/v1"
	"github.com/lunarianss/Luna/internal/pkg/code"
	"github.com/lunarianss/Luna/pkg/errors"
)

type SimplePromptTransform struct {
}

func (s *SimplePromptTransform) GetPromptStrAndRules(appMode model.AppMode, modelConfig *app.ModelConfigWithCredentialsEntity, prePrompt string, inputs map[string]interface{}, query string, context string, histories string) (string, map[string]interface{}, error) {

	var (
		variables = make(map[string]interface{})
	)

	promptTemplateConfig, err := s.GetPromptTemplate(appMode, modelConfig.Provider, modelConfig.Model, prePrompt, context == "", query == "", histories == "")

	if err != nil {
		return "", nil, err
	}

	for _, key := range promptTemplateConfig["custom_variable_keys"].([]string) {
		variables[key] = inputs[key]
	}

	for _, v := range promptTemplateConfig["special_variable_keys"].([]string) {
		if v == "#context" {
			variables["#context#"] = context
		} else if v == "#query#" {
			variables["#query#"] = query
		} else if v == "#histories#" {
			variables["#histories#"] = histories
		}
	}
	// todo 这里通过 template 去 format
	// promptTemplate, _ := promptTemplateConfig["prompt_template"].(utils.PromptTemplateParser)

	return "", promptTemplateConfig["prompt_rules"].(map[string]interface{}), nil

}
func (s *SimplePromptTransform) GetChatModelPromptMessage(appMode model.AppMode, prePrompt string, inputs map[string]interface{}, query string, context string, files []string, memory any, modelConfig *app.ModelConfigWithCredentialsEntity) ([]*message.PromptMessage, []string, error) {

	var promptMessages []*message.PromptMessage

	prompt, _, err := s.GetPromptStrAndRules(appMode, modelConfig, prePrompt, inputs, query, context, "")

	if err != nil {
		return nil, nil, err
	}

	if prompt != "" && query != "" {
		promptMessages = append(promptMessages, message.NewSystemMessage(prompt))
	}

	if query != "" {
		promptMessages = append(promptMessages, s.GetLastUserMessage(query, nil))
	} else {
		promptMessages = append(promptMessages, s.GetLastUserMessage(query, nil))
	}

	return promptMessages, nil, nil
}

func (s *SimplePromptTransform) GetLastUserMessage(prompt string, files []string) *message.PromptMessage {
	return message.NewUserMessage(prompt)
}

func (s *SimplePromptTransform) GetPrompt(appMode model.AppMode, promptTemplateEntity *app_config.PromptTemplateEntity, inputs map[string]interface{}, query string, files []string, context string, memory any, modelConfig *app.ModelConfigWithCredentialsEntity) ([]*message.PromptMessage, []string, error) {

	var (
		promptMessage []*message.PromptMessage
		stop          []string
		err           error
	)

	modelMode := modelConfig.Mode

	if modelMode == "chat" {
		promptMessage, stop, err = s.GetChatModelPromptMessage(appMode, promptTemplateEntity.SimplePromptTemplate, inputs, query, context, files, nil, modelConfig)

		if err != nil {
			return nil, nil, err
		}
	}

	return promptMessage, stop, nil
}

func (s *SimplePromptTransform) GetPromptTemplate(appMode model.AppMode, provider, model, prePrompt string, hasContext bool, queryInPrompt bool, withMemoryPrompt bool) (map[string]any, error) {

	var (
		customVariableKeys  []string
		specialVariableKeys []string
		prompt              string
		templatePromptRules string
	)

	promptRules, err := s.getPromptRole(appMode, provider, model)

	if err != nil {
		return nil, err
	}

	promptOrders := promptRules["system_prompt_orders"].([]interface{})

	for _, promptOrder := range promptOrders {
		if promptOrder.(string) == "context_prompt" && hasContext {
			prompt += promptRules["context_prompt"].(string)
			specialVariableKeys = append(specialVariableKeys, "#context#")
		} else if promptOrder == "pre_prompt" && prePrompt != "" {
			prompt += prePrompt + "\n"
			templateParser := &utils.PromptTemplateParser{
				Template: prePrompt,
			}
			customVariableKeys = templateParser.Exact()
		} else if promptOrder == "histories_prompt" && withMemoryPrompt {
			prompt += promptRules["histories_prompt"].(string)
			specialVariableKeys = append(specialVariableKeys, "#histories#")
		}
	}

	if queryInPrompt {
		templatePromptRules = promptRules["query_prompt"].(string)
		prompt += templatePromptRules
		specialVariableKeys = append(specialVariableKeys, "#query#")
	}

	return map[string]interface{}{
		"prompt_template": &utils.PromptTemplateParser{
			Template: prompt,
		},
		"custom_variable_keys":  customVariableKeys,
		"special_variable_keys": specialVariableKeys,
		"prompt_rules":          promptRules,
	}, nil

}

func (s *SimplePromptTransform) getPromptRole(appMode model.AppMode, provider, modelName string) (map[string]interface{}, error) {

	var (
		promptRoleMap map[string]interface{}
	)

	promptFileName := s.promptFileName(appMode, provider, modelName)

	_, fullFilePath, _, ok := runtime.Caller(0)

	if !ok {
		return nil, errors.WithCode(code.ErrRunTimeCaller, "Fail to get runtime caller info")
	}

	fileDir := filepath.Dir(fullFilePath)

	roleJsonPath := fmt.Sprintf("%s/prompt_templates/%s", fileDir, promptFileName)

	roleJsonContent, err := os.ReadFile(roleJsonPath)

	if err != nil {
		return nil, errors.WithCode(code.ErrRunTimeCaller, "Read file %s failed, Error: %+v", roleJsonPath, err)
	}

	if err := json.Unmarshal(roleJsonContent, &promptRoleMap); err != nil {
		return nil, errors.WithCode(code.ErrDecodingJSON, err.Error())
	}

	return promptRoleMap, nil

}
func (s *SimplePromptTransform) promptFileName(appMode model.AppMode, provider, modelName string) string {
	if appMode == model.COMPLETION {
		return "common_completion.json"
	} else {
		return "common_chat.json"
	}
}
