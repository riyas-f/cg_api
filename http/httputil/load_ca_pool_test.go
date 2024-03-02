package httputil

import (
	"os"
	"path"
	"testing"
)

func TestLoadCaPool(t *testing.T) {
	if os.Getenv("ROOT_CA_CERT") == "" {
		t.Errorf("Empty ROOT_CA_CERT environment")
		return
	}

	caCertPath := path.Join("../../", os.Getenv("ROOT_CA_CERT"))

	caPool := LoadRootCACertPool(caCertPath)
	t.Log(caPool)
	t.Log("Success")
}
