package options

type SystemOptions struct {
	AppBasePath                  string `mapstructure:"app-base-path"                  json:"-"`
	SsoEnforcedForSignin         bool   `mapstructure:"sso-enforced-for-signin"        json:"sso_enforced_for_signin"`
	SsoEnforcedForSigninProtocol string `mapstructure:"sso-enforced-for-signin-protocol" json:"sso_enforced_for_signin_protocol"`
	SsoEnforcedForWeb            bool   `mapstructure:"sso-enforced-for-web"           json:"sso_enforced_for_web"`
	SsoEnforcedForWebProtocol    string `mapstructure:"sso-enforced-for-web-protocol"  json:"sso_enforced_for_web_protocol"`
	EnableWebSsoSwitchComponent  bool   `mapstructure:"enable-web-sso-switch-component" json:"enable_web_sso_switch_component"`
	EnableEmailCodeLogin         bool   `mapstructure:"enable-email-code-login"        json:"enable_email_code_login"`
	EnableEmailPasswordLogin     bool   `mapstructure:"enable-email-password-login"    json:"enable_email_password_login"`
	EnableSocialOauthLogin       bool   `mapstructure:"enable-social-oauth-login"      json:"enable_social_oauth_login"`
	IsAllowRegister              bool   `mapstructure:"is-allow-register"              json:"is_allow_register"`
	IsAllowCreateWorkspace       bool   `mapstructure:"is-allow-create-workspace"      json:"is_allow_create_workspace"`
}

// NewJwtOptions creates a JwtOptions object with default parameters.
func NewSystemOptions() *SystemOptions {

	return &SystemOptions{}
}

// Validate checks validation of ServerRunOptions.
func (s *SystemOptions) Validate() []error {
	errors := []error{}

	return errors
}
