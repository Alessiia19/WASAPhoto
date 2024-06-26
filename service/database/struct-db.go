package database

import "time"

// User structure
type User struct {
	UserID   int    `json:"userID"`   // User's identifier
	Username string `json:"username"` // User's username
}

// Photo structure
type Photo struct {
	UserID        int       `json:"userID"`
	PhotoID       int       `json:"photoID"`
	Username      string    `json:"username"`
	ImageData     []byte    `json:"imageData"`
	UploadDate    time.Time `json:"uploadDate"`
	LikesCount    int       `json:"likesCount"`
	CommentsCount int       `json:"commentsCount"`
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

// Like structure
type Like struct {
	LikeID  int `json:"likeID"`
	UserID  int `json:"userID"`
	PhotoID int `json:"photoID"`
}

// Comment structure
type Comment struct {
	CommentID      int       `json:"commentID"`
	AuthorID       int       `json:"authorID"`
	AuthorUsername string    `json:"authorUsername"`
	PhotoID        int       `json:"photoID"`
	CommentText    string    `json:"commentText"`
	UploadDate     time.Time `json:"uploadDate"`
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
