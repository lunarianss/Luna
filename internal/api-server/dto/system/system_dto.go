package dto

type SystemConfigResponse struct {
	SsoEnforcedForSignin         bool   `json:"sso_enforced_for_signin"`
	SsoEnforcedForSigninProtocol string `json:"sso_enforced_for_signin_protocol"`
	SsoEnforcedForWeb            bool   `json:"sso_enforced_for_web"`
	SsoEnforcedForWebProtocol    string `json:"sso_enforced_for_web_protocol"`
	EnableWebSsoSwitchComponent  bool   `json:"enable_web_sso_switch_component"`
	EnableEmailCodeLogin         bool   `json:"enable_email_code_login"`
	EnableEmailPasswordLogin     bool   `json:"enable_email_password_login"`
	EnableSocialOauthLogin       bool   `json:"enable_social_oauth_login"`
	IsAllowRegister              bool   `json:"is_allow_register"`
	IsAllowCreateWorkspace       bool   `json:"is_allow_create_workspace"`
}

// SubscriptionModel 结构体
type SubscriptionModel struct {
	Plan     string `json:"plan"`
	Interval string `json:"interval"`
}

// BillingModel 结构体
type BillingModel struct {
	Enabled      bool              `json:"enabled"`
	Subscription SubscriptionModel `json:"subscription"`
}

// LimitationModel 结构体
type LimitationModel struct {
	Size  int `json:"size"`
	Limit int `json:"limit"`
}

// FeatureModel 结构体
type FeatureModel struct {
	Billing                   BillingModel    `json:"billing"`
	Members                   LimitationModel `json:"members"`
	Apps                      LimitationModel `json:"apps"`
	VectorSpace               LimitationModel `json:"vector_space"`
	AnnotationQuotaLimit      LimitationModel `json:"annotation_quota_limit"`
	DocumentsUploadQuota      LimitationModel `json:"documents_upload_quota"`
	DocsProcessing            string          `json:"docs_processing"`
	CanReplaceLogo            bool            `json:"can_replace_logo"`
	ModelLoadBalancingEnabled bool            `json:"model_load_balancing_enabled"`
	DatasetOperatorEnabled    bool            `json:"dataset_operator_enabled"`
}

// 初始化 FeatureModel 的函数，相当于 Python 中的默认值
func NewFeatureModel() *FeatureModel {
	return &FeatureModel{
		Billing: BillingModel{
			Enabled: false,
			Subscription: SubscriptionModel{
				Plan:     "sandbox",
				Interval: "",
			},
		},
		Members:                   LimitationModel{Size: 1, Limit: 1},
		Apps:                      LimitationModel{Size: 1, Limit: 10},
		VectorSpace:               LimitationModel{Size: 1, Limit: 5},
		AnnotationQuotaLimit:      LimitationModel{Size: 1, Limit: 10},
		DocumentsUploadQuota:      LimitationModel{Size: 1, Limit: 50},
		DocsProcessing:            "standard",
		CanReplaceLogo:            false,
		ModelLoadBalancingEnabled: false,
		DatasetOperatorEnabled:    false,
	}
}
