package extractor

import "fmt"

type ErrorHTTP struct {
	StatusCode int
	Message    string
}

func (e *ErrorHTTP) Error() string {
	return fmt.Sprintf("Error HTTP - Code: %d: \n Message: %s", e.StatusCode, e.Message)
}

type ErrorTechnical struct {
	Message string
	err     error
}

func (e *ErrorTechnical) Error() string {
	return fmt.Sprintf("Error Technical - Message: %s \n Error: %s", e.Message, e.err)
}
