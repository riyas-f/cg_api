package config

import (
	"crypto/tls"
	"os"
	"path/filepath"

	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
)

var (
	DB_PASSWORD string = os.Getenv("DB_PASSWORD")
	DB_USER     string = os.Getenv("DB_USER")
	DB_DATABASE string = os.Getenv("DB_DATABASE_NAME")
)

type Service struct {
	Host   string `json:"host"`
	Port   int    `json:"port,string"`
	Scheme string `json:"scheme"`
}
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

	Service struct {
		Auth    Service `json:"auth"`
		Session Service `json:"session"`
	} `json:"service"`

	Psgination struct {
		DefaultLimit int `json:"default_limit,string"`
		MaxLimit     int `json:"max_limit,string"`
	}

	*tls.Config
}

func ReadJSONConfiguration(path string) (*Config, error) {
	var config Config

	config.Database.Username = DB_USER
	config.Database.Password = DB_PASSWORD
	config.Database.Database = DB_DATABASE

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
