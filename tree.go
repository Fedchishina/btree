package btree

import (
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/exp/constraints"
)

type Tree[V constraints.Ordered] struct {
	storage NodeStorage[V]
	t       int
}

// NewTree is a function for creation empty tree
// - type V should be `ordered type` (`int`, `string`, `float` etc.)
// - param t is a min degree of b-tree. It can't be less than 2
// - param s is a storage when will be saved tree data
func NewTree[V constraints.Ordered](t int, s NodeStorage[V]) (*Tree[V], error) {
	if t < 2 {
		return nil, errors.New("t can't be less than 2")
	}

	return &Tree[V]{
		t:       t,
		storage: s,
	}, nil
}

// Exists is a function for searching key in Tree. If key exists in tree - returns true, else - returns false
// - param k should be `ordered type` (`int`, `string`, `float` etc.)
func (t *Tree[V]) Exists(k V) (bool, error) {
	root, err := t.storage.Read(RootName)
	if err != nil {
		return false, err
	}

	s, _, err := t.search(root, k)
	if err != nil {
		return false, err
	}

	return s != nil, nil
}

// Insert is a function for inserting element into Tree
// - param k should be `ordered type` (`int`, `string`, `float` etc.)
func (t *Tree[V]) Insert(k V) error {
	root, err := t.storage.Read(RootName)
	if err != nil {
		return err
	}

	if len(root.Keys) == t.maxKeysLength() {
		s := NewNode[V](t.t, RootName)
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

// maxKeysLength - internal function: return max amount of tree's keys in one Node
func (t *Tree[V]) maxKeysLength() int {
	return 2*t.t - 1
}

// insertNonFull - internal function for inserting key to a blank Node
func (t *Tree[V]) insertNonFull(n *Node[V], k V) error {
	i := 0
	for i < len(n.Keys) && k > n.Keys[i] {
		i++
	}

	if n.Leaf {
		n.insertKey(i, k)
		return t.storage.Write(n)
	}

	c, err := t.storage.Read(n.Children[i])
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
		c, err = t.storage.Read(n.Children[i])
	}

	return t.insertNonFull(c, k)
}

// splitChild - internal function for splitting Node with full amount of keys to two nodes
func (t *Tree[V]) splitChild(n, nodeToSplit *Node[V], i int) error {
	middleKey := nodeToSplit.Keys[t.t-1]
	n.insertKey(i, middleKey)

	newNode := newSplitNode(t.t, nodeToSplit, n.Name+strconv.Itoa(i+1))
	n.insertChild(i+1, newNode.Name)

	nodeToSplit.Name = n.Name + strconv.Itoa(i)
	nodeToSplit.Keys = nodeToSplit.Keys[:t.t-1]
	if !nodeToSplit.Leaf {
		nodeToSplit.Children = nodeToSplit.Children[:t.t]
	}

	if err := t.storage.Write(n); err != nil {
		return err
	}
	if err := t.storage.Write(newNode); err != nil {
		return err
	}
	if err := t.storage.Write(nodeToSplit); err != nil {
		return err
	}

	return nil
}

// search - search Node by key
func (t *Tree[V]) search(n *Node[V], k V) (*Node[V], int, error) {
	if n == nil {
		return nil, 0, nil
	}
	numKeys := len(n.Keys)

	i := 0
	for i < numKeys && k > n.Keys[i] {
		i++
	}

	if i < numKeys && k == n.Keys[i] {
		return n, i, nil
	}

	if n.Leaf {
		return nil, 0, nil
	}

	c, err := t.storage.Read(n.Children[i])
	if err != nil {
		return nil, 0, err
	}

	return t.search(c, k)
}

// Delete is a function for deleting Node by key in Tree
// - param k should be `ordered type` (`int`, `string`, `float` etc.)
// if Tree doesn't have this key - function returns an error
func (t *Tree[V]) Delete(k V) error {
	root, err := t.storage.Read(RootName)
	if err != nil {
		return err
	}

	n, i, err := t.search(root, k)
	if err != nil {
		return err
	}

	if n == nil {
		return errors.New(fmt.Sprintf("not found Node with key: %v", k))
	}

	if n.Leaf {
		n.deleteKeyByIndex(i)
		if err = t.storage.Write(n); err != nil {
			return err
		}
		return nil
	}

	childLeft, err := t.storage.Read(n.Children[i])
	if err != nil {
		return err
	}

	if len(childLeft.Keys) >= t.t {
		predecessor := childLeft.deleteMaxKey()
		n.Keys[i] = predecessor
		if err = t.storage.Write(n); err != nil {
			return err
		}

		if err = t.storage.Write(childLeft); err != nil {
			return err
		}

		return nil
	}

	childRight, err := t.storage.Read(n.Children[i+1])
	if err != nil {
		return err
	}
	if len(childRight.Keys) >= t.t {
		successor := childRight.deleteMinKey()
		n.Keys[i] = successor
		if err = t.storage.Write(n); err != nil {
			return err
		}

		if err = t.storage.Write(childRight); err != nil {
			return err
		}

		return nil
	}

	return t.mergeNodes(n, i)
}

// mergeNodes is an internal function for merging two nodes to one Node in Tree
func (t *Tree[V]) mergeNodes(n *Node[V], i int) error {
	leftChild, err := t.storage.Read(n.Children[i])
	if err != nil {
		return err
	}
	rightChild, err := t.storage.Read(n.Children[i+1])
	if err != nil {
		return err
	}

	leftChild.Keys = append(leftChild.Keys, rightChild.Keys...)

	if !leftChild.Leaf {
		leftChild.Children = append(leftChild.Children, rightChild.Children...)
	}

	if err = t.storage.Write(leftChild); err != nil {
		return err
	}

	n.Keys = append(n.Keys[:i], n.Keys[i+1:]...)
	n.Children = append(n.Children[:i+1], n.Children[i+2:]...)

	if len(n.Keys) == 0 {
		n.Keys = leftChild.Keys
		n.Children = leftChild.Children
		if len(n.Children) == 0 {
			n.Leaf = true
		}
		if err = t.storage.Delete(leftChild.Name); err != nil {
			return err
		}
	}

	if err = t.storage.Write(n); err != nil {
		return err
	}

	if err = t.storage.Delete(rightChild.Name); err != nil {
		return err
	}

	return nil
}
