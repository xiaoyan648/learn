package emailutil

import (
	"gopkg.in/gomail.v2"
)

// 线上公用邮箱配置
const (
	user = "notifications@qimao.com"
	pass = "YMW-9dt-V8s-ZHp"
	host = "smtp.exmail.qq.com"
	port = 465
)

func SendMail(mailTo []string, subject string, body string, attach []string, formSet string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", formSet)
	msg.SetHeader("To", mailTo...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)
	if len(attach) > 0 {
		for _, f := range attach {
			msg.Attach(f)
		}
	}

	dialer := gomail.NewDialer(host, port, user, pass)

	return dialer.DialAndSend(msg)
}
