package biz_entity

type ToolStaticConfiguration struct {
	Identity            *ToolIdentity    `json:"identity,omitempty" yaml:"identity,omitempty"`       // Optional identity of the tool
	Parameters          []*ToolParameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`   // Optional list of parameters
	Description         *ToolDescription `json:"description,omitempty" yaml:"description,omitempty"` // Optional description of the tool
	IsTeamAuthorization bool             `json:"is_team_authorization" yaml:"is_team_authorization"` // Whether team authorization is enabled
}
