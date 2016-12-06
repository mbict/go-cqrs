package database

import (
	"github.com/satori/go.uuid"
	"database/sql/driver"
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

