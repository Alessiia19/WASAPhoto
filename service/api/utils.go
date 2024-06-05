package api

import (
	"bytes"
	"image"
	"net/http"
	"regexp"
	"strings"
)

// --- AUTHENTICATION FUNCTIONS ---

// validateRequestingUser checks if the user is authorized.
func validateRequestingUser(userID, bearerToken string) int {
	if !isUserLoggedIn(bearerToken) {
		// The user is not authenticated.
		return http.StatusForbidden
	}

	if userID != bearerToken {
		// The user is not authorized.
		return http.StatusUnauthorized
	}

	// The user is authorized.
	return http.StatusOK
}

// extractBearer extracts the Bearer token from an authentication string.
func extractBearer(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

// isUserLoggedIn checks if the user is logged in.
// It returns true if the user is logged in, otherwise false.
func isUserLoggedIn(bearerToken string) bool {
	return bearerToken != ""
}

// --- USERNAME VALIDATION ---

// isValidUsername checks if the username meets the requirements defined in the OpenAPI specification.
func isValidUsername(username string) bool {
	// Check if the length of the username is within the specified range.
	if len(username) < 3 || len(username) > 16 {
		return false
	}

	// Check if the username contains only alphanumeric characters.
	match, _ := regexp.MatchString("^[a-zA-Z0-9]+$", username)
	return match
}

// --- PHOTO FORMAT VALIDATION ---

// CheckImageType verifica che il contenuto sia di tipo PNG o JPG.
func CheckImageType(data []byte) bool {
	contentType := http.DetectContentType(data)
	switch contentType {
	case "image/jpeg", "image/png":
		return true
	default:
		return false
	}
}

// ValidateImage verifica che il file sia effettivamente un'immagine PNG o JPG.
func ValidateImage(data []byte) bool {
	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return false
	}

	return format == "jpeg" || format == "png"
}
