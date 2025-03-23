package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/bootstrap"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/rabbitmq"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/aldotp/ecommerce-go-api/internal/core/service"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type worker struct {
	log             *zap.Logger
	rabbitMqService rabbitmq.RabbitMqInterface
	orderSvc        port.OrderService
	productSvc      port.ProductService
}

func NewOrderWorker(b *bootstrap.Bootstrap) *worker {
	return &worker{
		log:             b.Log,
		rabbitMqService: b.RabbitMQ,
		orderSvc:        service.NewOrderService(b.PaymentRepo, b.OrderRepo),
		productSvc:      service.NewProductService(b.ProductRepo, b.Cache),
	}
}

func (h *worker) UpdateStatusOrder() {
	ctx := context.Background()
	go h.WorkerUpdateStatus(ctx)
}

func (h *worker) WorkerUpdateStatus(ctx context.Context) {
	request := rabbitmq.RabbitMqConsumeRequest{
		QueueName:    consts.QueueUpdateStock,
		ConsumerName: fmt.Sprintf("go-%s", consts.QueueUpdateStock),
	}

	chClosedCh := make(chan *amqp.Error)

	msgs, err := h.rabbitMqService.Consume(request, chClosedCh)
	if err != nil {
		h.log.Error("failed to consume messages", zap.Error(err), zap.String("queue_name", request.QueueName), zap.String("from", "worker.send_email"))
		return
	}

	for {
		select {
		case amqErr := <-chClosedCh:
			h.log.Warn("channel closed by abnormal shutdown", zap.String("queue_name", request.QueueName), zap.Any("error", amqErr))
			time.Sleep(1 * time.Second)

			chClosedCh = make(chan *amqp.Error)
			msgs, err = h.rabbitMqService.Consume(request, chClosedCh)
			if err != nil {
				h.log.Error("failed to reconnect to RabbitMQ", zap.Error(err), zap.String("queue_name", request.QueueName))
				continue
			}

			h.log.Info("RabbitMQ channel reconnected", zap.String("queue_name", request.QueueName))

		case m := <-msgs:
			if m.Body == nil {
				_ = m.Ack(false)
				continue
			}

			h.log.Debug("message received", zap.String("queue_name", request.QueueName), zap.Any("message", string(m.Body)))

			var data dto.UpdateOrderStatus
			if err := json.Unmarshal(m.Body, &data); err != nil {
				h.log.Error("failed to unmarshal message body", zap.Error(err), zap.String("queue_name", request.QueueName))
				_ = m.Nack(false, false)
				continue
			}

			if h.rabbitMqService.IsClosed() {
				h.log.Warn("RabbitMQ channel closed, message will be requeued", zap.String("queue_name", request.QueueName))
				_ = m.Nack(false, true)
				continue
			}

			err := h.orderSvc.UpdateStatusOrder(ctx, data.OrderID, data.Status)
			if err != nil {
				h.log.Error("failed to update order status", zap.String("order_id", fmt.Sprintf("%d", data.OrderID)), zap.Error(err), zap.String("queue_name", request.QueueName), zap.Any("data", data))
				_ = m.Nack(false, true)
				continue
			} else {
				h.log.Info("update order status successfully", zap.String("order_id", fmt.Sprintf("%d", data.OrderID)), zap.String("queue_name", request.QueueName), zap.Any("data", data))
				_ = m.Ack(false)
			}
		}
	}
}
