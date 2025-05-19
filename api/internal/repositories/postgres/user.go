package postgres

import (
	"api/internal/user"
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Users struct {
	db    *DB
	table string
}

func NewUsers(db *DB) *Users {
	return &Users{
		db:    db,
		table: "users",
	}
}

func (u Users) Create(usr user.User) error {
	query, args, _ := builder.Insert(u.table).
		Columns("username", "login", "password").
		Values(usr.Username, usr.Login, usr.Password).
		ToSql()

	_, err := u.db.Exec(context.Background(), query, args...)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr); pgErr != nil && pgErr.Code == "23505" {
		return user.ErrAlreadyExists
	}
	return err
}

func (u Users) Get(login string) (user.User, error) {
	query, args, _ := builder.Select("*").
		From(u.table).
		Where(squirrel.Eq{"login": login}).
		ToSql()

	usr := user.User{}
	err := u.db.QueryRow(context.Background(), query, args...).
		Scan(&usr.Username, &usr.Login, &usr.Password)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.User{}, user.ErrNotFound
	}
	return usr, err
}

func (u Users) UpdateUsername(login, username string) error {
	query, args, _ := builder.Update(u.table).
		Set("username", username).
		Where(squirrel.Eq{"login": login}).
		ToSql()

	_, err := u.db.Exec(context.Background(), query, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.ErrNotFound
	}
	return err
}
