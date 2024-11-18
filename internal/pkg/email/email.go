package email

import (
	"sync"

	"github.com/lunarianss/Luna/internal/pkg/options"
	_email "github.com/lunarianss/Luna/pkg/email"
	"github.com/lunarianss/Luna/pkg/errors"
)

var (
	once     sync.Once
	EmailIns *_email.Mail
)

func GetEmailSMTPIns(opt *options.EmailOptions) (*_email.Mail, error) {

	var err error

	once.Do(func() {
		EmailIns, err = _email.NewMail("smtp", opt)
	})

	if err != nil || EmailIns == nil {
		return nil, errors.WithMessage(err, "failed to get email factory")
	}

	return EmailIns, nil
}
