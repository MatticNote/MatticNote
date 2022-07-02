package schemas

import (
	"database/sql"
	"github.com/segmentio/ksuid"
	"time"
)

type User struct {
	ID          ksuid.KSUID
	Username    sql.NullString
	Host        sql.NullString
	DisplayName sql.NullString
	Headline    sql.NullString
	Description sql.NullString
	CreatedAt   time.Time
	IsSilence   bool
	IsSuspend   bool
	IsModerator bool
	IsAdmin     bool
	DeletedAt   sql.NullTime
}
