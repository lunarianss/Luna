// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package email

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/smtp"
	"strings"

	"github.com/lunarianss/Luna/internal/pkg/options"
)

type Mail struct {
	Client          MailClient
	DefaultSendFrom string
}

// MailClient 接口，支持不同的邮件服务
type MailClient interface {
	Send(to string, subject string, html string, from string) error
}

// SMTPClient 使用 SMTP 发送邮件
type SMTPClient struct {
	Server   string
	Port     int
	Username string
	Password string
}

// Send 实现 MailClient 接口
func (s *SMTPClient) Send(to string, subject string, html string, from string) error {
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Server)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"From: " + s.Username + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" + html)

	addr := fmt.Sprintf("%s:%d", s.Server, s.Port)

	return s.SendMailUsingTLS(addr, auth, from, to, msg)
}

// NewMail 初始化邮件发送服务
func NewMail(mailType string, opt *options.EmailOptions) (*Mail, error) {
	m := &Mail{}

	switch mailType {
	case "smtp":
		client := &SMTPClient{
			Server:   opt.SMTPServer,
			Port:     opt.SMTPPort,
			Username: opt.SMTPFromEmail,
			Password: opt.SMTPPassword,
		}
		m.Client = client
		m.DefaultSendFrom = opt.SMTPFromEmail
	default:
		return nil, errors.New("unsupported mail type")
	}

	return m, nil
}

// Send 邮件发送函数
func (m *Mail) Send(to string, subject string, templateFile string, data map[string]interface{}, from string) error {
	if from == "" {
		from = m.DefaultSendFrom
	}
	if m.Client == nil {
		return errors.New("mail client is not initialized")
	}

	html, err := renderTemplate(templateFile, data)
	if err != nil {
		return err
	}

	return m.Client.Send(to, subject, html, from)
}

func (s *SMTPClient) SendMailUsingTLS(addr string, auth smtp.Auth, from string, to string, msg []byte) (err error) {
	c, err := s.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	tos := strings.Split(to, ";")
	for _, addr := range tos {
		if err = c.Rcpt(addr); err != nil {
			fmt.Print(err)
			return err
		}
	}
	w, err := c.Data()

	if err != nil {
		return err
	}
	_, err = w.Write(msg)

	if err != nil {
		return err
	}

	err = w.Close()

	if err != nil {
		return err
	}

	return c.Quit()
}
func (s *SMTPClient) Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// 渲染模板
func renderTemplate(templateFile string, data map[string]interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
