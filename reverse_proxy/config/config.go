package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type ServiceAPI struct {
	Endpoint string `json:"endpoint"`
	Host     string `json:"host"`
	Port     int    `json:"port,string"`
	Scheme   string `json:"scheme"`
}

type Config struct {
	ServiceName string `json:"service_name"`
	Version     string `json:"version"`

	Server struct {
		Host   string `json:"host"`
		Port   int    `json:"port,string"`
		Secure string `json:"secure"`
	} `json:"server"`

	Certificate struct {
		CertFile string `json:"certFile"`
		KeyFile  string `json:"KeyFile"`
	} `json:"certificate"`

	Services []ServiceAPI `json:"services"`
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

	data, err := io.ReadAll(configFile)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
