package service

import (
	"context"
	"time"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"
)

type CartService struct {
	CartItemRepo  port.CartItemRepository
	CartRepo      port.CartRepository
	OrderRepo     port.OrderRepository
	OrderItemRepo port.OrderItemRepository
	ProductRepo   port.ProductRepository
}

func NewCartService(
	cartItemRepo port.CartItemRepository,
	cartRepo port.CartRepository,
	orderRepo port.OrderRepository,
	orderItemRepo port.OrderItemRepository,
	productRepo port.ProductRepository,
) *CartService {
	return &CartService{
		CartItemRepo:  cartItemRepo,
		CartRepo:      cartRepo,
		OrderRepo:     orderRepo,
		OrderItemRepo: orderItemRepo,
		ProductRepo:   productRepo,
	}
}

func (s *CartService) GetCart(ctx context.Context, userID int) (dto.CartResponse, error) {
	var (
		response dto.CartResponse
	)

	cart, err := s.CartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return response, err
	}

	if cart == nil {
		return response, nil
	}

	cartItems, err := s.CartItemRepo.Finds(ctx, map[string]interface{}{"cart_id": cart.ID})
	if err != nil {
		return response, err
	}

	var totalPrice, quantity float64
	for _, item := range cartItems {
		product, err := s.ProductRepo.FindOne(ctx, item.ProductID)
		if err != nil {
			return response, err
		}

		totalPrice += product.Price * float64(item.Quantity)
		quantity += float64(item.Quantity)

		response.Items = append(response.Items, dto.CartItemResponse{
			Name:      product.Name,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
	}

	response.TotalProducts = len(response.Items)
	response.TotalItems = int(quantity)
	response.TotalPrice = totalPrice

	return response, nil
}

func (s *CartService) AddToCart(ctx context.Context, userID int, productID int, quantity int) error {
	cart, err := s.CartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if cart == nil {
		cart = &domain.Cart{
			UserID: userID,
		}
		if err := s.CartRepo.Store(ctx, cart); err != nil {
			return err
		}
	}

	product, err := s.ProductRepo.FindOne(ctx, productID)
	if err != nil {
		return err
	}

	existCartItem, err := s.CartItemRepo.FindOneByFilters(ctx, map[string]interface{}{"cart_id": cart.ID, "product_id": productID})
	if err != nil {
		return err
	}

	tNow := time.Now()
	cartItem := domain.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  quantity,
		CreatedAt: tNow,
		UpdatedAt: tNow,
	}

	if existCartItem == nil {
		err = s.CartItemRepo.Store(ctx, &cartItem)
		if err != nil {
			return err
		}
	} else {
		if product.Stock < quantity {
			return consts.ErrInsufficientStock
		}

		cartItem.Quantity += existCartItem.Quantity
		err := s.CartItemRepo.Update(ctx, existCartItem.ID, cartItem)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *CartService) RemoveFromCart(ctx context.Context, userID int, productID int) error {
	cart, err := s.CartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if cart == nil {
		return consts.ErrEmptyCart
	}

	existCartItem, err := s.CartItemRepo.FindOneByFilters(ctx, map[string]interface{}{"cart_id": cart.ID, "product_id": productID})
	if err != nil {
		return err
	}

	if existCartItem == nil {
		return consts.ErrDataNotFound
	}

	err = s.CartItemRepo.DeleteByProductID(ctx, productID)
	if err != nil {
		return err
	}

	return nil
}

func (s *CartService) UpdateCart(ctx context.Context, userID int, request dto.UpdateCartRequest) error {
	cart, err := s.CartRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if cart == nil {
		return consts.ErrEmptyCart
	}

	cartItem, err := s.CartItemRepo.FindOneByFilters(ctx, map[string]interface{}{"cart_id": cart.ID, "product_id": request.ProductID})
	if err != nil {
		return err
	}

	if cartItem == nil {
		return consts.ErrDataNotFound
	}

	product, err := s.ProductRepo.FindOne(ctx, request.ProductID)
	if err != nil {
		return err
	}

	if product.Stock < request.Quantity {
		return consts.ErrInsufficientStock
	}

	cartItem.Quantity = request.Quantity
	err = s.CartItemRepo.Update(ctx, cartItem.ID, *cartItem)
	if err != nil {
		return err
	}

	return nil
}
