package use_cases

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"jobsity-code-challenge/entities"
	"jobsity-code-challenge/use_cases/mocks/mocks"
)

func TestUseCase_Login(t *testing.T) {
	ctx := context.Background()
	t.Run("happy path", func(t *testing.T) {
		data := entities.LoginData{
			UserName: "userTest",
			Password: "12345",
		}
		hashPwd, _ := bcrypt.GenerateFromPassword([]byte(data.Password), 13)
		expected := entities.User{
			ID:        1,
			FirstName: "user",
			LastName:  "test",
			UserName:  "userTest",
			Password:  string(hashPwd),
		}
		repoMock := mocks.Repo{}
		repoMock.On("GetUserByUserName", ctx, data.UserName).Return(&expected, nil)
		u := UseCase{
			db: &repoMock,
		}
		got, err := u.Login(ctx, data)
		repoMock.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, expected, *got)
	})
	t.Run("invalid password", func(t *testing.T) {
		data := entities.LoginData{
			UserName: "userTest",
			Password: "12345",
		}
		hashPwd, _ := bcrypt.GenerateFromPassword([]byte("otherPwd"), 13)
		result := entities.User{
			ID:        1,
			FirstName: "user",
			LastName:  "test",
			UserName:  "userTest",
			Password:  string(hashPwd),
		}
		repoMock := mocks.Repo{}
		repoMock.On("GetUserByUserName", ctx, data.UserName).Return(&result, nil)
		u := UseCase{
			db: &repoMock,
		}
		got, err := u.Login(ctx, data)
		repoMock.AssertExpectations(t)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
	t.Run("user not exists", func(t *testing.T) {
		data := entities.LoginData{
			UserName: "userTest",
			Password: "12345",
		}
		repoMock := mocks.Repo{}
		repoMock.On("GetUserByUserName", ctx, data.UserName).Return(nil, nil)
		u := UseCase{
			db: &repoMock,
		}
		got, err := u.Login(ctx, data)
		repoMock.AssertExpectations(t)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
	t.Run("error on db", func(t *testing.T) {
		data := entities.LoginData{
			UserName: "userTest",
			Password: "12345",
		}
		repoMock := mocks.Repo{}
		repoMock.On("GetUserByUserName", ctx, data.UserName).Return(nil, errors.New("error"))
		u := UseCase{
			db: &repoMock,
		}
		got, err := u.Login(ctx, data)
		repoMock.AssertExpectations(t)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
