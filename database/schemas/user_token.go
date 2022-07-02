package schemas

import (
	"database/sql"
	"github.com/segmentio/ksuid"
)

type UserToken struct {
	ID        ksuid.KSUID
	Token     string
	UserId    *ksuid.KSUID
	ExpiredAt sql.NullTime
	IP        sql.NullString
}
