package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"jobsity-code-challenge/entities"
)

const (
	queryGetMessages = `SELECT t.* FROM (SELECT log.message, log.created_at, users.first_name, users.last_name, users.user_name
						FROM chat_log as log INNER JOIN users ON log.user_id = users.id ORDER BY log.created_at DESC LIMIT $1) t ORDER BY t.created_at ASC`
	queryInsertMessage = `INSERT INTO chat_log VALUES(DEFAULT, $1, $2, $3)`
)

func (d RepoDB) GetLastMessages(ctx context.Context, limit int) ([]entities.Message, error) {
	rows, err := d.db.QueryContext(ctx, queryGetMessages, limit)
	defer rows.Close()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "getting messages from DB")
	}
	var messages []entities.Message
	for rows.Next() {
		var message entities.Message
		if err = rows.Scan(
			&message.Message,
			&message.CreatedAt,
			&message.User.FirstName,
			&message.User.LastName,
			&message.User.UserName,
		); err != nil {
			return messages, fmt.Errorf("scanning row from DB error: %w", err)
		}
		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil {
		return messages, errors.Wrap(err, "reading rows from DB")
	}
	return messages, nil
}

func (d RepoDB) InsertMessage(msg entities.Message) error {
	_, err := d.db.Exec(queryInsertMessage, msg.User.ID, msg.Message, time.Now().UTC())
	if err != nil {
		return errors.Wrap(err, "inserting message into DB")
	}
	return nil
}
