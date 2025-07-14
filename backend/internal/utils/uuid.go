package utils

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

// UUID is a wrapper for uuid.UUID to support GORM BINARY(16) serialization
// Use utils.UUID in your models for BINARY(16) UUID columns
// Implements driver.Valuer and sql.Scanner

type UUID struct {
	uuid.UUID
}

func (u UUID) Value() (driver.Value, error) {
	if u.UUID == uuid.Nil {
		return nil, nil
	}
	return u.UUID[:], nil
}

func (u *UUID) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		if len(v) != 16 {
			return fmt.Errorf("invalid length for UUID: %d", len(v))
		}
		copy(u.UUID[:], v)
		return nil
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return err
		}
		u.UUID = parsed
		return nil
	case nil:
		u.UUID = uuid.Nil
		return nil
	default:
		return fmt.Errorf("cannot scan UUID from %T", src)
	}
}

// ParseUUID parses a string into a utils.UUID
func ParseUUID(id string) (UUID, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return UUID{}, err
	}
	return UUID{UUID: u}, nil
}
