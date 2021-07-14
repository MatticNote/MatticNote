package mn_misc

import (
	"database/sql/driver"
	"github.com/MatticNote/MatticNote/internal"
	"github.com/MatticNote/MatticNote/server/api/graphql/mn_type"
)

func Conv2Interface(value driver.Valuer) interface{} {
	val, _ := value.Value()
	return val
}

func ConvInternal2GQLType(ius *internal.UserStruct) mn_type.UserType {
	return mn_type.UserType{
		Uuid:           ius.Uuid,
		Username:       ius.Username,
		Host:           Conv2Interface(ius.Host),
		DisplayName:    Conv2Interface(ius.DisplayName),
		Summary:        Conv2Interface(ius.Summary),
		CreatedAt:      Conv2Interface(ius.CreatedAt),
		UpdatedAt:      Conv2Interface(ius.UpdatedAt),
		AcceptManually: ius.AcceptManually,
		IsBot:          ius.IsBot,
	}
}
