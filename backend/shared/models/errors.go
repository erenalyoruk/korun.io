package models

import "fmt"

type Error struct {
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code: %s, message: %s, err: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}
