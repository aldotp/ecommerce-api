package bootstrap

import (
	"log/slog"
	"os"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/auth/jwt"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/config"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/rabbitmq"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/storage/postgres"
	postgresRepo "github.com/aldotp/ecommerce-go-api/internal/adapter/storage/postgres/repository"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/storage/redis"
	"github.com/aldotp/ecommerce-go-api/pkg/logger"
)

func (b *Bootstrap) setConfig() {
	config, err := config.New()
	if err != nil {
		panic(err)
	}

	b.Config = config
}

func (b *Bootstrap) setLogger() {
	logger, err := logger.InitLogger(config.AppEnv())
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	b.Log = logger
}

func (b *Bootstrap) setPostgresDB() {
	db, err := postgres.New(b.ctx, &config.DB{
		Connection: config.DBConnection(),
		User:       config.DBUser(),
		Password:   config.DBPassword(),
		Host:       config.DBHost(),
		Port:       config.DBPort(),
		Name:       config.DBName(),
	})
	if err != nil {
		panic(err)
	}

	b.PostgresDB = db
}

func (b *Bootstrap) setJWTToken() {
	token, err := jwt.New()
	if err != nil {
		slog.Error("Error initializing token service", "error", err)
		os.Exit(1)
	}

	b.Token = token
}

func (b *Bootstrap) setCache() {
	cache, err := redis.New(b.ctx, &config.Redis{
		Addr:     config.RedisAddr(),
		Password: config.RedisPassword(),
	})
	if err != nil {
		panic(err)
	}

	b.Cache = cache
}

func (b *Bootstrap) setRabbitMQ() {
	mqConn, mqCh := rabbitmq.CreateConnection()
	b.RabbitMQ = rabbitmq.New(mqConn, mqCh, b.Log)
}

func (b *Bootstrap) setRestApiRepository() {
	b.UserRepo = postgresRepo.NewUserRepository(b.PostgresDB)
	b.OrderRepo = postgresRepo.NewOrderRepository(b.PostgresDB)
	b.ProductRepo = postgresRepo.NewProductRepository(b.PostgresDB)
	b.PaymentRepo = postgresRepo.NewPaymentRepository(b.PostgresDB)
	b.OrderItemRepo = postgresRepo.NewOrderItemRepository(b.PostgresDB)
	b.CartItemRepo = postgresRepo.NewCartItemRepository(b.PostgresDB)
	b.CartRepo = postgresRepo.NewCartRepository(b.PostgresDB)
	b.CategoryRepo = postgresRepo.NewCategoryRepository(b.PostgresDB)
	b.BalanceRepo = postgresRepo.NewBalanceRepository(b.PostgresDB)
}

func (b *Bootstrap) SetUpdateStatusConsumerRepository() {
	b.OrderRepo = postgresRepo.NewOrderRepository(b.PostgresDB)
	b.PaymentRepo = postgresRepo.NewPaymentRepository(b.PostgresDB)
	b.ProductRepo = postgresRepo.NewProductRepository(b.PostgresDB)
}

func (b *Bootstrap) SetExpiredPaymentConsumerRepository() {
	b.OrderRepo = postgresRepo.NewOrderRepository(b.PostgresDB)
	b.PaymentRepo = postgresRepo.NewPaymentRepository(b.PostgresDB)
	b.OrderItemRepo = postgresRepo.NewOrderItemRepository(b.PostgresDB)
	b.ProductRepo = postgresRepo.NewProductRepository(b.PostgresDB)
}
