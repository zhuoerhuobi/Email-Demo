package emailmodule

type FoxmailSender struct {
	EmailFrom  string
	SenderType string
}

func (s *FoxmailSender) GetSenderType() string {
	return "Foxmail"
}

func (s *FoxmailSender) SendEmail(subject, body string, to []string) error {
	// 实现使用Foxmail发送邮件的逻辑
	return nil
}
