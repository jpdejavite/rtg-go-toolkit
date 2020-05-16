package errors

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
