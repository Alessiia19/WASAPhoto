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
	GetMyStream(int) ([]CompletePhoto, error)
	GetUsers(string) ([]User, error)

	// utils
	GetPhotoUserID(int) (int, error)
	GetUserDetails(int) (User, error)
	GetFollowers(int) ([]User, error)
	GetFollowing(int) ([]User, error)
	GetUploadedPhotos(int) ([]CompletePhoto, error)

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
	_, err := db.Exec("PRAGMA foreign_key=ON;")
	if err != nil {
		return fmt.Errorf("error setting PRAGMA: %w", err)

	}

	usersQuery := `CREATE TABLE IF NOT EXISTS users (
				userid INTEGER PRIMARY KEY AUTOINCREMENT,
				username TEXT,
				UNIQUE (userid, username)
				);`

	_, err = db.Exec(usersQuery)

	if err != nil {
		return fmt.Errorf("error creating users structure: %w", err)
	}

	photosQuery := `CREATE TABLE IF NOT EXISTS photos (photoid INTEGER PRIMARY KEY AUTOINCREMENT,
					userid INTEGER,
					username TEXT,
					imageData BLOB,
					uploadDate DATETIME,
					likesCount INTEGER,
					commentsCount INTEGER,
					FOREIGN KEY(userid) REFERENCES user(userid) ON DELETE CASCADE);`
	_, err = db.Exec(photosQuery)
	if err != nil {
		return fmt.Errorf("error creating photos structure: %w", err)
	}

	followersQuery := `CREATE TABLE IF NOT EXISTS followers (
		userid INTEGER,
		followerid INTEGER,
		PRIMARY KEY (userid, followerid),
		FOREIGN KEY (userid) REFERENCES users(userid),
		FOREIGN KEY (followerid) REFERENCES users(userid)
		);`
	_, err = db.Exec(followersQuery)
	if err != nil {
		return fmt.Errorf("error creating followers table: %w", err)
	}

	followingQuery := `CREATE TABLE IF NOT EXISTS following (
            userid INTEGER,
            followingid INTEGER,
            PRIMARY KEY (userid, followingid),
            FOREIGN KEY (userid) REFERENCES users(userid),
            FOREIGN KEY (followingid) REFERENCES users(userid)
			);`

	_, err = db.Exec(followingQuery)
	if err != nil {
		return fmt.Errorf("error creating following table: %w", err)
	}

	bansQuery := `CREATE TABLE banned_users (
		userid INTEGER,
		banneduserid INTEGER,
		PRIMARY KEY (userid, banneduserid),
		FOREIGN KEY (userid) REFERENCES user(userid),
		FOREIGN KEY (bannedUserid) REFERENCES user(userid));`
	_, err = db.Exec((bansQuery))
	if err != nil {
		return fmt.Errorf("error creating ban structure: %w", err)
	}

	likesQuery := `CREATE TABLE likes (
		likeid INTEGER PRIMARY KEY AUTOINCREMENT,
		userid INTEGER,
		photoid INTEGER,
		FOREIGN KEY(userid) REFERENCES users(userid) ON DELETE CASCADE,
		FOREIGN KEY(photoid) REFERENCES photos(photoid) ON DELETE CASCADE
	);`
	_, err = db.Exec(likesQuery)
	if err != nil {
		return fmt.Errorf("error creating likes structure: %w", err)
	}
	commentsQuery := `CREATE TABLE IF NOT EXISTS comments (
		commentid INTEGER PRIMARY KEY AUTOINCREMENT,
		userid INTEGER,
		photoid INTEGER,
		upload_date DATETIME,
		commentText TEXT,
		FOREIGN KEY(userid) REFERENCES users(userid) ON DELETE CASCADE,
		FOREIGN KEY(photoid) REFERENCES photos(photoid) ON DELETE CASCADE
	);`

	_, err = db.Exec(commentsQuery)
	if err != nil {
		return fmt.Errorf("error creating comments structure: %w", err)
	}

	return nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
