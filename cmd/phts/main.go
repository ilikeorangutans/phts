package main

import (
	"context"
	"log"

	"github.com/ilikeorangutans/phts/pkg/server"
	"github.com/ilikeorangutans/phts/version"
	"github.com/spf13/viper"

	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

func parseConfig() server.Config {
	return server.Config{
		ServerURL:              viper.GetString("server_url"),
		AdminEmail:             viper.GetString("admin_email"),
		AdminPassword:          viper.GetString("admin_password"),
		Bind:                   viper.GetString("bind"),
		DatabaseHost:           viper.GetString("db_host"),
		DatabaseUser:           viper.GetString("db_user"),
		DatabasePassword:       viper.GetString("db_password"),
		DatabaseName:           viper.GetString("db_database"),
		DatabaseSSL:            viper.GetBool("db_ssl"),
		StorageEngine:          viper.GetString("storage_engine"),
		BucketName:             viper.GetString("minio_bucket"),
		MinioAccessKey:         viper.GetString("minio_access_key"),
		MinioSecretKey:         viper.GetString("minio_secret_key"),
		MinioEndpoint:          viper.GetString("minio_endpoint"),
		MinioUseSSL:            viper.GetBool("minio_use_ssl"),
		SmtpHost:               viper.GetString("smtp_host"),
		SmtpPort:               viper.GetInt("smtp_port"),
		SmtpUser:               viper.GetString("smtp_user"),
		SmtpPassword:           viper.GetString("smtp_password"),
		SmtpFrom:               viper.GetString("smtp_from"),
		FrontendStaticFilePath: viper.GetString("frontend_static_file_path"),
		AdminStaticFilePath:    viper.GetString("admin_static_file_path"),
	}
}

func setupEnvVars() {
	viper.SetEnvPrefix("phts")
	viper.AutomaticEnv()

	defaults := map[string]interface{}{
		"bind":        ":8080",
		"db_ssl":      false,
		"db_host":     "",
		"db_user":     "",
		"db_password": "",
		"db_database": "phts",

		"storage_engine": "file",

		"minio_bucket":     "",
		"minio_access_key": "",
		"minio_secret_key": "",
		"minio_endpoint":   "",
		"minio_use_ssl":    false,

		"frontend_static_file_path": "ui/dist/frontend/",
		"admin_static_file_path":    "ui/dist/admin/",
	}

	for key, value := range defaults {
		viper.SetDefault(key, value)
	}
}

func main() {
	log.Printf("phts starting up, version %s built on %s...", version.Sha, version.BuildTime)

	ctx := context.Background()

	setupEnvVars()

	config := parseConfig()
	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}
	main, err := server.NewMain(ctx, config)
	if err != nil {
		log.Fatalf("%s", err)
	}
	if err := main.Run(ctx); err != nil {
		log.Fatalf("error running:\n%+v", err)
	}
}
