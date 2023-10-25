package repo

import (
	"context"
	"database/sql"
	"github.com/evrone/go-clean-template/internal/auth/entity"
	"log"
)

type TokenRepo struct {
	db *sql.DB
}

func NewTokenRepo(db *sql.DB) *TokenRepo {
	return &TokenRepo{db}
}

func (t *TokenRepo) CreateUserToken(ctx context.Context, userToken entity.Token) error {
	_, err := t.db.Exec("INSERT INTO tokens (user_id, token) VALUES ($1, $2)",
		userToken.UserID, userToken.Token)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (t *TokenRepo) UpdateUserToken(ctx context.Context, userToken entity.Token) error {
	return nil

}
