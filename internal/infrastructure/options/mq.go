package options

type RocketMQOptions struct {
	Endpoint      []string `json:"endpoint" mapstructure:"endpoint"`
	GroupName     string   `json:"group-name" mapstructure:"group-name"`
	ProducerRetry int      `json:"producer-retry" mapstructure:"producer-retry"`
	Namespace     string   `json:"namespace" mapstructure:"namespace"`
	ConsumerRetry int      `json:"consumer-retry" mapstructure:"consumer-retry"`
	SecretKey     string   `json:"secret-key" mapstructure:"secret-key"`
	AccessKey     string   `json:"access-key" mapstructure:"access-key"`
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
