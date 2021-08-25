package ist

import (
	"github.com/MatticNote/MatticNote/misc"
	"github.com/google/uuid"
)

type NoteStruct struct {
	Uuid       uuid.UUID
	Author     *UserStruct
	CreatedAt  misc.NullTime
	Cw         misc.NullString
	Body       misc.NullString
	LocalOnly  bool
	Reply      *NoteStruct
	ReText     *NoteStruct
	Visibility string
}
