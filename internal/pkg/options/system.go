package options

type SystemOptions struct {
	AppBasePath                  string `mapstructure:"app-base-path"                  json:"-"`
	SsoEnforcedForSignin         bool   `mapstructure:"sso-enforced-for-signin"        json:"sso-enforced-for-signin"`
	SsoEnforcedForSigninProtocol string `mapstructure:"sso-enforced-for-signin-protocol" json:"sso-enforced-for-signin-protocol"`
	SsoEnforcedForWeb            bool   `mapstructure:"sso-enforced-for-web"           json:"sso-enforced-for-web"`
	SsoEnforcedForWebProtocol    string `mapstructure:"sso-enforced-for-web-protocol"  json:"sso-enforced-for-web-protocol"`
	EnableWebSsoSwitchComponent  bool   `mapstructure:"enable-web-sso-switch-component" json:"enable-web-sso-switch-component"`
	EnableEmailCodeLogin         bool   `mapstructure:"enable-email-code-login"        json:"enable-email-code-login"`
	EnableEmailPasswordLogin     bool   `mapstructure:"enable-email-password-login"    json:"enable-email-password-login"`
	EnableSocialOauthLogin       bool   `mapstructure:"enable-social-oauth-login"      json:"enable-social-oauth-login"`
	IsAllowRegister              bool   `mapstructure:"is-allow-register"              json:"is-allow-register"`
	IsAllowCreateWorkspace       bool   `mapstructure:"is-allow-create-workspace"      json:"is-allow-create-workspace"`
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
