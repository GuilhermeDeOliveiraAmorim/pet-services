package user

// EmailService abstrai o envio de emails.
type EmailService interface {
	SendPasswordResetEmail(to, resetLink string) error
}
