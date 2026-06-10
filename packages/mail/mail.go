package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

// ProviderType defines the email provider.
type ProviderType string

const (
	ProviderSMTP     ProviderType = "smtp"
	ProviderConsole  ProviderType = "console"
	ProviderSendGrid ProviderType = "sendgrid"
	ProviderMailgun  ProviderType = "mailgun"
)

// Config stores all configuration parameters for the mail system.
type Config struct {
	Provider     ProviderType
	Host         string // For SMTP
	Port         int    // For SMTP
	Username     string // For SMTP
	Password     string // For SMTP
	SenderEmail  string // Default sender email
	SenderName   string // Default sender name
	APIKey       string // For SendGrid/Mailgun
	Domain       string // For Mailgun (e.g., mail.yourdomain.com)
}

// Mailer is the common interface to send emails in the project.
type Mailer interface {
	Send(ctx context.Context, to []string, subject string, htmlBody string) error
}

// SMTPMailer sends mail via standard SMTP protocol.
type SMTPMailer struct {
	cfg Config
}

func NewSMTPMailer(cfg Config) *SMTPMailer {
	return &SMTPMailer{cfg: cfg}
}

func (m *SMTPMailer) Send(ctx context.Context, to []string, subject string, htmlBody string) error {
	var auth smtp.Auth
	if m.cfg.Username != "" || m.cfg.Password != "" {
		auth = smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	}

	from := fmt.Sprintf("%s <%s>", m.cfg.SenderName, m.cfg.SenderEmail)
	msg := []byte("To: " + strings.Join(to, ",") + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		htmlBody + "\r\n")

	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)
	return smtp.SendMail(addr, auth, m.cfg.SenderEmail, to, msg)
}

// ConsoleMailer only logs emails to console for development or testing.
type ConsoleMailer struct {
	log *zap.Logger
	cfg Config
}

func NewConsoleMailer(cfg Config, log *zap.Logger) *ConsoleMailer {
	return &ConsoleMailer{cfg: cfg, log: log}
}

func (m *ConsoleMailer) Send(ctx context.Context, to []string, subject string, htmlBody string) error {
	m.log.Info("Sending email (Console Mock)",
		zap.Strings("to", to),
		zap.String("from", m.cfg.SenderEmail),
		zap.String("subject", subject),
		zap.String("body", htmlBody),
	)
	return nil
}

// SendGridMailer sends mail via SendGrid API.
type SendGridMailer struct {
	cfg    Config
	client *http.Client
}

func NewSendGridMailer(cfg Config) *SendGridMailer {
	return &SendGridMailer{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (m *SendGridMailer) Send(ctx context.Context, to []string, subject string, htmlBody string) error {
	type emailUser struct {
		Email string `json:"email"`
		Name  string `json:"name,omitempty"`
	}
	type personalization struct {
		To []emailUser `json:"to"`
	}
	type content struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}
	type sendGridRequest struct {
		Personalizations []personalization `json:"personalizations"`
		From             emailUser         `json:"from"`
		Subject          string            `json:"subject"`
		Content          []content         `json:"content"`
	}

	tos := make([]emailUser, len(to))
	for i, email := range to {
		tos[i] = emailUser{Email: email}
	}

	reqBody := sendGridRequest{
		Personalizations: []personalization{{To: tos}},
		From:             emailUser{Email: m.cfg.SenderEmail, Name: m.cfg.SenderName},
		Subject:          subject,
		Content:          []content{{Type: "text/html", Value: htmlBody}},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+m.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("sendgrid returned status: %d", resp.StatusCode)
	}

	return nil
}

// MailgunMailer sends mail via Mailgun HTTP API.
type MailgunMailer struct {
	cfg    Config
	client *http.Client
}

func NewMailgunMailer(cfg Config) *MailgunMailer {
	return &MailgunMailer{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (m *MailgunMailer) Send(ctx context.Context, to []string, subject string, htmlBody string) error {
	apiURL := fmt.Sprintf("https://api.mailgun.net/v3/%s/messages", m.cfg.Domain)

	form := url.Values{}
	form.Set("from", fmt.Sprintf("%s <%s>", m.cfg.SenderName, m.cfg.SenderEmail))
	for _, email := range to {
		form.Add("to", email)
	}
	form.Set("subject", subject)
	form.Set("html", htmlBody)

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth("api", m.cfg.APIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mailgun returned status: %d", resp.StatusCode)
	}

	return nil
}

// NewMailer is a Factory function to create a Mailer matching the configuration.
func NewMailer(cfg Config, log *zap.Logger) (Mailer, error) {
	switch cfg.Provider {
	case ProviderSMTP:
		return NewSMTPMailer(cfg), nil
	case ProviderConsole:
		return NewConsoleMailer(cfg, log), nil
	case ProviderSendGrid:
		return NewSendGridMailer(cfg), nil
	case ProviderMailgun:
		return NewMailgunMailer(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported mail provider: %s", cfg.Provider)
	}
}
