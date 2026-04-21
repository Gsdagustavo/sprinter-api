package mail

type Sender interface {
	SendMail(to []string, subject, content string) error
}
