package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/info344-a17/challenges-KyleIWS/servers/gateway/models/users"
	"github.com/info344-a17/challenges-KyleIWS/servers/gateway/sessions"
)

//TODO: define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.

func (us *Ctx) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		sess := SessionState{}
		_, err := sessions.GetState(r, us.Key, us.SessionsStore, &sess)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not get session state %v", err), http.StatusUnauthorized)
		}
		prefix := r.URL.Query().Get("q")
		if len(prefix) > 0 {
			prefix = strings.ToLower(prefix)
			userObjects, err := us.RootTrieNode.GetN(prefix, 20)
			if err != nil {
				http.Error(w, fmt.Sprintf("error fetching users: %v", err), http.StatusBadRequest)
				return
			}
			users := us.UsersStore.ObjectIdsToUsers(userObjects)
			if err := json.NewEncoder(w).Encode(users); err != nil {
				http.Error(w, fmt.Sprintf("error returning new user json: %v", err), http.StatusInternalServerError)
				return
			}
		} else {
			if err := json.NewEncoder(w).Encode(make(map[string]string)); err != nil {
				http.Error(w, fmt.Sprintf("error returning new user json: %v", err), http.StatusInternalServerError)
				return
			}
		}
	case "POST":
		nu := &users.NewUser{}
		if err := json.NewDecoder(r.Body).Decode(nu); err != nil {
			http.Error(w, fmt.Sprintf("error decoding received json: %v", err), http.StatusBadRequest)
			return
		}
		// validate new user
		if err := nu.Validate(); err != nil {
			http.Error(w, fmt.Sprintf("error validating user: %v", err), http.StatusBadRequest)
			return
		}
		// make sure email isnt in usersotre already
		if _, err := us.UsersStore.GetByEmail(nu.Email); err == nil {
			http.Error(w, fmt.Sprintf("error email already exists"), http.StatusBadRequest)
			return
		}
		// make sure userstore doesnt have username already
		if _, err := us.UsersStore.GetByUserName(nu.UserName); err == nil {
			http.Error(w, fmt.Sprintf("error username already exists"), http.StatusBadRequest)
			return
		}
		// insert new user (they become user)
		user, err := us.UsersStore.Insert(nu)
		if err != nil {
			http.Error(w, fmt.Sprintf("error inserting user to database: %v", err), http.StatusInternalServerError)
			return
		}
		// begin a new session
		newSession := SessionState{
			AuthenticatedUser: user,
			TimeBegin:         time.Now(),
		}

		_, errBeginSession := sessions.BeginSession(us.Key, us.SessionsStore, &newSession, w)
		if errBeginSession != nil {
			http.Error(w, fmt.Sprintf("error generating session for user: %v", errBeginSession), http.StatusInternalServerError)
			return
		}
		// return with a http.StatusCreated and json encoded form of that created user
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, fmt.Sprintf("error returning new user json: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		// Somehow put the user/sessionn I have started into a request context?
		// Add this new user to the Trie
		us.RootTrieNode.Add(strings.ToLower(user.Email), user.ID)
		us.RootTrieNode.Add(strings.ToLower(user.UserName), user.ID)
		us.RootTrieNode.Add(strings.ToLower(user.FirstName), user.ID)
		us.RootTrieNode.Add(strings.ToLower(user.LastName), user.ID)
	default:
		http.Error(w, fmt.Sprintf("only accepts POST"), http.StatusMethodNotAllowed)
	}
}

func (us *Ctx) UsersMeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		sess := SessionState{}
		_, err := sessions.GetState(r, us.Key, us.SessionsStore, &sess)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not get session state %v", err), http.StatusUnauthorized)
		}
		user := sess.AuthenticatedUser
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, fmt.Sprintf("error returning user json: %v", err), http.StatusInternalServerError)
			return
		}
	case "PATCH":
		// Get the user from the request body
		sess := SessionState{}
		sessID, err := sessions.GetState(r, us.Key, us.SessionsStore, &sess)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not get session state %v", err), http.StatusBadRequest)
		}
		user := sess.AuthenticatedUser
		oldFirst := user.FirstName
		oldLast := user.LastName
		///////////
		// Get the updates from the request body as well
		upd := &users.Updates{}
		if err := json.NewDecoder(r.Body).Decode(upd); err != nil {
			http.Error(w, fmt.Sprintf("error decoding received json: %v", err), http.StatusBadRequest)
			return
		}
		// make sure the updates are valid
		if err := user.ApplyUpdates(upd); err != nil {
			http.Error(w, fmt.Sprintf("error applying updated updates: %v", err), http.StatusBadRequest)
			return
		}
		// update the user in our provided database
		if err := us.UsersStore.Update(user.ID, upd); err != nil {
			http.Error(w, fmt.Sprintf("could not update user: %v", err), http.StatusInternalServerError)
			return
		}
		us.SessionsStore.Save(sessID, &sess)
		///////////
		// Return updated user as json
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, fmt.Sprintf("error returning user json: %v", err), http.StatusInternalServerError)
			return
		}

		us.RootTrieNode.Delete(strings.ToLower(oldFirst), user.ID)
		us.RootTrieNode.Delete(strings.ToLower(oldLast), user.ID)
		us.RootTrieNode.Add(strings.ToLower(user.FirstName), user.ID)
		us.RootTrieNode.Add(strings.ToLower(user.LastName), user.ID)
	default:
		http.Error(w, fmt.Sprintf("only accepts GET and PATCH"), http.StatusMethodNotAllowed)
	}
}

func (sess *Ctx) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// decode given user
		cred := &users.Credentials{}
		if err := json.NewDecoder(r.Body).Decode(cred); err != nil {
			http.Error(w, fmt.Sprintf("error decoding received json: %v", err), http.StatusBadRequest)
			return
		}
		// get the user with the provided email
		u, err := sess.UsersStore.GetByEmail(cred.Email)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid credentials"), http.StatusUnauthorized)
			return
		}
		// authenticate with ther password
		if err := u.Authenticate(cred.Password); err != nil {
			http.Error(w, fmt.Sprintf("invalid credentials"), http.StatusUnauthorized)
			return
		}
		newSession := SessionState{
			AuthenticatedUser: u,
			TimeBegin:         time.Now(),
		}
		// begin a session
		if _, err := sessions.BeginSession(sess.Key, sess.SessionsStore, &newSession, w); err != nil {
			http.Error(w, fmt.Sprintf("error starting new session: %v", err), http.StatusInternalServerError)
			return
		}
		// return the user
		if err := json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, fmt.Sprintf("error returning new user json: %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("only accepts POST"), http.StatusMethodNotAllowed)
	}
}

func (sess *Ctx) SessionsMineHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		session := SessionState{}
		_, err := sessions.GetState(r, sess.Key, sess.SessionsStore, session)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not get session state %v", err), http.StatusUnauthorized)
		}
		if _, err := sessions.EndSession(r, sess.Key, sess.SessionsStore); err != nil {
			http.Error(w, fmt.Sprintf("error deleting state %v", err), http.StatusInternalServerError)
		}
	default:
		http.Error(w, fmt.Sprintf("only accepts DELETE"), http.StatusMethodNotAllowed)
	}
}
