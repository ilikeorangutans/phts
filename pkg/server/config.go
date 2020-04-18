package server

import "fmt"

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

func (c Config) DatabaseConnectionString() string {
	ssl := "enable"
	if !c.DatabaseSSL {
		ssl = "disable"
	}
	return fmt.Sprintf("user=%s host=%s password=%s dbname=%s sslmode=%s", c.DatabaseUser, c.DatabaseHost, c.DatabasePassword, c.DatabaseName, ssl)
}
