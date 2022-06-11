package main

import (
	config "check_https/pkg/config"
	"check_https/pkg/send"
	"crypto/tls"
	"fmt"
	"github.com/robfig/cron"
	"log"
	"net/http"
	"time"
)

var urls []string
var SendKey string
var sendMsg *send.Send

func initSetting() {
	setting, err := config.NewSetting()
	if err != nil {
		panic(err)
	}
	err = setting.ReadSection("urls", &urls)
	err = setting.ReadSection("send_key", &SendKey)
	sendMsg = send.NewSend(SendKey)
	return

}

func main() {
	initSetting()
	c := cron.New()
	err := c.AddFunc("45 0 * * * *", func() {
		for _, url := range urls {
			err := checkSSL(url)
			if err != nil {
				log.Printf("check %s https err : %v", url, err)
				return
			}
		}
	})
	if err != nil {
		log.Println(err)
		return
	}
	c.Start()
	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}

func checkSSL(url string) error {
	client := &http.Client{
		Transport: &http.Transport{
			// 注意如果证书已过期，那么只有在关闭证书校验的情况下链接才能建立成功
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		// 10s 超时后认为服务挂了
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	// 遍历所有证书
	for _, cert := range resp.TLS.PeerCertificates {
		if !cert.NotAfter.After(time.Now()) {
			msg := fmt.Sprintf("Website [%s] certificate has expired: %s", url, cert.NotAfter.Local().Format("2006-01-02 15:04:05"))
			err := sendMsg.SendMsg("has expired", msg)
			if err != nil {
				log.Println(err)
				return err
			}
			log.Println(msg)
			return nil
		}

		if cert.NotAfter.Sub(time.Now()) < 5*24*time.Hour {
			msg := fmt.Sprintf("Website [%s] certificate will expire, remaining time: %fh", url, cert.NotAfter.Sub(time.Now()).Hours())
			err := sendMsg.SendMsg("will expire", msg)
			if err != nil {
				log.Println(err)
				return err
			}
			log.Println(msg)
			return nil
		}
	}
	log.Printf("the %s https no expired\n", url)
	return nil
}
