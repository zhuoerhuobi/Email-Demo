package emailmodule

import "fmt"

type EmailSender interface {
	GetSenderType() string
	SendEmail(subject, body string, to []string) error
}

type EmailService struct {
	Senders []EmailSender
}

func (e *EmailService) SendEmail(subject, body string, to []string) error {
	for _, sender := range e.Senders {
		err := sender.SendEmail(subject, body, to)
		if err == nil {
			return nil
		}
		// 记录错误，继续尝试下一个服务商
		fmt.Println(sender.GetSenderType()+"发送失败", err)
	}
	return fmt.Errorf("所有邮件服务商发送失败")
}
