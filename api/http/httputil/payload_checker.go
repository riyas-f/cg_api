package httputil

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	mapset "github.com/deckarep/golang-set/v2"
)

// CheckParameterUnity check if all required json parameers are filled
//
// if a field with tag json and have a value empty, the function will
// return with a non nil value.
// This function will ignore field with "omitempty" tag
func CheckParametersUnity(v interface{}, requiredField []string) responseerror.HTTPCustomError {
	// get interface field
	s := reflect.ValueOf(v)

	if s.Kind() == reflect.Pointer {
		s = s.Elem()
	}

	typeS := s.Type()
	fieldMap := make(map[string][]string)

	for _, field := range requiredField {
		splitFieldName := strings.SplitN(field, ":", 2)

		// Add the field to recurse path
		if len(splitFieldName) > 1 {
			fieldMap[splitFieldName[0]] = append(fieldMap[splitFieldName[0]], splitFieldName[1])
			continue
		}

		// No need to recursively check the value
		field = splitFieldName[0]

		v := s.FieldByName(field)
		if v.IsValid() {
			// check if a field is empty
			if v.Type().Kind() == reflect.Slice {
				if v.Len() == 0 {
					f, _ := typeS.FieldByName(field)
					tag := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
					return responseerror.CreateBadRequestError(
						responseerror.MissingParameter,
						responseerror.MissingParameterMessage,
						map[string]string{
							"field": tag,
						},
					)
				}

				continue
			}

			if reflect.Zero(v.Type()).Interface() == v.Interface() {
				f, _ := typeS.FieldByName(field)
				tag := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
				return responseerror.CreateBadRequestError(
					responseerror.MissingParameter,
					responseerror.MissingParameterMessage,
					map[string]string{
						"field": tag,
					},
				)
			}
		}

	}

	// Recursively resolve field name and check it's value
	for structName, fieldNames := range fieldMap {
		v := s.FieldByName(structName).Interface()
		err := CheckParametersUnity(v, fieldNames)
		// Short circuit when found error
		if err != nil {
			return err
		}
	}

	return nil
}

// check if a request header match with the expected value
func CheckHeader(h http.Header, headerName []string, expectedValue []mapset.Set[string]) responseerror.HTTPCustomError {
	if len(headerName) != len(expectedValue) {
		return responseerror.CreateInternalServiceError(errors.New("name length isn't equal with value length"))
	}

	for i := 0; i < len(headerName); i++ {
		if !expectedValue[i].Contains(h.Get(headerName[i])) {
			return responseerror.CreateBadRequestError(
				responseerror.HeaderValueMistmatch,
				responseerror.HeaderValueMistmatchMessage,
				map[string]string{
					"name": headerName[i],
				},
			)
		}
	}

	return nil
}
