package config

import "github.com/spf13/viper"

// App related configuration
func AppName() string {
	return viper.GetString("APP_NAME")
}

func AppEnv() string {
	return viper.GetString("APP_ENV")
}

func AppVersion() string {
	return viper.GetString("APP_VERSION")
}

func SecretKey() string {
	return viper.GetString("SECRET_KEY")
}

func RefreshKey() string {
	return viper.GetString("REFRESH_KEY")
}
