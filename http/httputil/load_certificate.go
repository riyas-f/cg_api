package httputil

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	"github.com/youmark/pkcs8"
)

func LoadCertificate(certPath string, privateKeyPath string, passPath string) (*x509.Certificate, interface{}) {
	pw, err := os.ReadFile(passPath)

	if err != nil {
		log.Fatal(err)
	}

	publicKeyFile, err := os.ReadFile(certPath)

	if err != nil {
		log.Fatal(err)
	}

	pemBlock, _ := pem.Decode(publicKeyFile)

	if pemBlock == nil {
		log.Fatal("pem decode failed")
	}

	cert, err := x509.ParseCertificate(
		pemBlock.Bytes,
	)

	if err != nil {
		log.Fatal(err)
	}

	privateKeyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	pemBlock, _ = pem.Decode(privateKeyFile)

	if pemBlock == nil {
		log.Fatal("pem decode failed")
	}

	privateKey, err := pkcs8.ParsePKCS8PrivateKey(pemBlock.Bytes, pw)

	if err != nil {
		log.Fatal(err)
	}

	return cert, privateKey
}
