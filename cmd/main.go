package main

import (
	"context"
	"log"

	"github.com/aldotp/ecommerce-go-api/cmd/consumer"
	"github.com/aldotp/ecommerce-go-api/cmd/http"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/config"
	"github.com/spf13/cobra"

	_ "github.com/joho/godotenv/autoload"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer {your_token}" to authenticate

func main() {
	ctx := context.Background()
	config.LoadConfig()
	rootCmd := &cobra.Command{}

	// Restful API Command
	restCmd := cobra.Command{
		Use:   "rest",
		Short: "Rest is a command to start Restful Api server",
		Run: func(cmd *cobra.Command, args []string) {
			http.RunHTTPServer(ctx)
		},
	}

	// define consumer command
	consumerCmd := cobra.Command{
		Use:   "consumer",
		Short: "Consumer is a command to start consumer worker",
	}

	consumerExpiredPaymentCmd := cobra.Command{
		Use:   "expired_payment",
		Short: "Consumer is a command to start  UpdateStock consumer server",
		Run: func(cmd *cobra.Command, args []string) {
			consumer.RunExpiredPaymentConsumer(ctx)
		},
	}

	consumerUpdateStockCmd := cobra.Command{
		Use:   "update_status",
		Short: "Consumer is a command to start  UpdateStock consumer server",
		Run: func(cmd *cobra.Command, args []string) {
			consumer.RunUpdateStatusOrderConsumer(ctx)
		},
	}

	rootCmd.AddCommand(
		&restCmd,
		&consumerCmd,
	)

	consumerCmd.AddCommand(
		&consumerExpiredPaymentCmd,
		&consumerUpdateStockCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute command: %v", err)
	}

}
