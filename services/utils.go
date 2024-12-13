package services

import (
	"os"

	"github.com/samber/lo"
)

func GetEnv[T interface{}](key string, fallback T) T {
	var v interface{} = os.Getenv(key)
	var _v string = v.(string)
	if lo.IsEmpty(_v) {
		return fallback
	}

	return v.(T)
}

func GetIfPresent(value *string) *string {
	if value == nil {
		return lo.ToPtr("") // Default value if nil
	}
	return value
}
