package database

// User is the database object for a user entry.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// GetUsers returns a page of users.
var GetUsers = func(limit int, offset int) ([]User, error) {
	rows, err := Connection.Query("SELECT * FROM users LIMIT $1 OFFSET $2;", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users = []User{}
	for rows.Next() {
		var user User
		if err = rows.Scan(&user.ID, &user.Username); err != nil {
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

// GetUserID returns the user's id for a given username.
var GetUserID = func(username string) (int, error) {
	statement := "SELECT id FROM users WHERE username = $1;"
	var id int
	err := Connection.QueryRow(statement, username).Scan(&id)
	return id, err
}

// GetUser returns the user with a given id.
var GetUser = func(id int) (User, error) {
	statement := "SELECT * FROM users WHERE id = $1;"
	var user User
	err := Connection.QueryRow(statement, id).Scan(&user.ID, &user.Username)
	return user, err
}

// InsertUser inserts a new user into the users table.
var InsertUser = func(user User) (int, error) {
	statement := "INSERT INTO users (username) VALUES ($1) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, user.Username).Scan(&id)
	return id, err
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
