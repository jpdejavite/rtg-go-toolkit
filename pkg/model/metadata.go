package model

// MetaError metadata error to use in log
type MetaError struct {
	Err string `json:"error"`
}

// NewMetaError new metadata with error
func NewMetaError(err error) MetaError {
	return MetaError{
		Err: err.Error(),
	}
}
