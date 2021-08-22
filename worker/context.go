package worker

import (
	"errors"
	"fmt"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gocraft/work"
	"github.com/piprate/json-gold/ld"
	"log"
	"net/http"
)

type Context struct {
}

var (
	jldProc    = ld.NewJsonLdProcessor()
	jldOptions = ld.NewJsonLdOptions("")
	jldDoc     = map[string]interface{}{
		"@context": "https://www.w3.org/ns/activitystreams",
	}
)

func (c *Context) ProcessInbox(j *work.Job) error {
	data, ok := j.Args["data"]
	if !ok {
		return errors.New("no args: data")
	}
	doc, err := jldProc.Compact(data, jldDoc, jldOptions)
	if err != nil {
		log.Println("err: json-ld parse failed. ignore.")
		return nil
	}
	if len(doc) == 0 {
		log.Println("err: json-ld parsed, but nobody attributes. ignore.")
		return nil
	}

	apType, ok := doc["type"]
	if !ok {
		log.Println("err: type is not defined. ignore.")
		return nil
	}

	switch apType {
	case "Create":
		log.Println("Create activity")
	default:
		log.Println("err: unknown activity. ignore.")
	}

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
