package btree

import (
	"errors"
	"fmt"

	"golang.org/x/exp/constraints"
)

type Tree[V constraints.Ordered] struct {
	root *node[V]
	t    int
}

// New is a function for creation empty tree
// - type V should be `ordered type` (`int`, `string`, `float` etc.)
// - param t is a min degree of b-tree. It can't be less than 2
func New[V constraints.Ordered](t int) (*Tree[V], error) {
	if t < 2 {
		return nil, errors.New("t can't be less than 2")
	}

	n := newNode[V](t)

	return &Tree[V]{
		root: n,
		t:    t,
	}, nil
}

// Exists is a function for searching element in Tree. If element exists in tree - return true, else - false
// - param k should be `ordered type` (`int`, `string`, `float` etc.)
func (t *Tree[V]) Exists(k V) bool {
	return exists(t.root, k)
}

// Insert is a function for inserting element into Tree
// - param k should be `ordered type` (`int`, `string`, `float` etc.)
func (t *Tree[V]) Insert(k V) {
	r := t.root
	if r.n == 2*t.t-1 {
		s := newNode[V](t.t)
		t.root = s
		s.leaf = false
		s.children = append(s.children, r)
		t.splitChild(s, 0)
		t.insertNonFull(s, k)
	} else {
		t.insertNonFull(r, k)
	}
}

// Delete is a function for deleting key in Tree
// - param k should be `ordered type` (`int`, `string`, `float` etc.)
// if Tree doesn't have this key - function returns an error
func (t *Tree[V]) Delete(k V) error {
	n, i := search(t.root, k)

	if n == nil {
		return errors.New(fmt.Sprintf("not found node with key: %v", k))
	}

	if n.leaf {
		n.deleteKeyByIndex(i)
		return nil
	}

	if len(n.children[i].keys) >= t.t {
		predecessor := n.children[i].deleteMaxKey()
		n.keys[i] = predecessor
		return nil
	}
	if len(n.children[i+1].keys) >= t.t {
		successor := n.children[i+1].deleteMinKey()
		n.keys[i] = successor
		return nil
	}

	n.merge(i)

	return t.Delete(k)
}

// splitChild - internal function for splitting node with full amount of keys to two nodes
func (t *Tree[V]) splitChild(n *node[V], i int) {
	nodeToSplit := n.children[i]
	middleKey := nodeToSplit.keys[t.t-1]
	insertKeyToNode(n, i, middleKey)

	newNode := newSplitNode(t.t, nodeToSplit)
	n.children = insertNodeToNodes(n.children, i+1, newNode)

	nodeToSplit.n = t.t - 1
	nodeToSplit.keys = nodeToSplit.keys[:t.t-1]
	if !nodeToSplit.leaf {
		nodeToSplit.children = nodeToSplit.children[:t.t]
	}
}

// insertNonFull - internal function for inserting key to a blank node
func (t *Tree[V]) insertNonFull(n *node[V], k V) {
	i := insertIndex(n, k)

	if n.leaf {
		insertKeyToNode(n, i, k)
		return
	}

	if shouldSplitChildren(n, i, t.t) {
		t.splitChild(n, i)
		if k > n.keys[i] {
			i++
		}
	}

	t.insertNonFull(n.children[i], k)
}
