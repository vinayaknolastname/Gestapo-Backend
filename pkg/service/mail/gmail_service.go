package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"net/smtp"
	"strconv"

	"github.com/akmal4410/gestapo/internal/config"
	"github.com/akmal4410/gestapo/pkg/service/cache"
	"github.com/jordan-wright/email"
	"github.com/redis/go-redis/v9"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type GmailService struct {
	name                string
	senderEmailAdrress  string
	senderEmailPassword string
}

func NewGmailService(email *config.Email) EmailService {
	return &GmailService{
		name:                email.SenderName,
		senderEmailAdrress:  email.SenderAddress,
		senderEmailPassword: email.SenderPassword,
	}
}

func parseTemplate(email, templateType, content string) (int, *bytes.Buffer, error) {
	bodyTpl, err := template.ParseFiles(fmt.Sprintf("web/templates/%s.html", templateType))
	if err != nil {
		return 0, nil, err
	}
	otp := rand.Intn(900000) + 100000
	var body bytes.Buffer
	data := map[string]string{"otp": strconv.Itoa(otp), "email": email, "content": content}
	if err := bodyTpl.Execute(&body, data); err != nil {
		return 0, nil, err
	}
	return otp, &body, nil
}

func (sender *GmailService) SendOTP(to, subject, content string, redisCache cache.Cache) error {
	email := email.NewEmail()

	otp, htmlContent, err := parseTemplate(to, "email", content)
	if err != nil {
		return err
	}

	email.From = fmt.Sprintf("%s <%s>", sender.name, sender.senderEmailAdrress)
	email.To = []string{to}
	// email.Cc = cc
	// email.Bcc = bcc
	email.Subject = subject
	email.HTML = []byte(htmlContent.Bytes())

	// for _, f := range attachFiles {
	// 	_, err := email.AttachFile(f)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to attach file %s : %w", f, err)
	// 	}
	// }

	if err := redisCache.Delete(to); err != nil {
		if err != redis.Nil {
			return err
		}
	}

	if err := redisCache.Set(to, strconv.Itoa(otp)); err != nil {
		return err
	}
	smtpAuth := smtp.PlainAuth("", sender.senderEmailAdrress, sender.senderEmailPassword, smtpAuthAddress)
	return email.Send(smtpServerAddress, smtpAuth)
}

func (sender *GmailService) VerfiyOTP(user, otp string, redis cache.Cache) (bool, error) {
	cachedOtp, err := redis.Get(user)
	if err != nil {
		if err.Error() == "nil" {
			return false, nil
		}
		return false, err
	}
	// if the otp is valid, then delete it from the redis
	if cachedOtp == otp {
		if err := redis.Delete(user); err != nil {
			return false, err
		}
	}
	return cachedOtp == otp, nil
}
