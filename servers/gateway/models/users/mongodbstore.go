package users

import (
	"fmt"
	"strings"

	"github.com/info344-a17/challenges-KyleIWS/servers/gateway/indexes"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoStore struct {
	session *mgo.Session
	dbname  string
	colname string
}

func NewMongoStore(sess *mgo.Session, dbName string, collectionName string) *MongoStore {
	if sess == nil {
		panic("nil pointer passed for session")
	}
	return &MongoStore{
		session: sess,
		dbname:  dbName,
		colname: collectionName,
	}
}

func (ms *MongoStore) ObjectIdsToUsers(ids []bson.ObjectId) []*User {
	users := make([]*User, 0, len(ids))
	for _, objectId := range ids {
		user, err := ms.GetByID(objectId)
		if err == nil {
			users = append(users, user)
		}
	}
	return users
}

func (ms *MongoStore) LoadExistingUsers(root *indexes.TrieNode) (int, error) {
	result := &User{}
	count := 0
	col := ms.session.DB(ms.dbname).C(ms.colname)
	iterVal := col.Find(bson.M{}).Iter()
	for iterVal.Next(result) {
		count = count + 1
		root.Add(strings.ToLower(result.Email), result.ID)
		root.Add(strings.ToLower(result.UserName), result.ID)
		root.Add(strings.ToLower(result.FirstName), result.ID)
		root.Add(strings.ToLower(result.LastName), result.ID)
	}
	if err := iterVal.Err(); err != nil {
		return 0, err
	}
	return count, nil
}

//GetByID returns the User with the given ID
func (ms *MongoStore) GetByID(id bson.ObjectId) (*User, error) {
	result := &User{}
	col := ms.session.DB(ms.dbname).C(ms.colname)
	if err := col.Find(bson.M{"_id": id}).One(&result); err != nil {
		return InvalidUser, err
	}
	return result, nil
}

//GetByEmail returns the User with the given email
func (ms *MongoStore) GetByEmail(email string) (*User, error) {
	result := &User{}
	col := ms.session.DB(ms.dbname).C(ms.colname)
	if err := col.Find(bson.M{"email": email}).One(&result); err != nil {
		return InvalidUser, err
	}
	return result, nil
}

//GetByUserName returns the User with the given Username
func (ms *MongoStore) GetByUserName(username string) (*User, error) {
	result := &User{}
	col := ms.session.DB(ms.dbname).C(ms.colname)
	if err := col.Find(bson.M{"username": username}).One(&result); err != nil {
		return InvalidUser, err
	}
	return result, nil
}

//Insert converts the NewUser to a User, inserts
//it into the database, and returns it
func (ms *MongoStore) Insert(newUser *NewUser) (*User, error) {
	u, err := newUser.ToUser()
	if err != nil {
		return InvalidUser, fmt.Errorf("error getting user from new user: %v", err)
	}
	col := ms.session.DB(ms.dbname).C(ms.colname)
	if err := col.Insert(u); err != nil {
		return InvalidUser, fmt.Errorf("error inserting user to mongodb: %v", err)
	}
	return u, nil
}

//Update applies UserUpdates to the given user ID
func (ms *MongoStore) Update(userID bson.ObjectId, updates *Updates) error {
	updatedUser, err := ms.GetByID(userID)
	if err != nil {
		return err
	}
	errUp := updatedUser.ApplyUpdates(updates)
	if errUp != nil {
		return fmt.Errorf("error applying updates retrieved user: %v", errUp)
	}

	change := mgo.Change{
		Update:    bson.M{"$set": updatedUser},
		ReturnNew: true,
	}
	result := &User{}
	col := ms.session.DB(ms.dbname).C(ms.colname)
	if _, err := col.FindId(userID).Apply(change, result); err != nil {
		return fmt.Errorf("error updating record: %v", err)
	}
	return nil
}

//Delete deletes the user with the given ID
func (ms *MongoStore) Delete(userID bson.ObjectId) error {
	col := ms.session.DB(ms.dbname).C(ms.colname)
	if err := col.RemoveId(userID); err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}
