package api

import (
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
)

// User structure.
type User struct {
	UserID   int    `json:"userID"`   // User's identifier
	Username string `json:"username"` // User's username
}

// UserFromDatabase updates the current User struct with data from a database.User struct.
func (u *User) UserFromDatabase(user database.User) {
	u.UserID = user.UserID
	u.Username = user.Username
}

// UserToDatabase converts the current User struct to a database.User struct.
func (u *User) UserToDatabase() database.User {
	return database.User{
		UserID:   u.UserID,
		Username: u.Username,
	}
}

// Photo structure.
type Photo struct {
	UserID        int       `json:"userID"`
	PhotoID       int       `json:"photoID"`
	Username      string    `json:"username"`
	ImageData     []byte    `json:"imageData"`
	UploadDate    time.Time `json:"uploadDate"`
	LikesCount    int       `json:"likesCount"`
	CommentsCount int       `json:"commentsCount"`
}

// PhotoFromDatabase updates the current Photo struct with data from a database.Photo struct.
func (p *Photo) PhotoFromDatabase(photo database.Photo) {
	p.PhotoID = photo.PhotoID
	p.UserID = photo.UserID
	p.Username = photo.Username
	p.ImageData = photo.ImageData
	p.UploadDate = photo.UploadDate
	p.LikesCount = photo.LikesCount
	p.CommentsCount = photo.CommentsCount
}

// PhotoToDatabase converts the current Photo struct to a database.Photo struct.
func (p *Photo) PhotoToDatabase() database.Photo {
	return database.Photo{
		PhotoID:       p.PhotoID,
		UserID:        p.UserID,
		Username:      p.Username,
		ImageData:     p.ImageData,
		UploadDate:    p.UploadDate,
		LikesCount:    p.LikesCount,
		CommentsCount: p.CommentsCount,
	}
}

// CompletePhoto represents a photo object that includes the author's username, the image URL, the number of "likes" and comments,
// and details about users who have liked or commented, including the likes and comments themselves.
type CompletePhoto struct {
	UserID        int       `json:"userID"`
	PhotoID       int       `json:"photoID"`
	Username      string    `json:"username"`
	ImageData     []byte    `json:"imageData"`
	UploadDate    time.Time `json:"uploadDate"`
	LikesCount    int       `json:"likesCount"`
	Likes         []Like    `json:"likes"`
	CommentsCount int       `json:"commentsCount"`
	Comments      []Comment `json:"comments"`
}

// Profile structure that includes the number of "followers", "following" and photo uploaded, including their arrays
type Profile struct {
	UserID              int             `json:"userID"`              // User's identifier
	Username            string          `json:"username"`            // User's username
	Followers           []User          `json:"followers"`           // followers list
	Following           []User          `json:"following"`           // following list
	FollowersCount      int             `json:"followersCount"`      // followers number
	FollowingCount      int             `json:"followingCount"`      // following number
	UploadedPhotos      []CompletePhoto `json:"uploadedPhotos"`      // Photos array
	UploadedPhotosCount int             `json:"uploadedPhotosCount"` // Uploaded photos number
}

// Like structure.
type Like struct {
	LikeID  int `json:"likeID"`
	UserID  int `json:"userID"`
	PhotoID int `json:"photoID"`
}

// LikeFromDatabase updates the current Like struct with data from a database.Like struct.
func (l *Like) LikeFromDatabase(like database.Like) {
	l.LikeID = like.LikeID
	l.UserID = like.UserID
	l.PhotoID = like.PhotoID
}

// LikeToDatabase converts the current Like struct to a database.Like struct.
func (l *Like) LikeToDatabase() database.Like {
	return database.Like{
		LikeID:  l.LikeID,
		UserID:  l.UserID,
		PhotoID: l.PhotoID,
	}
}

// Comment structure.
type Comment struct {
	CommentID      int       `json:"commentID"`
	AuthorID       int       `json:"authorID"`
	AuthorUsername string    `json:"authorUsername"`
	PhotoID        int       `json:"photoID"`
	CommentText    string    `json:"commentText"`
	UploadDate     time.Time `json:"uploadDate"`
}

// CommentFromDatabase updates the current Comment struct with data from a database.Comment struct.
func (c *Comment) CommentFromDatabase(comment database.Comment) {
	c.CommentID = comment.CommentID
	c.AuthorID = comment.AuthorID
	c.AuthorUsername = comment.AuthorUsername
	c.PhotoID = comment.PhotoID
	c.CommentText = comment.CommentText
	c.UploadDate = comment.UploadDate
}

// CommentToDatabase converts the current Comment struct to a database.Comment struct.
func (c *Comment) CommentToDatabase() database.Comment {
	return database.Comment{
		CommentID:      c.CommentID,
		AuthorID:       c.AuthorID,
		AuthorUsername: c.AuthorUsername,
		PhotoID:        c.PhotoID,
		CommentText:    c.CommentText,
		UploadDate:     c.UploadDate,
	}
}
