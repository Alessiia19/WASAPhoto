package database

import (
	"database/sql"
	"errors"
	"fmt"
)

//          --- GET DATA FROM DATABASE FUNCTIONS ---

// GetPhotoUserID restituisce l'ID dell'utente che ha pubblicato la foto specificata.
func (db *appdbimpl) GetPhotoUserID(photoID int) (int, error) {
	var userID int
	err := db.c.QueryRow("SELECT userID FROM photos WHERE photoID = ?", photoID).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("error getting photo user ID: %w", err)
	}
	return userID, nil
}

// getUserDetails ottiene i dettagli dell'utente, compresi il nome utente e l'ID.
func (db *appdbimpl) getUserDetails(userID int) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT userID, username FROM users WHERE userID = ?", userID).
		Scan(&user.UserID, &user.Username)

	if errors.Is(err, sql.ErrNoRows) {
		return user, fmt.Errorf("user not found")
	} else if err != nil {
		return user, fmt.Errorf("error fetching user details: %w", err)
	}

	return user, nil
}

// getFollowers ottiene la lista di follower per l'utente specificato.
func (db *appdbimpl) getFollowers(userID int) ([]User, error) {
	var followers []User

	rows, err := db.c.Query("SELECT u.userID, u.username FROM users u JOIN followers f ON u.userID = f.followerID WHERE f.userID = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching followers: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var follower User
		if err := rows.Scan(&follower.UserID, &follower.Username); err != nil {
			return nil, fmt.Errorf("error scanning follower row: %w", err)
		}
		followers = append(followers, follower)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over follower rows: %w", err)
	}

	return followers, nil
}

// getFollowing ottiene la lista di utenti seguiti dall'utente specificato.
func (db *appdbimpl) getFollowing(userID int) ([]User, error) {
	var following []User

	rows, err := db.c.Query("SELECT u.userID, u.username FROM users u JOIN following f ON u.userID = f.followingID WHERE f.userID = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching following users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var followedUser User
		if err := rows.Scan(&followedUser.UserID, &followedUser.Username); err != nil {
			return nil, fmt.Errorf("error scanning following user row: %w", err)
		}
		following = append(following, followedUser)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over following user rows: %w", err)
	}

	return following, nil
}

// getUploadedPhotos ottiene la lista di foto caricate dall'utente specificato.
func (db *appdbimpl) getUploadedPhotos(userID int) ([]Photo, error) {
	var uploadedPhotos []Photo

	rows, err := db.c.Query("SELECT photoID, userID, imageData, uploadDate, likesCount, commentsCount FROM photos WHERE userID = ? ORDER BY uploadDate DESC", userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching uploaded photos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var photo Photo
		if err := rows.Scan(&photo.PhotoID, &photo.UserID, &photo.ImageData, &photo.UploadDate, &photo.LikesCount, &photo.CommentsCount); err != nil {
			return nil, fmt.Errorf("error scanning uploaded photo row: %w", err)
		}
		uploadedPhotos = append(uploadedPhotos, photo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over uploaded photo rows: %w", err)
	}

	return uploadedPhotos, nil
}
