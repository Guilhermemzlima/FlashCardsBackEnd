package errors

import (
	"fmt"
	"github.com/pkg/errors"
)

var ErrNotFound = &NotFound{Err: fmt.Errorf("not found")}
var ErrInvalidPayload = &InvalidPayload{Err: fmt.Errorf("invalid payload")}
var ErrInternalServer = &InternalServerError{Err: fmt.Errorf("internal server error")}
var ErrUnauthorized = &Unauthorized{Err: fmt.Errorf("unauthorized")}
var ErrBadRequest = &BadRequest{Err: fmt.Errorf("internal server error")}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func Wrap(err error) error {
	if _, ok := err.(stackTracer); ok {
		return err
	}

	return errors.WithStack(err)
}

func WrapWithMessage(err error, msg string) error {
	if _, ok := err.(stackTracer); ok {
		return errors.WithMessage(err, msg)
	}

	return errors.Wrap(err, msg)
}

func Cause(err error) error {
	return errors.Cause(err)
}

func New(msg string) error {
	return errors.New(msg)
}
