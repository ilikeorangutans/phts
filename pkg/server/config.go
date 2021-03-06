package server

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ilikeorangutans/phts/storage"
	"github.com/pkg/errors"
)

// Config is the server configuration
type Config struct {
	ServerURL           string
	AdminEmail          string
	AdminPassword       string
	InitialUser         string
	InitialUserPassword string
	Bind                string
	DatabaseHost        string
	DatabaseUser        string
	DatabasePassword    string
	DatabaseName        string
	DatabaseSSL         bool
	StorageEngine       string
	BucketName          string
	ProjectID           string
	MinioAccessKey      string
	MinioSecretKey      string
	MinioEndpoint       string
	MinioUseSSL         bool
	SmtpHost            string
	SmtpPort            int
	SmtpUser            string
	SmtpPassword        string
	SmtpFrom            string
	// FrontendStaticFilePath is the path where the frontend ui files can be found
	FrontendStaticFilePath string
	// AdminStaticFilePath is the path where the admin ui files can be found
	AdminStaticFilePath string
	// JWTSecret is used to encrypt JWT settings
	JWTSecret string
}

func (c Config) Validate() error {
	errors := ValidationErrors{}
	if c.ServerURL == "" {
		errors = append(errors, "PHTS_SERVER_URL not provided")
	}
	if c.Bind == "" {
		errors = append(errors, "PHTS_BIND not provided")
	}
	if c.DatabaseHost == "" {
		errors = append(errors, "PHTS_DB_HOST not provided")
	}
	if c.AdminEmail == "" || c.AdminPassword == "" {
		errors = append(errors, "admin email and password must be provided")
	}

	if c.JWTSecret == "" {
		log.Printf("No JWT secret set! This means sessions will be invalidated after server restarts.")
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

func (c Config) StorageBackend(ctx context.Context) (storage.Backend, error) {
	var backend storage.Backend
	var err error
	switch c.StorageEngine {
	case "gcs":
		backend, err = storage.NewGCSBackend(c.ProjectID, ctx, c.BucketName)
	case "minio":
		backend, err = storage.NewMinIOBackend(c.MinioEndpoint, c.MinioAccessKey, c.MinioSecretKey, c.BucketName, c.MinioUseSSL)
	default:
		backend = storage.NewFileBackend("tmp")
	}

	return backend, errors.Wrap(err, "could not instantiate storage backend")
}

type ValidationErrors []string

func (v ValidationErrors) Error() string {
	return strings.Join(v, ", ")
}
