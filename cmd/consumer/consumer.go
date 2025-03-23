package consumer

import (
	"context"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/bootstrap"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/consumer"
)

func RunUpdateStatusOrderConsumer(ctx context.Context) {
	b := bootstrap.NewBootstrap(ctx).BuildConsumerUpdateOrderStatusBootstrap()

	con := consumer.NewConsumer(b)
	con.Init()
	con.Start(con.UpdateStatusOrderConsumer)
}

func RunExpiredPaymentConsumer(ctx context.Context) {
	b := bootstrap.NewBootstrap(ctx).BuildConsumerExpiredPaymentBootstrap()

	con := consumer.NewConsumer(b)
	con.Init()
	con.Start(con.ExpiredPaymentConsumer)
}
