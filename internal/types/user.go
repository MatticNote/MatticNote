package types

import (
	"database/sql"
	"github.com/segmentio/ksuid"
)

type (
	User struct {
		ID            ksuid.KSUID
		Email         sql.NullString
		Username      sql.NullString
		Host          sql.NullString
		DisplayName   sql.NullString
		Headline      sql.NullString
		Description   sql.NullString
		CreatedAt     sql.NullTime
		EmailVerified sql.NullBool
		IsSilence     bool
		IsSuspend     bool
		IsActive      bool
		IsModerator   bool
		IsAdmin       bool
	}

	Note struct {
		ID        ksuid.KSUID
		Owner     ksuid.KSUID
		CW        sql.NullString
		Body      sql.NullString
		CreatedAt sql.NullTime
	}
)
