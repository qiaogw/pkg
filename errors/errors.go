package errors

import "github.com/cockroachdb/errors"

var (
	ErrExistsFile       = errors.New("exists file")
	ErrInvalidFieldName = errors.New("Invalid field name")
)

func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}
