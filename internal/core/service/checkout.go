package service

import (
	"context"
	"errors"
	"time"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"
)

type CheckoutService struct {
	ProductRepo   port.ProductRepository
	OrderRepo     port.OrderRepository
	OrderItemRepo port.OrderItemRepository
	CartRepo      port.CartRepository
	CartItemRepo  port.CartItemRepository
	PaymentRepo   port.PaymentRepository
}

func NewCheckoutService(
	productRepo port.ProductRepository,
	orderRepo port.OrderRepository,
	orderItemRepo port.OrderItemRepository,
	cartRepo port.CartRepository,
	cartItemRepo port.CartItemRepository,
	paymentRepo port.PaymentRepository,
) *CheckoutService {
	return &CheckoutService{
		ProductRepo:   productRepo,
		OrderRepo:     orderRepo,
		OrderItemRepo: orderItemRepo,
		CartRepo:      cartRepo,
		CartItemRepo:  cartItemRepo,
		PaymentRepo:   paymentRepo,
	}
}

func (s *CheckoutService) Checkout(ctx context.Context, userID int, paymentMethod string) (*dto.CheckoutResponse, error) {
	cart, err := s.CartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, consts.ErrEmptyCart
	}

	items, err := s.CartItemRepo.Finds(ctx, map[string]interface{}{"cart_id": cart.ID})
	if err != nil {
		return nil, err
	}

	var totalPrice float64
	for _, item := range items {
		product, err := s.ProductRepo.FindOne(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, errors.New("product not found")
		}
		if product.Stock < item.Quantity {
			return nil, errors.New("insufficient stock")
		}

		totalPrice += product.Price * float64(item.Quantity)
	}

	tNow := time.Now()
	order := &domain.Order{
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     "pending",
		CreatedAt:  tNow,
	}

	err = s.OrderRepo.Store(ctx, order)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		product, err := s.ProductRepo.FindOne(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}

		if err := s.OrderItemRepo.Store(ctx, &domain.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}); err != nil {
			return nil, err
		}

		newStock := product.Stock - item.Quantity
		if err := s.ProductRepo.UpdateStock(ctx, item.ProductID, newStock); err != nil {
			return nil, err
		}
	}

	expiredAt := tNow.Add(10 * time.Minute)
	err = s.PaymentRepo.Store(ctx, &domain.Payment{
		OrderID:       order.ID,
		PaymentMethod: paymentMethod,
		PaymentStatus: "pending",
		UpdatedAt:     tNow,
		CreatedAt:     tNow,
		ExpiredAt:     expiredAt,
	})
	if err != nil {
		return nil, err
	}

	if err := s.ClearCart(ctx, userID, cart.ID); err != nil {
		return nil, err
	}

	return &dto.CheckoutResponse{
		OrderID:       order.ID,
		PaymentMethod: paymentMethod,
		Total:         int(totalPrice),
	}, nil
}

func (s *CheckoutService) ClearCart(ctx context.Context, userID int, cartID int) error {
	err := s.CartRepo.DeleteByUserID(ctx, userID)
	if err != nil {
		return err
	}

	err = s.CartItemRepo.DeleteByCartID(ctx, cartID)
	if err != nil {
		return err
	}

	return nil
}
