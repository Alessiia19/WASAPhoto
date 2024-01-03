package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
)

// UpdateUsername aggiorna l'username dell'utente specificato nel database.
func (db *appdbimpl) UpdateUsername(userID int, newUsername string) error {
	// Verifica che l'username rispetti i requisiti desiderati.
	// (Aggiungi qui i controlli necessari, se vuoi farli)

	// Esegui l'aggiornamento dell'username nel database.
	result, err := db.c.Exec("UPDATE users SET username = ? WHERE userID = ?", newUsername, userID)
	if err != nil {
		// In caso di errore durante l'aggiornamento dell'username, restituisci l'errore.
		return fmt.Errorf("error updating username in database: %w", err)
	}

	// Verifica se l'utente è stato effettivamente aggiornato. - VALUTA SE LEVARE QUESTO CONTROLLO
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected after updating username: %w", err)
	}

	if rowsAffected == 0 {
		// Se nessuna riga è stata modificata, significa che l'utente con quell'ID non esiste.
		return sql.ErrNoRows
	}

	return nil
}

// FollowUser implementa l'operazione di follow nel database.
func (db *appdbimpl) FollowUser(userID, userIDToFollow int) error {

	// Controllo se l'utente sta cercando di eseguire l'operazione di follow su un utente che non esiste
	var existingUser int
	err := db.c.QueryRow("SELECT 1 FROM users WHERE userID = ?", userIDToFollow).Scan(&existingUser)
	if err == sql.ErrNoRows {
		return fmt.Errorf("The user you want to follow doesn't exists.")
	} else if err != nil {
		return fmt.Errorf("error checking existing user: %w", err)
	}

	// Controllo se l'utente è bloccato dall'utente che si sta cercando di seguire
	var isBanned int
	err = db.c.QueryRow("SELECT 1 FROM banned_users WHERE userID = ? AND bannedUserID = ?", userIDToFollow, userID).Scan(&isBanned)
	if err == nil {
		// L'utente è bloccato dall'altro utente
		return errors.New("you are banned by this user")
	} else if err != sql.ErrNoRows {
		// Altri errori diversi da ErrNoRows
		return fmt.Errorf("error checking if user is banned: %w", err)
	}

	// Controllo se l'utente ha bloccato l'altro utente
	var hasBanned int
	err = db.c.QueryRow("SELECT 1 FROM banned_users WHERE userID = ? AND bannedUserID = ?", userID, userIDToFollow).Scan(&hasBanned)
	if err == nil {
		// L'utente ha bloccato l'altro utente
		return errors.New("you have banned this user")
	} else if err != sql.ErrNoRows {
		// Altri errori diversi da ErrNoRows
		return fmt.Errorf("error checking if user is banned by you: %w", err)
	}

	// Controllo se l'utente segue già l'altro utente
	var existingFollower int
	err = db.c.QueryRow("SELECT 1 FROM followers WHERE userID = ? AND followerID = ?", userIDToFollow, userID).Scan(&existingFollower)
	if err == nil {
		// L'utente segue già l'altro utente
		return errors.New("already followed")
	} else if err != sql.ErrNoRows {
		// Altri errori diversi da ErrNoRows
		return fmt.Errorf("error checking existing follower: %w", err)
	}

	// Controllo se l'utente sta cercando di seguire se stesso
	if userID == userIDToFollow {
		return errors.New("cannot follow yourself")
	}

	// Aggiorna la tabella followers aggiungendo userID come follower di userIDToFollow
	_, err = db.c.Exec("INSERT INTO followers (userID, followerID) VALUES (?, ?)", userIDToFollow, userID)
	if err != nil {
		return fmt.Errorf("error updating followers table: %w", err)
	}

	// Aggiorna la tabella following aggiungendo userIDToFollow come following di userID
	_, err = db.c.Exec("INSERT INTO following (userID, followingID) VALUES (?, ?)", userID, userIDToFollow)
	if err != nil {
		return fmt.Errorf("error updating following table: %w", err)
	}

	return nil
}

// UnfollowUser implementa l'operazione di unfollow nel database.
func (db *appdbimpl) UnfollowUser(userID, followingID int) error {
	// Controllo se l'utente sta cercando di eseguire l'operazione di unfollow su un utente che non sta seguendo
	/*
		var existingFollower int
		err := db.c.QueryRow("SELECT 1 FROM followers WHERE userID = ? AND followerID = ?", followingID, userID).Scan(&existingFollower)
		if err == sql.ErrNoRows {
			return fmt.Errorf("You are trying to unfollow someone you don't follow.")
		} else if err != nil {
			return fmt.Errorf("error checking existing follower: %w", err)
		}
	*/

	// Rimuovi l'utente dai followers.
	_, err := db.c.Exec("DELETE FROM followers WHERE userID = ? AND followerID = ?", followingID, userID)
	if err != nil {
		return fmt.Errorf("error removing follower: %w", err)
	}

	// Rimuovi l'utente dai following.
	_, err = db.c.Exec("DELETE FROM following WHERE userID = ? AND followingID = ?", userID, followingID)
	if err != nil {
		return fmt.Errorf("error removing following: %w", err)
	}

	return nil
}

// BanUser implementa l'operazione di ban nel database.
func (db *appdbimpl) BanUser(UserID, bannedUserID int) error {
	// Verifica se l'utente da bloccare esiste nella tabella `users`.
	var existingUser int
	err := db.c.QueryRow("SELECT 1 FROM users WHERE userID = ?", bannedUserID).Scan(&existingUser)
	if err == sql.ErrNoRows {
		return fmt.Errorf("The user you want to ban doesn't exists.")
	} else if err != nil {
		return fmt.Errorf("error checking existing user: %w", err)
	}

	// Controllo se l'utente sta cercando di bannare se stesso
	if UserID == bannedUserID {
		return errors.New("cannot ban yourself")
	}

	// Controllo se l'utente è stato bannato da chi si sta cercando di bannare
	var isBanned int
	err = db.c.QueryRow("SELECT 1 FROM banned_users WHERE userID = ? AND bannedUserID = ?", bannedUserID, UserID).Scan(&isBanned)
	if err == nil {
		// L'utente è bloccato dall'altro utente
		return errors.New("you are banned by this user")
	} else if err != sql.ErrNoRows {
		// Altri errori diversi da ErrNoRows
		return fmt.Errorf("error checking if user is banned: %w", err)
	}

	// Controllo se l'utente ha bloccato l'altro utente
	var hasBanned int
	err = db.c.QueryRow("SELECT 1 FROM banned_users WHERE userID = ? AND bannedUserID = ?", UserID, bannedUserID).Scan(&hasBanned)
	if err == nil {
		// L'utente ha bloccato l'altro utente
		return errors.New("you have banned this user")
	} else if err != sql.ErrNoRows {
		// Altri errori diversi da ErrNoRows
		return fmt.Errorf("error checking if user is banned by you: %w", err)
	}

	// Utilizza la funzione UnfollowUser per gestire la rimozione dai followers e dai following.
	// L'utente bloccato smette automaticamente di essere sia follower che following dell'utente che lo blocca
	if err := db.UnfollowUser(UserID, bannedUserID); err != nil {
		return fmt.Errorf("error handling unfollow during ban operation: %w", err)
	}
	if err := db.UnfollowUser(bannedUserID, UserID); err != nil {
		return fmt.Errorf("error handling unfollow during ban operation: %w", err)
	}

	// Aggiorna la tabella dei bloccati.
	_, err = db.c.Exec("INSERT INTO banned_users (userID, bannedUserID) VALUES (?, ?)", UserID, bannedUserID)
	if err != nil {
		return fmt.Errorf("error updating banned_users table: %w", err)
	}

	return nil
}

// UnfollowUser implementa l'operazione di unfollow nel database.
func (db *appdbimpl) UnbanUser(userID, bannedUserID int) error {
	// Controllo se l'utente sta cercando di eseguire l'operazione di unbanning su un utente che non è stato bannato o non esiste
	var existingUser int
	err := db.c.QueryRow("SELECT 1 FROM banned_users WHERE userID = ? AND bannedUserID = ?", userID, bannedUserID).Scan(&existingUser)
	if err == sql.ErrNoRows {
		return fmt.Errorf("You are trying to unban someone who was not banned or doesn't exists.")
	} else if err != nil {
		return fmt.Errorf("error checking existing ban: %w", err)
	}

	// Rimuovi l'utente dai banned_users.
	_, err = db.c.Exec("DELETE FROM banned_users WHERE userID = ? AND bannedUserID = ?", userID, bannedUserID)
	if err != nil {
		return fmt.Errorf("error removing ban: %w", err)
	}

	return nil
}

// GetUserProfile restituisce i dettagli del profilo dell'utente specificato,
// comprese le foto caricate, il numero di follower, il numero di utenti seguiti
// e il numero totale di foto caricate.
func (db *appdbimpl) GetUserProfile(requestingUserID, requestedUserID int) (Profile, error) {
	var profile Profile

	// Controllo se l'utente che si vuole cercare esiste
	var existingUser int
	err := db.c.QueryRow("SELECT 1 FROM users WHERE userID = ?", requestedUserID).Scan(&existingUser)
	if err == sql.ErrNoRows {
		return profile, fmt.Errorf("The user you are searching doesn't exist")
	} else if err != nil {
		return profile, fmt.Errorf("error checking existing user: %w", err)
	}

	// Controllo se l'utente è bloccato dall'utente che si sta cercando
	var isBanned int
	err = db.c.QueryRow("SELECT 1 FROM banned_users WHERE userID = ? AND bannedUserID = ?", requestedUserID, requestingUserID).Scan(&isBanned)
	if err == nil {
		// L'utente è bloccato dall'altro utente
		return profile, errors.New("you are banned by this user")
	} else if err != sql.ErrNoRows {
		// Altri errori diversi da ErrNoRows
		return profile, fmt.Errorf("error checking if user is banned: %w", err)
	}

	// Ottenere dettagli utente (Username, Followers, Following, UploadedPhotosCount)
	user, err := db.getUserDetails(requestedUserID)
	if err != nil {
		return profile, err
	}

	// Ottenere la lista di follower
	followers, err := db.getFollowers(requestedUserID)
	if err != nil {
		return profile, err
	}

	// Ottenere la lista di utenti seguiti
	following, err := db.getFollowing(requestedUserID)
	if err != nil {
		return profile, err
	}

	// Ottenere la lista di foto caricate dall'utente
	uploadedPhotos, err := db.getUploadedPhotos(requestedUserID)
	if err != nil {
		return profile, err
	}

	// Costruire il profilo dell'utente
	profile = Profile{
		UserID:         user.UserID,
		Username:       user.Username,
		Followers:      followers,
		Following:      following,
		FollowersCount: len(followers),
		FollowingCount: len(following),
		UploadedPhotos: uploadedPhotos,
	}
	return profile, err
}

// GetMyStream restituisce lo stream dell'utente, costituito dalle foto pubblicate dai suoi following,
// comprese di dettagli quali il date-time in cui sono state pubblicate, il numero di like e di commenti.
func (db *appdbimpl) GetMyStream(userID int) ([]Photo, error) {
	// Ottieni la lista di utenti seguiti dall'utente specificato
	following, err := db.getFollowing(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting following users: %w", err)
	}

	var stream []Photo

	// Itera su ogni utente seguito per ottenere le foto dallo stream
	for _, followedUser := range following {
		// Ottieni la lista di foto caricate dall'utente seguito
		uploadedPhotos, err := db.getUploadedPhotos(followedUser.UserID)
		if err != nil {
			return nil, fmt.Errorf("error getting uploaded photos for user %d: %w", followedUser.UserID, err)
		}

		// Aggiungi le foto allo stream
		stream = append(stream, uploadedPhotos...)
	}

	// Ordina lo stream in ordine cronologico inverso
	sort.Slice(stream, func(i, j int) bool {
		return stream[i].UploadDate.After(stream[j].UploadDate)
	})

	return stream, nil
}
