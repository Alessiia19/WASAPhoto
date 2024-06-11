package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// CreatePhoto uploads a new photo to the database.
func (db *appdbimpl) CreatePhoto(p Photo) (Photo, error) {

	// Insert the new photo into the database.
	result, err := db.c.Exec("INSERT INTO photos (userid, username, imageData, uploadDate, likesCount, commentsCount) VALUES (?, ?, ?, ?, ?, ?)", p.UserID, p.Username, p.ImageData, p.UploadDate, p.LikesCount, p.CommentsCount)
	if err != nil {
		return p, fmt.Errorf("error creating photo in database: %w", err)
	}

	// Get the ID of the newly created photo.
	id, err := result.LastInsertId()
	if err != nil {
		return p, err
	}

	p.PhotoID = int(id)
	return p, nil
}

// LikePhoto adds a like to a photo in the database.
func (db *appdbimpl) LikePhoto(userID int, photoID int, l Like) error {

	// Check if the photo exists.
	var existingPhoto int
	err := db.c.QueryRow("SELECT 1 FROM photos WHERE photoid = ?", photoID).Scan(&existingPhoto)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.ErrNoRows // Photo not found
	} else if err != nil {
		return fmt.Errorf("error checking existing photo: %w", err)
	}

	// Check if the user has already liked the photo.
	var existingLike int
	err = db.c.QueryRow("SELECT 1 FROM likes WHERE userid = ? AND photoid = ?", userID, photoID).Scan(&existingLike)
	if err == nil {
		// The user has already liked the photo.
		return errors.New("already liked")
	}

	// Get the ID of the user who posted the photo.
	photoAuthorID, err := db.GetPhotoUserID(photoID)
	if err != nil {
		return err
	}

	// Check if the user who posted the photo has banned the current user.
	var isBanned int
	err = db.c.QueryRow("SELECT 1 FROM banned_users WHERE userid = ? AND banneduserid = ?", userID, photoAuthorID).Scan(&isBanned)
	if err == nil {
		return fmt.Errorf("cannot like a photo published by a user who has banned you")
	}

	// Increment the number of likes on the photo.
	_, err = db.c.Exec("UPDATE photos SET likesCount = likesCount + 1 WHERE photoid = ?", photoID)
	if err != nil {
		return fmt.Errorf("error updating likesCount in database: %w", err)
	}

	// Insert the like into the likes table.
	result, err := db.c.Exec("INSERT INTO likes (userID, photoID) VALUES (?, ?)", userID, photoID)
	if err != nil {
		return fmt.Errorf("error inserting like into database: %w", err)
	}

	// Get the ID of the newly created like.
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
	// Check if the photo exists.
	var existingPhoto int
	err := db.c.QueryRow("SELECT 1 FROM photos WHERE photoid = ?", photoID).Scan(&existingPhoto)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.ErrNoRows // Photo not found
	} else if err != nil {
		return fmt.Errorf("error checking existing photo: %w", err)
	}

	// Check if the like exists.
	var existingLike int
	err = db.c.QueryRow("SELECT 1 FROM likes WHERE likeid = ?", likeID).Scan(&existingLike)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.ErrNoRows // Like not found
	} else if err != nil {
		return fmt.Errorf("error checking existing like: %w", err)
	}

	// Decrement the number of likes on the photo.
	_, err = db.c.Exec("UPDATE photos SET likesCount = likesCount - 1 WHERE photoid = ?", photoID)
	if err != nil {
		return fmt.Errorf("error updating likesCount in database: %w", err)
	}

	// Remove the like from the likes table.
	_, err = db.c.Exec("DELETE FROM likes WHERE likeid = ? AND userid = ? AND photoid = ?", likeID, userID, photoID)
	if err != nil {
		return fmt.Errorf("error removing like from database: %w", err)
	}

	return nil
}

// CommentPhoto adds a comment to a photo in the database.
func (db *appdbimpl) CommentPhoto(userID, photoID int, authorUsername string, c Comment) (Comment, error) {
	// Check if the photo exists.
	var existingPhoto int
	err := db.c.QueryRow("SELECT 1 FROM photos WHERE photoid = ?", photoID).Scan(&existingPhoto)
	if errors.Is(err, sql.ErrNoRows) {
		return c, sql.ErrNoRows // Photo not found
	} else if err != nil {
		return c, fmt.Errorf("error checking existing photo: %w", err)
	}

	// Get the ID of the user who posted the photo.
	photoAuthorID, err := db.GetPhotoUserID(photoID)
	if err != nil {
		return c, err
	}

	// Check if the user who posted the photo has banned the current user.
	var isBanned int
	err = db.c.QueryRow("SELECT 1 FROM banned_users WHERE userid = ? AND banneduserid = ?", userID, photoAuthorID).Scan(&isBanned)
	if err == nil {
		return c, fmt.Errorf("cannot comment a photo published by a user who has banned you")
	}

	// Add the comment to the comments table.
	result, err := db.c.Exec("INSERT INTO comments (userid, username, photoid, commentText, uploadDate) VALUES (?, ?, ?, ?, ?)", userID, authorUsername, photoID, c.CommentText, c.UploadDate)
	if err != nil {
		return c, fmt.Errorf("error inserting comment into database: %w", err)
	}

	// Get the ID of the newly created comment.
	commentID, err := result.LastInsertId()
	if err != nil {
		return c, err
	}

	// Increment the number of comments on the photo.
	_, err = db.c.Exec("UPDATE photos SET commentsCount = commentsCount + 1 WHERE photoid = ?", photoID)
	if err != nil {
		return c, fmt.Errorf("error updating commentsCount in database: %w", err)
	}

	c.AuthorID = userID
	c.AuthorUsername = authorUsername
	c.CommentID = int(commentID)
	c.PhotoID = photoID

	return c, nil
}

// UncommentPhoto removes a specific comment from the specified photo in the database.
func (db *appdbimpl) UncommentPhoto(userID, photoID, commentID int) error {
	// Check if the photo exists.
	var existingPhoto int
	err := db.c.QueryRow("SELECT 1 FROM photos WHERE photoid = ?", photoID).Scan(&existingPhoto)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.ErrNoRows // Photo not found
	} else if err != nil {
		return fmt.Errorf("error checking existing photo: %w", err)
	}

	// Check if the comment exists and if the user who is trying to delete it is the author.
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

	// Decrement the number of comments on the photo.
	_, err = db.c.Exec("UPDATE photos SET commentsCount = commentsCount - 1 WHERE photoid = ?", photoID)
	if err != nil {
		return fmt.Errorf("error updating commentsCount in database: %w", err)
	}

	// Remove the comment from the comments table.
	_, err = db.c.Exec("DELETE FROM comments WHERE commentid = ? AND userid = ? AND photoid = ?", commentID, userID, photoID)
	if err != nil {
		return fmt.Errorf("error removing comment from database: %w", err)
	}

	return nil
}

// DeletePhoto removes a photo.
func (db *appdbimpl) DeletePhoto(userID, photoID int) error {

	// Check if the photo exists and belongs to the user.
	var existingPhotoUserID int
	err := db.c.QueryRow("SELECT userID FROM photos WHERE photoid = ?", photoID).Scan(&existingPhotoUserID)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.ErrNoRows // Photo not found
	} else if err != nil {
		return fmt.Errorf("error checking existing photo: %w", err)
	}

	if existingPhotoUserID != userID {
		return errors.New("cannot delete photos not published by you")
	}

	// Remove the photo from the photos table.
	_, err = db.c.Exec("DELETE FROM photos WHERE photoid = ?", photoID)
	if err != nil {
		return fmt.Errorf("error removing photo from database: %w", err)
	}

	// Remove likes associated with this photo.
	_, err = db.c.Exec("DELETE FROM likes WHERE photoid = ?", photoID)
	if err != nil {
		return fmt.Errorf("error removing photo's likes from database: %w", err)
	}

	// Remove comments associated with this photo.
	_, err = db.c.Exec("DELETE FROM comments WHERE photoid = ?", photoID)
	if err != nil {
		return fmt.Errorf("error removing photo's comments from database: %w", err)
	}

	return nil
}
