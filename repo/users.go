package repo

import (
	"database/sql"
	"strings"
	"time"
)

// User is the database object for a user entry.
type User struct {
	ID                  int            `json:"id"`
	AvatarId            int            `json:"avatarId"`
	Username            string         `json:"username"`
	Email               string         `json:"email"`
	PasswordHash        string         `json:"passwordHash"`
	CountryCode         string         `json:"countryCode"`
	XP                  int            `json:"xp"`
	IsPremium           bool           `json:"isPremium"`
	IsAdmin             bool           `json:"isAdmin"`
	PasswordResetToken  sql.NullString `json:"passwordResetToken"`
	PasswordResetExpiry sql.NullTime   `json:"passwordResetExpiry"`
	Joined              time.Time      `json:"joined"`
}

// UserDto is the dto for a user entry.
type UserDto struct {
	ID                      int       `json:"id"`
	AvatarId                int       `json:"avatarId"`
	AvatarName              string    `json:"avatarName"`
	AvatarDescription       string    `json:"avatarDescription"`
	AvatarPrimaryImageUrl   string    `json:"avatarPrimaryImageUrl"`
	AvatarSecondaryImageUrl string    `json:"avatarSecondaryImageUrl"`
	Username                string    `json:"username"`
	Email                   string    `json:"email"`
	CountryCode             string    `json:"countryCode"`
	Joined                  time.Time `json:"joined"`
}

type AuthUserDto struct {
	ID                      int            `json:"id"`
	AvatarId                int            `json:"avatarId"`
	AvatarName              string         `json:"avatarName"`
	AvatarDescription       string         `json:"avatarDescription"`
	AvatarPrimaryImageUrl   string         `json:"avatarPrimaryImageUrl"`
	AvatarSecondaryImageUrl string         `json:"avatarSecondaryImageUrl"`
	Username                string         `json:"username"`
	Email                   string         `json:"email"`
	PasswordHash            string         `json:"passwordHash"`
	CountryCode             string         `json:"countryCode"`
	XP                      int            `json:"xp"`
	IsPremium               bool           `json:"isPremium"`
	IsAdmin                 bool           `json:"isAdmin"`
	PasswordResetToken      sql.NullString `json:"passwordResetToken"`
	PasswordResetExpiry     sql.NullTime   `json:"passwordResetExpiry"`
	Joined                  time.Time      `json:"joined"`
}

type UpdateUserDto struct {
	AvatarId                int    `json:"avatarId" validate:"required"`
	AvatarName              string `json:"avatarName"`
	AvatarDescription       string `json:"avatarDescription"`
	AvatarPrimaryImageUrl   string `json:"avatarPrimaryImageUrl"`
	AvatarSecondaryImageUrl string `json:"avatarSecondaryImageUrl"`
	Username                string `json:"username" validate:"required,username"`
	Email                   string `json:"email" validate:"required,email"`
	CountryCode             string `json:"countryCode" validate:"required"`
	XP                      *int   `json:"xp" validate:"required"`
}

type UpdateUserXPDto struct {
	Score    int `json:"score"`
	MaxScore int `json:"maxScore"`
}

type TotalUsersDto struct {
	Day   string `json:"day"`
	Count int    `json:"count"`
}

// GetUsers returns a page of users.
var GetUsers = func(limit int, offset int) ([]UserDto, error) {
	rows, err := Connection.Query("SELECT u.id, a.id, a.name, a.description, a.primaryimageurl, a.secondaryimageurl, u.username, u.email, u.countrycode, u.joined FROM users u JOIN avatars a on a.id = u.avatarid LIMIT $1 OFFSET $2;", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users = []UserDto{}
	for rows.Next() {
		var user UserDto
		if err = rows.Scan(&user.ID, &user.AvatarId, &user.AvatarName, &user.AvatarDescription, &user.AvatarPrimaryImageUrl, &user.AvatarSecondaryImageUrl, &user.Username, &user.Email, &user.CountryCode, &user.Joined); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// GetFirstID returns the first ID for a given page.
var GetFirstUserID = func(offset int) (int, error) {
	statement := "SELECT id FROM users LIMIT 1 OFFSET $1;"
	var id int
	err := Connection.QueryRow(statement, offset).Scan(&id)
	return id, err
}

// GetUser returns the user with a given id.
var GetUser = func(id int) (UserDto, error) {
	statement := "SELECT u.id, a.id, a.name, a.description, a.primaryimageurl, a.secondaryimageurl, u.username, u.email, u.countrycode, u.joined FROM users u JOIN avatars a on a.id = u.avatarid WHERE u.id = $1;"
	var user UserDto
	err := Connection.QueryRow(statement, id).Scan(&user.ID, &user.AvatarId, &user.AvatarName, &user.AvatarDescription, &user.AvatarPrimaryImageUrl, &user.AvatarSecondaryImageUrl, &user.Username, &user.Email, &user.CountryCode, &user.Joined)
	return user, err
}

var GetAuthUser = func(id int) (AuthUserDto, error) {
	statement := "SELECT u.id, a.id, a.name, a.description, a.primaryimageurl, a.secondaryimageurl, u.username, u.email, u.passwordhash, u.countrycode, u.xp, u.ispremium, u.isadmin, u.passwordresettoken, u.passwordresetexpiry, u.joined FROM users u JOIN avatars a on a.id = u.avatarid WHERE u.id = $1;"
	var user AuthUserDto
	err := Connection.QueryRow(statement, id).Scan(&user.ID, &user.AvatarId, &user.AvatarName, &user.AvatarDescription, &user.AvatarPrimaryImageUrl, &user.AvatarSecondaryImageUrl, &user.Username, &user.Email, &user.PasswordHash, &user.CountryCode, &user.XP, &user.IsPremium, &user.IsAdmin, &user.PasswordResetToken, &user.PasswordResetExpiry, &user.Joined)
	return user, err
}

var GetAuthUserUsingEmail = func(email string) (AuthUserDto, error) {
	statement := "SELECT u.id, a.id, a.name, a.description, a.primaryimageurl, a.secondaryimageurl, u.username, u.email, u.passwordhash, u.countrycode, u.xp, u.ispremium, u.isadmin, u.passwordresettoken, u.passwordresetexpiry, u.joined FROM users u JOIN avatars a on a.id = u.avatarid WHERE u.email = $1;"
	var user AuthUserDto
	err := Connection.QueryRow(statement, email).Scan(&user.ID, &user.AvatarId, &user.AvatarName, &user.AvatarDescription, &user.AvatarPrimaryImageUrl, &user.AvatarSecondaryImageUrl, &user.Username, &user.Email, &user.PasswordHash, &user.CountryCode, &user.XP, &user.IsPremium, &user.IsAdmin, &user.PasswordResetToken, &user.PasswordResetExpiry, &user.Joined)
	return user, err
}

// InsertUser inserts a new user into the users table.
var InsertUser = func(user User) (int, error) {
	statement := "INSERT INTO users (avatarid, username, email, passwordHash, countrycode, xp, isPremium, isAdmin, joined) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, user.AvatarId, user.Username, user.Email, user.PasswordHash, user.CountryCode, 0, false, false, time.Now()).Scan(&id)
	return id, err
}

// UpdateUser updates a user entry.
var UpdateUser = func(userID int, user UpdateUserDto) error {
	statement := "UPDATE users set avatarid = $2, username = $3, email = $4, countryCode = $5, xp = $6 WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, userID, user.AvatarId, user.Username, user.Email, user.CountryCode, user.XP).Scan(&id)
}

func UpdateUserXP(userID, score, maxScore int) (int, error) {
	increase := calculateXPIncrease(score, maxScore)
	statement := "UPDATE users set xp = xp + $2 WHERE id = $1 RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, userID, increase).Scan(&id)
	return increase, err
}

func calculateXPIncrease(score, maxScore int) int {
	percent := (float64(score) / float64(maxScore)) * 100
	if maxScore <= 112 {
		if percent > 66.6 {
			return 3
		}
		if percent > 33.3 {
			return 2
		}
		return 1
	}

	if percent > 90 {
		return 5
	}
	if percent > 75 {
		return 4
	}
	if percent > 50 {
		return 3
	}
	if percent > 25 {
		return 2
	}
	return 1
}

// DeleteUser deletes a users scores, leaderboard entries and then the user entry in the users table.
var DeleteUser = func(userID int) error {
	leaderboardStatement := "DELETE FROM leaderboard WHERE userId = $1;"
	Connection.QueryRow(leaderboardStatement, userID)
	usersStatement := "DELETE FROM users WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(usersStatement, userID).Scan(&id)
}

var UsernameExists = func(username string) (bool, error) {
	statement := "SELECT COUNT(id) FROM users WHERE lower(username) = $1;"
	var count int
	err := Connection.QueryRow(statement, strings.ToLower(username)).Scan(&count)
	return count > 0, err
}

var AnotherUserWithUsername = func(id int, username string) (bool, error) {
	statement := "SELECT COUNT(id) FROM users WHERE id != $1 AND lower(username) = $2;"
	var count int
	err := Connection.QueryRow(statement, id, strings.ToLower(username)).Scan(&count)
	return count > 0, err
}

var EmailExists = func(email string) (bool, error) {
	statement := "SELECT COUNT(id) FROM users WHERE lower(email) = $1;"
	var count int
	err := Connection.QueryRow(statement, strings.ToLower(email)).Scan(&count)
	return count > 0, err
}

var AnotherUserWithEmail = func(id int, email string) (bool, error) {
	statement := "SELECT COUNT(id) FROM users WHERE id != $1 AND lower(email) = $2;"
	var count int
	err := Connection.QueryRow(statement, id, strings.ToLower(email)).Scan(&count)
	return count > 0, err
}

var SetPasswordResetValues = func(userID int, resetToken string, expiryDate time.Time) error {
	statement := "UPDATE users set passwordResetToken = $2, passwordResetExpiry = $3 WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, userID, resetToken, expiryDate).Scan(&id)
}

var ResetPassword = func(userID int, passwordHash string) error {
	statement := "UPDATE users set passwordhash = $2, passwordResetToken = null, passwordResetExpiry = null WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, userID, passwordHash).Scan(&id)
}

func GetLastWeekTotalUsers() ([]TotalUsersDto, error) {
	var result []TotalUsersDto
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, 0-i)
		count, err := getTotalUsersToDate(date)
		if err != nil {
			return nil, err
		}

		current := TotalUsersDto{
			Day:   date.Weekday().String(),
			Count: count,
		}
		result = append(result, current)
	}
	return result, nil
}

func getTotalUsersToDate(date time.Time) (int, error) {
	var count int
	statement := "SELECT COUNT(id) FROM users WHERE joined < $1;"
	err := Connection.QueryRow(statement, date).Scan(&count)
	return count, err
}
