package sessions

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	sessionID, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, ErrNoSessionID
	}
	store.Save(sessionID, sessionState)
	w.Header().Add(headerAuthorization, schemeBearer+sessionID.String())
	return sessionID, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	authHeader := r.Header.Get(headerAuthorization)
	if len(authHeader) == 0 {
		authHeader = r.URL.Query().Get("auth")
	}
	if len(authHeader) == 0 {
		return InvalidSessionID, ErrNoSessionID
	}
	if !strings.HasPrefix(authHeader, schemeBearer) {
		return InvalidSessionID, ErrInvalidScheme
	}
	authTokens := strings.Split(strings.Trim(authHeader, " "), " ")
	extractedID := authTokens[len(authTokens)-1]
	sessionID, err := ValidateID(extractedID, signingKey)
	if err != nil {
		//return the validation error.
		return InvalidSessionID, ErrInvalidID
	}
	return sessionID, nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	sessionID, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, ErrInvalidID
	}
	fmt.Println(sessionID.String())
	errorGet := store.Get(sessionID, &sessionState)
	if errorGet != nil {
		return InvalidSessionID, ErrStateNotFound
	}
	return sessionID, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	sessionID, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, ErrInvalidID
	}
	errorDelete := store.Delete(sessionID)
	if errorDelete != nil {
		return InvalidSessionID, ErrStateNotFound
	}
	return sessionID, nil
}
