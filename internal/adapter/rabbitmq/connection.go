package rabbitmq

import (
	"fmt"
	"log"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	mqConfig struct {
		Username string
		Password string
		Vhost    string
		Host     string
	}

	rabbitMqConn struct {
		mqConfig
	}
)

var (
	rmqConn *amqp.Connection
	rmqChan *amqp.Channel
)

func CreateConnection() (*amqp.Connection, *amqp.Channel) {
	conf := mqConfig{
		Host:     config.RabbitMQHost(),
		Username: config.RabbitMQUser(),
		Password: config.RabbitMQPassword(),
		Vhost:    config.RabbitMQVhost(),
	}

	rabbitMqConn := rabbitMqConn{mqConfig: conf}
	if (rmqConn == nil && rmqChan == nil) || rmqChan.IsClosed() || rmqConn.IsClosed() {
		rmqConn, rmqChan = rabbitMqConn.connect()
	}

	return rmqConn, rmqChan
}

func (conf rabbitMqConn) connect() (*amqp.Connection, *amqp.Channel) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s/%s", conf.Username, conf.Password, conf.Host, conf.Vhost)
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err)
		panic(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err)
		panic(err)
	}

	return conn, channel
}
