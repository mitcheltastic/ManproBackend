package email

import (
	"fmt"
	"log"
	"time"
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

// generateHTMLBody creates a professional, branded HTML email body for Intertwine.
func generateHTMLBody(code string) string {
	// Branding Colors:
	// Creamy White: #FDF8F0 (Main BG)
	// Sunny Yellow: #F9D748 (Accent)
	// Deep Brown: #5D4D37 (Text)
	// Soft Terracotta: #E49889 (Button)

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Intertwine Password Reset</title>
</head>
<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #FDF8F0;">
    <div style="max-width: 600px; margin: 20px auto; background-color: #FFFFFF; border-radius: 12px; border: 1px solid #E49889; overflow: hidden;">
        
        <!-- Header: Logo & Branding -->
        <div style="background-color: #F9D748; padding: 20px 30px; text-align: center; border-bottom: 4px solid #E49889;">
            <h1 style="margin: 0; color: #5D4D37; font-size: 24px; font-weight: bold;">
                Intertwine <span style="font-size: 18px;">&hearts;</span>
            </h1>
            <p style="margin: 5px 0 0; color: #5D4D37; font-size: 14px;">College Matchmaking, Simplified</p>
        </div>

        <!-- Body Content -->
        <div style="padding: 30px;">
            <h2 style="color: #5D4D37; font-size: 20px; margin-top: 0;">Password Reset Request</h2>
            
            <p style="color: #5D4D37; font-size: 16px; line-height: 1.5;">
                We received a request to reset the password for your Intertwine account.
                To ensure your security and allow you to regain access, please use the 6-digit verification code below.
            </p>

            <!-- Verification Code Block -->
            <div style="text-align: center; margin: 30px 0; padding: 15px; background-color: #FDF8F0; border: 1px solid #F9D748; border-radius: 8px;">
                <p style="margin: 0; color: #5D4D37; font-size: 14px;">Your Verification Code:</p>
                <p style="margin: 10px 0 0; color: #E49889; font-size: 32px; font-weight: bold; letter-spacing: 5px;">
                    %s
                </p>
            </div>
            
            <p style="color: #5D4D37; font-size: 16px; line-height: 1.5;">
                This code **expires in 15 minutes**. If you did not request a password reset, please ignore this email. Your password will remain unchanged.
            </p>

        </div>

        <!-- Footer -->
        <div style="padding: 20px 30px; background-color: #E49889; text-align: center; color: #FFFFFF; font-size: 12px;">
            <p style="margin: 0;">&copy; %d Intertwine | Find your perfect match.</p>
            <p style="margin: 5px 0 0;">Need help? Contact support@intertwine.com</p>
        </div>

    </div>
</body>
</html>
`, code, time.Now().Year())
}

// SendPasswordResetCode formats and sends the 6-digit code via email.
func (s *SMTPSender) SendPasswordResetCode(toEmail string, code string) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port) 
	
	// Authentication setup
	auth := smtp.PlainAuth("", s.username, s.password, s.host) 

	// Email content (MIME headers + body)
	// FIX: Use s.from (the formatted email address) in the Mail command.
	// FIX: Use the s.from variable in the 'From' header for the mail content.
	subject := fmt.Sprintf("Subject: Intertwine: Your Password Reset Code\r\nFrom: %s\r\n", s.from) 
	
	// CRITICAL CHANGE: Change Content-Type to text/html
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	
	// Use the HTML generation function
	htmlBody := generateHTMLBody(code)
	
	msg := []byte(subject + mime + htmlBody)
	
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
	// IMPORTANT: client.Mail() must use the exact authenticated email address, NOT the formatted string.
	// We use s.username here, which is the raw jokiinemail@gmail.com
	if err = client.Mail(s.username); err != nil {
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