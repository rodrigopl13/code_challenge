package use_cases

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"jobsity-code-challenge/entities"
)

func (u *UseCase) Login(ctx context.Context, data entities.LoginData) (*entities.User, error) {
	user, err := u.db.GetUserByUserName(ctx, data.UserName)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New(fmt.Sprintf("userName %s not valid", data.UserName))
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return nil, errors.New("password not valid")
	}
	user.Password = ""
	return user, nil
}
