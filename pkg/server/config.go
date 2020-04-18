package server

import (
	"fmt"
	"strings"
)

type Config struct {
	Bind             string
	DatabaseHost     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseSSL      bool
	StorageEngine    string
	BucketName       string
	ProjectID        string
	MinioAccessKey   string
	MinioSecretKey   string
	MinioEndpoint    string
	MinioUseSSL      bool
}

func (c Config) Validate() error {
	errors := ValidationErrors{}
	if c.Bind == "" {
		errors = append(errors, "PHTS_BIND not provided")
	}
	if c.DatabaseHost == "" {
		errors = append(errors, "PHTS_DB_HOST not provided")
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (c Config) DatabaseConnectionString() string {
	ssl := "enable"
	if !c.DatabaseSSL {
		ssl = "disable"
	}
	return fmt.Sprintf("user=%s host=%s password=%s dbname=%s sslmode=%s", c.DatabaseUser, c.DatabaseHost, c.DatabasePassword, c.DatabaseName, ssl)
}

type ValidationErrors []string

func (v ValidationErrors) Error() string {
	return strings.Join(v, ", ")
}
