package email

import (
	"bytes"
	"fmt"
	"html/template"
	"sync"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain        string
	Host          string
	Port          int
	Username      string
	Password      string
	Encryption    string
	FromAddress   string
	FromName      string
	Wait          *sync.WaitGroup
	MailerChannel chan Message
	ErrorChannel  chan error
	DoneChannel   chan bool
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
	Template    string
}

type VerifyMessageRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Data    string `json:"data"`
}

func CreateMail(wait *sync.WaitGroup) *Mail {
	return &Mail{
		ErrorChannel:  make(chan error),
		MailerChannel: make(chan Message, 100),
		DoneChannel:   make(chan bool),
		Domain:        "localhost",
		Host:          "localhost",
		Port:          1025,
		Encryption:    "none",
		FromName:      "info",
		FromAddress:   "info@company.com",
		Wait:          wait,
	}
}

func (email *Mail) ListenForMail() {
	for {
		select {
		case message := <-email.MailerChannel:
			go email.SendEmail(message, email.ErrorChannel)
		case <-email.ErrorChannel:
			continue
		case <-email.DoneChannel:
			return
		}
	}
}

func (email *Mail) SendEmail(msg Message, errorChannel chan error) {
	defer email.Wait.Done()

	if msg.Template == "" {
		msg.Template = "mail"
	}

	if msg.From == "" {
		msg.From = email.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = email.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := email.buildHTMLMessage(&msg)
	if err != nil {
		errorChannel <- err
	}

	plainTextMessage, err := email.buildPlainTextMessage(&msg)
	if err != nil {
		errorChannel <- err
	}

	server := mail.NewSMTPClient()
	server.Host = email.Host
	server.Port = email.Port
	server.Username = email.Username
	server.Password = email.Password
	server.Encryption = getEncryptionStandard(email.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		errorChannel <- err
	}

	emailToSend := mail.NewMSG()
	emailToSend.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)

	emailToSend.SetBody(mail.TextPlain, plainTextMessage)
	emailToSend.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, attach := range msg.Attachments {
			emailToSend.AddAttachment(attach)
		}
	}

	if err = emailToSend.Send(smtpClient); err != nil {
		errorChannel <- err
	}
}

func (mail *Mail) buildHTMLMessage(msg *Message) (string, error) {
	templateMessage := fmt.Sprintf("./cmd/web/templates/%s.html.gohtml", msg.Template)
	tmp, err := template.New("email-html").ParseFiles(templateMessage)
	if err != nil {
		return "", err
	}

	var templateBuffer bytes.Buffer
	if err = tmp.ExecuteTemplate(&templateBuffer, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := templateBuffer.String()
	formattedMessage, err = inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

func (mail *Mail) buildPlainTextMessage(msg *Message) (string, error) {
	templateMessage := fmt.Sprintf("./cmd/web/templates/%s.plain.gohtml", msg.Template)
	tmp, err := template.New("email-html").ParseFiles(templateMessage)
	if err != nil {
		return "", err
	}

	var templateBuffer bytes.Buffer
	if err = tmp.ExecuteTemplate(&templateBuffer, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := templateBuffer.String()

	return formattedMessage, nil
}

func getEncryptionStandard(encryption string) mail.Encryption {
	switch encryption {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}

func inlineCSS(formattedMessage string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(formattedMessage, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}
