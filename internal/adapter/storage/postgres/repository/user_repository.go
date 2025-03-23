package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/aldotp/ecommerce-go-api/internal/adapter/storage/postgres"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *postgres.DB
}

func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

// CreateUser creates a new user in the database
func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := ur.db.QueryBuilder.Insert("users").
		Columns("name", "email", "password", "role", "created_at", "updated_at").
		Values(user.Name, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id uint64) (*domain.User, error) {
	var user domain.User

	query := ur.db.QueryBuilder.Select("*").
		From("users").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, consts.ErrDataNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	query := ur.db.QueryBuilder.Select("*").
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, consts.ErrDataNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) ListUsers(ctx context.Context, offset, limit uint64) ([]domain.User, error) {
	var user domain.User
	var users []domain.User

	if limit == 0 {
		limit = 10
	}

	if offset == 0 {
		offset = 1
	}

	query := ur.db.QueryBuilder.Select("*").
		From("users").
		OrderBy("id").
		Limit(limit).
		Offset((offset - 1) * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ur.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {

	query := ur.db.QueryBuilder.Update("users").
		Set("name", sq.Expr("COALESCE(?, name)", nullString(user.Name))).
		Set("email", sq.Expr("COALESCE(?, email)", nullString(user.Email))).
		Set("password", sq.Expr("COALESCE(?, password)", nullString(user.Password))).
		Set("role", sq.Expr("COALESCE(?, role)", nullString(string(user.Role)))).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": user.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user by ID from the database
func (ur *UserRepository) DeleteUser(ctx context.Context, id uint64) error {
	query := ur.db.QueryBuilder.Delete("users").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = ur.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	var user domain.User

	query := ur.db.QueryBuilder.Select("*").
		From("users").
		Where(sq.Eq{"token": token}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, consts.ErrDataNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) ExistEmail(ctx context.Context, email string) (bool, error) {
	query := ur.db.QueryBuilder.Select("id").
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return false, err
	}

	var id uint64

	err = ur.db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
