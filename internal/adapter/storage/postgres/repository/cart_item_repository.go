package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/storage/postgres"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type CartItemRepository struct {
	db        *postgres.DB
	TableName string
}

func NewCartItemRepository(db *postgres.DB) *CartItemRepository {
	return &CartItemRepository{
		db:        db,
		TableName: "cart_items",
	}
}

// FindOne retrieves a single Categories by ID
func (r *CartItemRepository) FindOne(ctx context.Context, id int) (*domain.CartItem, error) {
	var cartItem domain.CartItem

	query := r.db.QueryBuilder.Select("id", "cart_id", "product_id", "quantity", "created_at", "updated_at").
		From(r.TableName).
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&cartItem.ID,
		&cartItem.CartID,
		&cartItem.ProductID,
		&cartItem.Quantity,
		&cartItem.CreatedAt,
		&cartItem.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &cartItem, nil
}

func (r *CartItemRepository) FindOneByFilters(ctx context.Context, filter map[string]interface{}) (*domain.CartItem, error) {
	query := r.db.QueryBuilder.Select("*").From(r.TableName)

	// Apply filters if provided
	for key, value := range filter {
		query = query.Where(sq.Eq{key: value})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var cartItem domain.CartItem
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&cartItem.ID,
		&cartItem.CartID,
		&cartItem.ProductID,
		&cartItem.Quantity,
		&cartItem.CreatedAt,
		&cartItem.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &cartItem, nil
}

func (r *CartItemRepository) Finds(ctx context.Context, filter map[string]interface{}) ([]domain.CartItem, error) {
	query := r.db.QueryBuilder.Select("*").From(r.TableName)

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

	var cartItems []domain.CartItem
	for rows.Next() {
		var cartItem domain.CartItem
		err := rows.Scan(
			&cartItem.ID,
			&cartItem.CartID,
			&cartItem.ProductID,
			&cartItem.Quantity,
			&cartItem.CreatedAt,
			&cartItem.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cartItems = append(cartItems, cartItem)
	}

	return cartItems, nil
}

// Store inserts a new Categories into the database
func (r *CartItemRepository) Store(ctx context.Context, data *domain.CartItem) error {
	query := r.db.QueryBuilder.Insert(r.TableName).
		Columns("cart_id", "product_id", "quantity", "created_at", "updated_at").
		Values(data.CartID, data.ProductID, data.Quantity, data.CreatedAt, data.UpdatedAt).
		Suffix("RETURNING id, cart_id, product_id, quantity, created_at, updated_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&data.ID,
		&data.CartID,
		&data.ProductID,
		&data.Quantity,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

// Update modifies an existing Categories in the database
func (r *CartItemRepository) Update(ctx context.Context, id int, updatedData domain.CartItem) error {
	query := r.db.QueryBuilder.Update(r.TableName).
		Set("quantity", updatedData.Quantity).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&updatedData.ID,
		&updatedData.CartID,
		&updatedData.ProductID,
		&updatedData.Quantity,
		&updatedData.CreatedAt,
		&updatedData.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *CartItemRepository) DeleteByCartID(ctx context.Context, cartID int) error {
	query := r.db.QueryBuilder.Delete(r.TableName).
		Where(sq.Eq{"cart_id": cartID})

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

func (r *CartItemRepository) DeleteByProductID(ctx context.Context, product_id int) error {
	query := r.db.QueryBuilder.Delete(r.TableName).
		Where(sq.Eq{"product_id": product_id})

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
