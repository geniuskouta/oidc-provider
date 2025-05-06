package repo

import (
	"errors"
	"fmt"
	"oidc/internal/oidc/domain"
	"oidc/internal/oidc/infra"
)

type UserRepository struct {
	db *infra.DB
}

func NewUserRepo(db *infra.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) Save(user *domain.User) error {
	if _, err := repo.db.Get(user.Email); err == nil {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	repo.db.Set(user.Email, user)
	return nil
}

func (repo *UserRepository) FindByEmail(email string) (*domain.User, error) {
	println("repo")
	userData, err := repo.db.Get(email)
	if err != nil {
		return nil, err
	}

	user, ok := userData.(*domain.User)
	if !ok {
		return nil, errors.New("failed to cast client data")
	}

	return user, nil
}
