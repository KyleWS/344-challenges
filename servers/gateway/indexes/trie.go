package indexes

import (
	"fmt"
	"sort"
	"sync"

	"gopkg.in/mgo.v2/bson"
)

const ErrKeyLength = "key length is zero"
const ErrRuneNotFound = "given rune is not a child of node"
const ErrValueNotFound = "given prefix not contained in trie"
const ErrValueAlreadyPresent = "given value is already in the trie"

type TrieNode struct {
	Children map[rune]*TrieNode
	Parent   *TrieNode
	Key      rune
	Values   map[bson.ObjectId]bool
	mx       sync.RWMutex
}

// RuneSlice is a helper datatype to help us sort
// a list of runes
type RuneSlice []rune

func (rs RuneSlice) Len() int {
	return len(rs)
}

func (rs RuneSlice) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs RuneSlice) Less(i, j int) bool {
	return rs[i] < rs[j]
}

///////// end of runeslice sort helper functions /////////////

// NewTrieNode creates a new node. Should be called with a rune to
// map as its key, and a node to store as its parent. The parent
// pointer enables easier delete logic.
func NewTrieNode(key rune, parent *TrieNode) *TrieNode {
	return &TrieNode{
		Key:      key,
		Parent:   parent,
		Children: make(map[rune]*TrieNode),
		Values:   make(map[bson.ObjectId]bool, 0),
	}
}

// helper function to see if the given node has the desired rune as a key.
// returns ErrRuneNotFound if the key is not in the node's children.
func (r *TrieNode) runifyFirstLetterOfKeyAndCheckForChild(key string) (*TrieNode, error) {
	runes := []rune(key)
	oneRune := runes[0]
	if r.Children != nil {
		if child, exists := r.Children[oneRune]; exists {
			return child, nil
		}
	}
	return nil, fmt.Errorf(ErrRuneNotFound)
}

// Adds a bson.ObjectID value to the given key index. \
// Returns ErrorValieAlreadyPresent if value is already
// in the Trie.
func (r *TrieNode) Add(key string, value bson.ObjectId) error {
	if len(key) == 0 {
		if _, exists := r.Values[value]; exists {
			return fmt.Errorf(ErrValueAlreadyPresent)
		}
		r.Values[value] = true
		return nil
	} else {
		childToGoTo, err := r.runifyFirstLetterOfKeyAndCheckForChild(key)
		// If the place we need to go does not exist...
		if err != nil && (err.Error() == ErrRuneNotFound) {
			// Pave the way with TrieNodes
			runes := []rune(key)
			childToGoTo = NewTrieNode(runes[0], r)
			r.Children[runes[0]] = childToGoTo
		}
		return childToGoTo.Add(key[1:], value)
	}
}

func (r *TrieNode) GetN(prefix string, n int) ([]bson.ObjectId, error) {
	if len(prefix) == 0 {
		// We are at the node we intended to get to, return Values
		obtainedBsons, _, err := SearchFromHere(r, n)
		if err != nil {
			return nil, err
		}
		return obtainedBsons, nil
	} else {
		childToGoTo, err := r.runifyFirstLetterOfKeyAndCheckForChild(prefix)
		if err != nil && (err.Error() == ErrRuneNotFound) {
			return nil, fmt.Errorf(ErrValueNotFound)
		}
		return childToGoTo.GetN(prefix[1:], n)
	}
}

func SearchFromHere(r *TrieNode, n int) ([]bson.ObjectId, int, error) {
	totalList := make([]bson.ObjectId, 0, len(r.Values)) // make n - len() at some point
	if len(r.Values) > 0 {
		sortedBsons := sortedBsonsFromSlice(r.Values)
		for _, bsonID := range sortedBsons {
			totalList = append(totalList, bsonID)
			n = n - 1
			if n <= 0 {
				return totalList, 0, nil
			}
		}
	}
	if len(r.Children) > 0 {
		sortedChildren := sortedRunesFromChildren(r.Children)
		for _, ru := range sortedChildren {
			bsonSlice, leafN, _ := SearchFromHere(r.Children[ru], n)
			totalList = append(totalList, bsonSlice...)
			n = leafN
			if n <= 0 {
				return totalList, 0, nil
			}
		}
	}
	return totalList, n, nil
}

// Remember to turn that shit to bson.ObjectID
func (r *TrieNode) Delete(key string, value bson.ObjectId) error {
	if len(key) == 0 {
		// We are where we want to delete said value
		if _, exists := r.Values[value]; !exists {
			return fmt.Errorf(ErrValueNotFound)
		}
		// If we make it here, we are good to delete our value
		delete(r.Values, value)
		// however we must consider the following cases:
		// // this value was one of many values in the map
		if len(r.Values) < 1 {
			if len(r.Children) == 0 {
				// if i have no children, lets delete upward
				r.Parent.deleteMe(r.Key)
			}
			// if this prefix has no more values,
			// but we still have children, then those obviously
			// have values because kyle is an expert programmer
			// so we just stop.
		}
		return nil
	} else {
		childToGoTo, err := r.runifyFirstLetterOfKeyAndCheckForChild(key)
		if err != nil && err.Error() == ErrRuneNotFound {
			return fmt.Errorf("given prefix not in tree: %v", err)
		}
		return childToGoTo.Delete(key[1:], value)
	}
}

func (r *TrieNode) deleteMe(key rune) {
	// I was called by my child.
	// First, I should delete that child.
	delete(r.Children, key)
	// Second, I should see if I myself should be deleted
	// // that means checking to see if I have values
	if len(r.Values) == 0 {
		// no values, how about children?
		if len(r.Children) == 0 {
			// alright, keep going. Am I the root?
			if r.Parent != nil {
				// I am not, lets do this
				r.Parent.deleteMe(r.Key)
			}
			// if I am the parent, and we already deleted that entry we
			// can stop here
		}
		// if there are other children, we are going to assume generously
		// kyle's programming skill and just stop. THose children _probably_
		// have other values in them
	}
	// if I now have values, then we definitely want to stop deleting.
}

func sortedRunesFromChildren(children map[rune]*TrieNode) []rune {
	runeList := make([]rune, 0, len(children))
	for ru := range children {
		runeList = append(runeList, ru)
	}
	sort.Sort(RuneSlice(runeList))
	return runeList
}

func sortedBsonsFromSlice(values map[bson.ObjectId]bool) []bson.ObjectId {
	// Make string slice
	bsonStringSlice := make([]string, 0, len(values))
	// add all the object ID hexes to it
	for IDs, _ := range values {
		bsonStringSlice = append(bsonStringSlice, IDs.Hex())
	}
	// sort that shit
	sort.Strings(bsonStringSlice)
	// add in order to new slice of bson.ObjectIdHex(s)
	bsonObjectIdSlice := make([]bson.ObjectId, 0, len(bsonStringSlice))
	for _, stringId := range bsonStringSlice {
		bsonObjectIdSlice = append(bsonObjectIdSlice, bson.ObjectIdHex(stringId))
	}
	return bsonObjectIdSlice
}

// func PrintTrie(root *TrieNode) {
// 	if root.Parent == nil {
// 		fmt.Println("Now printing Trie. Root has the following children.")
// 		fmt.Printf("children: %v\n", root.children())
// 	} else if root.Key != 0 {
// 		fmt.Printf("Looking at node: %s who has children %v\n", string(root.Key), root.children())
// 	}
// 	if root.Children != nil {
// 		for _, child := range root.Children {
// 			PrintTrie(child)
// 		}
// 	}
// }
//
// // helper function to print a given nodes children.
// func (r *TrieNode) children() string {
// 	if r.Children == nil {
// 		return "children is nil"
// 	} else if len(r.Children) == 0 {
// 		return "children is empty"
// 	} else {
// 		childString := ""
// 		for key, _ := range r.Children {
// 			childString += string(key) + " "
// 		}
// 		return childString
// 	}
// }
