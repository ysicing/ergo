package apiErr

import (
	"fmt"
)

func Err(message string, args ...interface{}) (err error) {
	return fmt.Errorf(message, args...)
}
