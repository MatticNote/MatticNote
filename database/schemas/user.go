package schemas

import (
	"database/sql"
	"github.com/segmentio/ksuid"
	"time"
)

type User struct {
	ID          ksuid.KSUID
	Username    sql.NullString
	Host        *ksuid.KSUID
	DisplayName sql.NullString
	Headline    sql.NullString
	Description sql.NullString
	CreatedAt   time.Time
	IsSilence   bool
	IsSuspend   bool
	IsModerator bool
	IsAdmin     bool
	PublicKey   *[]byte
	DeletedAt   sql.NullTime
}
