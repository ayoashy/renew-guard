package email

// EmailService defines the interface for sending emails
type EmailService interface {
	Send(to string, subject string, body string) error
	SendHTML(to string, subject string, htmlBody string) error
}

// EmailConfig holds configuration for email service
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// NewEmailService creates a new email service based on configuration
func NewEmailService(config EmailConfig) EmailService {
	return NewSMTPEmailService(config)
}
