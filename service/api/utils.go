package api

import (
	"net/http"
	"strings"
)

// --- AUTHENTICATION FUNCTIONS ---

// funzione di validazione dell'utente
func validateRequestingUser(identifier string, bearerToken string) int {
	// Controlla se l'utente è registrato, altrimenti ritorna errore
	if isNotLogged(bearerToken) {
		return http.StatusForbidden
	}

	// Controlla se l'user è autorizzato
	if identifier != bearerToken {
		return http.StatusUnauthorized
	}

	return 0
}

// estrae il token bearer
func extractBearer(auth string) string {
	var tokens = strings.Split(auth, " ")
	if len(tokens) == 2 {
		return strings.Trim(tokens[1], " ")
	}
	return ""
}

// Funzione che verifica se l'user è registrato
func isNotLogged(bearerToken string) bool {
	return bearerToken == ""
}
