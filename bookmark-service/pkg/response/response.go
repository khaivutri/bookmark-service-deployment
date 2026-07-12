package response

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type Message struct {
	Message string `json:"message"`
	Details []string `json:"details,omitempty"`
}

var (
	InternalServerErrResponse = Message{
		Message: "Internal Server Error",
		Details: nil,
	}
	InputErrResponse = Message{
		Message: "Invalid input",
		Details: nil,
	}
	CodeNotFoundResponse = Message{
        Message: "Code not found",
    }
)

func InputFieldError(err error) Message {
	if ok := errors.As(err, &validator.ValidationErrors{}); !ok {
		return InputErrResponse
	}

	var errs []string

	for _, err := range err.(validator.ValidationErrors) {
		errs = append(errs, err.Field() + " is invalid: ("+err.Tag()+") ")
	}

	return Message{
		Message: "Invalid input",
		Details: errs,
	}
}