package worker

import (
	"github.com/gocraft/work"
	"log"
)

type Context struct {
}

func (c *Context) ProcessInbox(j *work.Job) error {
	log.Println("it works!")
	return nil
}
