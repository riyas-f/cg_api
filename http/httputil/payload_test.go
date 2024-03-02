package httputil

import (
	"net/http"
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
)

type TestStruct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age,omitempty"`
}

func TestCheckPayloadParameters(t *testing.T) {
	data := &TestStruct{
		ID: "10",
	}

	err := CheckParametersUnity(data, []string{"ID", "Name", "Age"})

	t.Logf("Any error: %s", err)

	if err == nil {
		t.Errorf("Function expected to return non nil error, however it return nil")
	}
}

func TestHeaderCheck(t *testing.T) {
	h := http.Header{}
	h.Set("Content-Type", "application/xml")

	expectedValue := make([]mapset.Set[string], 1)
	expectedValue[0] = mapset.NewSet[string]("application/json")

	headerName := [1]string{"Content-Type"}

	// set := mapset.NewSet[string]("Test")

	// t.Log(reflect.TypeOf((set)))

	err := CheckHeader(h, headerName[:], expectedValue[:])

	t.Logf("Any error: %s", err)

	if err == nil {
		t.Errorf("Function expected to return an error, however it return nil")
	}

}
