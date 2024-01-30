package user

import "github.com/jackc/pgx/v5"

type UserRepo interface {
}

type userRepo struct {
	db *pgx.Conn
}

func NewUserRepo(db *pgx.Conn) *userRepo {
	return &userRepo{
		db: db,
	}
}
