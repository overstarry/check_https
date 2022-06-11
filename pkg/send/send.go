package send

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Send struct {
	SendKey string `json:"send_key"`
}

func NewSend(sendKey string) *Send {
	return &Send{SendKey: sendKey}
}

func (s *Send) SendMsg(title, msg string) error {
	res, err := http.Get(fmt.Sprintf("https://sctapi.ftqq.com/%s.send?title=%s&desp=%s", s.SendKey, title, msg))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	log.Println(string(body))
	return nil
}
