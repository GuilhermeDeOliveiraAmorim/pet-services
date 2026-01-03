package email

import (
	"errors"
	"fmt"
	"log/slog"
	"net/smtp"
)

// EmailServiceInterface is the common interface for both SMTP and Stub services
// (implements domainUser.EmailService)
type EmailServiceInterface interface {
	SendPasswordResetEmail(to, resetLink string) error
	SendEmailVerification(to, verificationLink string) error
}

// Config holds SMTP configuration.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	FromAddr string
	Logger   *slog.Logger
}

// SMTPService sends emails via SMTP.
type SMTPService struct {
	cfg Config
}

// NewSMTPService creates a new SMTP email service.
func NewSMTPService(cfg Config) *SMTPService {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	return &SMTPService{cfg: cfg}
}

// SendPasswordResetEmail sends a password reset email.
func (s *SMTPService) SendPasswordResetEmail(to, resetLink string) error {
	subject := "Redefinir Senha - Pet Services"
	body := fmt.Sprintf(`
Olá,

Você solicitou uma redefinição de senha. Clique no link abaixo para continuar:

%s

Este link expira em 1 hora.

Se você não solicitou isso, ignore este email.

Atenciosamente,
Pet Services Team
`, resetLink)

	return s.send(to, subject, body)
}

// SendEmailVerification sends an email verification email.
func (s *SMTPService) SendEmailVerification(to, verificationLink string) error {
	subject := "Verificar Email - Pet Services"
	body := fmt.Sprintf(`
Olá,

Obrigado por se cadastrar! Clique no link abaixo para verificar seu email:

%s

Este link expira em 24 horas.

Se você não criou essa conta, ignore este email.

Atenciosamente,
Pet Services Team
`, verificationLink)

	return s.send(to, subject, body)
}

func (s *SMTPService) send(to, subject, body string) error {
	if s.cfg.Host == "" || s.cfg.Port == 0 {
		s.cfg.Logger.Warn("smtp not configured, logging email instead",
			"to", to,
			"subject", subject,
		)
		fmt.Printf("📧 Email to %s\nSubject: %s\n%s\n", to, subject, body)
		return nil
	}

	if to == "" {
		return errors.New("email recipient cannot be empty")
	}

	auth := smtp.PlainAuth("", s.cfg.User, s.cfg.Password, s.cfg.Host)
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.cfg.FromAddr, to, subject, body)

	if err := smtp.SendMail(addr, auth, s.cfg.FromAddr, []string{to}, []byte(msg)); err != nil {
		s.cfg.Logger.Error("failed to send email",
			"to", to,
			"subject", subject,
			"error", err,
		)
		return fmt.Errorf("send email: %w", err)
	}

	s.cfg.Logger.Info("email sent successfully",
		"to", to,
		"subject", subject,
	)
	return nil
}

// StubEmailService logs emails to console (for development).
type StubEmailService struct {
	logger *slog.Logger
}

// NewStubEmailService creates a stub email service for development.
func NewStubEmailService(logger *slog.Logger) *StubEmailService {
	if logger == nil {
		logger = slog.Default()
	}
	return &StubEmailService{logger: logger}
}

// SendPasswordResetEmail logs the reset email.
func (s *StubEmailService) SendPasswordResetEmail(to, resetLink string) error {
	s.logger.Info("📧 [STUB] Password reset email",
		"to", to,
		"link", resetLink,
	)
	fmt.Printf("📧 [STUB] Password reset email sent to %s\n", to)
	fmt.Printf("   Reset link: %s\n", resetLink)
	return nil
}

// SendEmailVerification logs the verification email.
func (s *StubEmailService) SendEmailVerification(to, verificationLink string) error {
	s.logger.Info("📧 [STUB] Email verification",
		"to", to,
		"link", verificationLink,
	)
	fmt.Printf("📧 [STUB] Email verification sent to %s\n", to)
	fmt.Printf("   Verification link: %s\n", verificationLink)
	return nil
}
