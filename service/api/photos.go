package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// uploadPhoto uploads a new photo to the user's profile.
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the user ID from the path parameters.
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("uploadPhoto: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	authorizationStatus := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if authorizationStatus != http.StatusOK {
		w.WriteHeader(authorizationStatus)
		return
	}

	var photo Photo
	// Read the photo data from the request body.
	photo.ImageData, err = io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("uploadPhoto: error reading body content")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Image type checks (only png or jpeg is accepted).
	if !CheckImageType(photo.ImageData) {
		ctx.Logger.Error("uploadPhoto: unsupported image format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the username of the user.
	user, err := rt.db.GetUserDetails(userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("uploadPhoto: error retrieving user details")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update the photo data.
	photo.UserID = userID
	photo.Username = user.Username
	photo.UploadDate = time.Now()
	photo.LikesCount = 0
	photo.CommentsCount = 0

	// Create the photo in the database
	createdPhoto, err := rt.db.CreatePhoto(photo.PhotoToDatabase())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("uploadPhoto: Error creating photo in the database.")
		return
	}

	photo.PhotoFromDatabase(createdPhoto)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(photo)
}

// likePhoto adds a like to the specified photo.
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the user ID from the path parameters.
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("likePhoto: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	authorizationStatus := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if authorizationStatus != http.StatusOK {
		w.WriteHeader(authorizationStatus)
		return
	}

	// Extract the photo ID from the path parameters.
	photoID, err := strconv.Atoi(ps.ByName("photoid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("likePhoto: Invalid photo ID format.")
		return
	}

	var like Like
	// Like the photo
	if err := rt.db.LikePhoto(userID, photoID, like.LikeToDatabase()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// The photo does not exist, return a NotFound status.
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("likePhoto: Photo not found.")
			return
		}

		// Other errors
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("likePhoto: Error liking photo.")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UnlikePhoto removes a like from a photo.
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the user ID from the path parameters.
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("unlikePhoto: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	authorizationStatus := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if authorizationStatus != http.StatusOK {
		w.WriteHeader(authorizationStatus)
		return
	}

	// Extract the photo ID from the path parameters.
	photoID, err := strconv.Atoi(ps.ByName("photoid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("unlikePhoto: Invalid photo ID format.")
		return
	}

	// Extract the like ID from the path parameters.
	likeID, err := strconv.Atoi(ps.ByName("likeid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("unlikePhoto: Invalid like ID format.")
		return
	}

	// Unlike photo.
	if err := rt.db.UnlikePhoto(userID, photoID, likeID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If any of the input IDs do not exist, return a NotFound status.
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("unlikePhoto: Not found.")
			return
		}

		// Other errors.
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("unlikePhoto: Error removing like from a photo.")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// commentPhoto adds a comment to the specified photo.
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the user ID from the path parameters.
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("commentPhoto: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	authorizationStatus := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if authorizationStatus != http.StatusOK {
		w.WriteHeader(authorizationStatus)
		return
	}

	// Extract the photo ID from the path parameters.
	photoID, err := strconv.Atoi(ps.ByName("photoid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("commentPhoto: Invalid photo ID format.")
		return
	}

	var comment Comment
	// Extract the comment from the request body.
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("commentPhoto: Error decoding request body.")
		return
	}

	// Get the username of the user.
	user, err := rt.db.GetUserDetails(userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("commentPhoto: error retrieving user details")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comment.UploadDate = time.Now()

	// Comment photo
	newComment, err := rt.db.CommentPhoto(userID, photoID, user.Username, comment.CommentToDatabase())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("commentPhoto: Error commenting on photo.")
		return
	}

	// Update the user data with the information from the database.
	comment.CommentFromDatabase(newComment)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(comment)
}

// UncommentPhoto removes a comment from a photo.
func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the user ID from the path parameters.
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("uncommentPhoto: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	authorizationStatus := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if authorizationStatus != http.StatusOK {
		w.WriteHeader(authorizationStatus)
		return
	}

	// Extract the photo ID from the path parameters.
	photoID, err := strconv.Atoi(ps.ByName("photoid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("uncommentPhoto: Invalid photo ID format.")
		return
	}

	// Extract the comment ID from the path parameters.
	commentID, err := strconv.Atoi(ps.ByName("commentid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("uncommentPhoto: Invalid like ID format.")
		return
	}

	// Uncomment photo.
	if err := rt.db.UncommentPhoto(userID, photoID, commentID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If any of the input IDs do not exist, return a NotFound status.
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("uncommentPhoto: Not found.")
			return
		}

		// Altri Other errors.
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("uncommentPhoto: Error removing like from a photo.")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// deletePhoto removes a photo.
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the user ID from the path parameters.
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("deletePhoto: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	authorizationStatus := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if authorizationStatus != http.StatusOK {
		w.WriteHeader(authorizationStatus)
		return
	}

	// Extract the photo ID from the path parameters.
	photoID, err := strconv.Atoi(ps.ByName("photoid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("deletePhoto: Invalid photo ID format.")
		return
	}

	// Remove the photo from database.
	if err := rt.db.DeletePhoto(userID, photoID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// The photo does not exist, return a NotFound status.
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("deletePhoto: Photo not found.")
			return
		}

		// Other errors.
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("deletePhoto: Error removing photo.")
		return
	}

	w.WriteHeader(http.StatusOK)
}
