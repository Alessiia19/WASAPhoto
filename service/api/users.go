package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"io"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// setMyUserName aggiorna il nome utente dell'utente specificato.
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		// Se l'ID dell'utente non è un numero valido, restituisci un errore di formato.
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("setMyUserName: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	valid := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	var user User
	// Extract the username from the request body.
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		// If there is an error decoding the username, returns a Bad Request status.
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("setUsername: Invalid username format. Please follow the specified requirements.")
		return
	}

	// Aggiorna il nome utente nel database.
	if err := rt.db.UpdateUsername(userID, user.Username); err != nil {
		// Se c'è un errore nell'aggiornare l'username, restituisci un errore del server interno.
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("updateUsername: Error updating username in the database.")
		return
	}

	// Ritorna lo stato OK.
	w.WriteHeader(http.StatusOK)
}

// uploadPhoto carica una nuova foto sul profilo dell'utente.
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		// Se l'ID dell'utente non è un numero valido, restituisci un errore di formato.
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("uploadPhoto: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	valid := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	var photo Photo
	// Estrae l'immagine dal corpo della richiesta.
	if err := json.NewDecoder(r.Body).Decode(&photo); err != nil {
		// Se c'è un errore nella decodifica dell'immagine, restituisci uno stato di richiesta non valido.
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("uploadPhoto: Invalid JSON format.")
		return
	}

	// Leggi i dati dell'immagine dal corpo della richiesta.
	imageData, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("uploadPhoto: Error reading image data from request body.")
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(imageData))

	// Aggiorna i dati della foto
	photo.UserID = userID
	photo.UploadDate = time.Now()
	photo.LikesCount = 0
	photo.CommentsCount = 0

	// Create the photo in the database
	createdPhoto, err := rt.db.CreatePhoto(photo.PhotoToDatabase())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("uploadPhoto: Error creating photo in the database.")
		return
	}

	photo.PhotoFromDatabase(createdPhoto)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(photo)
}

// followUser aggiunge un utente alla lista di following dell'utente specificato.
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Estrai l'ID dell'utente che compie l'azione (follower) dal path.
	followerID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("followUser: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	valid := validateRequestingUser(strconv.Itoa(followerID), bearerToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Estrai l'ID dell'utente che si desidera seguire dalla richiesta.
	var followingUser User
	if err := json.NewDecoder(r.Body).Decode(&followingUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("followUser: Invalid request.")
		return
	}

	// Esegui l'operazione di follow nel database.
	if err := rt.db.FollowUser(followerID, followingUser.UserID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("followUser: Error following user in the database.")
		return
	}

	// Ritorna lo stato OK.
	w.WriteHeader(http.StatusOK)
}

// unfollowUser rimuove un utente dalla lista di following dell'utente specificato.
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Ottieni gli ID dal path
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("unfollowUser: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	valid := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	followingID, err := strconv.Atoi(ps.ByName("followingid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("unfollowUser: Invalid following ID format.")
		return
	}

	// Controllo se l'utente sta cercando di unfollow se stesso
	if userID == followingID {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("unfollowUser: Cannot unfollow yourself.")
		return
	}

	// Chiamare la funzione di database per eseguire l'operazione di unfollow
	if err := rt.db.UnfollowUser(userID, followingID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Se l'utente o il followingID non esiste, restituisci un errore 404.
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("unfollowUser: User or following user not found.")
			return
		}

		// Altri errori
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("unfollowUser: Error unfollowing user.")
		return
	}

	// Rispondi con uno status ok
	w.WriteHeader(http.StatusOK)
}

// banUser aggiunge un utente alla lista di utenti bloccati dall'utente specificato.
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Estrai l'ID dell'utente che compie l'azione dal path.
	UserID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("banUser: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	valid := validateRequestingUser(strconv.Itoa(UserID), bearerToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Estrai l'ID dell'utente che si desidera bloccare dalla richiesta.
	var bannedUser User
	if err := json.NewDecoder(r.Body).Decode(&bannedUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("banUser: Invalid request.")
		return
	}

	// Esegui l'operazione di ban nel database.
	if err := rt.db.BanUser(UserID, bannedUser.UserID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("banUser: Error banning user in the database.")
		return
	}

	// Ritorna lo stato OK.
	w.WriteHeader(http.StatusCreated)
}

// unfollowUser rimuove un utente dalla lista di following dell'utente specificato.
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Ottieni gli ID dal path
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("unbanUser: Invalid user ID format.")
		return
	}

	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	valid := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	bannedUserID, err := strconv.Atoi(ps.ByName("banneduserid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("unbanUser: Invalid banned user ID format.")
		return
	}

	// Controllo se l'utente sta cercando di unban se stesso
	if userID == bannedUserID {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error("unbanUser: Cannot unban yourself.")
		return
	}

	// Chiamare la funzione di database per eseguire l'operazione di unban
	if err := rt.db.UnbanUser(userID, bannedUserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Se l'utente o il bannedUserID non esiste, restituisci un errore 404.
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("unbanUser: User or banned user not found.")
			return
		}

		// Altri errori
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("unbanUser: Error unbanning user.")
		return
	}

	// Rispondi con uno status ok
	w.WriteHeader(http.StatusOK)
}

// getUserProfile ritorna il profilo dell'utente specificato
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Ottieni l'ID dell'user di cui si vuole visualizzare il profilo dal path
	requestedUserID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getUserProfile: Invalid user ID format.")
		return
	}

	// Ottieni l'ID dell'user che vuole effettuare l'operazione
	requestingUserID, err := strconv.Atoi(extractBearer(r.Header.Get("Authorization")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getUserProfile: error during authorization")
		return
	}

	// Chiamare la funzione del database per ottenere i dettagli del profilo dell'utente
	profile, err := rt.db.GetUserProfile(requestingUserID, requestedUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Se l'utente non esiste, restituisci un errore 404.
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("getUserProfile: Error getting user profile.")
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(profile)
}

// getMyStream ritorna lo stream dell'utente (foto dalle persone che l'utente segue)
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Ottieni l'ID dell'user dal path
	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("getMyStream: Invalid user ID format.")
		return
	}
	// Authorization
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	valid := validateRequestingUser(strconv.Itoa(userID), bearerToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Chiamare la funzione del database per ottenere lo stream dell'user
	stream, err := rt.db.GetMyStream(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Se l'utente non esiste, restituisci un errore 404.
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("getMyStream: Error getting user profile.")
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stream)
}
