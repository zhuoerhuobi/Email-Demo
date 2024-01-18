package main

import (
	"log"
)
import "Email-Demo/internal/emailmodule"

func main() {
	// 初始化服务商配置
	gmailSender := &emailmodule.GmailSender{
		EmailFrom:  "example@example.com",
		SenderType: "gmail",
	}
	foxmailSender := &emailmodule.FoxmailSender{
		EmailFrom:  "example@example.com",
		SenderType: "foxmail",
	}

	// 创建EmailService并添加服务商
	emailService := &emailmodule.EmailService{
		Senders: []emailmodule.EmailSender{gmailSender, foxmailSender},
	}

	// 使用emailService发送邮件
	err := emailService.SendEmail("主题", "正文", []string{"example@example.com"})
	if err != nil {
		log.Fatalf("发送邮件失败: %v", err)
	}
}
