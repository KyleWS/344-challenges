package users

import (
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var a *User = &User{
	Email:    "a@gmail.com",
	UserName: "aaa",
	ID:       bson.NewObjectId(),
}

var b *User = &User{
	Email:    "b@gmail.com",
	UserName: "bbb",
	ID:       bson.NewObjectId(),
}

var c *User = &User{
	Email:    "c@gmail.com",
	UserName: "ccc",
	ID:       bson.NewObjectId(),
}

func TestNewMemeStore(t *testing.T) {
	store := NewMemeStore(time.Hour, time.Minute)
	if store == nil {
		t.Errorf("error creating memestore, should not be nil")
	}
}

func TestGetByID(t *testing.T) {
	store := NewMemeStore(time.Hour, time.Minute)
	store.entries = []*User{a, b, c}
	u, err := store.GetByID(a.ID)
	if err != nil {
		t.Errorf("error when none was expected: %v", err)
	}
	if u != a {
		t.Errorf("error GetByID did not return correct user when valid ID provided")
	}
}

func TestGetByEmail(t *testing.T) {
	store := NewMemeStore(time.Hour, time.Minute)
	store.entries = []*User{a, b, c}
	u, err := store.GetByEmail(a.Email)
	if err != nil {
		t.Errorf("error when none was expected: %v", err)
	}
	if u != a {
		t.Errorf("error GetByEmail did not return correct user when valid email provided")
	}
}

func TestGetByUserName(t *testing.T) {
	store := NewMemeStore(time.Hour, time.Minute)
	store.entries = []*User{a, b, c}
	u, err := store.GetByUserName(a.UserName)
	if err != nil {
		t.Errorf("error when none was expected: %v", err)
	}
	if u != a {
		t.Errorf("error GetByUserName did not return correct user when valid username provided")
	}
}

func TestInsert(t *testing.T) {
	nu := &NewUser{
		Email:        "kyle@gmail.com",
		FirstName:    "kyle",
		LastName:     "ws",
		Password:     "123456",
		PasswordConf: "123456",
		UserName:     "xxcoolguy",
	}
	store := NewMemeStore(time.Hour, time.Minute)
	u, err := store.Insert(nu)
	if err != nil {
		t.Errorf("error inserting when none was expected: %v", err)
	}
	if _, err := store.GetByID(u.ID); err != nil {
		t.Errorf("error inserted user not found: %v", err)
	}
}

func TestUpdate(t *testing.T) {
	store := NewMemeStore(time.Hour, time.Minute)
	store.entries = []*User{a, b, c}
	up := &Updates{
		FirstName: "kyle",
		LastName:  "ws",
	}
	if err := store.Update(a.ID, up); err != nil {
		t.Errorf("error when none was expected: %v", err)
	}
	if a.FirstName != "kyle" {
		t.Errorf("error first name update failed name is %s should be %s", a.FirstName, "kyle")
	}
	if a.LastName != "ws" {
		t.Errorf("error last name update failed name is %s should be %s", a.LastName, "ws")
	}
	up = &Updates{
		FirstName: "",
		LastName:  "",
	}
	if err := store.Update(a.ID, up); err == nil {
		t.Errorf("error expected but none occured")
	}
}

func TestDelete(t *testing.T) {
	store := NewMemeStore(time.Hour, time.Minute)
	store.entries = []*User{a, b, c}
	if err := store.Delete(a.ID); err != nil {
		t.Errorf("error when none was expected: %v", err)
	}
	u, err := store.GetByID(a.ID)
	if err == nil && u != InvalidUser {
		t.Errorf("error still found user after deleting")
	}

}
