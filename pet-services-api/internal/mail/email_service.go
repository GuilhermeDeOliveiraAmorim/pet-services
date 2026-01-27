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
	<style>
		body { font-family: Arial, sans-serif; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #007bff; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; border: 1px solid #ddd; }
		.button { display: inline-block; margin-top: 20px; padding: 10px 20px; background-color: #007bff; color: white; text-decoration: none; border-radius: 5px; }
		.footer { text-align: center; margin-top: 20px; font-size: 12px; color: #666; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Bem-vindo ao Pet Services!</h1>
		</div>
		<div class="content">
			<p>Olá,</p>
			<p>Obrigado por se cadastrar no Pet Services. Para completar seu registro, você precisa verificar seu endereço de email.</p>
			<p>Clique no botão abaixo para verificar seu email:</p>
			<a href="%s" class="button">Verificar Email</a>
			<p>Ou copie e cole este link no seu navegador:</p>
			<p>%s</p>
			<p>Este link expira em 24 horas.</p>
			<p>Se você não criou esta conta, por favor ignore este email.</p>
		</div>
		<div class="footer">
			<p>&copy; 2026 Pet Services. Todos os direitos reservados.</p>
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
	<style>
		body { font-family: Arial, sans-serif; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #dc3545; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; border: 1px solid #ddd; }
		.button { display: inline-block; margin-top: 20px; padding: 10px 20px; background-color: #dc3545; color: white; text-decoration: none; border-radius: 5px; }
		.footer { text-align: center; margin-top: 20px; font-size: 12px; color: #666; }
		.warning { color: #dc3545; font-weight: bold; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Redefinição de Senha</h1>
		</div>
		<div class="content">
			<p>Olá,</p>
			<p>Recebemos uma solicitação para redefinir a senha da sua conta Pet Services.</p>
			<p>Clique no botão abaixo para criar uma nova senha:</p>
			<a href="%s" class="button">Redefinir Senha</a>
			<p>Ou copie e cole este link no seu navegador:</p>
			<p>%s</p>
			<p><span class="warning">Importante:</span> Este link expira em 1 hora. Se você não solicitou uma redefinição de senha, por favor ignore este email.</p>
			<p>Sua conta está segura enquanto você não clicar no link acima.</p>
		</div>
		<div class="footer">
			<p>&copy; 2026 Pet Services. Todos os direitos reservados.</p>
		</div>
	</div>
</body>
</html>
`, resetLink, resetLink)

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.user, s.password, s.host)
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
