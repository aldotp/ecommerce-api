package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/storage/postgres"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type PaymentRepository struct {
	db        *postgres.DB
	TableName string
}

func NewPaymentRepository(db *postgres.DB) *PaymentRepository {
	return &PaymentRepository{
		db:        db,
		TableName: "payments",
	}
}

func (r *PaymentRepository) Finds(ctx context.Context, filter map[string]interface{}) ([]domain.Payment, error) {
	query := r.db.QueryBuilder.Select("id, order_id, payment_method, payment_status, created_at, updated_at").From(r.TableName)

	// Apply filters if provided
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

	var payments []domain.Payment
	for rows.Next() {
		var payment domain.Payment
		err := rows.Scan(
			&payment.ID,
			&payment.OrderID,
			&payment.PaymentMethod,
			&payment.PaymentStatus,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *PaymentRepository) FindOne(ctx context.Context, id int) (*domain.Payment, error) {
	var payment domain.Payment

	query := r.db.QueryBuilder.Select("id, order_id, payment_method, payment_status, created_at, updated_at").
		From(r.TableName).
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.PaymentMethod,
		&payment.PaymentStatus,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &payment, nil
}

func (r *PaymentRepository) Store(ctx context.Context, data *domain.Payment) error {
	query := r.db.QueryBuilder.Insert(r.TableName).
		Columns("order_id", "payment_method", "payment_status", "created_at", "updated_at", "expired_at").
		Values(data.OrderID, data.PaymentMethod, data.PaymentStatus, data.CreatedAt, data.UpdatedAt, data.ExpiredAt).
		Suffix("RETURNING id, order_id, payment_method, payment_status, created_at, updated_at, expired_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&data.ID,
		&data.OrderID,
		&data.PaymentMethod,
		&data.PaymentStatus,
		&data.CreatedAt,
		&data.UpdatedAt,
		&data.ExpiredAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	}

	return nil
}

// Update modifies an existing Categories in the database
func (r *PaymentRepository) Update(ctx context.Context, order_id int, updatedData *domain.Payment) error {
	query := r.db.QueryBuilder.Update(r.TableName).
		Set("updated_at", time.Now()).
		Set("payment_method", sq.Expr("COALESCE(?, payment_method)", nullString(updatedData.PaymentMethod))).
		Set("payment_status", sq.Expr("COALESCE(?, payment_status)", nullString(updatedData.PaymentStatus))).
		Where(sq.Eq{"order_id": order_id}).
		Suffix("RETURNING id, order_id, payment_method, payment_status, created_at, updated_at, expired_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&updatedData.ID,
		&updatedData.OrderID,
		&updatedData.PaymentMethod,
		&updatedData.PaymentStatus,
		&updatedData.CreatedAt,
		&updatedData.UpdatedAt,
		&updatedData.ExpiredAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	}

	return nil
}

func (r *PaymentRepository) Delete(ctx context.Context, id int) error {
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

func (r *PaymentRepository) FindExpiredPayments(ctx context.Context, now time.Time) ([]*domain.Payment, error) {
	var payments []*domain.Payment

	query := r.db.QueryBuilder.Select("id", "order_id", "payment_method", "payment_status", "created_at", "expired_at").
		From(r.TableName).
		Where(sq.And{
			sq.Eq{"payment_status": "pending"},
			sq.Lt{"expired_at": now},
		})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var payment domain.Payment
		if err := rows.Scan(&payment.ID, &payment.OrderID, &payment.PaymentMethod, &payment.PaymentStatus, &payment.CreatedAt, &payment.ExpiredAt); err != nil {
			return nil, err
		}
		payments = append(payments, &payment)
	}

	return payments, nil
}

func (r *PaymentRepository) FindByUserIDandOrderID(ctx context.Context, userID int, orderID int) (*domain.Payment, error) {
	var payment domain.Payment

	query := r.db.QueryBuilder.Select("id", "order_id", "payment_method", "payment_status", "created_at", "updated_at", "expired_at").
		From(r.TableName).
		Where(sq.And{
			sq.Eq{"order_id": orderID},
		}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.PaymentMethod,
		&payment.PaymentStatus,
		&payment.CreatedAt,
		&payment.UpdatedAt,
		&payment.ExpiredAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &payment, nil
}
