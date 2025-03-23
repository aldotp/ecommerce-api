package config

import "github.com/spf13/viper"

// DB related configuration
func DBConnection() string {
	return viper.GetString("DB_CONNECTION")
}

func DBHost() string {
	return viper.GetString("DB_HOST")
}

func DBPort() string {
	return viper.GetString("DB_PORT")
}

func DBUser() string {
	return viper.GetString("DB_USER")
}

func DBPassword() string {
	return viper.GetString("DB_PASSWORD")
}

func DBName() string {
	return viper.GetString("DB_NAME")
}
