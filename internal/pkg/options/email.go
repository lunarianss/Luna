package options

import "github.com/spf13/pflag"

type EmailOptions struct {
	SMTPServer    string `json:"SMTP-server" mapstructure:"SMTP-server"`
	SMTPPort      int    `json:"SMTP-port" mapstructure:"SMTP-port"`
	SMTPUsername  string `json:"SMTP-username" mapstructure:"SMTP-username"`
	SMTPPassword  string `json:"SMTP-password" mapstructure:"SMTP-password"`
	SMTPFromEmail string `json:"SMTP-fromEmail" mapstructure:"SMTP-from-email"`
	TemplateDir   string `json:"template-dir" mapstructure:"template-dir"`
}

// Validate verifies flags passed to EmailOptions.
func (o *EmailOptions) Validate() []error {
	errs := []error{}

	return errs
}

// NewEmailOptions create a `zero` value instance.
func NewEmailOptions() *EmailOptions {
	return &EmailOptions{}
}

func (o *EmailOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.SMTPServer, "smtp.server", o.SMTPServer, "Hostname of your smtp server.")
	fs.StringVar(&o.SMTPUsername, "smtp.username", o.SMTPUsername, "Username of your smtp server.")
	fs.StringVar(&o.SMTPPassword, "smtp.password", o.SMTPPassword, "Password of your smtp server.")
	fs.IntVar(&o.SMTPPort, "smtp.port", o.SMTPPort, "The port the smtp server is listening on.")
}
