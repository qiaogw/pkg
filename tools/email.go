package tools

import (
	"github.com/astaxie/beego"
	"net/smtp"
	"strings"
)

/* SendMail
 *  to: example@example.com;example1@163.com;example2@sina.com.cn;...
 *  subject:The subject of mail
 *  body: The content of mail
 */
func SendMail(to string, subject string, body string) error {
	user := beego.AppConfig.String("mailfrom")
	password := beego.AppConfig.String("mailpassword")
	host := beego.AppConfig.String("mailhost")

	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var contentType string
	contentType = "Content-type:text/html;charset=utf-8"

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, sendTo, msg)
	return err
}
