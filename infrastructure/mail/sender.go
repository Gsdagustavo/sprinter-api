package mail

// Sender is a generic interface for email senders
type Sender interface {
	// SendMail defines a generic function to send an email
	SendMail(to []string, subject, content string) error
}
