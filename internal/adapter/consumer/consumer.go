package consumer

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/bootstrap"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/handler/worker"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/rabbitmq"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"
	"go.uber.org/zap"
)

type Operation func()

type consumer struct {
	bootstrap *bootstrap.Bootstrap
	rmq       rabbitmq.RabbitMqInterface
	log       *zap.Logger
}

type Consumer interface {
	Init()
	Start(Operation ...Operation)
	Stop() error

	UpdateStatusOrderConsumer()
	ExpiredPaymentConsumer()
}

func NewConsumer(b *bootstrap.Bootstrap) Consumer {
	return &consumer{
		bootstrap: b,
		rmq:       b.RabbitMQ,
		log:       b.Log.With(zap.String("from", "consumer")),
	}
}

func (c *consumer) Init() {
	rmqPreparations := []rabbitmq.Preparation{
		{
			IsBindingExchange: false,
			QueueName:         consts.QueueUpdateStock,
		},
		{
			Exchange: rabbitmq.RabbitMQExchange{
				Name: consts.ExchangeUpdateStock,
				Kind: "direct",
			},
		},
	}

	for _, p := range rmqPreparations {
		if p.Exchange.Name != "" {
			c.rmq.DeclareExchange(p.Exchange)
		}

		if p.QueueName != "" {
			c.rmq.DeclareQueue(p.QueueName)
		}

		if p.IsBindingExchange {
			c.rmq.BindingQueue(p.Exchange.Name, p.QueueName)
		}
	}

	c.log.Info("Consumer initialized...")
}

func (c *consumer) Start(operations ...Operation) {

	for _, op := range operations {
		op()
	}

	sigChan := make(chan os.Signal, 1)
	defer close(sigChan)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-sigChan
	c.Stop()

	c.log.Info("Consumer stopped...")

	os.Exit(0)
}

func (c *consumer) Stop() error {
	time.Sleep(1 * time.Second)

	c.log.Info("Consumer stopped...")
	return nil
}

func (c *consumer) UpdateStatusOrderConsumer() {
	c.log.Info("Consumer registered...", zap.String("job_name", "update_status"))

	worker.NewOrderWorker(c.bootstrap).UpdateStatusOrder()

}

func (c *consumer) ExpiredPaymentConsumer() {
	c.log.Info("Consumer registered...", zap.String("job_name", "expired_payment"))

	worker.NewPaymentWorker(c.bootstrap).Run()
}
