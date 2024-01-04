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

// likePhoto aggiunge un like alla foto specificata.
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Estrai l'ID dell'utente che vuole mettere mi piace a una foto
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		// Se l'ID dell'utente non è un numero valido, restituisci un errore di formato.
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("likePhoto: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	valid := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Estrai l'ID della foto dal path.
	photoID, err := strconv.Atoi(ps.ByName("photoid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("likePhoto: Invalid photo ID format.")
		return
	}

	var like Like
	// Esegui l'operazione di like nel database.
	if err := rt.db.LikePhoto(userID, photoID, like.LikeToDatabase()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Se la foto non esiste, restituisci un errore 404.
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("likePhoto: Photo not found.")
			return
		}

		// Altri errori
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("likePhoto: Error liking photo.")
		return
	}

	// Ritorna lo stato OK.
	w.WriteHeader(http.StatusCreated)
}

// UnlikePhoto removes a like from a photo.
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Estrai l'ID dell'utente che vuole rimuovere il mi piace da una foto
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		// Se l'ID dell'utente non è un numero valido, restituisci un errore di formato.
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("unlikePhoto: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	valid := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Estrai l'ID della foto dal path.
	photoID, err := strconv.Atoi(ps.ByName("photoid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("unlikePhoto: Invalid photo ID format.")
		return
	}

	likeID, err := strconv.Atoi(ps.ByName("likeid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("unlikePhoto: Invalid like ID format.")
		return
	}

	// Esegui l'operazione di unlike nel database.
	if err := rt.db.UnlikePhoto(userID, photoID, likeID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Se uno tra gli id ricercati non esiste restituisci un errore 404.
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("unlikePhoto: Not found.")
			return
		}

		// Altri errori
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("unlikePhoto: Error removing like from a photo.")
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
}

// commentPhoto aggiunge un commento a una foto specificata.
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Estrai l'ID dell'utente che vuole commentare una foto
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		// Se l'ID dell'utente non è un numero valido, restituisci un errore di formato.
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("commentPhoto: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	valid := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Estrai l'ID della foto dal path.
	photoID, err := strconv.Atoi(ps.ByName("photoid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("commentPhoto: Invalid photo ID format.")
		return
	}

	var comment Comment
	// Extract the comment from the request body.
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		// If there is an error decoding the comment, returns a Bad Request status.
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("commentPhoto: Error decoding request body.")
		return
	}

	newComment, err := rt.db.CommentPhoto(userID, photoID, comment.CommentToDatabase())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("commentPhoto: Error commenting on photo.")
		return
	}

	// Update the user data with the information from the database.
	comment.CommentFromDatabase(newComment)

	// Ritorna lo stato OK.
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(comment)
}

// UncommentPhoto removes a comment from a photo.
func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Estrai l'ID dell'utente che vuole rimuovere il commento da una foto
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		// Se l'ID dell'utente non è un numero valido, restituisci un errore di formato.
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("uncommentPhoto: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	valid := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Estrai l'ID della foto dal path.
	photoID, err := strconv.Atoi(ps.ByName("photoid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("uncommentPhoto: Invalid photo ID format.")
		return
	}

	commentID, err := strconv.Atoi(ps.ByName("commentid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("uncommentPhoto: Invalid like ID format.")
		return
	}

	// Esegui l'operazione di uncomment nel database.
	if err := rt.db.UncommentPhoto(userID, photoID, commentID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Se uno tra gli id ricercati non esiste restituisci un errore 404.
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("uncommentPhoto: Not found.")
			return
		}

		// Altri errori
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("uncommentPhoto: Error removing like from a photo.")
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
}

// deletePhoto removes a photo.
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Estrai l'ID dell'utente che vuole rimuovere la foto
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		// Se l'ID dell'utente non è un numero valido, restituisci un errore di formato.
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("deletePhoto: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	valid := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Estrai l'ID della foto dal path.
	photoID, err := strconv.Atoi(ps.ByName("photoid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("deletePhoto: Invalid photo ID format.")
		return
	}

	// Rimuovi la foto dal database.
	if err := rt.db.DeletePhoto(userID, photoID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Se la foto non esiste, restituisci un errore 404.
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("deletePhoto: Photo not found.")
			return
		}

		// Altri errori
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("deletePhoto: Error removing photo.")
		return
	}

	// Respond with success message
	w.WriteHeader(http.StatusOK)
}
