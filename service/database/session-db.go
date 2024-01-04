package database

func (db *appdbimpl) CreateUser(u User) (User, error) {

	var user User
	// Controlla se l'utente esiste già nel database.
	err := db.c.QueryRow("SELECT userID, username FROM users WHERE username = ?", u.Username).Scan(&user.UserID, &user.Username)
	if err == nil {
		// L'utente esiste già, restituisci i suoi dati.
		return user, nil
	}

	// Se l'utente non esiste, crea un nuovo utente.
	result, err := db.c.Exec("INSERT INTO users (username) VALUES (?)", u.Username)
	if err != nil {
		return user, err
	}

	// Ottieni l'ID dell'utente appena creato.
	id, err := result.LastInsertId()
	if err != nil {
		return u, err
	}

	// Aggiorna i dati dell'utente.
	u.UserID = int(id)
	return u, nil
}
