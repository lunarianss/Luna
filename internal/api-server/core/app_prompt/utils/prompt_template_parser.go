// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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
