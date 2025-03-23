package bootstrap

import (
	"context"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/config"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/rabbitmq"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/storage/postgres"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"

	"go.uber.org/zap"
)

type Bootstrap struct {
	ctx        context.Context
	Log        *zap.Logger
	Config     *config.Config
	PostgresDB *postgres.DB
	RabbitMQ   rabbitmq.RabbitMqInterface

	UserRepo      port.UserRepository
	OrderRepo     port.OrderRepository
	ProductRepo   port.ProductRepository
	PaymentRepo   port.PaymentRepository
	OrderItemRepo port.OrderItemRepository
	CartItemRepo  port.CartItemRepository
	CartRepo      port.CartRepository
	CategoryRepo  port.CategoryRepository
	BalanceRepo   port.BalanceRepository

	Token port.TokenInterface
	Cache port.CacheInterface
}

func NewBootstrap(ctx context.Context) *Bootstrap {
	return &Bootstrap{
		ctx: ctx,
	}
}
