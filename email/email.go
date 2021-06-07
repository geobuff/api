package email

import (
	"fmt"
	"os"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var SendResetToken = func(email, resetLink string) (*rest.Response, error) {
	from := mail.NewEmail(os.Getenv("EMAIL_NAME"), os.Getenv("EMAIL_ADDRESS"))
	subject := "Password Reset Request"
	to := mail.NewEmail("User", email)
	plainTextContent := fmt.Sprintf("Hi there,\n\nBelow is the link to reset the password for your account:\n%s\n\nIf you did not request a password reset please disregard this email.\n\nFrom,\nThe GeoBuff Team", resetLink)
	htmlContent := fmt.Sprintf("<div><p>Hi there,</p><p>Below is the link to reset the password for your account:</p><p>%s</p><p>If you did not request a password reset please disregard this email.</p><p>From,</p><p>The GeoBuff Team</p></div>", resetLink)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	return client.Send(message)
}
