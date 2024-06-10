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
	err := db.c.QueryRow("SELECT userid FROM photos WHERE photoid = ?", photoID).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("error getting photo user ID: %w", err)
	}
	return userID, nil
}

// getUserDetails retrieves user details (userid, username)
func (db *appdbimpl) GetUserDetails(userID int) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT userid, username FROM users WHERE userid = ?", userID).
		Scan(&user.UserID, &user.Username)

	if errors.Is(err, sql.ErrNoRows) {
		return user, fmt.Errorf("user not found")
	} else if err != nil {
		return user, fmt.Errorf("error fetching user details: %w", err)
	}

	return user, nil
}

// getFollowers retrieves the list of followers for the specified user.
func (db *appdbimpl) GetFollowers(userID int) ([]User, error) {
	var followers []User

	rows, err := db.c.Query("SELECT u.userid, u.username FROM users u JOIN followers f ON u.userid = f.followerid WHERE f.userid = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching followers: %w", err)
	}
	defer rows.Close() // Ensure the rows are closed after the query.

	// Iterate over the rows to extract each follower's data.
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

// getFollowing retrieves the list of users followed for the specified user.
func (db *appdbimpl) GetFollowing(userID int) ([]User, error) {
	var following []User

	rows, err := db.c.Query("SELECT u.userid, u.username FROM users u JOIN following f ON u.userid = f.followingid WHERE f.userid = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching following users: %w", err)
	}
	defer rows.Close() // Ensure the rows are closed after the query.

	// Iterate over the rows to extract each following's data.
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

// getLikes retrieves the list of likes for the specified photo.
func (db *appdbimpl) GetLikes(photoID int) ([]Like, error) {
	var likes []Like
	rows, err := db.c.Query("SELECT likeid, userid, photoid FROM likes WHERE photoid = ?", photoID)
	if err != nil {
		return nil, fmt.Errorf("error fetching likes: %w", err)
	}
	defer rows.Close() // Ensure the rows are closed after the query.

	// Iterate over the rows to extract each like's data.
	for rows.Next() {
		var like Like
		if err := rows.Scan(&like.LikeID, &like.UserID, &like.PhotoID); err != nil {
			return nil, fmt.Errorf("error scanning like row: %w", err)
		}
		likes = append(likes, like)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over likes rows: %w", err)
	}

	return likes, nil
}

// getComments retrieves the list of comments for the specified photo.
func (db *appdbimpl) GetComments(photoID int) ([]Comment, error) {
	var comments []Comment
	rows, err := db.c.Query("SELECT commentid, userid, username, photoid, commentText, uploadDate FROM comments WHERE photoid = ?", photoID)
	if err != nil {
		return nil, fmt.Errorf("error fetching comments: %w", err)
	}
	defer rows.Close() // Ensure the rows are closed after the query.

	// Iterate over the rows to extract each comment's data.
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.CommentID, &comment.AuthorID, &comment.AuthorUsername, &comment.PhotoID, &comment.CommentText, &comment.UploadDate); err != nil {
			return nil, fmt.Errorf("error scanning comment row: %w", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over comment rows: %w", err)
	}

	return comments, nil
}

// getUploadedPhotos retrieves the list of photos uploaded by the user specified.
func (db *appdbimpl) GetUploadedPhotos(userID int) ([]CompletePhoto, error) {
	var uploadedPhotos []CompletePhoto

	// Fetch all photos uploaded by the user, ordered by upload date in descending order.
	rows, err := db.c.Query("SELECT photoid, userid, username, imageData, uploadDate, likesCount, commentsCount FROM photos WHERE userid = ? ORDER BY uploadDate DESC", userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching uploaded photos: %w", err)
	}
	defer rows.Close() // Ensure the rows are closed after the query.

	// Iterate over the query results to read each photo's data.
	for rows.Next() {
		var photo CompletePhoto
		if err := rows.Scan(&photo.PhotoID, &photo.UserID, &photo.Username, &photo.ImageData, &photo.UploadDate, &photo.LikesCount, &photo.CommentsCount); err != nil {
			return nil, fmt.Errorf("error scanning uploaded photo row: %w", err)
		}

		// Retrieve the list of likes for each photo.
		photo.Likes, err = db.GetLikes(photo.PhotoID)
		if err != nil {
			return nil, fmt.Errorf("error fetching likes for photoID %d: %w", photo.PhotoID, err)
		}

		// Retrieve the list of comments for each photo.
		photo.Comments, err = db.GetComments(photo.PhotoID)
		if err != nil {
			return nil, fmt.Errorf("error fetching comments for photoID %d: %w", photo.PhotoID, err)
		}

		uploadedPhotos = append(uploadedPhotos, photo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over uploaded photo rows: %w", err)
	}

	return uploadedPhotos, nil
}
