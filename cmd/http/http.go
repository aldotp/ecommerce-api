package http

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/bootstrap"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/handler/http"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/router"
	"github.com/aldotp/ecommerce-go-api/internal/core/service"
)

func RunHTTPServer(ctx context.Context) {
	f := bootstrap.NewBootstrap(ctx).BuildRestBootstrap()

	// Services
	userService := service.NewUserService(f.UserRepo, f.Cache, f.Token, f.Log, f.BalanceRepo)
	authService := service.NewAuthService(f.UserRepo, f.Token, f.Log)
	productService := service.NewProductService(f.ProductRepo, f.Cache)
	categoryService := service.NewCategoryService(f.CategoryRepo, f.Cache)
	cartService := service.NewCartService(f.CartItemRepo, f.CartRepo, f.OrderRepo, f.OrderItemRepo, f.ProductRepo)
	checkoutService := service.NewCheckoutService(f.ProductRepo, f.OrderRepo, f.OrderItemRepo, f.CartRepo, f.CartItemRepo, f.PaymentRepo)
	balanceService := service.NewBalanceService(f.BalanceRepo, f.Cache)
	paymentService := service.NewPaymentService(f.PaymentRepo, f.OrderRepo, f.RabbitMQ, f.BalanceRepo, balanceService)
	orderService := service.NewOrderService(f.PaymentRepo, f.OrderRepo)

	// Handlers
	userHandler := http.NewUserHandler(userService, f.Log)
	authHandler := http.NewAuthHandler(authService, f.Log)
	productHandler := http.NewProductHandler(productService, f.Log)
	categoryHandler := http.NewCategoryHandler(categoryService, f.Log)
	cartHandler := http.NewCartHandler(cartService, f.Log)
	checkoutHandler := http.NewCheckoutHandler(checkoutService, f.Log)
	paymentHandler := http.NewPaymentHandler(paymentService, f.Log)
	orderHandler := http.NewOrderHandler(orderService, f.Log)
	balanceHandler := http.NewBalanceHandler(balanceService, f.Log)

	// HTTP server
	routes, err := router.NewRouter(
		f.Token,
		authHandler,
		userHandler,
		productHandler,
		categoryHandler,
		cartHandler,
		checkoutHandler,
		paymentHandler,
		orderHandler,
		balanceHandler,
	)
	if err != nil {
		slog.Error("Error creating router", "error", err)
	}

	// Start server
	listenAddr := fmt.Sprintf("%s:%s", f.Config.HTTP.URL, f.Config.HTTP.Port)
	slog.Info("Starting the HTTP server", "listen_address", listenAddr)
	err = routes.Serve(listenAddr)
	if err != nil {
		slog.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}
