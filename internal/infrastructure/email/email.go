// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package email

import (
	"sync"

	_email "github.com/lunarianss/Luna/infrastructure/email"
	"github.com/lunarianss/Luna/infrastructure/errors"
	"github.com/lunarianss/Luna/internal/infrastructure/options"
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
