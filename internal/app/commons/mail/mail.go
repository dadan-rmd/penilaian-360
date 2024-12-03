package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"path"
	"runtime"
	"text/template"

	"github.com/rs/zerolog/log"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = "587"
const CONFIG_SENDER_NAME = "KBR <kbrprimedev@gmail.com>"
const CONFIG_AUTH_EMAIL = "kbrprimedev@gmail.com"
const CONFIG_AUTH_PASSWORD = "pssamdqtneaxwfwt"

func SendVerifyMail(to []string, name, id string) (err error) {
	auth := smtp.PlainAuth("", CONFIG_AUTH_EMAIL, CONFIG_AUTH_PASSWORD, CONFIG_SMTP_HOST)

	_, file, _, _ := runtime.Caller(0)
	rootPath := path.Join(file, "../../../../../")
	t, err := template.ParseFiles(rootPath + "/assets/template/verifikasi.html")
	if err != nil {
		log.Error().Msg("open template verifikasi email")
		return
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Verifikasi Email \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Username string
		Url      string
	}{
		Username: name,
		Url:      os.Getenv("URL_VERIFIKASI") + id,
	})

	// Sending email.
	err = smtp.SendMail(CONFIG_SMTP_HOST+":"+CONFIG_SMTP_PORT, auth, CONFIG_AUTH_EMAIL, to, body.Bytes())
	if err != nil {
		log.Error().Msg("error SendMail : " + err.Error())
		return
	}
	return
}

func SendOtpMail(to []string, name, otp string) (err error) {
	auth := smtp.PlainAuth("", CONFIG_AUTH_EMAIL, CONFIG_AUTH_PASSWORD, CONFIG_SMTP_HOST)

	_, file, _, _ := runtime.Caller(0)
	rootPath := path.Join(file, "../../../../../")
	t, err := template.ParseFiles(rootPath + "/assets/template/otp.html")
	if err != nil {
		log.Error().Msg("open template otp email")
		return
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Kode Verifikasi OTP \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Username string
		Otp      string
	}{
		Username: name,
		Otp:      otp,
	})

	// Sending email.
	err = smtp.SendMail(CONFIG_SMTP_HOST+":"+CONFIG_SMTP_PORT, auth, CONFIG_AUTH_EMAIL, to, body.Bytes())
	if err != nil {
		log.Error().Msg("error SendMail : " + err.Error())
		return
	}
	return
}

func SendReportComment(to []string, detailEpisode, urlComment, comment string) (err error) {
	auth := smtp.PlainAuth("", CONFIG_AUTH_EMAIL, CONFIG_AUTH_PASSWORD, CONFIG_SMTP_HOST)

	_, file, _, _ := runtime.Caller(0)
	rootPath := path.Join(file, "../../../../../")
	t, err := template.ParseFiles(rootPath + "/assets/template/report-comment.html")
	if err != nil {
		log.Error().Msg("open template report-comment")
		return
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Konfirmasi Pelanggaran Aturan Penggunaan - Tindakan Diperlukan\n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		UrlEpisode string
		Url        string
		Comment    string
	}{
		UrlEpisode: os.Getenv("HOST_KBRPRIME") + "/podcast/" + detailEpisode,
		Url:        os.Getenv("HOST_API_ANALYTICS_KBRPRIME") + urlComment,
		Comment:    comment,
	})

	// Sending email.
	err = smtp.SendMail(CONFIG_SMTP_HOST+":"+CONFIG_SMTP_PORT, auth, CONFIG_AUTH_EMAIL, to, body.Bytes())
	if err != nil {
		log.Error().Msg("error SendMail : " + err.Error())
		return
	}
	return
}
