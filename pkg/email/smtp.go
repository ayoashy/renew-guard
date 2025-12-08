package email

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"
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
	return s.sendMultipart(to, subject, body, body)
}

func (s *SMTPEmailService) SendHTML(to string, subject string, htmlBody string) error {
	// Create plain text version from HTML (simple strip tags approach)
	plainText := s.htmlToPlainText(htmlBody)
	return s.sendMultipart(to, subject, plainText, htmlBody)
}

// htmlToPlainText converts HTML to plain text (basic implementation)
func (s *SMTPEmailService) htmlToPlainText(html string) string {
	// Remove HTML tags and decode common entities
	text := html
	text = strings.ReplaceAll(text, "<br>", "\n")
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = strings.ReplaceAll(text, "<br />", "\n")
	text = strings.ReplaceAll(text, "</p>", "\n\n")
	text = strings.ReplaceAll(text, "</div>", "\n")
	
	// Simple tag removal
	for strings.Contains(text, "<") && strings.Contains(text, ">") {
		start := strings.Index(text, "<")
		end := strings.Index(text, ">")
		if end > start {
			text = text[:start] + text[end+1:]
		} else {
			break
		}
	}
	
	// Clean up extra whitespace
	lines := strings.Split(text, "\n")
	var cleaned []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	
	return strings.Join(cleaned, "\n")
}

// generateMessageID creates a unique Message-ID header
func (s *SMTPEmailService) generateMessageID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("<%s@%s>", base64.URLEncoding.EncodeToString(b), s.config.SMTPHost)
}

// sendMultipart sends email with both plain text and HTML versions
func (s *SMTPEmailService) sendMultipart(to string, subject string, plainBody string, htmlBody string) error {
	from := s.config.FromEmail
	fromName := s.config.FromName
	
	// Connect with timeout
	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()
	
	// Set deadline for the entire operation
	conn.SetDeadline(time.Now().Add(30 * time.Second))
	
	// Create SMTP client
	client, err := smtp.NewClient(conn, s.config.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()
	
	// Say hello
	if err := client.Hello("localhost"); err != nil {
		return fmt.Errorf("failed to send EHLO: %w", err)
	}
	
	// Start TLS if available
	if ok, _ := client.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: s.config.SMTPHost}
		if err := client.StartTLS(config); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}
	
	// Authenticate
	if err := client.Auth(s.auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	
	// Set sender and recipient
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}
	
	// Get data writer
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	
	// Build multipart message
	boundary := s.generateBoundary()
	
	// Write headers
	headers := textproto.MIMEHeader{}
	headers.Set("From", fmt.Sprintf("%s <%s>", fromName, from))
	headers.Set("To", to)
	headers.Set("Subject", subject)
	headers.Set("MIME-Version", "1.0")
	headers.Set("Content-Type", fmt.Sprintf("multipart/alternative; boundary=\"%s\"", boundary))
	headers.Set("Date", time.Now().Format(time.RFC1123Z))
	headers.Set("Message-ID", s.generateMessageID())
	headers.Set("X-Mailer", "RenewGuard/1.0")
	
	// Write headers to message
	for k, v := range headers {
		fmt.Fprintf(writer, "%s: %s\r\n", k, v[0])
	}
	fmt.Fprintf(writer, "\r\n")
	
	// Write plain text part
	fmt.Fprintf(writer, "--%s\r\n", boundary)
	fmt.Fprintf(writer, "Content-Type: text/plain; charset=UTF-8\r\n")
	fmt.Fprintf(writer, "Content-Transfer-Encoding: 7bit\r\n\r\n")
	fmt.Fprintf(writer, "%s\r\n\r\n", plainBody)
	
	// Write HTML part
	fmt.Fprintf(writer, "--%s\r\n", boundary)
	fmt.Fprintf(writer, "Content-Type: text/html; charset=UTF-8\r\n")
	fmt.Fprintf(writer, "Content-Transfer-Encoding: 7bit\r\n\r\n")
	fmt.Fprintf(writer, "%s\r\n\r\n", htmlBody)
	
	// Close boundary
	fmt.Fprintf(writer, "--%s--\r\n", boundary)
	
	// Close writer
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}
	
	return nil
}

// generateBoundary creates a unique boundary for multipart messages
func (s *SMTPEmailService) generateBoundary() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("===============%s==", base64.StdEncoding.EncodeToString(b))
}
