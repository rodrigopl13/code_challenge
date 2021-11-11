package use_cases

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"jobsity-code-challenge/entities"
)

const (
	stock            = "stock"
	csvNumFields     = 8
	stockBotResponse = "%s quote is $%s per share"
)

var commands = map[string]bool{
	stock: true,
}

func (u *UseCase) GetMessages(ctx context.Context, limit int) ([]entities.Message, error) {
	return u.db.GetLastMessages(ctx, limit)
}

func (u *UseCase) ProcessMessage(message string) ([]byte, error) {
	var msg entities.Message
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		return []byte(message), errors.Wrap(err, "unmarshalling json message")
	}

	if strings.HasPrefix(msg.Message, "/") {
		return u.runCommand(message, msg)
	}
	err := u.db.InsertMessage(msg)
	if err != nil {
		return []byte(message), errors.Wrap(err, "inserting message in database")
	}

	return []byte(message), nil
}

func (u *UseCase) runCommand(message string, msg entities.Message) ([]byte, error) {
	res := entities.Message{
		ID: 0,
		User: entities.User{
			ID:        0,
			FirstName: "Chat-Bot",
			UserName:  "chat-bot",
		},
		CreatedAt: time.Now().UTC(),
	}
	c, v, err := parseCommand(msg.Message)
	if err != nil {
		res.Message = err.Error()
		b, _ := json.Marshal(res)
		return b, errors.Wrap(err, message)
	}
	if !isKnownCommand(c) {
		res.Message = fmt.Sprintf("unknown command: %s", c)
		b, _ := json.Marshal(res)
		return b, errors.New(res.Message)
	}
	price, err := u.callStooqAPI(v)
	if err != nil {
		res.Message = "something went wrong calling stooq!"
		b, _ := json.Marshal(res)
		return b, errors.Wrap(err, "calling stooq API")
	}
	res.Message = fmt.Sprintf(stockBotResponse, strings.ToUpper(v), price)
	err = u.publishStockMessage(res)
	if err != nil {
		res.Message = "error executing command"
		b, _ := json.Marshal(res)
		return b, errors.Wrap(err, "publishing message to rabbitmq")
	}
	return nil, nil
}

func (u *UseCase) publishStockMessage(msg entities.Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return u.rabbit.Publish(u.rabbit.GetStockQueueName(), b)
}

func parseCommand(message string) (command, value string, err error) {
	v := strings.Split(message[1:], "=")
	if len(v) != 2 {
		return "", "", fmt.Errorf("invalid command structure")
	}
	command = strings.ToLower(strings.TrimSpace(v[0]))
	value = strings.ToLower(strings.TrimSpace(v[1]))
	return
}

func isKnownCommand(command string) bool {
	return commands[command]
}

func (u *UseCase) callStooqAPI(stockCode string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(u.stooqUrlString, stockCode))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.FieldsPerRecord = csvNumFields
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}
	return records[1][6], nil
}
