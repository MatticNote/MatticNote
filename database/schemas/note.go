package schemas

import (
	"database/sql"
	"github.com/segmentio/ksuid"
	"time"
)

type Note struct {
	ID        ksuid.KSUID
	Owner     *ksuid.KSUID
	CW        sql.NullString
	Body      sql.NullString
	ReplyID   *ksuid.KSUID
	RetextID  *ksuid.KSUID
	CreatedAt time.Time
}
