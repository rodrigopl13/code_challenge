package use_cases

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"jobsity-code-challenge/entities"
	"jobsity-code-challenge/use_cases/mocks/mocks"
)

func TestUseCase_SignUp(t *testing.T) {
	ctx := context.Background()
	t.Run("happy path", func(t *testing.T) {
		data := entities.User{
			FirstName: "test",
			LastName:  "test",
			UserName:  "test_user",
			Password:  "test_pwd",
		}
		repoMock := mocks.Repo{}
		repoMock.On("CreateNewUSer", ctx, data).Return(nil)
		u := UseCase{
			db: &repoMock,
		}

		assert.NoError(t, u.SignUp(ctx, data))
		repoMock.AssertExpectations(t)
	})
	t.Run("error on db", func(t *testing.T) {
		data := entities.User{
			FirstName: "test",
			LastName:  "test",
			UserName:  "test_user",
			Password:  "test_pwd",
		}
		repoMock := mocks.Repo{}
		repoMock.On("CreateNewUSer", ctx, data).Return(errors.New("error"))
		u := UseCase{
			db: &repoMock,
		}

		assert.Error(t, u.SignUp(ctx, data))
		repoMock.AssertExpectations(t)
	})
}
