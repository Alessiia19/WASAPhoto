package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// UpdateUsername updates the username of the specified user in the database.
func (db *appdbimpl) UpdateUsername(userID int, newUsername string) error {
	// Converts the provided username to lowercase for comparison.
	newUsernameLower := strings.ToLower(newUsername)

	// Check if the username is already in use by another user.
	var existingUserID int
	err := db.c.QueryRow("SELECT userid FROM users WHERE LOWER(username) = ?", newUsernameLower).Scan(&existingUserID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		// Return an error if there is an issue during the username search.
		return fmt.Errorf("error checking existing username: %w", err)
	}
	// Check if the found user ID does not match the current user's ID, indicating the username is already taken by someone else.
	if existingUserID != 0 && existingUserID != userID {
		return fmt.Errorf("username %s already in use by another user", newUsername)
	}

	// Update the username in the database
	_, err = db.c.Exec("UPDATE users SET username = ? WHERE userid = ?", newUsername, userID)
	if err != nil {
		return fmt.Errorf("error updating username in database: %w", err)
	}

	return nil
}

// GetUserProfile retrieves the details of the specified user's profile,
// including uploaded photos, the number of followers, the number of users followed,
// and the total number of photos uploaded.
func (db *appdbimpl) GetUserProfile(requestingUserID, requestedUserID int) (Profile, error) {
	var profile Profile

	// Check if the user to be searched exists.
	var existingUser int
	err := db.c.QueryRow("SELECT 1 FROM users WHERE userid = ?", requestedUserID).Scan(&existingUser)
	if errors.Is(err, sql.ErrNoRows) {
		return profile, fmt.Errorf("the user you are searching doesn't exist")
	} else if err != nil {
		return profile, fmt.Errorf("error checking existing user: %w", err)
	}

	// Check if the user is banned by the user being searched.
	var isBanned int
	err = db.c.QueryRow("SELECT 1 FROM banned_users WHERE userid = ? AND banneduserid = ?", requestedUserID, requestingUserID).Scan(&isBanned)
	if err == nil {
		// The user is banned by the searched user.
		return profile, errors.New("you are banned by this user")
	}

	// Retrieve user details (userid, username)
	user, err := db.GetUserDetails(requestedUserID)
	if err != nil {
		return profile, err
	}

	// Retrieve the list of followers.
	followers, err := db.GetFollowers(requestedUserID)
	if err != nil {
		return profile, err
	}

	// Retrieve the list of users followed.
	following, err := db.GetFollowing(requestedUserID)
	if err != nil {
		return profile, err
	}

	// Retrieve the list of photos uploaded by the user.
	uploadedPhotos, err := db.GetUploadedPhotos(requestedUserID)
	if err != nil {
		return profile, err
	}

	// Construct the user's profile with the gathered data.
	profile = Profile{
		UserID:              user.UserID,
		Username:            user.Username,
		Followers:           followers,
		Following:           following,
		FollowersCount:      len(followers),
		FollowingCount:      len(following),
		UploadedPhotos:      uploadedPhotos,
		UploadedPhotosCount: len(uploadedPhotos),
	}
	return profile, err
}

// FollowUser adds a user to the specified user's following list.
func (db *appdbimpl) FollowUser(userID, userIDToFollow int) error {

	// Check if the user being followed exists.
	var existingUser int
	err := db.c.QueryRow("SELECT 1 FROM users WHERE userid = ?", userIDToFollow).Scan(&existingUser)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("the user you want to follow doesn't exists")
	} else if err != nil {
		return fmt.Errorf("error checking existing user: %w", err)
	}

	// Check if the user is banned by the user they are trying to follow.
	var isBanned int
	err = db.c.QueryRow("SELECT 1 FROM banned_users WHERE userid = ? AND banneduserid = ?", userIDToFollow, userID).Scan(&isBanned)
	if err == nil {
		// The user is banned by the other user.
		return errors.New("you are banned by this user")
	}

	// Check if the user has banned the other user.
	var hasBanned int
	err = db.c.QueryRow("SELECT 1 FROM banned_users WHERE userid = ? AND banneduserid = ?", userID, userIDToFollow).Scan(&hasBanned)
	if err == nil {
		// The user has banned the other user.
		return errors.New("you have banned this user")
	}

	// Check if the user already follows the other user.
	var existingFollower int
	err = db.c.QueryRow("SELECT 1 FROM followers WHERE userid = ? AND followerid = ?", userIDToFollow, userID).Scan(&existingFollower)
	if err == nil {
		// The user already follows the other user.
		return errors.New("already followed")
	}

	// Check if the user is trying to follow themselves.
	if userID == userIDToFollow {
		return errors.New("cannot follow yourself")
	}

	// Update the followers table by adding userID as a follower of userIDToFollow.
	_, err = db.c.Exec("INSERT INTO followers (userID, followerID) VALUES (?, ?)", userIDToFollow, userID)
	if err != nil {
		return fmt.Errorf("error updating followers table: %w", err)
	}

	// Update the following table by adding userIDToFollow as a following of userID.
	_, err = db.c.Exec("INSERT INTO following (userID, followingID) VALUES (?, ?)", userID, userIDToFollow)
	if err != nil {
		return fmt.Errorf("error updating following table: %w", err)
	}

	return nil
}

// UnfollowUser removes a user from the specified user's following list.
func (db *appdbimpl) UnfollowUser(userID, followingID int) error {

	// Check if the user is attempting to unfollow someone they are not currently following.
	var existingFollower int
	err := db.c.QueryRow("SELECT 1 FROM followers WHERE userid = ? AND followerid = ?", followingID, userID).Scan(&existingFollower)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("you are trying to unfollow someone you don't follow")
	} else if err != nil {
		return fmt.Errorf("error checking existing follower: %w", err)
	}

	// Remove the user from the followers table.
	_, err = db.c.Exec("DELETE FROM followers WHERE userid = ? AND followerid = ?", followingID, userID)
	if err != nil {
		return fmt.Errorf("error removing follower: %w", err)
	}

	// Remove the user from the following table.
	_, err = db.c.Exec("DELETE FROM following WHERE userid = ? AND followingid = ?", userID, followingID)
	if err != nil {
		return fmt.Errorf("error removing following: %w", err)
	}

	return nil
}

// BanUser adds a user to the specified user's banned list.
func (db *appdbimpl) BanUser(userID, bannedUserID int) error {

	// Verify if the user to be banned exists.
	var existingUser int
	err := db.c.QueryRow("SELECT 1 FROM users WHERE userid = ?", bannedUserID).Scan(&existingUser)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("the user you want to ban doesn't exists")
	} else if err != nil {
		return fmt.Errorf("error checking existing user: %w", err)
	}

	// Check if the user is trying to ban themselves.
	if userID == bannedUserID {
		return errors.New("cannot ban yourself")
	}

	// Check if the user is banned by the user they are trying to ban.
	var isBanned int
	err = db.c.QueryRow("SELECT 1 FROM banned_users WHERE userid = ? AND banneduserid = ?", bannedUserID, userID).Scan(&isBanned)
	if err == nil {
		// The user is banned by the other user.
		return errors.New("you are banned by this user")
	}

	// Check if the user has already banned the other user.
	var hasBanned int
	err = db.c.QueryRow("SELECT 1 FROM banned_users WHERE userid = ? AND banneduserid = ?", userID, bannedUserID).Scan(&hasBanned)
	if err == nil {
		// User already banned.
		return errors.New("you have banned this user")
	}

	// Handle the removal from followers and following if necessary.
	// The banned user will automatically stop being a follower and following of the user who bans them.
	var follows int
	err = db.c.QueryRow("SELECT 1 FROM followers WHERE userid = ? AND followerid = ?", userID, bannedUserID).Scan(&follows)
	if err == nil {
		// Remove the banned user from followers if they are following the user executing the ban.
		if err := db.UnfollowUser(bannedUserID, userID); err != nil {
			return fmt.Errorf("error handling unfollow during ban operation: %w", err)
		}
	}

	// Check if the user is following the user to be banned.
	var followed int
	err = db.c.QueryRow("SELECT 1 FROM following WHERE userid = ? AND followingid = ?", userID, bannedUserID).Scan(&followed)
	if err == nil {
		// Remove the initiating user from following the user to be banned.
		if err := db.UnfollowUser(userID, bannedUserID); err != nil {
			return fmt.Errorf("error handling unfollow during ban operation: %w", err)
		}
	}

	// Update the banned_users table.
	_, err = db.c.Exec("INSERT INTO banned_users (userid, banneduserid) VALUES (?, ?)", userID, bannedUserID)
	if err != nil {
		return fmt.Errorf("error updating banned_users table: %w", err)
	}

	return nil
}

// UnbanUser removes a user from the specified user's banned list.
func (db *appdbimpl) UnbanUser(userID, bannedUserID int) error {
	// Check if the user is trying to unban someone who is not currently banned or does not exist.
	var existingUser int
	err := db.c.QueryRow("SELECT 1 FROM banned_users WHERE userid = ? AND banneduserid = ?", userID, bannedUserID).Scan(&existingUser)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("you are trying to unban someone who was not banned or doesn't exists")
	} else if err != nil {
		return fmt.Errorf("error checking existing ban: %w", err)
	}

	// Remove the user from the banned_users table.
	_, err = db.c.Exec("DELETE FROM banned_users WHERE userid = ? AND banneduserid = ?", userID, bannedUserID)
	if err != nil {
		return fmt.Errorf("error removing ban: %w", err)
	}

	return nil
}

// GetMyStream returns the stream of the user, consisting of photos posted by their following,
// including details such as the date-time they were posted, the number of likes, and comments.
func (db *appdbimpl) GetMyStream(userID int) ([]CompletePhoto, error) {
	// Retrieve the list of users that the specified user is following.
	following, err := db.GetFollowing(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting following users: %w", err)
	}

	var stream []CompletePhoto

	// Iterate over each followed user to obtain photos from their stream.
	for _, followedUser := range following {
		// Retrieve the list of photos uploaded by the followed user.
		uploadedPhotos, err := db.GetUploadedPhotos(followedUser.UserID)
		if err != nil {
			return nil, fmt.Errorf("error getting uploaded photos for user %d: %w", followedUser.UserID, err)
		}

		// Add the photos to the stream.
		stream = append(stream, uploadedPhotos...)
	}

	// Sort the stream in reverse chronological order.
	sort.Slice(stream, func(i, j int) bool {
		return stream[i].UploadDate.After(stream[j].UploadDate)
	})

	return stream, nil
}

// GetUsers searches for users by username that starts with the specified substring.
func (db *appdbimpl) GetUsers(userID int, usernameSubstring string) ([]User, error) {
	var users []User

	// Define SQL query to find users whose usernames contain the specified substring
	// and who are not banned by the user making the request.
	query := "SELECT userid, username FROM users WHERE username LIKE ? AND userid NOT IN (SELECT userid FROM banned_users WHERE banneduserid = ?)"
	rows, err := db.c.Query(query, usernameSubstring+"%", userID)
	if err != nil {
		return nil, fmt.Errorf("error querying users by username substring: %w", err)
	}
	defer rows.Close() // Ensure that the rows are closed after operations are complete.

	// Iterate over the query results and extract user data.
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserID, &user.Username); err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}

// GetBanStatus checks if one user has been banned by another user.
func (db *appdbimpl) GetBanStatus(userID, userToCheckID int) (bool, error) {
	var hasBanned int
	// Check if the specified user (userID) has banned the other user (userToCheckID).
	err := db.c.QueryRow("SELECT 1 FROM banned_users WHERE userid = ? AND banneduserid = ?", userID, userToCheckID).Scan(&hasBanned)
	if err == nil {
		// The user has banned the other user.
		return true, nil
	} else if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return false, err
}
