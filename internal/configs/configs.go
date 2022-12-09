package configs

import (
	"fmt"
)

type Error struct {
	Field string
}

func (e *Error) Error() string {
	return fmt.Sprintf("config error(parsing %s[%T])", e.Field, e.Field)
}

func Get[T any](cfg map[string]any, key string) (a T, err error) {
	a, ok := cfg[key].(T)
	if !ok {
		err = &Error{key}
		return
	}
	return
}
