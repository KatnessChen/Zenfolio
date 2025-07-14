package utils

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

// UUID is a wrapper for uuid.UUID that handles conversion to/from VARCHAR(36) in database
// This allows us to use standard uuid.UUID throughout the codebase while handling database serialization
type UUID struct {
	uuid.UUID
}

// NewUUID creates a new UUID
func NewUUID() UUID {
	return UUID{UUID: uuid.New()}
}

// ParseUUID parses a string into a UUID
func ParseUUID(s string) (UUID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return UUID{}, err
	}
	return UUID{UUID: parsed}, nil
}

// MustParseUUID parses a string into a UUID, panics on error
func MustParseUUID(s string) UUID {
	return UUID{UUID: uuid.MustParse(s)}
}

// Value implements driver.Valuer interface for database storage as VARCHAR(36)
func (u UUID) Value() (driver.Value, error) {
	if u.UUID == uuid.Nil {
		return nil, nil
	}
	return u.UUID.String(), nil
}

// Scan implements sql.Scanner interface for reading from VARCHAR(36) database column
func (u *UUID) Scan(value interface{}) error {
	if value == nil {
		u.UUID = uuid.Nil
		return nil
	}

	switch v := value.(type) {
	case string:
		if v == "" {
			u.UUID = uuid.Nil
			return nil
		}
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("invalid UUID string: %w", err)
		}
		u.UUID = parsed
		return nil
	case []byte:
		if len(v) == 0 {
			u.UUID = uuid.Nil
			return nil
		}
		return u.Scan(string(v))
	default:
		return fmt.Errorf("cannot scan UUID from type %T", value)
	}
}

// String returns the string representation of the UUID
func (u UUID) String() string {
	return u.UUID.String()
}

// IsNil returns true if the UUID is nil
func (u UUID) IsNil() bool {
	return u.UUID == uuid.Nil
}

// ToStandardUUID converts to standard uuid.UUID
func (u UUID) ToStandardUUID() uuid.UUID {
	return u.UUID
}

// FromStandardUUID creates a UUID from standard uuid.UUID
func FromStandardUUID(id uuid.UUID) UUID {
	return UUID{UUID: id}
}

// MarshalJSON implements json.Marshaler
func (u UUID) MarshalJSON() ([]byte, error) {
	if u.UUID == uuid.Nil {
		return []byte("null"), nil
	}
	return []byte(`"` + u.UUID.String() + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (u *UUID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		u.UUID = uuid.Nil
		return nil
	}

	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid UUID JSON format")
	}

	str := string(data[1 : len(data)-1])
	if str == "" {
		u.UUID = uuid.Nil
		return nil
	}

	parsed, err := uuid.Parse(str)
	if err != nil {
		return fmt.Errorf("invalid UUID in JSON: %w", err)
	}

	u.UUID = parsed
	return nil
}
