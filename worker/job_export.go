package worker

import (
	"errors"
	"github.com/gocraft/work"
	"github.com/google/uuid"
	"log"
)

func (c Context) ExportData(j *work.Job) error {
	targetUserRaw, ok := j.Args["targetUser"]
	if !ok {
		return errors.New("no args: targetUser")
	}
	targetUser, ok := targetUserRaw.(uuid.UUID)
	if !ok {
		return errors.New("cannot convert to UUID")
	}
	log.Printf("Export requested: %s\n", targetUser.String())

	// todo: ユーザデータのエクスポート処理をする

	return nil
}
