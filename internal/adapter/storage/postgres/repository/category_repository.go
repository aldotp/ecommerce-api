package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/storage/postgres"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type CategoryRepository struct {
	db        *postgres.DB
	TableName string
}

func NewCategoryRepository(db *postgres.DB) *CategoryRepository {
	return &CategoryRepository{
		db:        db,
		TableName: "categories",
	}
}

func (r *CategoryRepository) Finds(ctx context.Context, filter map[string]interface{}) ([]domain.Category, error) {
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

	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// FindOne retrieves a single Categories by ID
func (r *CategoryRepository) FindOne(ctx context.Context, id int) (*domain.Category, error) {
	var category domain.Category

	query := r.db.QueryBuilder.Select("id, name").
		From(r.TableName).
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&category.ID,
		&category.Name,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

// Store inserts a new Categories into the database
func (r *CategoryRepository) Store(ctx context.Context, data *domain.Category) error {
	query := r.db.QueryBuilder.Insert(r.TableName).
		Columns("name").
		Values(data.Name).
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
func (r *CategoryRepository) Update(ctx context.Context, id int, updatedData domain.Category) error {
	query := r.db.QueryBuilder.Update(r.TableName).
		Set("name", sq.Expr("COALESCE(?, name)", updatedData.Name)).
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
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a Categories by ID from the database
func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
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
