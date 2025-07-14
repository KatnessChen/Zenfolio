package types

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

// BinaryUUID is a custom type that stores UUIDs as BINARY(16) in MySQL
type BinaryUUID uuid.UUID

// Scan implements the sql.Scanner interface for reading from database
func (u *BinaryUUID) Scan(value interface{}) error {
	if value == nil {
		*u = BinaryUUID(uuid.Nil)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		if len(v) == 16 {
			// Binary UUID from database
			parsed, err := uuid.FromBytes(v)
			if err != nil {
				return err
			}
			*u = BinaryUUID(parsed)
			return nil
		} else if len(v) == 36 {
			// String UUID (fallback)
			parsed, err := uuid.ParseBytes(v)
			if err != nil {
				return err
			}
			*u = BinaryUUID(parsed)
			return nil
		}
		return fmt.Errorf("invalid UUID byte length: %d", len(v))
	case string:
		// String UUID
		parsed, err := uuid.Parse(v)
		if err != nil {
			return err
		}
		*u = BinaryUUID(parsed)
		return nil
	default:
		return fmt.Errorf("unsupported UUID type: %T", value)
	}
}

// Value implements the driver.Valuer interface for writing to database
func (u BinaryUUID) Value() (driver.Value, error) {
	if uuid.UUID(u) == uuid.Nil {
		return nil, nil
	}
	// Convert to uuid.UUID and get bytes for BINARY(16) storage
	uuidVal := uuid.UUID(u)
	return uuidVal[:], nil
}

// String returns the string representation of the UUID
func (u BinaryUUID) String() string {
	return uuid.UUID(u).String()
}

// UUID converts BinaryUUID back to uuid.UUID
func (u BinaryUUID) UUID() uuid.UUID {
	return uuid.UUID(u)
}

// MarshalJSON implements json.Marshaler interface
func (u BinaryUUID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + u.String() + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (u *BinaryUUID) UnmarshalJSON(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("invalid UUID JSON")
	}
	// Remove quotes
	str := string(data[1 : len(data)-1])
	parsed, err := uuid.Parse(str)
	if err != nil {
		return err
	}
	*u = BinaryUUID(parsed)
	return nil
}

// NewBinaryUUID creates a new BinaryUUID
func NewBinaryUUID() BinaryUUID {
	return BinaryUUID(uuid.New())
}

// ParseBinaryUUID parses a string UUID into BinaryUUID
func ParseBinaryUUID(s string) (BinaryUUID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return BinaryUUID{}, err
	}
	return BinaryUUID(parsed), nil
}

// MustParseBinaryUUID parses a string UUID into BinaryUUID, panics on error
func MustParseBinaryUUID(s string) BinaryUUID {
	parsed := uuid.MustParse(s)
	return BinaryUUID(parsed)
}
