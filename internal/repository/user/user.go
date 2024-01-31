package user

import (
	"context"
	"errors"
	"inditilla/internal/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	SaveUser(context.Context, entity.UserSignupForm) (int, error)
	// GetUserByEmail(context.Context, string) *entity.UserEntity
}

type userRepo struct {
	db *pgx.Conn
}

func NewUserRepo(db *pgx.Conn) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) SaveUser(ctx context.Context, u entity.UserSignupForm) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 15)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO users (first_name, last_name, email, hashed_password)
		VALUES ($1, $2, $3, $4) RETURNING id`

	var id int

	err = r.db.QueryRow(ctx, query, u.FirstName, u.LastName, u.Email, hashedPassword).Scan(&id)
	if err != nil {
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == "23505" {
				return 0, entity.ErrDuplicateEmail
			}
		}

		return 0, err
	}

	return id, nil
}
