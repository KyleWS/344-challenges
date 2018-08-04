package users

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Email     string        `json:"email"`
	PassHash  []byte        `json:"-"` //stored, but not encoded to clients
	UserName  string        `json:"userName"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	PhotoURL  string        `json:"photoURL"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

const ErrPassShort = "error password is too short"
const ErrPassMatch = "error passwords should match"
const ErrUsernameEmpty = "error username should not be empty"

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	if _, err := mail.ParseAddress(nu.Email); err != nil {
		return fmt.Errorf("error parsing email: %v", err)
	}
	//- Password must be at least 6 characters
	if len(nu.Password) < 6 {
		return fmt.Errorf(ErrPassShort)
	}
	if strings.Compare(nu.Password, nu.PasswordConf) != 0 {
		return fmt.Errorf(ErrPassMatch)
	}
	if len(nu.UserName) < 1 {
		return fmt.Errorf(ErrUsernameEmpty)
	}

	return nil
}

//ToUser converts the NewUser to a User, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {
	user := &User{
		Email:     nu.Email,
		FirstName: nu.FirstName,
		LastName:  nu.LastName,
		UserName:  nu.UserName,
	}
	md5Hasher := md5.New()
	email := strings.ToLower(strings.Trim(nu.Email, " "))
	io.WriteString(md5Hasher, email)
	gravatarURL := gravatarBasePhotoURL + string(md5Hasher.Sum(nil))
	user.PhotoURL = gravatarURL
	user.ID = bson.NewObjectId()
	err := user.SetPassword(nu.Password)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	return user, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put betweeen the names
func (u *User) FullName() string {
	if len(u.FirstName) == 0 || len(u.LastName) == 0 {
		return u.FirstName + u.LastName
	}
	return u.FirstName + " " + u.LastName
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	if len(password) == 0 {
		return fmt.Errorf("error password should not be empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return fmt.Errorf("error could not genereate hash for password")
	}
	u.PassHash = hash
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
	if err != nil {
		return fmt.Errorf("error password incorrect")
	}
	return nil
}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	if len(updates.FirstName) < 1 {
		return fmt.Errorf("error first name must not be empty")
	}
	if len(updates.LastName) < 1 {
		return fmt.Errorf("error last name must not be empty")
	}
	u.FirstName = updates.FirstName
	u.LastName = updates.LastName
	return nil
}
