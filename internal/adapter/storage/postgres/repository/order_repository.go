package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/storage/postgres"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type OrderRepository struct {
	db        *postgres.DB
	TableName string
}

func NewOrderRepository(db *postgres.DB) *OrderRepository {
	return &OrderRepository{
		db:        db,
		TableName: "orders",
	}
}

func (r *OrderRepository) Finds(ctx context.Context, filter map[string]interface{}) ([]domain.Order, error) {
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

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.TotalPrice,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// FindOne retrieves a single Categories by ID
func (r *OrderRepository) FindOne(ctx context.Context, id int, userID int) (*domain.Order, error) {
	var Order domain.Order

	query := r.db.QueryBuilder.Select("id", "user_id", "total_price", "status", "created_at").
		From(r.TableName).
		Where(sq.Eq{"id": id}).
		Where(sq.Eq{"user_id": userID}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&Order.ID,
		&Order.UserID,
		&Order.TotalPrice,
		&Order.Status,
		&Order.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &Order, nil
}

// Store inserts a new Categories into the database
func (r *OrderRepository) Store(ctx context.Context, data *domain.Order) error {
	query := r.db.QueryBuilder.Insert(r.TableName).
		Columns("user_id", "total_price", "status", "created_at").
		Values(data.UserID, data.TotalPrice, data.Status, data.CreatedAt).
		Suffix("RETURNING id, user_id, total_price, status, created_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&data.ID,
		&data.UserID,
		&data.TotalPrice,
		&data.Status,
		&data.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

// Update modifies an existing Categories in the database
func (r *OrderRepository) Update(ctx context.Context, id int, updatedData *domain.Order) error {
	query := r.db.QueryBuilder.Update(r.TableName).
		Set("status", sq.Expr("COALESCE(?, status)", updatedData.Status)).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, user_id, total_price, status, created_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&updatedData.ID,
		&updatedData.UserID,
		&updatedData.TotalPrice,
		&updatedData.Status,
		&updatedData.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) Delete(ctx context.Context, id int) error {
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
