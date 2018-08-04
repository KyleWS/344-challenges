package users

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type MemeStore struct {
	entries []*User
}

func NewMemeStore(sessionDuration time.Duration, purgeInterval time.Duration) *MemeStore {
	return &MemeStore{
		entries: []*User{},
	}
}

//GetByID returns the User with the given ID
func (m *MemeStore) GetByID(id bson.ObjectId) (*User, error) {
	for _, u := range m.entries {
		if u.ID == id {
			return u, nil
		}
	}
	return InvalidUser, ErrUserNotFound
}

//GetByEmail returns the User with the given email
func (m *MemeStore) GetByEmail(email string) (*User, error) {
	for _, u := range m.entries {
		if u.Email == email {
			return u, nil
		}
	}
	return InvalidUser, ErrUserNotFound
}

//GetByUserName returns the User with the given Username
func (m *MemeStore) GetByUserName(username string) (*User, error) {
	for _, u := range m.entries {
		if u.UserName == username {
			return u, nil
		}
	}
	return InvalidUser, ErrUserNotFound
}

//Insert converts the NewUser to a User, inserts
//it into the database, and returns it
func (m *MemeStore) Insert(newUser *NewUser) (*User, error) {
	u, err := newUser.ToUser()
	if err != nil {
		return InvalidUser, err
	}
	m.entries = append(m.entries, u)
	return u, nil
}

//Update applies UserUpdates to the given user ID
func (m *MemeStore) Update(userID bson.ObjectId, updates *Updates) error {
	u, err := m.GetByID(userID)
	if err != nil {
		return err
	}
	if err := u.ApplyUpdates(updates); err != nil {
		return err
	}
	return nil
}

//Delete deletes the user with the given ID
func (m *MemeStore) Delete(userID bson.ObjectId) error {
	for index, u := range m.entries {
		if u.ID == userID {
			m.entries = append(m.entries[:index], m.entries[index+1:]...)
			return nil
		}
	}
	return ErrUserNotFound
}
