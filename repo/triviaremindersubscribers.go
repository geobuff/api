package repo

import (
	"database/sql"
	"time"
)

type TriviaReminderSubscriber struct {
	ID    int       `json:"id"`
	Email string    `json:"email"`
	Added time.Time `json:"added"`
}

type CreateTriviaReminderSubscriberDto struct {
	Email string `json:"email"`
}

func TriviaReminderSubscriberExists(email string) (bool, error) {
	statement := "SELECT id from triviaReminderSubscribers WHERE email = $1;"
	var id int
	err := Connection.QueryRow(statement, email).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func InsertTriviaReminderSubscriber(entry CreateTriviaReminderSubscriberDto) error {
	statement := "INSERT INTO triviaReminderSubscribers (email, added) VALUES ($1, $2) RETURNING id;"
	var id int
	return Connection.QueryRow(statement, entry.Email, time.Now()).Scan(&id)
}

func DeleteTriviaReminderSubscriber(subscriberID int) error {
	statement := "DELETE FROM triviaReminderSubscribers where id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, subscriberID).Scan(&id)
}
