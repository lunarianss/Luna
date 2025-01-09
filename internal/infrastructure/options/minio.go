package options

type MinioOptions struct {
	Endpoint  string `json:"endpoint" mapstructure:"endpoint"`
	SecretKey string `json:"secret_key" mapstructure:"secret-key"`
	AccessKey string `json:"access_key" mapstructure:"access-key"`
	UseSSL    bool   `json:"use_ssl" mapstructure:"use-ssl"`
	Bucket    string `json:"bucket" mapstructure:"bucket"`
}

// NewJwtOptions creates a JwtOptions object with default parameters.
func NewMinioOptions() *MinioOptions {
	return &MinioOptions{}
}

// Validate verifies flags passed to MySQLOptions.
func (o *MinioOptions) Validate() []error {
	errs := []error{}

	return errs
}
