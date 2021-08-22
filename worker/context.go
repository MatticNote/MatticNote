package worker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gocraft/work"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

type Context struct {
}

func (c *Context) ProcessInbox(j *work.Job) error {
	docRaw, ok := j.Args["doc"]
	if !ok {
		return errors.New("no args: data")
	}
	doc, ok := docRaw.(map[string]interface{})
	if !ok {
		return errors.New("cannot convert to map[string]interface{}")
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
	to := j.ArgString("to")
	if err := j.ArgError(); err != nil {
		return err
	}
	bodyRaw, ok := j.Args["bodyRaw"]
	if !ok {
		return errors.New("no arg: bodyRaw")
	}
	fromUuidRaw, ok := j.Args["fromUuid"]
	if !ok {
		return errors.New("no arg: fromUuid")
	}
	fromUuid, ok := fromUuidRaw.(uuid.UUID)
	if !ok {
		return errors.New("arg type error: fromUuid")
	}

	body, err := json.Marshal(bodyRaw)
	if err != nil {
		return errors.New("json parse failed")
	}

	// todo: urlは引数から取る
	req, err := http.NewRequest(http.MethodPost, to, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("MatticNote/%s", internal.Version))
	req.Header.Set("Content-Type", "application/activity+json")

	signer, err := misc.GetHttpSignatureMethod()
	if err != nil {
		return err
	}

	privateKey, err := internal.GetUserPrivateKey(fromUuid)
	if err != nil {
		return err
	}

	err = signer.SignRequest(
		privateKey,
		fmt.Sprintf("%s/activity/user/%s", config.Config.Server.Endpoint, fromUuid.String()),
		req,
		body,
	)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode >= 500 {
		return errors.New("returned server error response")
	}

	return nil
}
