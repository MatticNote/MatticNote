package worker

import (
	"errors"
	"github.com/gocraft/work"
	"log"
)

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
