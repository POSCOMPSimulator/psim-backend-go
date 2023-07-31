package mailer

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	From     string
	Password string
	SMTPHost string
	SMTPPort string
	Auth     smtp.Auth
}

func NewMailer(from, password, smtpHost, smtpPort string) *Mailer {
	auth := smtp.PlainAuth("", from, password, smtpHost)
	return &Mailer{
		From:     from,
		Password: password,
		SMTPHost: smtpHost,
		SMTPPort: smtpPort,
		Auth:     auth,
	}
}

func (m *Mailer) SendMail(to []string, subject string, message []byte) error {
	msg := []byte("To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		string(message) + "\r\n")
	err := smtp.SendMail(m.SMTPHost+":"+m.SMTPPort, m.Auth, m.From, to, msg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Email sent successfully!")
	return nil
}

func (m *Mailer) SendVerificationMail(to []string, verificationCode string) error {
	subject := "Verifique sua conta no PSIM"
	message := []byte("Obrigado por se inscrever no PSIM!\r\n\r\n" +
		"Por favor, verifique seu endereço de e-mail inserindo o seguinte código em nosso site:\r\n\r\n" +
		verificationCode + "\r\n\r\n" +
		"Atenciosamente,\r\n" +
		"A equipe do PSIM")
	err := m.SendMail(to, subject, message)
	if err != nil {
		return err
	}
	fmt.Println("E-mail de verificação enviado com sucesso!")
	return nil
}

func (m *Mailer) SendRecoverMail(to []string, resetLink string) error {
	subject := "Recuperação de senha do PSIM"
	message := []byte("Olá!\r\n\r\n" +
		"Recebemos uma solicitação para redefinir a senha da sua conta no PSIM.\r\n\r\n" +
		"Para redefinir sua senha, clique no seguinte link:\r\n\r\n" +
		resetLink + "\r\n\r\n" +
		"Se você não solicitou uma redefinição de senha, pode ignorar este e-mail.\r\n\r\n" +
		"Atenciosamente,\r\n" +
		"A equipe do PSIM")
	err := m.SendMail(to, subject, message)
	if err != nil {
		return err
	}
	fmt.Println("E-mail de recuperação de senha enviado com sucesso!")
	return nil
}
