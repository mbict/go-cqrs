package mysql

import (
	"database/sql/driver"
	"github.com/google/uuid"
)

//mysql optimized ID key, we store it as a bytes VARBINARY(16) instead of string
type UUID struct {
	uuid.UUID
}

func MysqlUUID(uuid uuid.UUID) UUID {
	return UUID{
		UUID: uuid,
	}
}

//optimized way of storing uuid in mysql
func (u UUID) Value() (driver.Value, error) {
	return u.UUID, nil
}
