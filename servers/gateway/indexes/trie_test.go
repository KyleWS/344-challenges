package indexes

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestNewTrieNode(t *testing.T) {
	node := NewTrieNode(32, nil)
	if node.Key != 32 {
		t.Errorf("error setting new trie node to rune alue 32")
	}
}

func TestOneGetN(t *testing.T) {
	head := NewTrieNode(0, nil)
	world := bson.NewObjectId()
	head.Add("hello", world)
	_, err := head.GetN("hello", 1)
	if err != nil {
		t.Errorf("error got error when none was expected, when getting single value from trie %v", err)
	}
}

func TestManyGetN(t *testing.T) {
	head := NewTrieNode(0, nil)
	hey := bson.NewObjectId()
	ethan := bson.NewObjectId()
	howsit := bson.NewObjectId()
	going := bson.NewObjectId()
	checkSlice := []bson.ObjectId{hey, ethan, howsit, going}
	head.Add("hello", hey)
	head.Add("hello", ethan)
	head.Add("hello", howsit)
	head.Add("hello", going)
	retrievedSlice, err := head.GetN("hello", 4)
	if err != nil {
		t.Errorf("error when getting multiple values from trie %v", err)
	}
	if len(retrievedSlice) < 4 {
		t.Errorf("error should not have retrieved slice of length less than 4")
	}
	for i, val := range retrievedSlice {
		if val != checkSlice[i] {
			t.Errorf("error of retrieved slice is wrong in many get n")
		}
	}
}

func TestEmptyGetN(t *testing.T) {
	head := NewTrieNode(0, nil)
	retrievedSlice, err := head.GetN("hello", 100)
	if err != nil && err.Error() != ErrValueNotFound {
		t.Errorf("error when getting multiple values from trie %v", err)
	}
	if len(retrievedSlice) > 0 {
		t.Errorf("error should be empty slice")
	}
}

func TestSingleBranchGetN(t *testing.T) {
	head := NewTrieNode(0, nil)
	hey := bson.NewObjectId()
	ethan := bson.NewObjectId()
	howsit := bson.NewObjectId()
	going := bson.NewObjectId()
	checkSlice := []bson.ObjectId{hey, ethan, howsit, going}
	head.Add("hello", hey)
	head.Add("helloo", ethan)
	head.Add("hellooo", howsit)
	head.Add("helloooo", going)
	retrievedSlice, err := head.GetN("hello", 4)
	if err != nil {
		t.Errorf("error when getting multiple values from trie %v", err)
	}
	if len(retrievedSlice) < 4 {
		t.Errorf("error should not have retrieved slice of length less than 4")
	}
	for i, val := range retrievedSlice {
		if val != checkSlice[i] {
			t.Errorf("error of retrieved slice is wrong in many get n")
		}
	}
}

func TestManyBranchGetN(t *testing.T) {
	head := NewTrieNode(0, nil)
	hey := bson.NewObjectId()
	ethan := bson.NewObjectId()
	howsit := bson.NewObjectId()
	going := bson.NewObjectId()
	checkSlice := []bson.ObjectId{hey, ethan, howsit, going}
	head.Add("hello", hey)
	head.Add("helloa", ethan)
	head.Add("hellob", howsit)
	head.Add("helloc", going)
	retrievedSlice, err := head.GetN("hello", 4)
	if err != nil {
		t.Errorf("error when getting multiple values from trie %v", err)
	}
	if len(retrievedSlice) < 4 {
		t.Errorf("error should not have retrieved slice of length less than 4")
	}
	for i, val := range retrievedSlice {
		if val != checkSlice[i] {
			t.Errorf("error of retrieved slice is wrong in many get n")
		}
	}
}

func TestAddSingleValue(t *testing.T) {
	head := NewTrieNode(0, nil)
	world := bson.NewObjectId()
	if err := head.Add("hello", world); err != nil {
		t.Errorf("error when adding to trie %v", err)
	}
	retrievedSlice, _ := head.GetN("hello", 2)
	if retrievedSlice == nil {
		t.Errorf("error retrieved bson slice should not be nil")
	}
	if len(retrievedSlice) == 0 {
		t.Errorf("error retrieved bson slice should not be empty")
	}
	if retrievedSlice[0] != world {
		t.Errorf("error retrieved bson slice does not match")
	}
}

func TestAddManyValuesOneBranch(t *testing.T) {
	head := NewTrieNode(0, nil)
	world := bson.NewObjectId()
	otherWorld := bson.NewObjectId()
	underWorld := bson.NewObjectId()
	underWear := bson.NewObjectId()
	checkSlice := []bson.ObjectId{world, otherWorld, underWorld, underWear}
	head.Add("hello", world)
	head.Add("helloo", otherWorld)
	head.Add("hellooo", underWorld)
	head.Add("helloooo", underWear)

	retrievedSlice, _ := head.GetN("hello", 4)
	if retrievedSlice == nil {
		t.Errorf("error retrieved bson slice should not be nil")
	}
	if len(retrievedSlice) < 4 {
		t.Errorf("error retrieved bson slice should not be less than 4 elements long")
	}
	for i, retrieved := range retrievedSlice {
		if retrieved != checkSlice[i] {
			t.Errorf("error did not get matching elements: %v should be the same as %v", retrieved, checkSlice[i])
		}
	}
}

func TestAddManyValuesManyBranches(t *testing.T) {
	head := NewTrieNode(0, nil)
	world := bson.NewObjectId()
	otherWorld := bson.NewObjectId()
	underWorld := bson.NewObjectId()
	underWear := bson.NewObjectId()
	checkSlice := []bson.ObjectId{world, otherWorld, underWorld, underWear}
	head.Add("hello", world)
	head.Add("helloa", otherWorld)
	head.Add("hellob", underWorld)
	head.Add("helloc", underWear)

	retrievedSlice, _ := head.GetN("hello", 4)
	if retrievedSlice == nil {
		t.Errorf("error retrieved bson slice should not be nil")
	}
	if len(retrievedSlice) < 4 {
		t.Errorf("error retrieved bson slice should not be less than 4 elements long")
	}
	for i, retrieved := range retrievedSlice {
		if retrieved != checkSlice[i] {
			t.Errorf("error did not get matching elements: %v should be the same as %v", retrieved, checkSlice[i])
		}
	}
}

func TestAddSameValueTwiceInSamePrefix(t *testing.T) {
	head := NewTrieNode(0, nil)
	world := bson.NewObjectId()
	head.Add("hello", world)
	head.Add("hello", world)

	retrievedSlice, _ := head.GetN("hello", 2)
	if retrievedSlice == nil {
		t.Errorf("error retrieved bson slice should not be nil")
	}
	if len(retrievedSlice) > 1 {
		t.Errorf("error there should only be one value when adding the same value twice")
	}
}

func TestDeleteOne(t *testing.T) {
	head := NewTrieNode(0, nil)
	hey := bson.NewObjectId()
	head.Add("hello", hey)
	head.Delete("hello", hey)
	retrievedSlice, err := head.GetN("hello", 10)
	if err != nil && err.Error() != ErrValueNotFound {
		t.Errorf("getting after deleting should get %s", ErrValueNotFound)
	}
	if len(retrievedSlice) > 0 {
		t.Errorf("retrieved slice should be empty")
	}
	if err := head.Delete("hello", hey); err == nil {
		t.Errorf("should not be able to delete same value twice")
	}
}

func TestDeleteOneOfMany(t *testing.T) {
	head := NewTrieNode(0, nil)
	hey := bson.NewObjectId()
	ethan := bson.NewObjectId()
	howsit := bson.NewObjectId()
	going := bson.NewObjectId()
	checkSlice := []bson.ObjectId{ethan, howsit, going}
	head.Add("hello", hey)
	head.Add("hello", ethan)
	head.Add("hello", howsit)
	head.Add("hello", going)
	head.Delete("hello", hey)
	retrievedSlice, err := head.GetN("hello", 4)
	if err != nil {
		t.Errorf("error when getting multiple values from trie %v", err)
	}
	if len(retrievedSlice) < 3 {
		t.Errorf("error should not have retrieved slice of length less than 3")
	}
	for i, val := range retrievedSlice {
		if val != checkSlice[i] {
			t.Errorf("error when comparing slice to deleted version fo slice")
		}
	}
}

func TestDeleteTail(t *testing.T) {
	head := NewTrieNode(0, nil)
	world := bson.NewObjectId()
	underWear := bson.NewObjectId()
	head.Add("a", world)
	head.Add("abc", underWear)
	head.Delete("abc", underWear)

	retrievedSlice, _ := head.GetN("a", 4)
	if retrievedSlice == nil {
		t.Errorf("error retrieved bson slice should not be nil")
	}
	if head.Children['a'].Children['b'] != nil {
		t.Errorf("error when deleting a tail node, the intermediate nodes between it and the next value holding parent should be delteed")
	}
	if len(retrievedSlice) > 1 {
		t.Errorf("should not get more than one value back after deleting one value from two-value trie")
	}
}
