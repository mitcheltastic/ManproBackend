package email

import (
	"fmt"
	"log"
	"net/smtp"
	"crypto/tls"
)

// Sender defines the interface for sending email notifications.
type Sender interface {
	SendPasswordResetCode(toEmail string, code string) error
}

// SMTPSender is the concrete implementation of the Sender interface using SMTP.
type SMTPSender struct {
	host     string
	port     string
	username string
	password string
	from     string
}

// NewSMTPSender creates a new SMTPSender instance.
func NewSMTPSender(host, port, user, pass, from string) Sender {
	return &SMTPSender{
		host:     host,
		port:     port,
		username: user,
		password: pass,
		from:     from,
	}
}

// SendPasswordResetCode formats and sends the 6-digit code via email.
func (s *SMTPSender) SendPasswordResetCode(toEmail string, code string) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port) 
	
	// Authentication setup
	// CRITICAL FIX: Reverting the host parameter back to s.host (smtp.gmail.com).
	// This error means the server requires the host name for the final handshake validation.
	auth := smtp.PlainAuth("", s.username, s.password, s.host) 

	// Email content (MIME headers + body)
	subject := "Subject: Your Password Reset Code\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/plain; charset=\"UTF-8\";\r\n\r\n"
	body := fmt.Sprintf("Your 6-digit password reset code is: %s. This code will expire in 15 minutes. Please use it immediately to reset your password.", code)
	
	msg := []byte(subject + mime + body)
	
	// 1. Establish the UNENCRYPTED connection
	client, err := smtp.Dial(addr)
	if err != nil {
		log.Printf("FATAL EMAIL ERROR: Failed to dial SMTP server (Check port 587 access): %v", err)
		return fmt.Errorf("failed to dial SMTP server: %w", err)
	}

	defer client.Close()

	// 2. Upgrade the connection to TLS/SSL (STARTTLS)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName: s.host,
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		log.Printf("FATAL EMAIL ERROR: Failed to initiate STARTTLS handshake: %v", err)
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// 3. Authenticate the client using App Password
	if err = client.Auth(auth); err != nil {
		log.Printf("FATAL EMAIL ERROR: SMTP authentication failed. Check App Password/User: %v", err)
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}
	
	// 4. Send the email (using low-level commands)
	if err = client.Mail(s.from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err = client.Rcpt(toEmail); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}
	
	// 5. Quit the session
	if err = client.Quit(); err != nil {
		log.Printf("Warning: Failed to quit SMTP session: %v", err)
	}

	log.Printf("SUCCESS: Password reset code sent to %s.", toEmail)
	return nil
}