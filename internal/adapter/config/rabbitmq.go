package config

import "github.com/spf13/viper"

func RabbitMQHost() string {
	return viper.GetString("RABBITMQ_HOST")
}

func RabbitMQUser() string {
	return viper.GetString("RABBITMQ_USER")
}

func RabbitMQPassword() string {
	return viper.GetString("RABBITMQ_PASSWORD")
}

func RabbitMQVhost() string {
	return viper.GetString("RABBITMQ_VHOST")
}
