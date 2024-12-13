// @alchemy replace package dao
package gorm

import (
	"fmt"
	"reflect"
)

// SetIfPresent sets a field of a struct if the provided value is not nil.
func SetIfPresent(target interface{}, fieldName string, value interface{}) error {
	// Ensure the target is a pointer to a struct
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}

	// Get the underlying struct and its type
	v = v.Elem()
	field := v.FieldByName(fieldName)

	// Check if the field exists and is settable
	if !field.IsValid() {
		return fmt.Errorf("field '%s' does not exist", fieldName)
	}
	if !field.CanSet() {
		return fmt.Errorf("field '%s' cannot be set", fieldName)
	}

	// Check if the value is not nil
	if value == nil {
		return nil
	}

	// Set the field value if types match
	val := reflect.ValueOf(value)
	if field.Type() != val.Type() {
		return fmt.Errorf("type mismatch: field '%s' is of type %s, but value is of type %s",
			fieldName, field.Type(), val.Type())
	}

	field.Set(val)
	return nil
}
