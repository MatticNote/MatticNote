package worker

import (
	"fmt"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gocraft/work"
	"log"
	"net/http"
)

type Context struct {
}

func (c *Context) ProcessInbox(j *work.Job) error {
	log.Println("it works!")
	return nil
}

func (c *Context) Deliver(j *work.Job) error {
	// todo: urlは引数から取る
	req, err := http.NewRequest(http.MethodPost, "https://example.com/inbox", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("MatticNote/%s", internal.Version))
	req.Header.Set("Content-Type", "application/activity+json")

	// todo: このあとにSignRequestで秘密鍵、鍵ID（URL）、↑のHttpRequest、送信するデータの本文([]byte形式)で処理
	_, err = misc.GetHttpSignatureMethod()
	if err != nil {
		return err
	}

	// todo: この辺りでHTTP通信をする

	return nil
}
