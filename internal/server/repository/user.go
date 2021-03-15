package repository

import (
	"fmt"
	"go-template/internal/server/model"

	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Get(userID string) (*model.User, error) {
	query := `SELECT id, name FROM user where id = ?`
	user := model.User{}
	if err := r.db.Get(&user, query, userID); err != nil {
		return nil, fmt.Errorf("Get User failed. userId: %v, error: %w", userID, err)
	}
	return &user, nil
}
