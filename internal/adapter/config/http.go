package config

import "github.com/spf13/viper"

// HTTP related configuration
func HTTPPort() string {
	return viper.GetString("HTTP_PORT")
}

func HTTPURL() string {
	return viper.GetString("HTTP_URL")
}

func HTTPAllowedOrigins() string {
	return viper.GetString("HTTP_ALLOWED_ORIGINS")
}
