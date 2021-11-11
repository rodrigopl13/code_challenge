package controllers

import (
	"context"

	"jobsity-code-challenge/entities"
)

type useCase interface {
	Login(ctx context.Context, data entities.LoginData) (*entities.User, error)
	SignUp(ctx context.Context, data entities.User) error
	GetMessages(ctx context.Context, limit int) ([]entities.Message, error)
}

type Tokenizer interface {
	GenerateToken(user entities.User) (string, error)
}

type Controller struct {
	usecase   useCase
	wsServer  *WsServer
	tokenizer Tokenizer
}

func New(u useCase, ws *WsServer, t Tokenizer) Controller {
	return Controller{usecase: u, wsServer: ws, tokenizer: t}
}
