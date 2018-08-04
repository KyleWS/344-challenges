package users

import (
	"fmt"
	"testing"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var nuTest1 = &NewUser{
	Email:        "kyle@gmail.com",
	FirstName:    "Kyle",
	LastName:     "Williams-Smith",
	Password:     "123456",
	PasswordConf: "123456",
	UserName:     "narutohokage123",
}

var nuTest2 = &NewUser{
	Email:        "casey@gmail.com",
	FirstName:    "Casey",
	LastName:     "Williams-Smith",
	Password:     "123456",
	PasswordConf: "123456",
	UserName:     "theRealStuartReges",
}

var nuTest3 = &NewUser{
	Email:        "elmo@gmail.com",
	FirstName:    "Elmo",
	LastName:     "TheRedTrollThing",
	Password:     "123456",
	PasswordConf: "123456",
	UserName:     "hugme",
}

func GetNewMongoStore() (*MongoStore, error) {
	sess, err := mgo.Dial("127.0.0.1")
	if err != nil {
		return nil, fmt.Errorf("error connecting to local mongodb: %v", err)
	}
	ms := NewMongoStore(sess, "dbtest", "coltest")
	return ms, nil
}

func ClearCollection(ms *MongoStore) error {
	col := ms.session.DB(ms.dbname).C(ms.colname)
	if _, err := col.RemoveAll(bson.M{}); err != nil {
		return fmt.Errorf("error could not clear collection")
	}
	return nil
}

func TestMongoGetByID(t *testing.T) {
	ms, err := GetNewMongoStore()
	if err != nil {
		t.Errorf("error connecting to db: %v", err)
	}
	if err := ClearCollection(ms); err != nil {
		t.Errorf("error clearing database before test: %v", err)
	}
	user1, _ := ms.Insert(nuTest1)
	ms.Insert(nuTest2)
	ms.Insert(nuTest3)
	result, err := ms.GetByID(user1.ID)
	if err != nil {
		t.Errorf("error trying to get users from mongodb: %v", err)
	} else if result.ID != user1.ID {
		t.Errorf("error retrieved incorrect user")
	}
}

func TestMongoGetByEmail(t *testing.T) {
	ms, err := GetNewMongoStore()
	if err != nil {
		t.Errorf("error connecting to db: %v", err)
	}
	if err := ClearCollection(ms); err != nil {
		t.Errorf("error clearing database before test: %v", err)
	}
	user1, _ := ms.Insert(nuTest1)
	ms.Insert(nuTest2)
	ms.Insert(nuTest3)
	result, err := ms.GetByEmail(user1.Email)
	fmt.Println(result)
	if err != nil {
		t.Errorf("error trying to get users from mongodb: %v", err)
	}
	if result.ID != user1.ID {
		t.Errorf("error retrieved incorrect user")
	}
}

func TestMongoGetByUserName(t *testing.T) {
	ms, err := GetNewMongoStore()
	if err != nil {
		t.Errorf("error connecting to db: %v", err)
	}
	if err := ClearCollection(ms); err != nil {
		t.Errorf("error clearing database before test: %v", err)
	}
	user1, _ := ms.Insert(nuTest1)
	ms.Insert(nuTest2)
	ms.Insert(nuTest3)
	result, err := ms.GetByUserName(user1.UserName)
	if err != nil {
		t.Errorf("error trying to get users from mongodb: %v", err)
	}
	if result.ID != user1.ID {
		t.Errorf("error retrieved incorrect user")
	}
}

func TestMongoInsert(t *testing.T) {
	ms, err := GetNewMongoStore()
	if err != nil {
		t.Errorf("error connecting to db: %v", err)
	}
	if err := ClearCollection(ms); err != nil {
		t.Errorf("error clearing database before test: %v", err)
	}
	_, errInsert := ms.Insert(nuTest1)
	if errInsert != nil {
		t.Errorf("error inserting user into mongodb: %v", errInsert)
	}
}

func TestMongoUpdate(t *testing.T) {
	ms, err := GetNewMongoStore()
	if err != nil {
		t.Errorf("error connecting to db: %v", err)
	}
	if err := ClearCollection(ms); err != nil {
		t.Errorf("error clearing database before test: %v", err)
	}
	user1, _ := ms.Insert(nuTest1)
	ms.Insert(nuTest2)
	ms.Insert(nuTest3)

	updates := &Updates{
		FirstName: "Mr.Kyle",
		LastName:  "Williamson-ShmancyPantsy",
	}
	if err := ms.Update(user1.ID, updates); err != nil {
		t.Errorf("error updating user record: %v", err)
	}
	user, err := ms.GetByID(user1.ID)
	if user.FirstName != "Mr.Kyle" {
		t.Errorf("error first name was not updated properly")
	}
	if user.LastName != "Williamson-ShmancyPantsy" {
		t.Errorf("error last name was not updated properly")
	}
}

func TestMongoDelete(t *testing.T) {
	ms, err := GetNewMongoStore()
	if err != nil {
		t.Errorf("error connecting to db: %v", err)
	}
	if err := ClearCollection(ms); err != nil {
		t.Errorf("error clearing database before test: %v", err)
	}
	user1, _ := ms.Insert(nuTest1)
	ms.Insert(nuTest2)
	ms.Insert(nuTest3)
	if err := ms.Delete(user1.ID); err != nil {
		t.Errorf("unexpected error when deleting user: %v", err)
	}
	_, errDelete := ms.GetByID(user1.ID)
	if errDelete == nil {
		t.Errorf("did not get error when one was expected. User should have been deleted.")
	}
}
