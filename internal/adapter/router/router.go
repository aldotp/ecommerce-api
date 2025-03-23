package router

import (
	_ "github.com/aldotp/ecommerce-go-api/docs"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/config"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/handler/http"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/middleware"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/aldotp/ecommerce-go-api/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	*gin.Engine
}

func NewRouter(
	token port.TokenInterface,
	authHandler *http.AuthHandler,
	userHandler *http.UserHandler,
	productHandler *http.ProductHandler,
	categoryHandler *http.CategoryHandler,
	cartHandler *http.CartHandler,
	checkoutHandler *http.CheckoutHandler,
	paymentHandler *http.PaymentHandler,
	orderHandler *http.OrderHandler,
	balanceHandler *http.BalanceHandler,
) (*Router, error) {

	// Set Gin mode
	if config.AppEnv() == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Custom validator
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		if err := v.RegisterValidation("user_role", userRoleValidator); err != nil {
			return nil, err
		}
	}

	util.Index(router, config.AppVersion(), config.AppName())
	util.Metrics(router)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API Routes
	api := router.Group("/api")
	v1 := api.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh-token", authHandler.RefreshToken)
		}

		register := v1.Group("/register")
		{
			register.POST("/", userHandler.Register)
		}

		profile := v1.Group("/profile")
		{
			authUser := profile.Group("/").Use(middleware.AuthMiddleware(token))
			{
				authUser.GET("/", userHandler.GetProfile)
			}
		}

		user := v1.Group("/users")
		{
			authUser := user.Group("/").Use(middleware.AuthMiddleware(token))
			{
				admin := authUser.Use(middleware.AdminMiddleware())
				{
					authUser.GET("/", userHandler.ListUsers)
					authUser.GET("/:id", userHandler.GetUser)
					admin.PUT("/:id", userHandler.UpdateUser)
					admin.DELETE("/:id", userHandler.DeleteUser)
				}
			}
		}

		product := v1.Group("/products")
		{
			authUser := product.Group("/").Use(middleware.AuthMiddleware(token))
			{

				authUser.GET("/", productHandler.ListProducts)
				authUser.GET("/:id", productHandler.GetProduct)

				admin := authUser.Use(middleware.AdminMiddleware())
				{
					admin.POST("/", productHandler.CreateProduct)
					admin.DELETE("/:id", productHandler.DeleteProduct)
					admin.PUT("/:id", productHandler.UpdateProduct)
				}

			}
		}

		category := v1.Group("/categories")
		{
			authUser := category.Group("/").Use(middleware.AuthMiddleware(token))
			{

				authUser.GET("/", categoryHandler.ListCategory)
				authUser.GET("/:id", categoryHandler.GetCategory)

				admin := authUser.Use(middleware.AdminMiddleware())
				{
					admin.POST("/", categoryHandler.CreateCategory)
					admin.DELETE("/:id", categoryHandler.DeleteCategory)
					admin.PUT("/:id", categoryHandler.UpdateCategory)
				}

			}
		}

		cart := v1.Group("/carts")
		{
			authUser := cart.Group("/").Use(middleware.AuthMiddleware(token))
			{
				authUser.POST("/", cartHandler.AddToCart)
				authUser.GET("/", cartHandler.ViewCart)
				authUser.DELETE("/", cartHandler.RemoveFromCart)
				authUser.PUT("/", cartHandler.UpdateCart)
			}
		}

		checkout := v1.Group("/checkout")
		{
			authUser := checkout.Group("/").Use(middleware.AuthMiddleware(token))
			{
				authUser.POST("/", checkoutHandler.Checkout)
			}
		}

		payment := v1.Group("/payments")
		{
			authUser := payment.Group("/").Use(middleware.AuthMiddleware(token))
			{
				authUser.POST("/pay", paymentHandler.Pay)
			}
		}

		order := v1.Group("/orders")
		{
			authUser := order.Group("/").Use(middleware.AuthMiddleware(token))
			{
				authUser.GET("", orderHandler.GetOrders)
				authUser.GET("/:id", orderHandler.GetOrderDetail)
			}
		}

		balance := v1.Group("/balance")
		{
			authUser := balance.Group("/").Use(middleware.AuthMiddleware(token))
			{
				authUser.POST("/deposit", balanceHandler.Deposit)
				authUser.POST("/transfer", balanceHandler.Transfer)
				authUser.GET("", balanceHandler.CheckBalance)
				authUser.POST("/withdraw", balanceHandler.Withdraw)
			}
		}
	}

	return &Router{
		router,
	}, nil
}

// Serve starts the HTTP server
func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
