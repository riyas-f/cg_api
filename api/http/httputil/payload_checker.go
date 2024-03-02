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
	s := reflect.ValueOf(v).Elem()
	typeS := s.Type()
	// typeOfS := s.Type()

	// for i := 0; i < typeOfS.NumField(); i++ {
	// 	field := typeOfS.Field(i)
	// 	jsonTag := field.Tag.Get("json")

	// 	// Gatekeep conditional
	// 	if jsonTag == "-" || jsonTag == "" {
	// 		continue
	// 	}

	// 	if x := strings.SplitAfter(jsonTag, ","); len(x) > 1 {
	// 		if x[1] == "omitempty" {
	// 			continue
	// 		}
	// 	}

	// 	a := s.Field(i).Interface()

	for _, field := range requiredField {
		v := s.FieldByName(field)
		if v.IsValid() {
			// check if a field is empty
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
