package todo

import "fmt"

// Error contains error information.
type Error struct {
	Err     error
	Message string
	Code    int
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v: %v", e.Message, e.Err.Error())
}
