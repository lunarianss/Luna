package biz_entity

import common "github.com/lunarianss/Luna/internal/api-server/domain/provider/entity/biz_entity/common_relation"

type CredentialsType string

const (
	CredentialsTypeSecretInput CredentialsType = "secret-input"
	CredentialsTypeTextInput   CredentialsType = "text-input"
	CredentialsTypeSelect      CredentialsType = "select"
	CredentialsTypeBoolean     CredentialsType = "boolean"
)

type ToolLabelEnum string

const (
	ToolLabelSearch        ToolLabelEnum = "search"
	ToolLabelImage         ToolLabelEnum = "image"
	ToolLabelVideos        ToolLabelEnum = "videos"
	ToolLabelWeather       ToolLabelEnum = "weather"
	ToolLabelFinance       ToolLabelEnum = "finance"
	ToolLabelDesign        ToolLabelEnum = "design"
	ToolLabelTravel        ToolLabelEnum = "travel"
	ToolLabelSocial        ToolLabelEnum = "social"
	ToolLabelNews          ToolLabelEnum = "news"
	ToolLabelMedical       ToolLabelEnum = "medical"
	ToolLabelProductivity  ToolLabelEnum = "productivity"
	ToolLabelEducation     ToolLabelEnum = "education"
	ToolLabelBusiness      ToolLabelEnum = "business"
	ToolLabelEntertainment ToolLabelEnum = "entertainment"
	ToolLabelUtilities     ToolLabelEnum = "utilities"
	ToolLabelOther         ToolLabelEnum = "other"
)

type ToolProviderIdentity struct {
	Author      string             `json:"author" yaml:"author"`
	Name        string             `json:"name" yaml:"name"`
	Description *common.I18nObject `json:"description" yaml:"description"`
	Icon        string             `json:"icon" yaml:"icon"`
	Label       *common.I18nObject `json:"label" yaml:"label"`
	Tags        []ToolLabelEnum    `json:"tags" yaml:"tags"`
}

type ToolCredentialsOption struct {
	Value string             `json:"value" yaml:"value"`
	Label *common.I18nObject `json:"label" yaml:"label"`
}

type ToolProviderCredentials struct {
	Name        string                   `json:"name" yaml:"name"`
	Type        CredentialsType          `json:"type" yaml:"type"`
	Required    bool                     `json:"required" yaml:"required"`
	Default     interface{}              `json:"default" yaml:"default"`
	Options     []*ToolCredentialsOption `json:"options" yaml:"options"`
	Label       *common.I18nObject       `json:"label" yaml:"label"`
	Help        *common.I18nObject       `json:"help" yaml:"help"`
	URL         string                   `json:"url" yaml:"url"`
	Placeholder *common.I18nObject       `json:"placeholder" yaml:"placeholder"`
}
