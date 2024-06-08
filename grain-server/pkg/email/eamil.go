package emailx

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"time"
)

type MailServer struct {
	p        *email.Pool
	userName string
	password string
	host     string
	smtp     string
	port     string
}

func NewMailServer(userName string, password string, host string, smtp string, port string) (*MailServer, error) {
	e := &MailServer{userName: userName, password: password, host: host, smtp: smtp, port: port}
	return e, e.initEmailServer()
}

func (e *MailServer) initEmailServer() error {
	var err error
	e.p, err = email.NewPool(
		fmt.Sprintf("%s:%s", e.smtp, e.port),
		10,
		smtp.PlainAuth("", e.userName, e.password, e.host),
	)
	return err
}

func (e *MailServer) SendEmail(se *email.Email) error {
	if err := e.p.Send(se, 10*time.Second); err != nil {
		return err
	}
	return nil
}
