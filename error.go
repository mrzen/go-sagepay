package sagepay

import (
	"fmt"
	"strings"
)

// Error represents an error from the sage API
type Error struct {
	Code        int    `json:"code"`
	Property    string `json:"property"`
	Description string `json:"description"`
	UserMessage string `json:"clientMessage"`
}

// ErrorResponse represents the response given by an Error
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

func (e ErrorResponse) Error() string {
	msgs := make([]string, len(e.Errors))

	for i, e := range e.Errors {
		msgs[i] = e.Error()
	}

	return strings.Join(msgs, "\n")
}

func (e Error) Error() string {
	return fmt.Sprintf("Error(%d): %s (%s, %s)", e.Code, e.UserMessage, e.Property, e.Description)
}
