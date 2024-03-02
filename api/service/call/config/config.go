package config

import (
	"crypto/tls"
	"os"
	"path/filepath"

	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
)

type Config struct {
	ServiceName string `json:"service_name"`
	Version     string `json:"version"`
	Database    struct {
		Host     string `json:"host"`
		Port     int    `json:"port,string"`
		Username string `json:"username"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"database"`

	Server struct {
		Host   string `json:"host"`
		Port   int    `json:"port,string"`
		Secure string `json:"secure"`
	} `json:"server"`

	*tls.Config
}

func ReadJSONConfiguration(path string) (*Config, error) {
	var config Config

	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exePath := filepath.Dir(exe)

	configFile, err := os.Open(filepath.Join(exePath, "", path))
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	err = jsonutil.DecodeJSON(configFile, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
