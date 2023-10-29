package btree

import (
	"errors"
	"strconv"

	"golang.org/x/exp/constraints"
)

// RootName - is name of root node
const RootName = "0"

type TreeStorage[V constraints.Ordered] struct {
	storage storage[V]
	t       int
}

// NewTreeStorage is a function for creation empty tree
// - type V should be `ordered type` (`int`, `string`, `float` etc.)
// - param t is a min degree of b-tree. It can't be less than 2
// - param name is a name of storage when will be saved tree data
func NewTreeStorage[V constraints.Ordered](t int, name string) (*TreeStorage[V], error) {
	if t < 2 {
		return nil, errors.New("t can't be less than 2")
	}

	storage, err := NewFileStorage[V](name, t)
	if err != nil {
		return nil, err
	}

	return &TreeStorage[V]{
		t:       t,
		storage: storage,
	}, nil
}

// Exists is a function for searching element in TreeStorage. If element exists in tree - return true, else - false
// - param k should be `ordered type` (`int`, `string`, `float` etc.)
func (t *TreeStorage[V]) Exists(k V) (bool, error) {
	root, err := t.storage.readNode(RootName)
	if err != nil {
		return false, err
	}

	s, err := t.search(root, k)
	if err != nil {
		return false, err
	}

	return s != nil, nil
}

// Insert is a function for inserting element into TreeStorage
// - param k should be `ordered type` (`int`, `string`, `float` etc.)
func (t *TreeStorage[V]) Insert(k V) error {
	root, err := t.storage.readNode(RootName)
	if err != nil {
		return err
	}

	if len(root.Keys) == t.maxKeysLength() {
		s := newNodeStorage[V](t.t, RootName)
		s.Leaf = false
		s.Children = append(s.Children, RootName+RootName)
		if err := t.splitChild(s, root, 0); err != nil {
			return err
		}

		if err := t.insertNonFull(s, k); err != nil {
			return err
		}

		return nil
	}

	if err := t.insertNonFull(root, k); err != nil {
		return err
	}

	return nil
}

// maxKeysLength - internal function: return max amount of tree's keys in one node
func (t *TreeStorage[V]) maxKeysLength() int {
	return 2*t.t - 1
}

// insertNonFull - internal function for inserting key to a blank nodeStorage
func (t *TreeStorage[V]) insertNonFull(n *nodeStorage[V], k V) error {
	i := 0
	for i < len(n.Keys) && k > n.Keys[i] {
		i++
	}

	if n.Leaf {
		insertKeyToNodeStorage(n, i, k)
		return t.storage.writeNode(n)
	}

	c, err := t.storage.readNode(n.Children[i])
	if err != nil {
		return err
	}

	reReadChildren := false
	if len(c.Keys) == t.maxKeysLength() {
		if err := t.splitChild(n, c, i); err != nil {
			return err
		}
		if i < len(n.Keys) && k > n.Keys[i] {
			i++
			reReadChildren = true
		}
	}

	if reReadChildren {
		c, err = t.storage.readNode(n.Children[i])
	}

	return t.insertNonFull(c, k)
}

// splitChild - internal function for splitting nodeStorage with full amount of keys to two nodes
func (t *TreeStorage[V]) splitChild(n, nodeToSplit *nodeStorage[V], i int) error {
	middleKey := nodeToSplit.Keys[t.t-1]
	insertKeyToNodeStorage(n, i, middleKey)

	newNode := newSplitNodeStorage(t.t, nodeToSplit, n.Name+strconv.Itoa(i+1))
	insertChild(n, i+1, newNode.Name)

	nodeToSplit.Name = n.Name + strconv.Itoa(i)
	nodeToSplit.Keys = nodeToSplit.Keys[:t.t-1]
	if !nodeToSplit.Leaf {
		nodeToSplit.Children = nodeToSplit.Children[:t.t]
	}

	if err := t.storage.writeNode(n); err != nil {
		return err
	}
	if err := t.storage.writeNode(newNode); err != nil {
		return err
	}
	if err := t.storage.writeNode(nodeToSplit); err != nil {
		return err
	}

	return nil
}

// search - search nodeStorage by key
func (t *TreeStorage[V]) search(n *nodeStorage[V], k V) (*nodeStorage[V], error) {
	if n == nil {
		return nil, nil
	}
	numKeys := len(n.Keys)

	i := 0
	for i < numKeys && k > n.Keys[i] {
		i++
	}

	if i < numKeys && k == n.Keys[i] {
		return n, nil
	}

	if n.Leaf {
		return nil, nil
	}

	c, err := t.storage.readNode(n.Children[i])
	if err != nil {
		return nil, err
	}

	return t.search(c, k)
}
