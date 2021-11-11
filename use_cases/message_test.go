package use_cases

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"jobsity-code-challenge/entities"
	"jobsity-code-challenge/use_cases/mocks/mocks"
)

func TestUseCase_GetMessages(t *testing.T) {
	ctx := context.Background()
	t.Run("happy path", func(t *testing.T) {
		limit := 10
		expected := []entities.Message{
			{
				ID: 1,
				User: entities.User{
					ID:        1,
					FirstName: "user",
					LastName:  "test",
					UserName:  "userTest",
					Password:  "somePwd",
				},
				Message:   "test message",
				CreatedAt: time.Time{},
			},
		}
		repoMock := mocks.Repo{}
		repoMock.On("GetLastMessages", ctx, limit).Return(expected, nil)
		u := UseCase{
			db: &repoMock,
		}
		got, err := u.GetMessages(ctx, limit)
		repoMock.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
	})
	t.Run("none message in database", func(t *testing.T) {
		limit := 10
		repoMock := mocks.Repo{}
		repoMock.On("GetLastMessages", ctx, limit).Return(nil, nil)
		u := UseCase{
			db: &repoMock,
		}
		got, err := u.GetMessages(ctx, limit)
		repoMock.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})
	t.Run("error in database", func(t *testing.T) {
		limit := 10
		repoMock := mocks.Repo{}
		repoMock.On("GetLastMessages", ctx, limit).Return(nil, errors.New("error"))
		u := UseCase{
			db: &repoMock,
		}
		got, err := u.GetMessages(ctx, limit)
		repoMock.AssertExpectations(t)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestUseCase_ProcessMessage(t *testing.T) {
	t.Run("happy path, with no command", func(t *testing.T) {
		created := time.Now().Format(time.RFC3339)
		tNorm, _ := time.Parse(time.RFC3339, created)
		message := entities.Message{
			User: entities.User{
				ID:        1,
				FirstName: "user",
				LastName:  "test",
				UserName:  "userTest",
			},
			Message:   "some message",
			CreatedAt: tNorm,
		}
		repoMock := mocks.Repo{}
		repoMock.On("InsertMessage", message).Return(nil)
		u := UseCase{
			db: &repoMock,
		}
		b, _ := json.Marshal(message)
		got, err := u.ProcessMessage(string(b))
		repoMock.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Equal(t, b, got)
	})
	t.Run("happy path, with command", func(t *testing.T) {
		created := time.Now().Format(time.RFC3339)
		tNorm, _ := time.Parse(time.RFC3339, created)
		message := entities.Message{
			User: entities.User{
				ID:        1,
				FirstName: "user",
				LastName:  "test",
				UserName:  "userTest",
			},
			Message:   "/stock=stockCode",
			CreatedAt: tNorm,
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`Symbol,Date,Time,Open,High,Low,Close,Volume
test,2021-11-10,01:00:00,0,0,0,100.50,0`))
		}))
		defer server.Close()
		brokerMock := mocks.Broker{}
		brokerMock.On("Publish", "testQueue", mock.Anything).Return(nil)
		brokerMock.On("GetStockQueueName").Return("testQueue")
		u := UseCase{
			rabbit:         &brokerMock,
			stooqUrlString: server.URL + "/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv",
		}
		b, _ := json.Marshal(message)
		got, err := u.ProcessMessage(string(b))
		brokerMock.AssertExpectations(t)
		assert.NoError(t, err)
		assert.Nil(t, got)

	})
	t.Run("error inserting message in database, no command", func(t *testing.T) {
		created := time.Now().Format(time.RFC3339)
		tNorm, _ := time.Parse(time.RFC3339, created)
		message := entities.Message{
			User: entities.User{
				ID:        1,
				FirstName: "user",
				LastName:  "test",
				UserName:  "userTest",
			},
			Message:   "some message",
			CreatedAt: tNorm,
		}
		repoMock := mocks.Repo{}
		repoMock.On("InsertMessage", message).Return(errors.New("error"))
		u := UseCase{
			db: &repoMock,
		}
		b, _ := json.Marshal(message)
		got, err := u.ProcessMessage(string(b))
		repoMock.AssertExpectations(t)
		assert.Error(t, err)
		assert.Equal(t, b, got)
	})
	t.Run("error sending message to rabbitmq", func(t *testing.T) {
		created := time.Now().Format(time.RFC3339)
		tNorm, _ := time.Parse(time.RFC3339, created)
		message := entities.Message{
			User: entities.User{
				ID:        1,
				FirstName: "user",
				LastName:  "test",
				UserName:  "userTest",
			},
			Message:   "/stock=stockCode",
			CreatedAt: tNorm,
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`Symbol,Date,Time,Open,High,Low,Close,Volume
test,2021-11-10,01:00:00,0,0,0,100.50,0`))
		}))
		defer server.Close()
		brokerMock := mocks.Broker{}
		brokerMock.On("Publish", "testQueue", mock.Anything).Return(errors.New("error"))
		brokerMock.On("GetStockQueueName").Return("testQueue")
		u := UseCase{
			rabbit:         &brokerMock,
			stooqUrlString: server.URL + "/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv",
		}
		b, _ := json.Marshal(message)
		got, err := u.ProcessMessage(string(b))
		brokerMock.AssertExpectations(t)
		assert.Error(t, err)
		assert.Contains(t, string(got), "error executing command")
	})
	t.Run("error command not known", func(t *testing.T) {
		created := time.Now().Format(time.RFC3339)
		tNorm, _ := time.Parse(time.RFC3339, created)
		message := entities.Message{
			User: entities.User{
				ID:        1,
				FirstName: "user",
				LastName:  "test",
				UserName:  "userTest",
			},
			Message:   "/command=stockCode",
			CreatedAt: tNorm,
		}
		u := UseCase{}
		b, _ := json.Marshal(message)
		got, err := u.ProcessMessage(string(b))
		assert.Error(t, err)
		assert.Contains(t, string(got), "unknown command: command")
	})
	t.Run("error command message wrong structure", func(t *testing.T) {
		created := time.Now().Format(time.RFC3339)
		tNorm, _ := time.Parse(time.RFC3339, created)
		message := entities.Message{
			User: entities.User{
				ID:        1,
				FirstName: "user",
				LastName:  "test",
				UserName:  "userTest",
			},
			Message:   "/stock stockCode",
			CreatedAt: tNorm,
		}
		u := UseCase{}
		b, _ := json.Marshal(message)
		got, err := u.ProcessMessage(string(b))
		assert.Error(t, err)
		assert.Contains(t, string(got), "invalid command structure")
	})
	t.Run("error calling stooq API", func(t *testing.T) {
		created := time.Now().Format(time.RFC3339)
		tNorm, _ := time.Parse(time.RFC3339, created)
		message := entities.Message{
			User: entities.User{
				ID:        1,
				FirstName: "user",
				LastName:  "test",
				UserName:  "userTest",
			},
			Message:   "/stock=stockCode",
			CreatedAt: tNorm,
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`Server Error`))
		}))
		defer server.Close()
		u := UseCase{
			stooqUrlString: server.URL + "/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv",
		}
		b, _ := json.Marshal(message)
		got, err := u.ProcessMessage(string(b))
		assert.Error(t, err)
		assert.Contains(t, string(got), "something went wrong calling stooq!")
	})
}
