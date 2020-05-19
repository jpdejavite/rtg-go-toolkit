package errors

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// CustomError is a struct to custom error
type CustomError struct {
	Code    string
	Message string
}

// NotAuthorizedError not authorized custom error
var NotAuthorizedError = New("Not_authorized", "Not authorized")

// New is a creator ErrorWrapper struct
func New(code, message string) error {
	return CustomError{
		Code:    code,
		Message: message,
	}
}

// Error return error message field value
func (e CustomError) Error() string {
	return e.Message

}

// HandleGraphqlError handle graphql with custom error
func HandleGraphqlError() graphql.ErrorPresenterFunc {
	return func(ctx context.Context, err error) *gqlerror.Error {
		if customError, ok := err.(CustomError); ok {
			gqlError := &gqlerror.Error{
				Message:    customError.Message,
				Extensions: map[string]interface{}{"code": customError.Code},
			}
			return gqlError
		}

		return graphql.DefaultErrorPresenter(ctx, err)
	}
}
