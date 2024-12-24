package options

type WeaviateOptions struct {
	Endpoint string `json:"endpoint" mapstructure:"endpoint"`
	Schema   string `json:"schema" mapstructure:"schema"`
	ApiKey   string `json:"api-key" mapstructure:"api-key"`
}

// NewJwtOptions creates a JwtOptions object with default parameters.
func NewWeaviateOptions() *WeaviateOptions {
	return &WeaviateOptions{}
}

// Validate verifies flags passed to MySQLOptions.
func (o *WeaviateOptions) Validate() []error {
	errs := []error{}

	return errs
}
