package worker

import (
	"github.com/MatticNote/MatticNote/misc"
	"github.com/gocraft/work"
	"log"
)

type Context struct {
}

func (c *Context) ProcessInbox(j *work.Job) error {
	log.Println("it works!")
	return nil
}

func (c *Context) Deliver(j *work.Job) error {
	_, err := misc.GetHttpSignatureMethod()
	if err != nil {
		return err
	}
	return nil
}
