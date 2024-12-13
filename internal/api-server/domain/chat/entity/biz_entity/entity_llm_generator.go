package biz_entity

import "strings"

const RULE_CONFIG_PROMPT_GENERATE_TEMPLATE = `
Here is a task description for which I would like you to create a high-quality prompt template for:
<task_description>
{{TaskDescription}}
</task_description>
Based on task description, please create a well-structured prompt template that another AI could use to consistently complete the task. The prompt template should include:
- Descriptive variable names surrounded by {{ }} (two curly brackets) to indicate where the actual values will be substituted in. Choose variable names that clearly indicate the type of value expected. Variable names have to be composed of number, english alphabets and underline and nothing else. 
- Clear instructions for the AI that will be using this prompt, demarcated with <instructions> tags. The instructions should provide step-by-step directions on how to complete the task using the input variables. Also Specifies in the instructions that the output should not contain any xml tag. 
- Relevant examples if needed to clarify the task further, demarcated with <example> tags. Do not use curly brackets any other than in <instruction> section. 
- Any other relevant sections demarcated with appropriate XML tags like <input>, <output>, etc.
- Use the same language as task description. 
- Output in PLACEHOLDER_XML and start with <instruction>
Please generate the full prompt template and output only the prompt template.
`

func GetRuleConfigPromptGenerateTemplate() string {
	return strings.ReplaceAll(RULE_CONFIG_PROMPT_GENERATE_TEMPLATE, "PLACEHOLDER_XML", "``` xml ```")
}

const RULE_CONFIG_PARAMETER_GENERATE_TEMPLATE = `
I need to extract the following information from the input text. The <information to be extracted> tag specifies the 'type', 'description' and 'required' of the information to be extracted. 
<information to be extracted>
variables name bounded two double curly brackets. Variable name has to be composed of number, english alphabets and underline and nothing else. 
</information to be extracted>

Step 1: Carefully read the input and understand the structure of the expected output.
Step 2: Extract relevant parameters from the provided text based on the name and description of object. 
Step 3: Structure the extracted parameters to JSON object as specified in <structure>.
Step 4: Ensure that the list of variable_names is properly formatted and valid. The output should not contain any XML tags. Output an empty list if there is no valid variable name in input text. 

### Structure
Here is the structure of the expected output, I should always follow the output structure. 
["variable_name_1", "variable_name_2"]

### Input Text
Inside <text></text> XML tags, there is a text that I should extract parameters and convert to a JSON object.
<text>
{{.InputText}}
</text>

### Answer
I should always output a valid list. Output nothing other than the list of variable_name. Output an empty list if there is no variable name in input text.
`

const RULE_CONFIG_STATEMENT_GENERATE_TEMPLATE = `
<instruction>
Step 1: Identify the purpose of the chatbot from the variable {{.TaskDescription}} and infer chatbot's tone  (e.g., friendly, professional, etc.) to add personality traits. 
Step 2: Create a coherent and engaging opening statement.
Step 3: Ensure the output is welcoming and clearly explains what the chatbot is designed to do. Do not include any XML tags in the output.
Please use the same language as the user's input language. If user uses chinese then generate opening statement in chinese,  if user uses english then generate opening statement in english. 
Example Input: 
Provide customer support for an e-commerce website
Example Output: 
Welcome! I'm here to assist you with any questions or issues you might have with your shopping experience. Whether you're looking for product information, need help with your order, or have any other inquiries, feel free to ask. I'm friendly, helpful, and ready to support you in any way I can.
<Task>
Here is the task description: {{.InputText}}

You just need to generate the output
`
