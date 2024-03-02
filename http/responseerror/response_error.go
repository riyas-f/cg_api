package responseerror

import (
	"bytes"
	"fmt"
	"html/template"
)

type errorType string
type errorMessageTemplate string

// type NamedArgs map[string]string

type ResponseError struct {
	Code    int
	Message string
	Name    string
}

// Create error message from specified template and arguments
//
// {{.key}} in template string will be replaced by the value in namedArgs with the same key
// as specified in the template
// if the key in template doesn't exist in namedArgs, it will be replaced with <no value>.
// make sure to match the namedArgs with the key specified with the template. If the template
// doesn't require any arguments, use nil as the namedArgs, it will return the template
//
// below is example for message formatting:
//
//	tmp1 := "Hello {{.name}}"
//	args := map[string][string]{"name": "aditya"}
//	ParseMessage(tmp1, args) // "Hello aditya"
func ParseMessage(tmp errorMessageTemplate, namedArgs map[string]string) string {
	if namedArgs == nil {
		return string(tmp)
	}

	var tp1 bytes.Buffer

	tmp1, err := template.New("message").Parse(string(tmp))

	if err != nil {
		panic(err)
	}

	if err = tmp1.Execute(&tp1, namedArgs); err != nil {
		panic(err)
	}

	return (tp1.String())
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("Status Code: %d, Message: %s", e.Code, e.Message)
}

func (e *ResponseError) Get() *ResponseError {
	return e
}

type HTTPCustomError interface {
	Get() *ResponseError
	// Implement error interface
	Error() string
}
