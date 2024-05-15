package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// CreatePhoto carica una nuova foto nel database.
func (db *appdbimpl) CreatePhoto(p Photo) (Photo, error) {

	// Esegui l'inserimento della nuova foto nel database.
	result, err := db.c.Exec("INSERT INTO photos (userid, username, imageData, uploadDate, likesCount, commentsCount) VALUES (?, ?, ?, ?, ?, ?)", p.UserID, p.Username, p.ImageData, p.UploadDate, p.LikesCount, p.CommentsCount)
	if err != nil {
		return p, fmt.Errorf("error creating photo in database: %w", err)
	}

	// Ottieni l'ID della foto appena creata.
	id, err := result.LastInsertId()
	if err != nil {
		return p, err
	}

	p.PhotoID = int(id)
	return p, nil
}

// LikePhoto implementa l'operazione di mettere like a una foto nel database.
func (db *appdbimpl) LikePhoto(userID int, photoID int, l Like) error {
	// Verifica se la foto esiste.
	var existingPhoto int
	err := db.c.QueryRow("SELECT 1 FROM photos WHERE photoid = ?", photoID).Scan(&existingPhoto)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.ErrNoRows // Foto non trovata
	} else if err != nil {
		return fmt.Errorf("error checking existing photo: %w", err)
	}

	// Controllo se l'utente ha già messo mi piace alla foto
	var existingLike int
	err = db.c.QueryRow("SELECT 1 FROM likes WHERE userid = ? AND photoid = ?", userID, photoID).Scan(&existingLike)
	if err == nil {
		// L'utente ha già messo mi piace
		return errors.New("already liked")
	}

	// Ottenere l'ID dell'utente che ha pubblicato la foto
	photoAuthorID, err := db.GetPhotoUserID(photoID)
	if err != nil {
		return err
	}

	// Controllo se l'utente sta cercando di mettere mi piace ad una sua foto - VALUTA SE TOGLIERLO
	if userID == photoAuthorID {
		return errors.New("cannot like your own photos")
	}

	// Verifica se l'utente che ha pubblicato la foto ha bannato l'utente corrente.
	var isBanned int
	err = db.c.QueryRow("SELECT 1 FROM banned_users WHERE userid = ? AND banneduserid = ?", userID, photoAuthorID).Scan(&isBanned)
	if err != nil {
		return fmt.Errorf("error checking if user is banned: %w", err)
	}

	if isBanned == 1 {
		return errors.New("cannot like a photo published by a user who has banned you")
	}

	// Incrementa il numero di likes della foto.
	_, err = db.c.Exec("UPDATE photos SET likesCount = likesCount + 1 WHERE photoid = ?", photoID)
	if err != nil {
		return fmt.Errorf("error updating likesCount in database: %w", err)
	}

	// Inserisci il like nella tabella likes.
	result, err := db.c.Exec("INSERT INTO likes (userID, photoID) VALUES (?, ?)", userID, photoID)
	if err != nil {
		return fmt.Errorf("error inserting like into database: %w", err)
	}

	// Ottieni l'ID del like appena creato.
	likeID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	l.UserID = userID
	l.PhotoID = photoID
	l.LikeID = int(likeID)

	return nil
}

// UnlikePhoto removes a specific like from the specified photo in the database.
func (db *appdbimpl) UnlikePhoto(userID, photoID, likeID int) error {
	// Check if the photo exists
	var existingPhoto int
	err := db.c.QueryRow("SELECT 1 FROM photos WHERE photoid = ?", photoID).Scan(&existingPhoto)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.ErrNoRows // Photo not found
	} else if err != nil {
		return fmt.Errorf("error checking existing photo: %w", err)
	}

	// Check if the like exists
	var existingLike int
	err = db.c.QueryRow("SELECT 1 FROM likes WHERE likeid = ?", likeID).Scan(&existingLike)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.ErrNoRows // Like not found
	} else if err != nil {
		return fmt.Errorf("error checking existing like: %w", err)
	}

	// Decrement the number of likes on the photo
	_, err = db.c.Exec("UPDATE photos SET likesCount = likesCount - 1 WHERE photoid = ?", photoID)
	if err != nil {
		return fmt.Errorf("error updating likesCount in database: %w", err)
	}

	// Remove the like from the likes table
	_, err = db.c.Exec("DELETE FROM likes WHERE likeid = ? AND userid = ? AND photoid = ?", likeID, userID, photoID)
	if err != nil {
		return fmt.Errorf("error removing like from database: %w", err)
	}

	return nil
}

// CommentPhoto aggiunge un commento a una foto nel database.
func (db *appdbimpl) CommentPhoto(userID, photoID int, c Comment) (Comment, error) {
	// Verifica se la foto esiste.
	var existingPhoto int
	err := db.c.QueryRow("SELECT 1 FROM photos WHERE photoid = ?", photoID).Scan(&existingPhoto)
	if errors.Is(err, sql.ErrNoRows) {
		return c, sql.ErrNoRows // Foto non trovata
	} else if err != nil {
		return c, fmt.Errorf("error checking existing photo: %w", err)
	}

	// Ottenere l'ID dell'utente che ha pubblicato la foto
	photoAuthorID, err := db.GetPhotoUserID(photoID)
	if err != nil {
		return c, err
	}

	// Verifica se l'utente che ha pubblicato la foto ha bannato l'utente corrente.
	var isBanned int
	err = db.c.QueryRow("SELECT 1 FROM banned_users WHERE userid = ? AND banneduserid = ?", userID, photoAuthorID).Scan(&isBanned)
	if err != nil {
		return c, fmt.Errorf("error checking if user is banned: %w", err)
	}

	if isBanned == 1 {
		return c, errors.New("cannot comment a photo published by a user who has banned you")
	}

	// Aggiungi il commento alla tabella dei commenti.
	result, err := db.c.Exec("INSERT INTO comments (userID, photoID, commentText) VALUES (?, ?, ?)", userID, photoID, c.CommentText)
	if err != nil {
		return c, fmt.Errorf("error inserting comment into database: %w", err)
	}

	// Ottieni l'ID del commento appena creato.
	commentID, err := result.LastInsertId()
	if err != nil {
		return c, err
	}

	// Incrementa il numero di commenti della foto.
	_, err = db.c.Exec("UPDATE photos SET commentsCount = commentsCount + 1 WHERE photoid = ?", photoID)
	if err != nil {
		return c, fmt.Errorf("error updating commentsCount in database: %w", err)
	}

	c.AuthorID = userID
	c.CommentID = int(commentID)
	c.PhotoID = photoID
	return c, nil
}

// UncommentPhoto removes a specific comment from the specified photo in the database.
func (db *appdbimpl) UncommentPhoto(userID, photoID, commentID int) error {
	// Check if the photo exists
	var existingPhoto int
	err := db.c.QueryRow("SELECT 1 FROM photos WHERE photoid = ?", photoID).Scan(&existingPhoto)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.ErrNoRows // Photo not found
	} else if err != nil {
		return fmt.Errorf("error checking existing photo: %w", err)
	}

	// Check if the comment exists and if the user who is trying to delete it is the author
	var commentAuthorID int
	err = db.c.QueryRow("SELECT userID FROM comments WHERE commentid = ? AND photoid = ?", commentID, photoID).Scan(&commentAuthorID)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.ErrNoRows // Comment not found
	} else if err != nil {
		return fmt.Errorf("error checking existing comment: %w", err)
	}

	if commentAuthorID != userID {
		return errors.New("cannot delete comments not published by you")
	}

	// Decrement the number of comments on the photo
	_, err = db.c.Exec("UPDATE photos SET commentsCount = commentsCount - 1 WHERE photoid = ?", photoID)
	if err != nil {
		return fmt.Errorf("error updating commentsCount in database: %w", err)
	}

	// Remove the comment from the comments table
	_, err = db.c.Exec("DELETE FROM comments WHERE commentid = ? AND userid = ? AND photoid = ?", commentID, userID, photoID)
	if err != nil {
		return fmt.Errorf("error removing comment from database: %w", err)
	}

	return nil
}

// DeletePhoto rimuove una foto.
func (db *appdbimpl) DeletePhoto(userID, photoID int) error {
	// Verifica se la foto esiste e appartiene all'utente
	var existingPhotoUserID int
	err := db.c.QueryRow("SELECT userID FROM photos WHERE photoid = ?", photoID).Scan(&existingPhotoUserID)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.ErrNoRows // Foto non trovata
	} else if err != nil {
		return fmt.Errorf("errore durante il controllo dell'esistenza della foto: %w", err)
	}

	if existingPhotoUserID != userID {
		return errors.New("cannot delete photos not published by you")
	}

	// Rimuovi la foto dalla tabella photos
	_, err = db.c.Exec("DELETE FROM photos WHERE photoid = ?", photoID)
	if err != nil {
		return fmt.Errorf("errore durante la rimozione della foto dal database: %w", err)
	}

	// Rimuovi i like associati a questa foto
	_, err = db.c.Exec("DELETE FROM likes WHERE photoid = ?", photoID)
	if err != nil {
		return fmt.Errorf("errore durante la rimozione dei like associati alla foto: %w", err)
	}

	// Rimuovi i commenti associati a questa foto
	_, err = db.c.Exec("DELETE FROM comments WHERE photoid = ?", photoID)
	if err != nil {
		return fmt.Errorf("errore durante la rimozione dei commenti associati alla foto: %w", err)
	}

	return nil
}
