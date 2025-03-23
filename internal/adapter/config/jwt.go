package config

import "github.com/spf13/viper"

// Token related configuration
func TokenDuration() string {
	return viper.GetString("TOKEN_DURATION")
}

func AccessTokenExpired() int {
	return viper.GetInt("ACCESS_TOKEN_EXPIRED")
}

func RefreshTokenExpired() int {
	return viper.GetInt("REFRESH_TOKEN_EXPIRED")
}
