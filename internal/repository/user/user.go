package user

import (
	"context"
	"errors"
	"inditilla/internal/entity"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	SaveUser(context.Context, entity.UserSignupForm) (int, error)
	Authenticate(context.Context, string, string) (entity.UserEntity, error)
	Exists(context.Context, string) (bool, error)
	GetById(context.Context, int) (entity.UserEntity, error)
	Update(context.Context, *entity.UserEntity) error
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

func (r *userRepo) Authenticate(ctx context.Context, email string, password string) (entity.UserEntity, error) {
	user := entity.UserEntity{}
	var hashedPassword []byte

	query := `SELECT * FROM users WHERE email=$1`

	err := r.db.QueryRow(ctx, query, email).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &hashedPassword, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.UserEntity{}, entity.ErrInvalidCredentials
		} else {
			return entity.UserEntity{}, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return entity.UserEntity{}, entity.ErrInvalidCredentials
		} else {
			return entity.UserEntity{}, err
		}
	}

	return user, nil
}

func (r *userRepo) Exists(ctx context.Context, email string) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(
		SELECT true 
		FROM users 
		WHERE users.email = $1 
		)`

	err := r.db.QueryRow(ctx, query, email).Scan(&exists)

	if errors.Is(err, pgx.ErrNoRows) {
		return false, entity.ErrNoRecord
	}

	return exists, err
}

func (r *userRepo) GetById(ctx context.Context, id int) (entity.UserEntity, error) {
	user := entity.UserEntity{}
	var hashedPassword []byte

	query := `SELECT * FROM users WHERE id=$1`
	err := r.db.QueryRow(ctx, query, id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &hashedPassword, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.UserEntity{}, entity.ErrNoRecord
		}
		return entity.UserEntity{}, err
	}

	user.Password = string(hashedPassword)

	return user, nil
}

func (r *userRepo) Update(ctx context.Context, user *entity.UserEntity) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 15)
	if err != nil {
		return err
	}

	query := `
		UPDATE users 
		SET first_name = $1, last_name = $2, email = $3, hashed_password = $4
		WHERE id = $5
		`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = r.db.Exec(ctx, query, user.FirstName, user.LastName, user.Email, hashedPassword, user.Id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.ErrEditConflict
		}
		return err
	}

	return nil
}
