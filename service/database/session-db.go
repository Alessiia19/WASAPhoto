package database

func (db *appdbimpl) CreateUser(u User) (User, error) {

	var user User
	// Check if the user already exists in the database.
	err := db.c.QueryRow("SELECT userid, username FROM users WHERE username = ?", u.Username).Scan(&user.UserID, &user.Username)
	if err == nil {
		// The user already exists, return their data.
		return user, nil
	}

	// If the user doesn't exist, create a new user.
	result, err := db.c.Exec("INSERT INTO users (username) VALUES (?)", u.Username)
	if err != nil {
		return user, err
	}

	// Get the ID of the newly created user.
	id, err := result.LastInsertId()
	if err != nil {
		return u, err
	}

	// Update the user's data.
	u.UserID = int(id)
	return u, nil
}
