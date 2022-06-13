package types

import (
	"github.com/segmentio/ksuid"
	"time"
)

type User struct {
	ID          ksuid.KSUID
	Email       *string
	Password    *string
	Username    *string
	Host        *string
	DisplayName *string
	Description *string
	CreatedAt   time.Time
}
