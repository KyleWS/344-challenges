package handlers

import (
	"time"

	"github.com/info344-a17/challenges-KyleIWS/servers/gateway/models/users"
)

//TODO: define a session state struct for this web server
//see the assignment description for the fields you should include
//remember that other packages can only see exported fields!

type SessionState struct {
	TimeBegin         time.Time
	AuthenticatedUser *users.User
}
