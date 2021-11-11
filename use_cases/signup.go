package use_cases

import (
	"context"

	"jobsity-code-challenge/entities"
)

func (u *UseCase) SignUp(ctx context.Context, data entities.User) error {
	return u.db.CreateNewUSer(ctx, data)
}
