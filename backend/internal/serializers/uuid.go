package serializers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"gorm.io/gorm/schema"
)

// UUIDBinarySerializer handles serialization of uuid.UUID to BINARY(16)
type UUIDBinarySerializer struct{}

// Scan implements the schema.SerializerInterface
func (UUIDBinarySerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	if dbValue == nil {
		return nil
	}

	var result uuid.UUID
	switch v := dbValue.(type) {
	case []byte:
		if len(v) == 16 {
			// Binary UUID from database
			result, err = uuid.FromBytes(v)
			if err != nil {
				return err
			}
		} else if len(v) == 36 {
			// String UUID (fallback)
			result, err = uuid.ParseBytes(v)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("invalid UUID byte length: %d", len(v))
		}
	case string:
		// String UUID
		result, err = uuid.Parse(v)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported UUID type: %T", dbValue)
	}

	field.ReflectValueOf(ctx, dst).Set(reflect.ValueOf(result))
	return nil
}

// Value implements the schema.SerializerInterface
func (UUIDBinarySerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	if fieldValue == nil {
		return nil, nil
	}

	switch v := fieldValue.(type) {
	case uuid.UUID:
		if v == uuid.Nil {
			return nil, nil
		}
		// Return as 16-byte slice for BINARY(16) storage
		return v[:], nil
	case *uuid.UUID:
		if v == nil || *v == uuid.Nil {
			return nil, nil
		}
		return (*v)[:], nil
	default:
		return nil, fmt.Errorf("unsupported UUID type: %T", fieldValue)
	}
}
