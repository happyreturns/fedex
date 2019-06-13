package models

import "fmt"

// ReplyError is used when the severity of the API reply indicates an error
type ReplyError struct {
	Message  string
	Severity string
}

func (r ReplyError) Error() string {
	return fmt.Sprintf("reply got status %s with error: %s", r.Severity, r.Message)
}
