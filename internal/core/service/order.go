package service

import (
	"context"

	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"
)

type OrderService struct {
	OrderRepo   port.OrderRepository
	PaymentRepo port.PaymentRepository
}

func NewOrderService(paymentRepo port.PaymentRepository, orderRepo port.OrderRepository) *OrderService {
	return &OrderService{
		PaymentRepo: paymentRepo,
		OrderRepo:   orderRepo,
	}
}

func (s *OrderService) UpdateStatusOrder(ctx context.Context, orderID int, status string) error {
	var order domain.Order

	switch status {
	case consts.Paid:
		s.updateStatusOnProccess(orderID, &order)
	}

	return s.OrderRepo.Update(ctx, orderID, &order)
}

func (s *OrderService) updateStatusOnProccess(orderID int, order *domain.Order) {
	*order = domain.Order{
		ID:     orderID,
		Status: consts.Paid,
	}
}

func (s *OrderService) GetOrder(ctx context.Context, orderID int, userID int) (*domain.Order, error) {
	return s.OrderRepo.FindOne(ctx, orderID, userID)
}

func (s *OrderService) ListOrders(ctx context.Context, userId int) ([]domain.Order, error) {
	filter := map[string]interface{}{}

	if userId != 0 {
		filter["user_id"] = userId
	}

	return s.OrderRepo.Finds(ctx, filter)
}
