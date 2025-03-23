package service

import (
	"context"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/rabbitmq"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"
)

type PaymentService struct {
	OrderRepo   port.OrderRepository
	PaymentRepo port.PaymentRepository
	BalanceRepo port.BalanceRepository
	BalanceSvc  port.BalanceService
	rabbitmq    rabbitmq.RabbitMqInterface
}

func NewPaymentService(paymentRepo port.PaymentRepository, orderRepo port.OrderRepository, rabbitmq rabbitmq.RabbitMqInterface, balanceRepo port.BalanceRepository, balanceSvc port.BalanceService) *PaymentService {
	return &PaymentService{
		PaymentRepo: paymentRepo,
		OrderRepo:   orderRepo,
		rabbitmq:    rabbitmq,
		BalanceRepo: balanceRepo,
		BalanceSvc:  balanceSvc,
	}
}

func (s *PaymentService) MakePayment(ctx context.Context, userID int, orderID int) error {

	existPayment, err := s.PaymentRepo.FindByUserIDandOrderID(ctx, userID, orderID)
	if err != nil {
		return err
	}

	if existPayment.PaymentStatus == "completed" {
		return nil
	}

	order, err := s.OrderRepo.FindOne(ctx, orderID, userID)
	if err != nil {
		return err
	}

	if order == nil {
		return consts.ErrDataNotFound
	}

	if existPayment.PaymentMethod == "balance" {
		balance, err := s.BalanceRepo.GetBalance(ctx, uint64(userID))
		if err != nil {
			return err
		}

		if balance < float64(order.TotalPrice) {
			return consts.ErrInsufficientBalance
		}

		err = s.BalanceSvc.Withdraw(ctx, uint64(userID), float64(order.TotalPrice))
		if err != nil {
			return err
		}

	}

	payment := domain.Payment{
		OrderID:       orderID,
		PaymentStatus: "completed",
	}

	err = s.PaymentRepo.Update(ctx, orderID, &payment)
	if err != nil {
		return err
	}

	err = s.rabbitmq.Publish(ctx, rabbitmq.RabbitMqPublishRequest{
		QueueName: consts.QueueUpdateStock,
		Messages: dto.UpdateOrderStatus{
			OrderID: orderID,
			Status:  consts.Paid,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
