package utils

type PromptTemplateParser struct {
	Template         string
	WithVariableTmpl bool
	Regex            string
	VariableKeys     []string
}

func (p *PromptTemplateParser) Exact() []string {

	return nil
}
