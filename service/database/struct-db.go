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
	ImageData     []byte    `json:"imageData"`
	UploadDate    time.Time `json:"uploadDate"`
	LikesCount    int       `json:"likesCount"`
	CommentsCount int       `json:"commentsCount"`
	//WhoLiked      []int     `json:"whoLiked"`
	//Comments    []Comment `json:"comments"`
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
	UserID         int     `json:"userID"`         // User's identifier
	Username       string  `json:"username"`       // User's username
	Followers      []User  `json:"followers"`      // followers list
	Following      []User  `json:"following"`      // following list
	FollowersCount int     `json:"followersCount"` // followers number
	FollowingCount int     `json:"followingCount"` // following number
	UploadedPhotos []Photo `json:"uploadedPhotos"`
}
