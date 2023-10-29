package btree

import "golang.org/x/exp/constraints"

// nodeStorage is the structure of TreeStorage's Node.
// nodeStorage's name is a name of TreeStorage's Node
// nodeStorage's keys is an array of ordered keys (each key has ordered type)
// nodeStorage's children is an array of nodeStorage names (children of this nodeStorage)
// leaf is a sign: nodeStorage is leaf or not
type nodeStorage[V constraints.Ordered] struct {
	Name     string
	Keys     []V
	Children []string
	Leaf     bool
}

// newNode - internal function for creating empty nodeStorage
func newNodeStorage[V constraints.Ordered](t int, name string) *nodeStorage[V] {
	return &nodeStorage[V]{
		Name:     name,
		Keys:     make([]V, 0, 2*t-1),
		Children: make([]string, 0, 2*t),
		Leaf:     true,
	}
}

// insertKeyToNodeStorage - insert key to nodeStorage on the i-position in key's array
func insertKeyToNodeStorage[V constraints.Ordered](n *nodeStorage[V], i int, k V) {
	n.Keys = append(n.Keys, k)
	copy(n.Keys[i+1:], n.Keys[i:])
	n.Keys[i] = k
}

// insertChild - insert child to children of nodeStorage on the i-position
func insertChild[V constraints.Ordered](n *nodeStorage[V], i int, child string) {
	n.Children = append(n.Children, child)
	copy(n.Children[i+1:], n.Children[i:])
	n.Children[i] = child
}

// newSplitNodeStorage - internal function. Create additional nodeStorage from nodeToSplit and returns it
func newSplitNodeStorage[V constraints.Ordered](t int, nodeToSplit *nodeStorage[V], name string) *nodeStorage[V] {
	n := newNodeStorage[V](t, name)
	n.Leaf = nodeToSplit.Leaf
	n.Keys = append(n.Keys, nodeToSplit.Keys[t:]...)

	if !nodeToSplit.Leaf {
		n.Children = append(n.Children, nodeToSplit.Children[t:]...)
	}

	return n
}
