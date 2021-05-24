package repo

import (
	"database/sql"
	"time"
)

// User is the database object for a user entry.
type User struct {
	ID                  int            `json:"id"`
	Username            string         `json:"username"`
	Email               string         `json:"email"`
	PasswordHash        string         `json:"passwordHash"`
	CountryCode         string         `json:"countryCode"`
	XP                  int            `json:"xp"`
	IsAdmin             bool           `json:"isAdmin"`
	IsPremium           bool           `json:"isPremium"`
	PasswordResetToken  sql.NullString `json:"passwordResetToken"`
	PasswordResetExpiry sql.NullTime   `json:"passwordResetExpiry"`
}

// UserDto is the dto for a user entry.
type UserDto struct {
	ID                  int            `json:"id"`
	Username            string         `json:"username"`
	Email               string         `json:"email"`
	CountryCode         string         `json:"countryCode"`
	XP                  int            `json:"xp"`
	IsAdmin             bool           `json:"isAdmin"`
	IsPremium           bool           `json:"isPremium"`
	PasswordResetToken  sql.NullString `json:"passwordResetToken"`
	PasswordResetExpiry sql.NullTime   `json:"passwordResetExpiry"`
}

// GetUsers returns a page of users.
var GetUsers = func(limit int, offset int) ([]UserDto, error) {
	rows, err := Connection.Query("SELECT id, username, email, countrycode, xp, isadmin, ispremium, passwordresettoken, passwordresetexpiry FROM users LIMIT $1 OFFSET $2;", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users = []UserDto{}
	for rows.Next() {
		var user UserDto
		if err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.CountryCode, &user.XP, &user.IsAdmin, &user.IsPremium, &user.PasswordResetToken, &user.PasswordResetExpiry); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// GetFirstID returns the first ID for a given page.
var GetFirstID = func(limit int, offset int) (int, error) {
	statement := "SELECT id FROM users LIMIT $1 OFFSET $2;"
	var id int
	err := Connection.QueryRow(statement, limit, offset).Scan(&id)
	return id, err
}

// GetUser returns the user with a given id.
var GetUser = func(id int) (UserDto, error) {
	statement := "SELECT id, username, email, countrycode, xp, isadmin, ispremium, passwordresettoken, passwordresetexpiry FROM users WHERE id = $1;"
	var user UserDto
	err := Connection.QueryRow(statement, id).Scan(&user.ID, &user.Username, &user.Email, &user.CountryCode, &user.XP, &user.IsAdmin, &user.IsPremium, &user.PasswordResetToken, &user.PasswordResetExpiry)
	return user, err
}

// GetUserUsingEmail returns the user with a given email.
var GetUserUsingEmail = func(email string) (User, error) {
	statement := "SELECT * FROM users WHERE email = $1;"
	var user User
	err := Connection.QueryRow(statement, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CountryCode, &user.XP, &user.IsAdmin, &user.IsPremium, &user.PasswordResetToken, &user.PasswordResetExpiry)
	return user, err
}

// InsertUser inserts a new user into the users table.
var InsertUser = func(user User) (User, error) {
	statement := "INSERT INTO users (username, email, passwordHash, countrycode, xp, isAdmin, isPremium) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;"
	var newUser User
	err := Connection.QueryRow(statement, user.Username, user.Email, user.PasswordHash, user.CountryCode, 0, false, false).Scan(&newUser.ID, &newUser.Username, &newUser.Email, &newUser.PasswordHash, &newUser.CountryCode, &newUser.XP, &newUser.IsAdmin, &newUser.IsPremium, &user.PasswordResetToken, &user.PasswordResetExpiry)
	return newUser, err
}

// UpdateUser updates a user entry.
var UpdateUser = func(user User) error {
	statement := "UPDATE users set username = $2, email = $3, countryCode = $4, xp = $5 WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, user.ID, user.Username, user.Email, user.CountryCode, user.XP).Scan(&id)
}

// DeleteUser deletes a users scores, leaderboard entries and then the user entry in the users table.
var DeleteUser = func(userID int) error {
	scoresStatement := "DELETE FROM scores WHERE userId = $1;"
	Connection.QueryRow(scoresStatement, userID)
	leaderboardStatement := "DELETE FROM countries_leaderboard WHERE userId = $1;"
	Connection.QueryRow(leaderboardStatement, userID)
	usersStatement := "DELETE FROM users WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(usersStatement, id).Scan(&id)
}

var UsernameExists = func(username string) (bool, error) {
	statement := "SELECT COUNT(id) FROM users WHERE username = $1;"
	var count int
	err := Connection.QueryRow(statement, username).Scan(&count)
	return count > 0, err
}

var EmailExists = func(email string) (bool, error) {
	statement := "SELECT COUNT(id) FROM users WHERE email = $1;"
	var count int
	err := Connection.QueryRow(statement, email).Scan(&count)
	return count > 0, err
}

var SetPasswordResetValues = func(userID int, resetToken string, expiryDate time.Time) error {
	statement := "UPDATE users set passwordResetToken = $2, passwordResetExpiry = $3 WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, userID, resetToken, expiryDate).Scan(&id)
}

func ResetPassword(userID int, passwordHash string) error {
	statement := "UPDATE users set passwordhash = $2, passwordResetToken = null, passwordResetExpiry = null WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, userID, passwordHash).Scan(&id)
}

func UpgradeSubscription(userID int) error {
	statement := "UPDATE users set isPremium = true id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, userID).Scan(&id)
}
