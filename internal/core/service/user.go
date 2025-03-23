package service

import (
	"context"
	"time"

	"github.com/aldotp/ecommerce-go-api/internal/adapter/dto"
	"github.com/aldotp/ecommerce-go-api/internal/core/domain"
	"github.com/aldotp/ecommerce-go-api/internal/core/port"
	"github.com/aldotp/ecommerce-go-api/pkg/consts"
	"github.com/aldotp/ecommerce-go-api/pkg/util"
	"go.uber.org/zap"
)

type UserService struct {
	cache       port.CacheInterface
	repo        port.UserRepository
	token       port.TokenInterface
	balanceRepo port.BalanceRepository
	log         *zap.Logger
}

func NewUserService(repo port.UserRepository, cache port.CacheInterface, token port.TokenInterface, log *zap.Logger, balanceRepo port.BalanceRepository) *UserService {
	return &UserService{
		cache:       cache,
		repo:        repo,
		token:       token,
		log:         log,
		balanceRepo: balanceRepo,
	}
}

// Register creates a new user
func (us *UserService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {

	exist, err := us.repo.ExistEmail(ctx, user.Email)
	if err != nil {
		us.log.Error(err.Error())
		return nil, consts.ErrInternal
	}

	if exist {
		return nil, consts.ErrEmailAlreadyExist
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		us.log.Error(err.Error())
		return nil, consts.ErrInternal
	}

	tNow := time.Now()

	user.Password = hashedPassword
	user.Role = domain.Customer
	user.CreatedAt = tNow
	user.UpdatedAt = tNow
	user, err = us.repo.CreateUser(ctx, user)
	if err != nil {
		us.log.Error(err.Error())
		if err == consts.ErrConflictingData {
			return nil, err
		}
		return nil, consts.ErrInternal
	}

	err = us.balanceRepo.Store(ctx, &domain.Balance{
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		us.log.Error(err.Error())
		return nil, consts.ErrInternal
	}

	return user, nil
}

// GetUser gets a user by ID
func (us *UserService) GetUser(ctx context.Context, id uint64) (*domain.User, error) {
	var user *domain.User

	cacheKey := util.GenerateCacheKey("user", id)
	cachedUser, err := us.cache.Get(ctx, cacheKey)
	if err == nil {
		err := util.Deserialize(cachedUser, &user)
		if err != nil {
			return nil, consts.ErrInternal
		}
		return user, nil
	}

	user, err = us.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == consts.ErrDataNotFound {
			return nil, err
		}
		return nil, consts.ErrInternal
	}

	userSerialized, err := util.Serialize(user)
	if err != nil {
		return nil, consts.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, consts.ErrInternal
	}

	return user, nil
}

// ListUsers lists all users
func (us *UserService) ListUsers(ctx context.Context, page, limit uint64) ([]domain.User, error) {
	var users []domain.User

	params := util.GenerateCacheKeyParams(page, limit)
	cacheKey := util.GenerateCacheKey("users", params)

	cachedUsers, err := us.cache.Get(ctx, cacheKey)
	if err == nil {
		err := util.Deserialize(cachedUsers, &users)
		if err != nil {
			return nil, consts.ErrInternal
		}

		return users, nil
	}

	users, err = us.repo.ListUsers(ctx, page, limit)
	if err != nil {
		return nil, consts.ErrInternal
	}

	usersSerialized, err := util.Serialize(users)
	if err != nil {
		return nil, consts.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, usersSerialized, time.Minute*10)
	if err != nil {
		return nil, consts.ErrInternal
	}

	return users, nil
}

// UpdateUser updates a user's name, email, and password
func (us *UserService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUser, err := us.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		if err == consts.ErrDataNotFound {
			return nil, err
		}
		return nil, consts.ErrInternal
	}

	emptyData := user.Name == "" &&
		user.Email == "" &&
		user.Password == "" &&
		user.Role == ""
	sameData := existingUser.Name == user.Name &&
		existingUser.Email == user.Email &&
		existingUser.Role == user.Role
	if emptyData || sameData {
		return nil, consts.ErrNoUpdatedData
	}

	var hashedPassword string

	if user.Password != "" {
		hashedPassword, err = util.HashPassword(user.Password)
		if err != nil {
			return nil, consts.ErrInternal
		}
	}

	user.Password = hashedPassword

	_, err = us.repo.UpdateUser(ctx, user)
	if err != nil {
		if err == consts.ErrConflictingData {
			return nil, err
		}
		return nil, consts.ErrInternal
	}

	cacheKey := util.GenerateCacheKey("user", user.ID)

	err = us.cache.Delete(ctx, cacheKey)
	if err != nil {
		return nil, consts.ErrInternal
	}

	userSerialized, err := util.Serialize(user)
	if err != nil {
		return nil, consts.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, consts.ErrInternal
	}

	err = us.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return nil, consts.ErrInternal
	}

	return user, nil
}

// DeleteUser deletes a user by ID
func (us *UserService) DeleteUser(ctx context.Context, id uint64) error {
	_, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == consts.ErrDataNotFound {
			return err
		}
		return consts.ErrInternal
	}

	cacheKey := util.GenerateCacheKey("user", id)

	err = us.cache.Delete(ctx, cacheKey)
	if err != nil {
		return consts.ErrInternal
	}

	err = us.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return consts.ErrInternal
	}

	return us.repo.DeleteUser(ctx, id)
}

// GetUser gets a user by ID
func (us *UserService) GetProfile(ctx context.Context, id uint64) (*dto.GetProfile, error) {
	var profile *dto.GetProfile

	cacheKey := util.GenerateCacheKey("user_profile", id)
	cachedUser, err := us.cache.Get(ctx, cacheKey)
	if err == nil {
		err := util.Deserialize(cachedUser, &profile)
		if err != nil {
			return nil, consts.ErrInternal
		}
		return profile, nil
	}

	user, err := us.repo.GetUserByID(ctx, id)
	if err != nil {
		if err == consts.ErrDataNotFound {
			return nil, err
		}
		return nil, consts.ErrInternal
	}

	profile = &dto.GetProfile{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	profileSerialized, err := util.Serialize(profile)
	if err != nil {
		return nil, consts.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, profileSerialized, 0)
	if err != nil {
		return nil, consts.ErrInternal
	}

	return profile, nil
}
