package errors

import "encoding/json"

// GraphqlErrorExtensions graphql custom error extension
type GraphqlErrorExtensions struct {
	Code string `json:"code"`
}

// GraphqlError graphql custom error
type GraphqlError struct {
	Message    string                 `json:"message"`
	Extensions GraphqlErrorExtensions `json:"extensions"`
}

// GraphqlErrors graphql errors struct
type GraphqlErrors struct {
	Errors []GraphqlError `json:"errors"`
}

// NewGraphqlError build graphql error
func NewGraphqlError(code string, message string) GraphqlErrors {
	return GraphqlErrors{
		Errors: []GraphqlError{
			GraphqlError{
				Message: message,
				Extensions: GraphqlErrorExtensions{
					Code: code,
				},
			},
		},
	}
}

// NewGraphqlErrorToJSON build graphql error and convert to json string
func NewGraphqlErrorToJSON(code string, message string) string {
	jsonB, _ := json.Marshal(NewGraphqlError(code, message))
	return string(jsonB)
}
