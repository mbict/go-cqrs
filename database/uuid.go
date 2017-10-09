package database

import (
	"database/sql/driver"
	"github.com/satori/go.uuid"
)

//mysql specific fixuture
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
	return u.Bytes(), nil
}
