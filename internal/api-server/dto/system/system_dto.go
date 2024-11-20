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
