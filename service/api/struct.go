package api

import (
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
)

// User structure
type User struct {
	UserID   int    `json:"userID"`   // User's identifier
	Username string `json:"username"` // User's username
}

func (u *User) UserFromDatabase(user database.User) {
	u.UserID = user.UserID
	u.Username = user.Username
}

func (u *User) UserToDatabase() database.User {
	return database.User{
		UserID:   u.UserID,
		Username: u.Username,
	}
}

// Photo rappresenta un oggetto foto che include il nome dell'autore, l'URL dell'immagine, il numero di
// "mi piace" e commenti, e dettagli sugli utenti che hanno messo "mi piace" o commentato, compresi
// i commenti stessi.
type Photo struct {
	UserID        int       `json:"userID"`
	PhotoID       int       `json:"photoID"`
	ImageData     []byte    `json:"imageData"`
	UploadDate    time.Time `json:"uploadDate"`
	LikesCount    int       `json:"likesCount"`
	CommentsCount int       `json:"commentsCount"`
}

func (p *Photo) PhotoFromDatabase(photo database.Photo) {
	p.PhotoID = photo.PhotoID
	p.UserID = photo.UserID
	p.ImageData = photo.ImageData
	p.UploadDate = photo.UploadDate
	p.LikesCount = photo.LikesCount
	p.CommentsCount = photo.CommentsCount
}

func (p *Photo) PhotoToDatabase() database.Photo {
	return database.Photo{
		PhotoID:       p.PhotoID,
		UserID:        p.UserID,
		ImageData:     p.ImageData,
		UploadDate:    p.UploadDate,
		LikesCount:    p.LikesCount,
		CommentsCount: p.CommentsCount,
	}
}

type Profile struct {
	UserID         int     `json:"userID"`         // User's identifier
	Username       string  `json:"username"`       // User's username
	Followers      []User  `json:"followers"`      // followers list
	Following      []User  `json:"following"`      // following list
	FollowersCount int     `json:"followersCount"` // followers number
	FollowingCount int     `json:"followingCount"` // following number
	UploadedPhotos []Photo `json:"uploadedPhotos"`
}

type Like struct {
	LikeID  int `json:"likeID"`
	UserID  int `json:"userID"`
	PhotoID int `json:"photoID"`
}

func (l *Like) LikeFromDatabase(like database.Like) {
	l.LikeID = like.LikeID
	l.UserID = like.UserID
	l.PhotoID = like.PhotoID
}

func (l *Like) LikeToDatabase() database.Like {
	return database.Like{
		LikeID:  l.LikeID,
		UserID:  l.UserID,
		PhotoID: l.PhotoID,
	}
}

// Comment rappresenta un commento su una foto.
type Comment struct {
	CommentID   int    `json:"commentID"`
	AuthorID    int    `json:"authorID"`
	PhotoID     int    `json:"photoID"`
	CommentText string `json:"commentText"`
}

func (c *Comment) CommentFromDatabase(comment database.Comment) {
	c.CommentID = comment.CommentID
	c.AuthorID = comment.AuthorID
	c.PhotoID = comment.PhotoID
	c.CommentText = comment.CommentText
}

func (c *Comment) CommentToDatabase() database.Comment {
	return database.Comment{
		CommentID:   c.CommentID,
		AuthorID:    c.AuthorID,
		PhotoID:     c.PhotoID,
		CommentText: c.CommentText,
	}
}
