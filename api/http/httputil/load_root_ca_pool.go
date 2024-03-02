package httputil

import (
	"crypto/x509"
	"fmt"
	"os"
)

func LoadRootCACertPool(rootCAPath string) *x509.CertPool {
	fmt.Printf("loading ca cert from : %s\n", rootCAPath)

	rootCA, err := os.ReadFile(rootCAPath)

	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(rootCA)
	return certPool

}
