package worker

import (
	"context"
	"log"
	"time"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/bootstrap"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
)

type PaymentWorker struct {
	PaymentRepo   port.PaymentRepository
	OrderRepo     port.OrderRepository
	OrderItemRepo port.OrderItemRepository
	ProductRepo   port.ProductRepository
}

func NewPaymentWorker(b *bootstrap.Bootstrap) *PaymentWorker {
	return &PaymentWorker{
		PaymentRepo:   b.PaymentRepo,
		OrderRepo:     b.OrderRepo,
		OrderItemRepo: b.OrderItemRepo,
		ProductRepo:   b.ProductRepo,
	}
}

func (w *PaymentWorker) Run() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.cancelExpiredPayments()
		}
	}
}

func (w *PaymentWorker) cancelExpiredPayments() {
	ctx := context.Background()
	now := time.Now()

	payments, err := w.PaymentRepo.FindExpiredPayments(ctx, now)
	if err != nil {
		log.Println("Error finding expired payments:", err)
		return
	}

	for _, payment := range payments {
		log.Println("Cancelling payment for Order ID:", payment.OrderID)
		go func(data *domain.Payment) {
			err := w.ProcessExpiredPayments(ctx, data)
			if err != nil {
				log.Println("Error processing expired payment:", err)
			}
		}(payment)
	}
}

func (w *PaymentWorker) ProcessExpiredPayments(ctx context.Context, payment *domain.Payment) error {
	err := w.PaymentRepo.Update(ctx, payment.ID, &domain.Payment{PaymentStatus: "failed"})
	if err != nil {
		log.Println("Error updating payment status:", err)
	}

	items, err := w.OrderItemRepo.Finds(ctx, map[string]interface{}{"order_id": payment.OrderID})
	if err != nil {
		log.Println("Error finding order items:", err)
	}

	for _, item := range items {
		product, err := w.ProductRepo.FindOne(ctx, item.ProductID)
		if err != nil {
			log.Println("Error finding product:", err)

		}

		newStock := product.Stock + item.Quantity
		err = w.ProductRepo.UpdateStock(ctx, item.ProductID, newStock)
		if err != nil {
			log.Println("Error updating stock:", err)

		}
	}

	err = w.OrderRepo.Update(ctx, payment.OrderID, &domain.Order{Status: "cancelled"})
	if err != nil {
		log.Println("Error updating order status:", err)
	}

	return nil

}
