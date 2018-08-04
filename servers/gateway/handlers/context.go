package handlers

import (
	"github.com/info344-a17/challenges-KyleIWS/servers/gateway/indexes"
	"github.com/info344-a17/challenges-KyleIWS/servers/gateway/models/users"
	"github.com/info344-a17/challenges-KyleIWS/servers/gateway/sessions"
)

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store

type Ctx struct {
	Key           string
	SessionsStore sessions.Store
	UsersStore    users.Store
	RootTrieNode  *indexes.TrieNode
}
