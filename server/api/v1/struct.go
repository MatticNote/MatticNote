package v1

import (
	"github.com/MatticNote/MatticNote/misc"
)

type (
	v1UserRes struct {
		Uuid           string          `json:"uuid"`
		Username       string          `json:"username"`
		Host           misc.NullString `json:"host"`
		DisplayName    misc.NullString `json:"display_name"`
		Summary        misc.NullString `json:"summary"`
		CreatedAt      misc.NullTime   `json:"created_at"`
		UpdatedAt      misc.NullTime   `json:"updated_at"`
		AcceptManually bool            `json:"accept_manually"`
		IsBot          bool            `json:"is_bot"`
	}
)
