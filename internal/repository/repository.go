package repository

import (
	"inditilla/internal/repository/user"

	"github.com/jackc/pgx/v5"
)

type Repositories struct {
	User user.UserRepo
}

func New(db *pgx.Conn) *Repositories {
	return &Repositories{
		User: user.NewUserRepo(db),
	}
}
