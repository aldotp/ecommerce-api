package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/storage/postgres"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type ProductRepository struct {
	db        *postgres.DB
	TableName string
}

func NewProductRepository(db *postgres.DB) *ProductRepository {
	return &ProductRepository{
		db:        db,
		TableName: "products",
	}
}

// Finds retrieves multiple products based on the provided filter
func (r *ProductRepository) Finds(ctx context.Context, filter map[string]interface{}) ([]domain.Product, error) {
	query := r.db.QueryBuilder.Select("*").From(r.TableName)

	for key, value := range filter {
		if key == "search" {
			query = query.Where(sq.Or{
				sq.Like{"LOWER(name)": "%" + value.(string) + "%"},
			})
			continue
		}

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

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CategoryID,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// FindOne retrieves a single product by ID
func (r *ProductRepository) FindOne(ctx context.Context, id int) (*domain.Product, error) {
	var product domain.Product

	query := r.db.QueryBuilder.Select("*").
		From(r.TableName).
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CategoryID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

// Store inserts a new product into the database
func (r *ProductRepository) Store(ctx context.Context, data *domain.Product) error {
	query := r.db.QueryBuilder.Insert(r.TableName).
		Columns("name", "description", "price", "stock", "category_id", "created_at", "updated_at").
		Values(data.Name, data.Description, data.Price, data.Stock, data.CategoryID, time.Now(), time.Now()).
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

// Update modifies an existing product in the database
func (r *ProductRepository) Update(ctx context.Context, id int, updatedData domain.Product) error {
	query := r.db.QueryBuilder.Update(r.TableName).
		Set("name", sq.Expr("COALESCE(?, name)", updatedData.Name)).
		Set("price", sq.Expr("COALESCE(?, price)", updatedData.Price)).
		Set("stock", sq.Expr("COALESCE(?, stock)", updatedData.Stock)).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&updatedData.ID,
		&updatedData.Name,
		&updatedData.Price,
		&updatedData.Stock,
		&updatedData.CreatedAt,
		&updatedData.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a product by ID from the database
func (r *ProductRepository) Delete(ctx context.Context, id int) error {
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

func (r *ProductRepository) UpdateStock(ctx context.Context, productID, newStock int) error {
	query := r.db.QueryBuilder.Update(r.TableName).
		Set("stock", newStock).
		Where(sq.Eq{"id": productID})

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
