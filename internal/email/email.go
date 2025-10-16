package email

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
	"os"
)

type EmailService struct {
	SMTPHost string
	SMTPPort string
	From     string
	Password string
}

func NewEmailService() *EmailService {
	return &EmailService{
		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPPort: os.Getenv("SMTP_PORT"),
		From:     os.Getenv("SMTP_FROM"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
}

// GenerateVerificationCode generates a random 6-digit code
func GenerateVerificationCode() (string, error) {
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// SendVerificationCode sends a verification code to the specified email
func (e *EmailService) SendVerificationCode(to, code string) error {
	auth := smtp.PlainAuth("", e.From, e.Password, e.SMTPHost)

	subject := "Subject: Fastnet VPN - Verification Code\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
				<h2 style="color: #333;">Fastnet VPN Verification</h2>
				<p>Your verification code is:</p>
				<div style="background-color: #f4f4f4; padding: 15px; text-align: center; font-size: 32px; font-weight: bold; letter-spacing: 5px; margin: 20px 0;">
					%s
				</div>
				<p>This code will expire in 10 minutes.</p>
				<p>If you didn't request this code, please ignore this email.</p>
			</div>
		</body>
		</html>
	`, code)

	message := []byte(subject + mime + body)

	addr := fmt.Sprintf("%s:%s", e.SMTPHost, e.SMTPPort)
	err := smtp.SendMail(addr, auth, e.From, []string{to}, message)
	if err != nil {
		return err
	}

	return nil
}