package mail

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailService interface {
	SendVerificationEmail(to, token string) error
	SendPasswordResetEmail(to, token string) error
}

type SMTPEmailService struct {
	host          string
	port          string
	user          string
	password      string
	from          string
	verifyBaseURL string
	resetBaseURL  string
}

func NewSMTPEmailService(host, port, user, password, from, verifyBaseURL, resetBaseURL string) *SMTPEmailService {
	return &SMTPEmailService{
		host:          host,
		port:          port,
		user:          user,
		password:      password,
		from:          from,
		verifyBaseURL: verifyBaseURL,
		resetBaseURL:  resetBaseURL,
	}
}

func (s *SMTPEmailService) SendVerificationEmail(to, token string) error {
	subject := "Verifique seu email - Pet Services"
	verifyLink := fmt.Sprintf("%s?token=%s", s.verifyBaseURL, token)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<style>
		body { margin: 0; padding: 0; background-color: #f4f7fb; font-family: -apple-system, BlinkMacSystemFont, Segoe UI, Roboto, Arial, sans-serif; color: #1f2937; }
		.wrapper { width: 100%%; padding: 28px 12px; }
		.container { max-width: 620px; margin: 0 auto; background-color: #ffffff; border: 1px solid #e5e7eb; border-radius: 14px; overflow: hidden; }
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #0f766e 0%%, #06b6d4 100%%); }
		.header h1 { margin: 0; font-size: 28px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.content p { margin: 0 0 14px 0; }
		.callout { background-color: #f0fdfa; border: 1px solid #99f6e4; border-radius: 10px; padding: 14px; margin: 18px 0; }
		.button-wrap { text-align: center; margin: 24px 0; }
		.button { display: inline-block; padding: 12px 26px; background-color: #0f766e; color: #ffffff !important; text-decoration: none; border-radius: 999px; font-weight: 600; }
		.link-box { word-break: break-all; background-color: #f9fafb; border: 1px solid #e5e7eb; border-radius: 10px; padding: 12px; font-size: 14px; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Bem-vindo ao Pet Services</h1>
			</div>
			<div class="content">
				<p>Ola,</p>
				<p>Obrigado por se cadastrar. Para ativar sua conta, confirme seu endereco de email clicando no botao abaixo.</p>

				<div class="button-wrap">
					<a href="%s" class="button">Verificar e-mail</a>
				</div>

				<div class="callout">
					<p><strong>Importante:</strong> este link expira em 24 horas.</p>
					<p>Se voce nao criou esta conta, pode ignorar este email com seguranca.</p>
				</div>

				<p>Se o botao nao funcionar, copie e cole este link no navegador:</p>
				<div class="link-box">%s</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, verifyLink, verifyLink)

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendPasswordResetEmail(to, token string) error {
	subject := "Redefinir sua senha - Pet Services"
	resetLink := fmt.Sprintf("%s?token=%s", s.resetBaseURL, token)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<style>
		body { margin: 0; padding: 0; background-color: #f4f7fb; font-family: -apple-system, BlinkMacSystemFont, Segoe UI, Roboto, Arial, sans-serif; color: #1f2937; }
		.wrapper { width: 100%%; padding: 28px 12px; }
		.container { max-width: 620px; margin: 0 auto; background-color: #ffffff; border: 1px solid #e5e7eb; border-radius: 14px; overflow: hidden; }
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #b91c1c 0%%, #ef4444 100%%); }
		.header h1 { margin: 0; font-size: 28px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.content p { margin: 0 0 14px 0; }
		.callout { background-color: #fef2f2; border: 1px solid #fecaca; border-radius: 10px; padding: 14px; margin: 18px 0; }
		.button-wrap { text-align: center; margin: 24px 0; }
		.button { display: inline-block; padding: 12px 26px; background-color: #b91c1c; color: #ffffff !important; text-decoration: none; border-radius: 999px; font-weight: 600; }
		.link-box { word-break: break-all; background-color: #f9fafb; border: 1px solid #e5e7eb; border-radius: 10px; padding: 12px; font-size: 14px; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Redefinicao de senha</h1>
			</div>
			<div class="content">
				<p>Ola,</p>
				<p>Recebemos uma solicitacao para redefinir a senha da sua conta Pet Services.</p>

				<div class="button-wrap">
					<a href="%s" class="button">Redefinir senha</a>
				</div>

				<div class="callout">
					<p><strong>Importante:</strong> este link expira em 1 hora.</p>
					<p>Se voce nao solicitou a redefinicao, ignore este email. Sua conta permanece segura.</p>
				</div>

				<p>Se o botao nao funcionar, copie e cole este link no navegador:</p>
				<div class="link-box">%s</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, resetLink, resetLink)

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) sendEmail(to, subject, body string) error {
	var auth smtp.Auth
	if s.user != "" || s.password != "" {
		auth = smtp.PlainAuth("", s.user, s.password, s.host)
	}
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	message := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=utf-8\r\n\r\n%s",
		s.from, to, subject, body,
	)

	if err := smtp.SendMail(addr, auth, s.from, []string{to}, []byte(message)); err != nil {
		return fmt.Errorf("erro ao enviar email: %w", err)
	}

	return nil
}

func GetEmailServiceFromEnv() EmailService {
	return NewSMTPEmailService(
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
		os.Getenv("SMTP_FROM"),
		os.Getenv("EMAIL_VERIFY_BASE_URL"),
		os.Getenv("PASSWORD_RESET_BASE_URL"),
	)
}
