package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App   *App
		Token *Token
		DB    *DB
		HTTP  *HTTP
	}

	App struct {
		Name string
		Env  string
	}

	Token struct {
		Duration string
	}

	Redis struct {
		Addr     string
		Password string
	}

	DB struct {
		Connection string
		Host       string
		Port       string
		User       string
		Password   string
		Name       string
	}

	HTTP struct {
		Env            string
		URL            string
		Port           string
		AllowedOrigins string
		ReadTimeout    string
		WriteTimeout   string
		MaxHeaderBytes string
	}

	GCS struct {
		Credential string
		BucketName string
	}

	MongoDB struct {
		User     string
		Password string
		Host     string
		Port     string
		DBName   string
	}
)

// New creates a new container instance
func New() (*Config, error) {
	if os.Getenv("APP_ENV") != "production" {
		LoadConfig()
	}

	app := &App{
		Name: AppName(),
		Env:  AppEnv(),
	}

	token := &Token{
		Duration: TokenDuration(),
	}

	db := &DB{
		Connection: DBConnection(),
		Host:       DBHost(),
		Port:       DBPort(),
		User:       DBUser(),
		Password:   DBPassword(),
		Name:       DBName(),
	}

	http := &HTTP{
		Env:            AppEnv(),
		URL:            HTTPURL(),
		Port:           HTTPPort(),
		AllowedOrigins: HTTPAllowedOrigins(),
	}

	return &Config{
		app,
		token,
		db,
		http,
	}, nil
}

func LoadConfig() {
	viper.AutomaticEnv()

	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using system environment variables")
	} else {
		log.Println("Using .env file:", viper.ConfigFileUsed())
	}
}
