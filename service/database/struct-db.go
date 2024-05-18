package database

import "time"

// User structure
type User struct {
	UserID   int    `json:"userID"`   // User's identifier
	Username string `json:"username"` // User's username
}

// Photo rappresenta un oggetto foto che include il nome dell'autore, l'URL dell'immagine, il numero di
// "mi piace" e commenti, e dettagli sugli utenti che hanno messo "mi piace" o commentato, compresi
// i commenti stessi.
type Photo struct {
	UserID        int       `json:"userID"`
	PhotoID       int       `json:"photoID"`
	Username      string    `json:"username"`
	ImageData     []byte    `json:"imageData"`
	UploadDate    time.Time `json:"uploadDate"`
	LikesCount    int       `json:"likesCount"`
	CommentsCount int       `json:"commentsCount"`
}

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

type Like struct {
	LikeID  int `json:"likeID"`
	UserID  int `json:"userID"`
	PhotoID int `json:"photoID"`
}

// Comment rappresenta un commento su una foto.
type Comment struct {
	CommentID   int    `json:"commentID"`
	AuthorID    int    `json:"authorID"`
	PhotoID     int    `json:"photoID"`
	CommentText string `json:"commentText"`
}

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
