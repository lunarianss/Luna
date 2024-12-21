package options

type RocketMQOptions struct {
	Endpoint      []string `json:"endpoint" mapstructure:"endpoint"`
	GroupName     string   `json:"group-name" mapstructure:"group-name"`
	ProducerRetry int      `json:"producer-retry" mapstructure:"producer-retry"`
}

// NewJwtOptions creates a JwtOptions object with default parameters.
func NewRocketMQOptions() *RocketMQOptions {
	return &RocketMQOptions{}
}

// Validate verifies flags passed to MySQLOptions.
func (o *RocketMQOptions) Validate() []error {
	errs := []error{}

	return errs
}
