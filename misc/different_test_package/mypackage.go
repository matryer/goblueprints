package mypackage

import "net/mail"

type MailMan struct{}

func (m *MailMan) Send(subject, body string, to ...*mail.Address) {

}

func New() *MailMan {
	return &MailMan{}
}
