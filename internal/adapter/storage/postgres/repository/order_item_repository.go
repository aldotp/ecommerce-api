package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/storage/postgres"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type OrderItemRepository struct {
	db        *postgres.DB
	TableName string
}

func NewOrderItemRepository(db *postgres.DB) *OrderItemRepository {
	return &OrderItemRepository{
		db:        db,
		TableName: "order_items",
	}
}

// FindOne retrieves a single Categories by ID
func (r *OrderItemRepository) FindOne(ctx context.Context, id int) (*domain.OrderItem, error) {
	var orderItem domain.OrderItem

	query := r.db.QueryBuilder.Select("id", "order_id", "product_id", "quantity", "price").
		From(r.TableName).
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&orderItem.ID,
		&orderItem.OrderID,
		&orderItem.ProductID,
		&orderItem.Quantity,
		&orderItem.Price,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &orderItem, nil
}

// Store inserts a new Categories into the database
func (r *OrderItemRepository) Store(ctx context.Context, data *domain.OrderItem) error {
	query := r.db.QueryBuilder.Insert(r.TableName).
		Columns("order_id", "product_id", "quantity", "price").
		Values(data.OrderID, data.ProductID, data.Quantity, data.Price).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(&data.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderItemRepository) Update(ctx context.Context, id int, updatedData domain.OrderItem) error {
	query := r.db.QueryBuilder.Update(r.TableName).
		Set("status", sq.Expr("COALESCE(?, quantity)", updatedData.Quantity)).
		Set("price", sq.Expr("COALESCE(?, price)", updatedData.Price)).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&updatedData.ID,
		&updatedData.OrderID,
		&updatedData.ProductID,
		&updatedData.Quantity,
		&updatedData.Price,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderItemRepository) Delete(ctx context.Context, id int) error {
	query := r.db.QueryBuilder.Delete(r.TableName).
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderItemRepository) Finds(ctx context.Context, filter map[string]interface{}) ([]domain.OrderItem, error) {
	query := r.db.QueryBuilder.Select("*").From(r.TableName)

	for key, value := range filter {
		query = query.Where(sq.Eq{key: value})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderItems []domain.OrderItem
	for rows.Next() {
		var orderItem domain.OrderItem
		err := rows.Scan(
			&orderItem.ID,
			&orderItem.OrderID,
			&orderItem.ProductID,
			&orderItem.Quantity,
			&orderItem.Price,
		)
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, orderItem)
	}

	return orderItems, nil
}
