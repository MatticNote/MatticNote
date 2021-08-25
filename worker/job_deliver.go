package worker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MatticNote/MatticNote/activitypub"
	"github.com/MatticNote/MatticNote/config"
	"github.com/MatticNote/MatticNote/internal/ist"
	"github.com/MatticNote/MatticNote/internal/user"
	"github.com/MatticNote/MatticNote/internal/version"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gocraft/work"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"time"
)

func (c *Context) Deliver(j *work.Job) error {
	to := j.ArgString("to")
	if err := j.ArgError(); err != nil {
		return err
	}
	bodyRaw, ok := j.Args["body"]
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

	req, err := http.NewRequest(http.MethodPost, to, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("MatticNote/%s", version.Version))
	req.Header.Set("Content-Type", "application/activity+json")

	signer, err := misc.GetHttpSignatureMethod()
	if err != nil {
		return err
	}

	privateKey, err := user.GetUserPrivateKey(fromUuid)
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

func (c *Context) NotePreJob(j *work.Job) error {
	visibility := j.ArgString("visibility")
	if j.ArgError() != nil {
		return j.ArgError()
	}
	var createdNote *ist.NoteStruct
	if createdNoteRaw, ok := j.Args["createdNote"]; !ok {
		log.Println("err: createdNote is not defined")
		return nil
	} else {
		createdNote = createdNoteRaw.(*ist.NoteStruct)
	}
	var authorUuid uuid.UUID
	if authorUuidRaw, ok := j.Args["authorUuid"]; !ok {
		log.Println("err: authorUuid is not defined")
		return nil
	} else {
		authorUuid = authorUuidRaw.(uuid.UUID)
	}

	activity := activitypub.RenderNoteActivity(createdNote)
	switch strings.ToUpper(visibility) {
	case "PUBLIC", "UNLISTED", "FOLLOWER":
		followerInbox, err := user.GetUserFollowerInbox(authorUuid)
		if err != nil {
			return err
		}
		if len(followerInbox) > 0 {
			for _, inbox := range followerInbox {
				_, err = Enqueue.Enqueue(
					JobDeliver,
					work.Q{
						"to":       inbox,
						"body":     activity,
						"fromUuid": authorUuid,
					},
				)
				if err != nil {
					return err
				}
			}

		}
	case "DIRECT":
		// todo: ダイレクトによるdeliver
	}

	return nil
}
