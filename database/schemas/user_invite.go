package schemas

import (
	"database/sql"
	"github.com/segmentio/ksuid"
)

type UserInvite struct {
	ID        ksuid.KSUID
	Owner     *ksuid.KSUID
	Code      string
	Count     sql.NullInt32
	ExpiredAt sql.NullTime
}
