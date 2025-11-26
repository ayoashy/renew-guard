package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type SMTPEmailService struct {
	config EmailConfig
	auth   smtp.Auth
}

func NewSMTPEmailService(config EmailConfig) *SMTPEmailService {
	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)
	return &SMTPEmailService{
		config: config,
		auth:   auth,
	}
}

func (s *SMTPEmailService) Send(to string, subject string, body string) error {
	return s.send(to, subject, body, "text/plain")
}

func (s *SMTPEmailService) SendHTML(to string, subject string, htmlBody string) error {
	return s.send(to, subject, htmlBody, "text/html")
}

func (s *SMTPEmailService) send(to string, subject string, body string, contentType string) error {
	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = fmt.Sprintf("%s; charset=UTF-8", contentType)

	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body

	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)

	// Create TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.config.SMTPHost,
	}

	// Connect to SMTP server
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		// Try without TLS
		return s.sendWithoutTLS(to, []byte(message))
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// Authenticate
	if err := client.Auth(s.auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Set sender
	if err := client.Mail(s.config.FromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipient
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send message
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}

func (s *SMTPEmailService) sendWithoutTLS(to string, message []byte) error {
	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)
	return smtp.SendMail(addr, s.auth, s.config.FromEmail, []string{to}, message)
}
