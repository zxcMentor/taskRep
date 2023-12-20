package repository

import (
	"awesomeProject/internal/models/user"
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user user.User) error
	GetByID(ctx context.Context, id string) (user.User, error)
	Update(ctx context.Context, user user.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]user.User, error)
}

type UserRepo struct {
	DB *sql.DB
	qb squirrel.StatementBuilderType
}

func NewDBPostgres(db *sql.DB) *UserRepo {
	return &UserRepo{
		DB: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *UserRepo) Create(ctx context.Context, user user.User) error {
	query, args, err := r.qb.
		Insert("users").
		Columns("username", "email", "password", "created_at", "updated_at").
		Values(user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.DB.ExecContext(ctx, query, args...)
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (user.User, error) {
	var user user.User

	query, args, err := r.qb.
		Select("id", "username", "email", "password", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return user, err
	}

	err = r.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepo) Update(ctx context.Context, user user.User) error {
	query, args, err := r.qb.
		Update("users").
		Set("username", user.Username).
		Set("email", user.Email).
		Set("password", user.Password).
		Set("updated_at", user.UpdatedAt).
		Where(squirrel.Eq{"id": user.ID}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.DB.ExecContext(ctx, query, args...)
	return err
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	query, args, err := r.qb.
		Update("users").
		Set("deleted_at", time.Now()).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.DB.ExecContext(ctx, query, args...)
	return err
}

func (r *UserRepo) List(ctx context.Context, limit, offset int) ([]user.User, error) {
	var users []user.User

	query, args, err := r.qb.
		Select("id", "username", "email", "created_at", "updated_at").
		From("users").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user user.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
