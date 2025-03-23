package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/storage/postgres"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type CartRepository struct {
	db        *postgres.DB
	TableName string
}

func NewCartRepository(db *postgres.DB) *CartRepository {
	return &CartRepository{
		db:        db,
		TableName: "carts",
	}
}

func (r *CartRepository) FindOne(ctx context.Context, id int) (*domain.Cart, error) {
	var Cart domain.Cart

	query := r.db.QueryBuilder.Select("id", "user_id", "created_at", "updated_at").
		From(r.TableName).
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&Cart.ID,
		&Cart.UserID,
		&Cart.CreatedAt,
		&Cart.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &Cart, nil
}

func (r *CartRepository) FindByUserID(ctx context.Context, userID int) (*domain.Cart, error) {
	var cart domain.Cart

	query := r.db.QueryBuilder.Select("id", "user_id", "created_at", "updated_at").
		From(r.TableName).
		Where(sq.Eq{"user_id": userID}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &cart, nil
}

// Store inserts a new Categories into the database
func (r *CartRepository) Store(ctx context.Context, data *domain.Cart) error {
	query := r.db.QueryBuilder.Insert(r.TableName).
		Columns("user_id", "created_at", "updated_at").
		Values(data.UserID, time.Now(), time.Now()).
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

// Update modifies an existing Categories in the database
func (r *CartRepository) Update(ctx context.Context, id int, updatedData domain.Cart) error {
	query := r.db.QueryBuilder.Update(r.TableName).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&updatedData.ID,
		&updatedData.UserID,
		&updatedData.CreatedAt,
		&updatedData.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *CartRepository) DeleteByUserID(ctx context.Context, userID int) error {
	query := r.db.QueryBuilder.Delete(r.TableName).
		Where(sq.Eq{"user_id": userID})

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
