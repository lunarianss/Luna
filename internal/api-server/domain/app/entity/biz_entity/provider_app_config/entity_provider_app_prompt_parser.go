package biz_entity

import (
	"fmt"
	"regexp"
)

var (
	// Regular expressions to match template variables.
	RegEx                = regexp.MustCompile(`\{\{([a-zA-Z_][a-zA-Z0-9_]{0,29}|#histories#|#query#|#context#)\}\}`)
	WithVariableTplRegex = regexp.MustCompile(`\{\{([a-zA-Z_][a-zA-Z0-9_]{0,29}|#[a-zA-Z0-9_]{1,50}\.[a-zA-Z0-9_\.\-]{1,100}#|#histories#|#query#|#context#)\}\}`)
)

type promptTemplateParser struct {
	Template        string
	withVariableTpl bool
	Regex           *regexp.Regexp
	VariableKeys    []string
}

type IPromptTemplateParser interface {
	Extract() []string
	Format(inputs map[string]interface{}, removeTemplateVariables bool) string
}

func NewPromptTemplateParse(template string, withVariableTpl bool) *promptTemplateParser {

	var regExp *regexp.Regexp

	if withVariableTpl {
		regExp = WithVariableTplRegex
	} else {
		regExp = RegEx
	}

	return &promptTemplateParser{
		Template:        template,
		Regex:           regExp,
		withVariableTpl: withVariableTpl,
	}
}

// extract finds all variables in the template that match the regex.
func (p *promptTemplateParser) Extract() []string {
	matches := p.Regex.FindAllStringSubmatch(p.Template, -1)
	var keys []string
	for _, match := range matches {
		if len(match) > 1 {
			keys = append(keys, match[1])
		}
	}
	return keys
}

// Format replaces template variables in the template with the provided inputs.
func (p *promptTemplateParser) Format(inputs map[string]interface{}, removeTemplateVariables bool) string {
	replacer := func(match string) string {
		key := match[2 : len(match)-2]
		value, found := inputs[key]

		if found {
			if removeTemplateVariables {
				if strValue, ok := value.(string); ok {
					return p.removeTemplateVariablesFromText(strValue)
				}
			}
			return fmt.Sprintf("%v", value)
		}
		return match
	}
	return regexp.MustCompile(`\{\{.*?\}\}`).ReplaceAllStringFunc(p.Template, replacer)
}

func (p *promptTemplateParser) removeTemplateVariablesFromText(text string) string {
	var regex *regexp.Regexp
	if p.withVariableTpl {
		regex = WithVariableTplRegex
	} else {
		regex = RegEx
	}
	return regex.ReplaceAllString(text, `{\1}`)
}
