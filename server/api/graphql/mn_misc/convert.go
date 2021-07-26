package mn_misc

import (
	"database/sql/driver"
)

func Conv2Interface(value driver.Valuer) interface{} {
	val, _ := value.Value()
	return val
}
