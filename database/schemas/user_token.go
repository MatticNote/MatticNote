package schemas

import (
	"database/sql"
	"github.com/segmentio/ksuid"
	"time"
)

type UserToken struct {
	ID        ksuid.KSUID
	Token     string
	UserId    *ksuid.KSUID
	ExpiredAt sql.NullTime
	IP        sql.NullString
	CreatedAt time.Time
}
