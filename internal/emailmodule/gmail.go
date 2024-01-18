package emailmodule

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
	"os"
	"strings"
)

type GmailSender struct {
	EmailFrom  string
	SenderType string
}

func (s *GmailSender) GetSenderType() string {
	return "Gmail"
}

func (s *GmailSender) SendEmail(subject, body string, to []string) error {
	ctx := context.Background()

	b, err := os.ReadFile("config/static/credentials.json")
	if err != nil {
		log.Fatalf("无法读取凭据文件: %v", err)
	}

	// 如果修改了scopes，删除之前保存的token.json。
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("无法解析凭据文件: %v", err)
	}

	srv, err := gmail.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, getToken(config))))
	if err != nil {
		log.Fatalf("无法检索Gmail client: %v", err)
	}

	// 创建并发送邮件
	var message gmail.Message
	emailTo := strings.Join(to, ",")
	messageStr := []byte(
		"From: " + s.EmailFrom + "\r\n" +
			"To:  " + emailTo + "\r\n" +
			"Subject: " + subject + "\r\n\r\n" +
			body)

	message.Raw = base64.URLEncoding.EncodeToString(messageStr)
	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Fatalf("无法发送邮件: %v", err)
	}

	fmt.Println("邮件已发送！")
	return nil
}

func getToken(config *oauth2.Config) *oauth2.Token {
	// Token文件存储了用户的token。
	const tokenFile = "config/static/token.json"
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return tok
}

// getTokenFromWeb 从Google的web流程中请求一个Token。
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("请在浏览器中打开以下链接，并授权应用程序: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("无法读取授权码: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("无法从web中获取token: %v", err)
	}
	return tok
}

// tokenFromFile 从文件中读取Token。
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// saveToken 将Token保存到文件中。
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("保存Token文件到: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("无法缓存oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
