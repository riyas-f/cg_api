package jsonutil

import (
	"encoding/json"
	"io"
)

func EncodeToJson(v interface{}) ([]byte, error) {
	// response := base.Response{
	// 	Status:  responseStatus,
	// 	Message: responseMessage,
	// 	Data:    data,
	// }

	jsonResponse, err := json.Marshal(v)

	if err != nil {
		return nil, err
	}

	return jsonResponse, nil
}

func DecodeJSON(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)

	var err error

	for {
		err = decoder.Decode(&v)
		if err != nil {

			if err == io.EOF {
				err = nil
			}

			break
		}
	}

	return err
}
