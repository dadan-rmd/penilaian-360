package mail

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"text/template"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
	"gopkg.in/gomail.v2"
)

func SendEvaluation(to []string, cc []string, evaluatedName, name, deadline string) error {
	fmt.Println(to, cc, evaluatedName, name, deadline)
	// Get the root path
	_, file, _, _ := runtime.Caller(0)
	rootPath := path.Join(file, "../../../../../")
	tmplPath := rootPath + "/assets/template/evaluation.html"
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open email template")
		return err
	}

	// Generate email body
	var body bytes.Buffer
	err = t.Execute(&body, struct {
		EvaluatedName string
		Name          string
		Deadline      string
		URL           string
	}{
		EvaluatedName: evaluatedName,
		Name:          name,
		Deadline:      deadline,
		URL:           os.Getenv("URL_EVALUATION"),
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute email template")
		return err
	}

	// Setup email
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("CONFIG_SENDER_NAME")+" <"+os.Getenv("CONFIG_AUTH_EMAIL")+">")
	m.SetHeader("To", to...)
	if len(cc) > 0 {
		m.SetHeader("Cc", cc...)
	}
	m.SetHeader("Subject", "Verifikasi Email")
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(os.Getenv("CONFIG_SMTP_HOST"), cast.ToInt(os.Getenv("CONFIG_SMTP_PORT")), os.Getenv("CONFIG_AUTH_EMAIL"), os.Getenv("CONFIG_AUTH_PASSWORD"))

	// Send email
	if err := d.DialAndSend(m); err != nil {
		log.Error().Err(err).Msg("Failed to send email")
		return err
	}

	log.Info().Msg("Email sent successfully!")
	return nil
}
