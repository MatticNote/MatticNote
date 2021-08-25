package ist

import (
	"encoding/pem"
	"github.com/MatticNote/MatticNote/misc"
	"github.com/google/uuid"
)

type UserStruct struct {
	Uuid           uuid.UUID
	Username       string
	Host           misc.NullString
	DisplayName    misc.NullString
	Summary        misc.NullString
	CreatedAt      misc.NullTime
	UpdatedAt      misc.NullTime
	IsSilence      bool
	AcceptManually bool
	IsBot          bool
	IsSuspend      bool
	PublicKey      *pem.Block
}

type UserRelationStruct struct {
	Following     bool
	FollowPending bool
	Follows       bool
	Muting        bool
	Blocking      bool
	Blocked       bool
}

type LocalUserStruct struct {
	UserStruct
	Email          string
	AcceptManually bool
	IsSuperuser    bool
	PrivateKey     *pem.Block
}
