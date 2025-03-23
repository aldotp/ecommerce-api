package config

import "github.com/spf13/viper"

// Redis related configuration
func RedisAddr() string {
	return viper.GetString("REDIS_ADDR")
}
func RedisPassword() string {
	return viper.GetString("REDIS_PASSWORD")
}
