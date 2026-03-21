package mail

import (
	"fmt"
	"html"
	"net/smtp"
	"os"
)

type EmailService interface {
	SendVerificationEmail(to, token string) error
	SendWelcomeAfterVerificationEmail(to, name string) error
	SendLoginBlockedAlertEmail(to, name, reason string) error
	SendPasswordResetEmail(to, token string) error
	SendPasswordChangedAlertEmail(to, name string) error
	SendPasswordResetSuccessEmail(to, name string) error
	SendAccountDeactivatedEmail(to, name string) error
	SendAccountReactivatedEmail(to, name string) error
	SendAccountDeletedEmail(to, name string) error
	SendRequestCreatedEmail(to, providerName, ownerName, petName, serviceName, requestID string) error
	SendRequestCreatedOwnerConfirmationEmail(to, ownerName, providerName, petName, serviceName, requestID string) error
	SendRequestAcceptedEmail(to, ownerName, providerName, petName, requestID string) error
	SendRequestRejectedEmail(to, ownerName, providerName, petName, reason, requestID string) error
	SendRequestCompletedEmail(to, ownerName, providerName, petName, requestID string) error
	SendReviewReminderEmail(to, ownerName, providerName, petName, requestID string) error
	SendReviewReceivedEmail(to, providerName, ownerName string, rating float64, comment string) error
	SendAdoptionGuardianProfileApprovedEmail(to, name string) error
	SendAdoptionGuardianProfileRejectedEmail(to, name, reason string) error
	SendAdoptionApplicationSubmittedEmail(to, applicantName, petName string) error
	SendAdoptionApplicationReceivedGuardianEmail(to, guardianName, applicantName, petName string) error
	SendAdoptionApplicationApprovedEmail(to, applicantName, petName, guardianContact string) error
	SendAdoptionApplicationRejectedEmail(to, applicantName, petName string) error
	SendAdoptionApplicationWithdrawnGuardianEmail(to, guardianName, applicantName, petName string) error
	SendPetAdoptedGuardianEmail(to, guardianName, petName string) error
	SendPetAdoptedApplicantEmail(to, applicantName, petName string) error
	SendPetAdoptedRejectedApplicantsEmail(to, petName string) error
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

func (s *SMTPEmailService) SendWelcomeAfterVerificationEmail(to, name string) error {
	subject := "Conta ativada com sucesso - Pet Services"

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
		.callout { background-color: #f0fdfa; border: 1px solid #99f6e4; border-radius: 10px; padding: 14px; margin: 18px 0; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Conta ativada</h1>
			</div>
			<div class="content">
				<p>Ola %s,</p>
				<p>Seu e-mail foi verificado e sua conta esta ativa.</p>
				<div class="callout">
					<p>Agora voce ja pode completar seu perfil e aproveitar os servicos da plataforma.</p>
				</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(name))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendLoginBlockedAlertEmail(to, name, reason string) error {
	subject := "Alerta de seguranca: tentativa de login bloqueada - Pet Services"

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
		.header h1 { margin: 0; font-size: 26px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.callout { background-color: #fef2f2; border: 1px solid #fecaca; border-radius: 10px; padding: 14px; margin: 18px 0; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Tentativa de login bloqueada</h1>
			</div>
			<div class="content">
				<p>Ola %s,</p>
				<p>Identificamos uma tentativa de login que foi bloqueada.</p>
				<div class="callout">
					<p><strong>Motivo:</strong> %s</p>
					<p>Se nao foi voce, recomendamos revisar a seguranca da sua conta.</p>
				</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(name), html.EscapeString(reason))

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

func (s *SMTPEmailService) SendPasswordChangedAlertEmail(to, name string) error {
	subject := "Alerta de seguranca: senha alterada - Pet Services"

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #1d4ed8 0%%, #3b82f6 100%%); }
		.header h1 { margin: 0; font-size: 28px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.callout { background-color: #eff6ff; border: 1px solid #bfdbfe; border-radius: 10px; padding: 14px; margin: 18px 0; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Senha alterada</h1>
			</div>
			<div class="content">
				<p>Ola %s,</p>
				<p>Este e um alerta de seguranca confirmando que a senha da sua conta foi alterada com sucesso.</p>
				<div class="callout">
					<p>Se voce nao reconhece esta alteracao, redefina sua senha imediatamente e entre em contato com o suporte.</p>
				</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(name))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendPasswordResetSuccessEmail(to, name string) error {
	subject := "Senha redefinida com sucesso - Pet Services"

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #15803d 0%%, #22c55e 100%%); }
		.header h1 { margin: 0; font-size: 28px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.callout { background-color: #f0fdf4; border: 1px solid #86efac; border-radius: 10px; padding: 14px; margin: 18px 0; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Senha redefinida</h1>
			</div>
			<div class="content">
				<p>Ola %s,</p>
				<p>Sua senha foi redefinida com sucesso.</p>
				<div class="callout">
					<p>Se voce nao solicitou esta redefinicao, proteja sua conta imediatamente e entre em contato com o suporte.</p>
				</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(name))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendAccountDeactivatedEmail(to, name string) error {
	subject := "Conta desativada - Pet Services"

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #6b7280 0%%, #9ca3af 100%%); }
		.header h1 { margin: 0; font-size: 28px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.callout { background-color: #f3f4f6; border: 1px solid #d1d5db; border-radius: 10px; padding: 14px; margin: 18px 0; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Conta desativada</h1>
			</div>
			<div class="content">
				<p>Ola %s,</p>
				<p>Sua conta foi desativada com sucesso.</p>
				<div class="callout">
					<p>Todos os tokens ativos foram revogados por seguranca. Quando desejar, voce pode reativar a conta.</p>
				</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(name))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendAccountReactivatedEmail(to, name string) error {
	subject := "Conta reativada - Pet Services"

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
		.callout { background-color: #f0fdfa; border: 1px solid #99f6e4; border-radius: 10px; padding: 14px; margin: 18px 0; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Conta reativada</h1>
			</div>
			<div class="content">
				<p>Ola %s,</p>
				<p>Sua conta foi reativada com sucesso.</p>
				<div class="callout">
					<p>Agora voce ja pode fazer login novamente e continuar usando a plataforma.</p>
				</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(name))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendAccountDeletedEmail(to, name string) error {
	subject := "Conta removida - Pet Services"

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #7f1d1d 0%%, #b91c1c 100%%); }
		.header h1 { margin: 0; font-size: 28px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.callout { background-color: #fef2f2; border: 1px solid #fecaca; border-radius: 10px; padding: 14px; margin: 18px 0; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Conta removida</h1>
			</div>
			<div class="content">
				<p>Ola %s,</p>
				<p>Sua conta foi removida com sucesso.</p>
				<div class="callout">
					<p>Se esta acao nao foi realizada por voce, entre em contato com o suporte imediatamente.</p>
				</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(name))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendRequestCreatedEmail(to, providerName, ownerName, petName, serviceName, requestID string) error {
	subject := "Nova solicitacao recebida - Pet Services"

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
		.details { background-color: #f0fdfa; border: 1px solid #99f6e4; border-radius: 10px; padding: 14px; margin-top: 16px; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Nova solicitacao</h1>
			</div>
			<div class="content">
				<p>Ola %s,</p>
				<p>Voce recebeu uma nova solicitacao de servico.</p>
				<div class="details">
					<p><strong>Cliente:</strong> %s</p>
					<p><strong>Pet:</strong> %s</p>
					<p><strong>Servico:</strong> %s</p>
					<p><strong>ID da solicitacao:</strong> %s</p>
				</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(providerName), html.EscapeString(ownerName), html.EscapeString(petName), html.EscapeString(serviceName), html.EscapeString(requestID))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendRequestCreatedOwnerConfirmationEmail(to, ownerName, providerName, petName, serviceName, requestID string) error {
	subject := "Solicitacao enviada com sucesso - Pet Services"

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #1d4ed8 0%%, #3b82f6 100%%); }
		.header h1 { margin: 0; font-size: 28px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.details { background-color: #eff6ff; border: 1px solid #bfdbfe; border-radius: 10px; padding: 14px; margin-top: 16px; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Solicitacao enviada</h1>
			</div>
			<div class="content">
				<p>Ola %s,</p>
				<p>Sua solicitacao foi enviada com sucesso para o prestador.</p>
				<div class="details">
					<p><strong>Prestador:</strong> %s</p>
					<p><strong>Pet:</strong> %s</p>
					<p><strong>Servico:</strong> %s</p>
					<p><strong>ID da solicitacao:</strong> %s</p>
				</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(ownerName), html.EscapeString(providerName), html.EscapeString(petName), html.EscapeString(serviceName), html.EscapeString(requestID))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendRequestAcceptedEmail(to, ownerName, providerName, petName, requestID string) error {
	subject := "Sua solicitacao foi aceita - Pet Services"

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #15803d 0%%, #22c55e 100%%); }
		.header h1 { margin: 0; font-size: 28px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.details { background-color: #f0fdf4; border: 1px solid #86efac; border-radius: 10px; padding: 14px; margin-top: 16px; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Solicitacao aceita</h1>
			</div>
			<div class="content">
				<p>Ola %s,</p>
				<p>Sua solicitacao foi aceita com sucesso.</p>
				<div class="details">
					<p><strong>Prestador:</strong> %s</p>
					<p><strong>Pet:</strong> %s</p>
					<p><strong>ID da solicitacao:</strong> %s</p>
				</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(ownerName), html.EscapeString(providerName), html.EscapeString(petName), html.EscapeString(requestID))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendRequestRejectedEmail(to, ownerName, providerName, petName, reason, requestID string) error {
	subject := "Sua solicitacao foi rejeitada - Pet Services"

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
		.details { background-color: #fef2f2; border: 1px solid #fecaca; border-radius: 10px; padding: 14px; margin-top: 16px; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Solicitacao rejeitada</h1>
			</div>
			<div class="content">
				<p>Ola %s,</p>
				<p>Sua solicitacao foi rejeitada.</p>
				<div class="details">
					<p><strong>Prestador:</strong> %s</p>
					<p><strong>Pet:</strong> %s</p>
					<p><strong>Motivo:</strong> %s</p>
					<p><strong>ID da solicitacao:</strong> %s</p>
				</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(ownerName), html.EscapeString(providerName), html.EscapeString(petName), html.EscapeString(reason), html.EscapeString(requestID))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendRequestCompletedEmail(to, ownerName, providerName, petName, requestID string) error {
	subject := "Sua solicitacao foi concluida - Pet Services"

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #1d4ed8 0%%, #3b82f6 100%%); }
		.header h1 { margin: 0; font-size: 28px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.details { background-color: #eff6ff; border: 1px solid #bfdbfe; border-radius: 10px; padding: 14px; margin-top: 16px; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Solicitacao concluida</h1>
			</div>
			<div class="content">
				<p>Ola %s,</p>
				<p>Seu atendimento foi concluido com sucesso.</p>
				<div class="details">
					<p><strong>Prestador:</strong> %s</p>
					<p><strong>Pet:</strong> %s</p>
					<p><strong>ID da solicitacao:</strong> %s</p>
				</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(ownerName), html.EscapeString(providerName), html.EscapeString(petName), html.EscapeString(requestID))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendReviewReminderEmail(to, ownerName, providerName, petName, requestID string) error {
	subject := "Lembrete: avalie seu atendimento - Pet Services"

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #1d4ed8 0%%, #3b82f6 100%%); }
		.header h1 { margin: 0; font-size: 28px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.details { background-color: #eff6ff; border: 1px solid #bfdbfe; border-radius: 10px; padding: 14px; margin-top: 16px; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Lembrete de avaliacao</h1>
			</div>
			<div class="content">
				<p>Ola %s,</p>
				<p>Seu atendimento foi concluido e sua avaliacao e muito importante para a comunidade.</p>
				<div class="details">
					<p><strong>Prestador:</strong> %s</p>
					<p><strong>Pet:</strong> %s</p>
					<p><strong>ID da solicitacao:</strong> %s</p>
				</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(ownerName), html.EscapeString(providerName), html.EscapeString(petName), html.EscapeString(requestID))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendReviewReceivedEmail(to, providerName, ownerName string, rating float64, comment string) error {
	subject := "Voce recebeu uma nova review - Pet Services"

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
		.details { background-color: #f0fdfa; border: 1px solid #99f6e4; border-radius: 10px; padding: 14px; margin-top: 16px; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Nova review recebida</h1>
			</div>
			<div class="content">
				<p>Ola %s,</p>
				<p>Seu perfil recebeu uma nova avaliacao.</p>
				<div class="details">
					<p><strong>Cliente:</strong> %s</p>
					<p><strong>Nota:</strong> %.1f/5</p>
					<p><strong>Comentario:</strong> %s</p>
				</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(providerName), html.EscapeString(ownerName), rating, html.EscapeString(comment))

	return s.sendEmail(to, subject, body)
}

// === ADOPTION MODULE EMAILS ===

func (s *SMTPEmailService) SendAdoptionGuardianProfileApprovedEmail(to, name string) error {
	subject := "🎉 Seu perfil foi aprovado! Comece a publicar anúncios"

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #10b981 0%%, #059669 100%%); }
		.header h1 { margin: 0; font-size: 28px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.cta { text-align: center; margin-top: 28px; }
		.cta a { background-color: #10b981; color: #ffffff; padding: 14px 32px; text-decoration: none; border-radius: 8px; display: inline-block; font-weight: 600; }
		.cta a:hover { background-color: #059669; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Bem-vindo! 🎉</h1>
			</div>
			<div class="content">
				<p>Olá <strong>%s</strong>,</p>
				<p>Excelentes notícias! Sua candidatura para ser responsável por adoção foi <strong>APROVADA</strong>!</p>
				<p>Agora você pode:</p>
				<ul>
					<li>Publicar pets para adoção</li>
					<li>Gerenciar candidaturas de interessados</li>
					<li>Acompanhar adoções bem-sucedidas</li>
				</ul>
				<p>Que família quer ajudar a dar um novo lar a um pet? 🐾</p>
				<div class="cta">
					<a href="https://petservices.com/adoption/listings/new">Publicar Primeiro Anúncio</a>
				</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(name))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendAdoptionGuardianProfileRejectedEmail(to, name, reason string) error {
	subject := "Informações sobre sua candidatura de perfil"

	reasonHTML := ""
	if reason != "" {
		reasonHTML = fmt.Sprintf(`<p><strong>Motivo:</strong> %s</p>`, html.EscapeString(reason))
	}

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #ef4444 0%%, #dc2626 100%%); }
		.header h1 { margin: 0; font-size: 28px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Atualização</h1>
			</div>
			<div class="content">
				<p>Olá <strong>%s</strong>,</p>
				<p>Obrigado por se candidatar para ser responsável por adoção em Pet Services.</p>
				<p>Infelizmente, sua candidatura foi recusada neste momento.</p>
				%s
				<p>Isso não significa que você não seja uma boa família! Muitas vezes, é apenas uma questão de timing ou circunstâncias.</p>
				<p><strong>Próximas etapas:</strong></p>
				<ul>
					<li>Entre em contato conosco em <strong>support@petservices.com</strong></li>
					<li>Podemos discutir como melhorar sua candidatura</li>
					<li>Você pode tentar novamente após 30 dias com informações atualizadas</li>
				</ul>
				<p>Continue acompanhando nossos anúncios - existem muitos pets esperando por um lar!</p>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(name), reasonHTML)

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendAdoptionApplicationSubmittedEmail(to, applicantName, petName string) error {
	subject := "✓ Sua candidatura foi enviada"

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #3b82f6 0%%, #2563eb 100%%); }
		.header h1 { margin: 0; font-size: 28px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.details { background-color: #eff6ff; border: 1px solid #bfdbfe; border-radius: 10px; padding: 14px; margin-top: 16px; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Candidatura Enviada!</h1>
			</div>
			<div class="content">
				<p>Olá <strong>%s</strong>,</p>
				<p>Sua candidatura foi enviada com sucesso! 🎉</p>
				<div class="details">
					<p><strong>Pet:</strong> %s</p>
					<p><strong>Status:</strong> Aguardando resposta do responsável</p>
				</div>
				<p>O responsável analisará sua candidatura e entrará em contato em breve se tudo correr bem.</p>
				<p><strong>Enquanto espera:</strong></p>
				<ul>
					<li>👉 Continue navegando por outros anúncios</li>
					<li>🔔 Receba notificações de novos pets</li>
					<li>❓ Qualquer dúvida? Contate support@petservices.com</li>
				</ul>
				<p>Desejamos boa sorte! 🐾</p>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(applicantName), html.EscapeString(petName))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendAdoptionApplicationReceivedGuardianEmail(to, guardianName, applicantName, petName string) error {
	subject := "📬 Nova candidatura para " + petName

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #f59e0b 0%%, #d97706 100%%); }
		.header h1 { margin: 0; font-size: 28px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.candidate-info { background-color: #fef3c7; border: 1px solid #fcd34d; border-radius: 10px; padding: 14px; margin-top: 16px; }
		.cta { text-align: center; margin-top: 28px; }
		.cta a { background-color: #f59e0b; color: #ffffff; padding: 14px 32px; text-decoration: none; border-radius: 8px; display: inline-block; font-weight: 600; }
		.cta a:hover { background-color: #d97706; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Nova Candidatura! 📬</h1>
			</div>
			<div class="content">
				<p>Olá <strong>%s</strong>,</p>
				<p>Excelentes notícias! Uma nova pessoa se candidatou a <strong>%s</strong>!</p>
				<div class="candidate-info">
					<p><strong>Candidato:</strong> %s</p>
					<p><strong>Interesse:</strong> Adotar %s</p>
				</div>
				<p>Revise a candidatura e decida se essa é a família perfeita para %s.</p>
				<p><strong>⏰ Importante:</strong> Responda em até 7 dias para manter a candidatura ativa e demonstrar interesse.</p>
				<div class="cta">
					<a href="https://petservices.com/adoption/listings">Ver Candidatura</a>
				</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(guardianName), html.EscapeString(petName), html.EscapeString(applicantName), html.EscapeString(petName), html.EscapeString(petName))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendAdoptionApplicationApprovedEmail(to, applicantName, petName, guardianContact string) error {
	subject := "🎉 Parabéns! Sua candidatura foi aprovada!"

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #10b981 0%%, #059669 100%%); }
		.header h1 { margin: 0; font-size: 32px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.steps { background-color: #f0fdf4; border: 1px solid #dcfce7; border-radius: 10px; padding: 14px; margin-top: 16px; }
		.steps ol { padding-left: 20px; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Parabéns!</h1>
			</div>
			<div class="content">
				<p>Olá <strong>%s</strong>,</p>
				<p>Excelentes notícias! Sua candidatura para <strong>%s</strong> foi <strong>APROVADA</strong>! 🎉</p>
				<p>O responsável aceitou sua candidatura! Prepare-se para receber seu novo companheiro.</p>
				<div class="steps">
					<h3>Próximos passos:</h3>
					<ol>
						<li><strong>Entre em contato:</strong> %s</li>
						<li><strong>Agende um encontro</strong> para conhecer %s pessoalmente</li>
						<li><strong>Conheça o pet</strong> e confira se é o match perfeito</li>
						<li><strong>Assine os documentos</strong> de adoção</li>
						<li><strong>Leve seu novo companheiro para casa!</strong> 🐾</li>
					</ol>
				</div>
				<p>Que maravilhoso momento! Bem-vindo à comunidade de adotantes Pet Services!</p>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(applicantName), html.EscapeString(petName), html.EscapeString(guardianContact), html.EscapeString(petName))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendAdoptionApplicationRejectedEmail(to, applicantName, petName string) error {
	subject := "Atualização sobre sua candidatura"

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #6366f1 0%%, #4f46e5 100%%); }
		.header h1 { margin: 0; font-size: 28px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.encouragement { background-color: #fef2f2; border: 1px solid #fee2e2; border-radius: 10px; padding: 14px; margin-top: 16px; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Atualização</h1>
			</div>
			<div class="content">
				<p>Olá <strong>%s</strong>,</p>
				<p>Obrigado pelo seu interesse em <strong>%s</strong>!</p>
				<p>Infelizmente, sua candidatura não foi selecionada nesta ocasião. O responsável recebeu outras candidaturas que estiveram mais alinhadas com as necessidades do pet.</p>
				<p>Mas isso não significa que você não seja uma ótima família! 💪</p>
				<div class="encouragement">
					<h3>Não desista!</h3>
					<ul>
						<li>👀 Continue navegando - existem muitos pets esperando por você</li>
						<li>🔔 Ative notificações para receber novos anúncios</li>
						<li>💬 Contate support@petservices.com se tiver dúvidas</li>
					</ul>
				</div>
				<p>O pet certo para você aparecerá em breve. Continue acompanhando!</p>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(applicantName), html.EscapeString(petName))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendPetAdoptedGuardianEmail(to, guardianName, petName string) error {
	subject := "🎉 Pet adotado com sucesso!"

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #10b981 0%%, #059669 100%%); }
		.header h1 { margin: 0; font-size: 32px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.celebration { background-color: #f0fdf4; border: 1px solid #dcfce7; border-radius: 10px; padding: 14px; margin-top: 16px; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Sucesso! 🎉</h1>
			</div>
			<div class="content">
				<p>Olá <strong>%s</strong>,</p>
				<p><strong>%s</strong> foi adotado com sucesso! Que notícia maravilhosa!</p>
				<div class="celebration">
					<p>🐾 Você ajudou a dar um novo lar a um pet.</p>
					<p>💝 Obrigado por fazer parte da Pet Services!</p>
				</div>
				<p><strong>E agora?</strong></p>
				<ul>
					<li>📸 Compartilhe histórias de sucesso com a comunidade</li>
					<li>🌟 Veja o histórico de adoções realizadas</li>
					<li>📢 Publique um novo pet se tiver mais para oferecer</li>
				</ul>
				<p>Continue ajudando pets a encontrarem seus novos lares!</p>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(guardianName), html.EscapeString(petName))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendPetAdoptedApplicantEmail(to, applicantName, petName string) error {
	subject := "🎉 Bem-vindo ao seu novo companheiro!"

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #10b981 0%%, #059669 100%%); }
		.header h1 { margin: 0; font-size: 32px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.resources { background-color: #f0fdf4; border: 1px solid #dcfce7; border-radius: 10px; padding: 14px; margin-top: 16px; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Bem-vindo! 🎉</h1>
			</div>
			<div class="content">
				<p>Olá <strong>%s</strong>,</p>
				<p>Parabéns! <strong>%s</strong> agora é parte da sua família! 🐾</p>
				<p>Este é o começo de uma linda história juntos. Bem-vindo à comunidade de adotantes Pet Services!</p>
				<div class="resources">
					<h3>Recursos úteis:</h3>
					<ul>
						<li>📚 Guia de incorporação do pet na família</li>
						<li>🏥 Dicas de saúde e cuidados básicos</li>
						<li>🐾 Comunidade de adotantes para trocar experiências</li>
						<li>📞 Suporte veterinário disponível</li>
					</ul>
				</div>
				<p>Queremos ver como está indo! Sinta-se à vontade para compartilhar atualizações e fotos.</p>
				<p>Felicidades! 💝</p>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(applicantName), html.EscapeString(petName))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendPetAdoptedRejectedApplicantsEmail(to, petName string) error {
	subject := "Atualização: Candidatura para " + petName

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #6366f1 0%%, #4f46e5 100%%); }
		.header h1 { margin: 0; font-size: 28px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.message { background-color: #fef3c7; border: 1px solid #fcd34d; border-radius: 10px; padding: 16px; margin-top: 16px; }
		.suggestions { background-color: #eff6ff; border: 1px solid #bfdbfe; border-radius: 10px; padding: 14px; margin-top: 16px; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Atualização de Candidatura</h1>
			</div>
			<div class="content">
				<p>Olá!</p>
				<div class="message">
					<p><strong>%s</strong> já foi adotado! 🐾</p>
					<p>Infelizmente, a candidatura não foi selecionada nesta oportunidade, mas isso não significa que o pet perfeito para sua família não está por vir!</p>
				</div>
				<div class="suggestions">
					<h3>O que fazer agora?</h3>
					<ul>
						<li>✨ Explore outros anúncios de animais disponíveis</li>
						<li>📌 Mantenha seu perfil atualizado para futuras seleções</li>
						<li>❤️ Deixe sua candidatura em aberto para novos pets que chegarem</li>
						<li>💌 Receba notificações de pets que combinam com seu perfil</li>
					</ul>
				</div>
				<p>A Pet Services conecta você com o companheiro perfeito. Continue conosco e continue procurando! 🐕🐈</p>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(petName))

	return s.sendEmail(to, subject, body)
}

func (s *SMTPEmailService) SendAdoptionApplicationWithdrawnGuardianEmail(to, guardianName, applicantName, petName string) error {
	subject := "Candidato retirou a candidatura - Pet Services"

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
		.header { padding: 28px 24px; text-align: center; color: #ffffff; background: linear-gradient(135deg, #6b7280 0%%, #9ca3af 100%%); }
		.header h1 { margin: 0; font-size: 26px; line-height: 1.2; }
		.content { padding: 28px 24px; font-size: 16px; line-height: 1.6; }
		.content p { margin: 0 0 14px 0; }
		.details { background-color: #f3f4f6; border: 1px solid #d1d5db; border-radius: 10px; padding: 14px; margin: 18px 0; }
		.callout { background-color: #eff6ff; border: 1px solid #bfdbfe; border-radius: 10px; padding: 14px; margin: 18px 0; }
		.footer { text-align: center; padding: 20px 16px 26px 16px; font-size: 12px; color: #6b7280; }
	</style>
</head>
<body>
	<div class="wrapper">
		<div class="container">
			<div class="header">
				<h1>Candidatura retirada</h1>
			</div>
			<div class="content">
				<p>Ola %s,</p>
				<p>Informamos que um candidato retirou sua candidatura ao seu anuncio de adocao.</p>
				<div class="details">
					<p><strong>Candidato:</strong> %s</p>
					<p><strong>Anuncio:</strong> %s</p>
				</div>
				<div class="callout">
					<p>Seu anuncio continua ativo e outros candidatos podem se inscrever. Acesse o painel para acompanhar as candidaturas restantes.</p>
				</div>
			</div>
			<div class="footer">
				&copy; 2026 Pet Services. Todos os direitos reservados.
			</div>
		</div>
	</div>
</body>
</html>
`, html.EscapeString(guardianName), html.EscapeString(applicantName), html.EscapeString(petName))

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
