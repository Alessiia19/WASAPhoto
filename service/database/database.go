/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	CreateUser(User) (User, error)
	UpdateUsername(int, string) error
	CreatePhoto(Photo) (Photo, error)
	FollowUser(int, int) error
	UnfollowUser(int, int) error
	BanUser(int, int) error
	UnbanUser(int, int) error
	LikePhoto(int, int, Like) error
	UnlikePhoto(int, int, int) error
	CommentPhoto(int, int, Comment) (Comment, error)
	UncommentPhoto(int, int, int) error
	DeletePhoto(int, int) error
	GetUserProfile(int, int) (Profile, error)
	GetMyStream(int) ([]Photo, error)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string

	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='user';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = createTables(db)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func createTables(db *sql.DB) error {
	db.Exec("PRAGMA foreign_key=ON;")
	userQuery := `CREATE TABLE IF NOT EXISTS users (
				userID INTEGER PRIMARY KEY AUTOINCREMENT,
				username TEXT,
				UNIQUE (userID, username)
				);`

	_, err := db.Exec(userQuery)

	if err != nil {
		return fmt.Errorf("error creating users structure: %w", err)
	}

	// VEDI SE TOGLIERE IMAGEDATA
	photoQuery := `CREATE TABLE photos (photoID INTEGER PRIMARY KEY AUTOINCREMENT,
					userID INTEGER,
					imageData BLOB,
					uploadDate DATETIME,
					likesCount INTEGER,
					commentsCount INTEGER,
					FOREIGN KEY(userID) REFERENCES user(userID) ON DELETE CASCADE);`
	_, err = db.Exec(photoQuery)
	if err != nil {
		return fmt.Errorf("error creating photos structure: %w", err)
	}

	followersQuery := `CREATE TABLE IF NOT EXISTS followers (
		userID INTEGER,
		followerID INTEGER,
		PRIMARY KEY (userID, followerID),
		FOREIGN KEY (userID) REFERENCES users(userID),
		FOREIGN KEY (followerID) REFERENCES users(userID)
		);`
	_, err = db.Exec(followersQuery)
	if err != nil {
		return fmt.Errorf("error creating followers table: %w", err)
	}

	followingQuery := `CREATE TABLE IF NOT EXISTS following (
            userID INTEGER,
            followingID INTEGER,
            PRIMARY KEY (userID, followingID),
            FOREIGN KEY (userID) REFERENCES users(userID),
            FOREIGN KEY (followingID) REFERENCES users(userID)
			);`

	_, err = db.Exec(followingQuery)
	if err != nil {
		return fmt.Errorf("error creating following table: %w", err)
	}

	banQuery := `CREATE TABLE banned_users (
		userID INTEGER,
		bannedUserID INTEGER,
		PRIMARY KEY (userID, bannedUserID),
		FOREIGN KEY (userID) REFERENCES user(userID),
		FOREIGN KEY (bannedUserID) REFERENCES user(userID));`
	_, err = db.Exec((banQuery))
	if err != nil {
		return fmt.Errorf("error creating ban structure: %w", err)
	}

	likesQuery := `CREATE TABLE likes (
		likeID INTEGER PRIMARY KEY AUTOINCREMENT,
		userID INTEGER,
		photoID INTEGER,
		FOREIGN KEY(userID) REFERENCES users(userID) ON DELETE CASCADE,
		FOREIGN KEY(photoID) REFERENCES photos(photoID) ON DELETE CASCADE
	);`
	_, err = db.Exec(likesQuery)
	if err != nil {
		return fmt.Errorf("error creating likes structure: %w", err)
	}
	commentQuery := `CREATE TABLE IF NOT EXISTS comments (
		commentID INTEGER PRIMARY KEY AUTOINCREMENT,
		userID INTEGER,
		photoID INTEGER,
		upload_date DATETIME,
		commentText TEXT,
		FOREIGN KEY(userID) REFERENCES users(userID) ON DELETE CASCADE,
		FOREIGN KEY(photoID) REFERENCES photos(photoID) ON DELETE CASCADE
	);`

	_, err = db.Exec(commentQuery)
	if err != nil {
		return fmt.Errorf("error creating comments structure: %w", err)
	}

	/*
		commentQuery := `CREATE TABLE comment (commentId INTEGER PRIMARY KEY AUTOINCREMENT,
						commentText TEXT,
						upload_date DATETIME,
						userID INTEGER,
						photoId INTEGER,
						FOREIGN KEY(userID) REFERENCES user(userID) ON DELETE CASCADE,
						FOREIGN KEY(photoId) REFERENCES photo(photoId) ON DELETE CASCADE);`
		_, err = db.Exec(commentQuery)
		if err != nil {
			return fmt.Errorf("error creating comment structure: %w", err)
		}
		likeQuery := `CREATE TABLE like (userID INTEGER,
						likedPhotoId INTEGER,
						PRIMARY KEY (userID, likedPhotoId),
						FOREIGN KEY (userID) REFERENCES user(userID),
						FOREIGN KEY (likedPhotoId) REFERENCES photo(photoId));`
		_, err = db.Exec(likeQuery)
		if err != nil {
			return fmt.Errorf("error creating like structure: %w", err)
		}
	*/

	return nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
