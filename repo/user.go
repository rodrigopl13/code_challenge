package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"jobsity-code-challenge/entities"
)

const (
	queryGetUserByUserName = `SELECT * FROM users WHERE user_name = $1`
	queryInsertNewUser     = `INSERT INTO users VALUES(DEFAULT, $1, $2, $3, $4, $5)`
)

func (d RepoDB) GetUserByUserName(ctx context.Context, userName string) (*entities.User, error) {
	var user entities.User
	row := d.db.QueryRowContext(ctx, queryGetUserByUserName, userName)
	var createdAt time.Time
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.UserName, &user.Password, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, fmt.Sprintf("getting user %s from DB", userName))
	}
	return &user, nil
}

func (d RepoDB) CreateNewUSer(ctx context.Context, user entities.User) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 13)
	if err != nil {
		return errors.Wrap(err, "generating password hash")
	}
	_, err = d.db.ExecContext(ctx, queryInsertNewUser, user.FirstName, user.LastName, user.UserName, string(hashPassword), time.Now().UTC())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("inserting user %s on DB", user.UserName))
	}
	return nil
}
