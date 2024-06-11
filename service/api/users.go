package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// setMyUserName updates the username of the specified user.
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		// If the user ID is not valid, returns an error.
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("setMyUserName: Invalid user ID")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	authorizationStatus := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if authorizationStatus != http.StatusOK {
		w.WriteHeader(authorizationStatus)
		return
	}

	var user User
	// Extract the username from the request body.
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("setMyUserName: Invalid username.")
		return
	}

	// Check if the username meets the requirements.
	if !isValidUsername(user.Username) {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("setMyUserName: Invalid username format. Please follow the specified requirements.")
		return
	}

	// Update the username in the database.
	if err := rt.db.UpdateUsername(userID, user.Username); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("setMyUserName: Error updating username in the database.")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// getUserProfile returns the profile of the specified user.
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the ID of the user whose profile is to be viewed from the path.
	requestedUserID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getUserProfile: Invalid user ID format.")
		return
	}

	// Extract the ID of the user making the request.
	requestingUserID, err := strconv.Atoi(extractBearer(r.Header.Get("Authorization")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getUserProfile: error during authorization")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	authorizationStatus := validateRequestingUser(strconv.Itoa(requestingUserID), bearerToken)
	if authorizationStatus != http.StatusOK {
		w.WriteHeader(authorizationStatus)
		return
	}

	// Call the database function to get the user profile details.
	profile, err := rt.db.GetUserProfile(requestingUserID, requestedUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return a 404 error if the user does not exist.
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getUserProfile: Error getting user profile.")
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(profile)
}

// followUser adds a user to the specified user's following list.
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the ID of the user making the request.
	followerID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("followUser: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	authorizationStatus := validateRequestingUser(strconv.Itoa(followerID), bearerToken)
	if authorizationStatus != http.StatusOK {
		w.WriteHeader(authorizationStatus)
		return
	}

	// Extract the user ID of the user to be followed from the request body.
	var followingUser User
	if err := json.NewDecoder(r.Body).Decode(&followingUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("followUser: Invalid request.")
		return
	}

	// Follow the user.
	if err := rt.db.FollowUser(followerID, followingUser.UserID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("followUser: Error following user in the database.")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// unfollowUser removes a user from the specified user's following list.
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the ID of the user making the request.
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("unfollowUser: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	authorizationStatus := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if authorizationStatus != http.StatusOK {
		w.WriteHeader(authorizationStatus)
		return
	}

	// Extract the user ID of the user to be unfollowed.
	followingID, err := strconv.Atoi(ps.ByName("followingid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("unfollowUser: Invalid following ID format.")
		return
	}

	// Check if the user is trying to unfollow themselves.
	if userID == followingID {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("unfollowUser: Cannot unfollow yourself.")
		return
	}

	// Unfollow the user.
	if err := rt.db.UnfollowUser(userID, followingID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("unfollowUser: User or following user not found.")
			return
		}

		// Other errors.
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("unfollowUser: Error unfollowing user.")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// banUser adds a user to the specified user's banned list.
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the ID of the user making the request.
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("banUser: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	authorizationStatus := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if authorizationStatus != http.StatusOK {
		w.WriteHeader(authorizationStatus)
		return
	}

	// Extract the user ID of the user to be banned from the request body.
	var bannedUser User
	if err := json.NewDecoder(r.Body).Decode(&bannedUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("banUser: Invalid request.")
		return
	}

	// Ban the user
	if err := rt.db.BanUser(userID, bannedUser.UserID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("banUser: Error banning user in the database.")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// unBanUser removes a user from the specified user's banned list.
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the ID of the user making the request.
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("unbanUser: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	authorizationStatus := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if authorizationStatus != http.StatusOK {
		w.WriteHeader(authorizationStatus)
		return
	}

	// Extract the user ID of the banned user.
	bannedUserID, err := strconv.Atoi(ps.ByName("banneduserid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("unbanUser: Invalid banned user ID format.")
		return
	}

	// Check if the user is trying to unban themselves.
	if userID == bannedUserID {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("unbanUser: Cannot unban yourself.")
		return
	}

	// Unban the user.
	if err := rt.db.UnbanUser(userID, bannedUserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return a NotFound status if either the user or the banned user does not exist.
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("unbanUser: User or banned user not found.")
			return
		}

		// Other errors.
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("unbanUser: Error unbanning user.")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// getMyStream returns the stream of the user, consisting of photos from people the user follows.
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the ID of the user making the request.
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getMyStream: Invalid user ID format.")
		return
	}
	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	authorizationStatus := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if authorizationStatus != http.StatusOK {
		w.WriteHeader(authorizationStatus)
		return
	}

	// Call the database function to get the user's stream.
	stream, err := rt.db.GetMyStream(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return a NotFound status if the user does not exist.
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getMyStream: Error getting user profile.")
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stream)
}

// getUsers searches for users whose usernames starts with a specified substring.
func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the ID of the user making the request from the authorization.
	userIDstr := extractBearer(r.Header.Get("Authorization"))

	// Convert the user ID from string to integer.
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getUsers: Invalid user ID")
		return
	}

	// Retrieve the search substring from the URL query parameter "username".
	query := r.URL.Query().Get("username")
	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("getUsers: Query 'username' is missing")
		return
	}

	// Call the database function to get users matching the query substring.
	users, err := rt.db.GetUsers(userID, query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getUsers: Error fetching users from database")
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}

// getBanStatus checks if a user has been banned by the currently logged-in user.
func (rt *_router) getBanStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the ID of the user making the request.
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getBanStatus: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	authorizationStatus := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if authorizationStatus != http.StatusOK {
		w.WriteHeader(authorizationStatus)
		return
	}

	// Extract the ID of the user to check.
	userToCheckID, err := strconv.Atoi(ps.ByName("banneduserid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getBanStatus: Invalid user to check ID format.")
		return
	}

	// Call the database function to get the ban status.
	isBanned, err := rt.db.GetBanStatus(userID, userToCheckID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getBanStatus: Error while checking ban status")
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(isBanned)
}
