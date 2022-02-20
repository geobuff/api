package repo

import (
	"database/sql"
	"time"
)

type NewsletterSubscriber struct {
	ID    int       `json:"id"`
	Email string    `json:"email"`
	Added time.Time `json:"added"`
}

type CreateNewsletterSubscriberDto struct {
	Email string `json:"email"`
}

func NewsletterSubscriberExists(email string) (bool, error) {
	statement := "SELECT id from newsletterSubscribers WHERE email = $1;"
	var id int
	err := Connection.QueryRow(statement, email).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func InsertNewsletterSubscriber(entry CreateNewsletterSubscriberDto) error {
	statement := "INSERT INTO newsletterSubscribers (email, added) VALUES ($1, $2) RETURNING id;"
	var id int
	return Connection.QueryRow(statement, entry.Email, time.Now()).Scan(&id)
}

func DeleteNewsletterSubscriber(subscriberID int) error {
	statement := "DELETE FROM newsletterSubscribers where id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, subscriberID).Scan(&id)
}
