package use_cases

import (
	"context"

	"jobsity-code-challenge/entities"
)

//go:generate mockery --name Repo --filename ./mocks/repo.go --outpkg mocks
type Repo interface {
	GetUserByUserName(ctx context.Context, userName string) (*entities.User, error)
	CreateNewUSer(ctx context.Context, user entities.User) error
	GetLastMessages(ctx context.Context, limit int) ([]entities.Message, error)
	InsertMessage(msg entities.Message) error
}

//go:generate mockery --name Broker --filename ./mocks/broker.go --outpkg mocks
type Broker interface {
	GetStockQueueName() string
	Publish(key string, msg []byte) error
}

type UseCase struct {
	db             Repo
	rabbit         Broker
	stooqUrlString string
}

func New(db Repo, rabbit Broker, stooqUrlString string) *UseCase {
	return &UseCase{db: db, rabbit: rabbit, stooqUrlString: stooqUrlString}
}
