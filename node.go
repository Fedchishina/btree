package btree

import "golang.org/x/exp/constraints"

// Node is the structure of Tree's Node.
// Name is a name of Tree's node
// Keys is an array of ordered keys (each key has ordered type)
// Children is an array of Node names (children of this Node)
// Leaf is a sign: Node is leaf or not
type Node[V constraints.Ordered] struct {
	Name     string
	Keys     []V
	Children []string
	Leaf     bool
}

// NewNode - internal function for creating empty Node
func NewNode[V constraints.Ordered](t int, name string) *Node[V] {
	return &Node[V]{
		Name:     name,
		Keys:     make([]V, 0, 2*t-1),
		Children: make([]string, 0, 2*t),
		Leaf:     true,
	}
}

// insertKey - insert key to Node on the i-position in key's array
func (n *Node[V]) insertKey(i int, k V) {
	n.Keys = append(n.Keys, k)
	copy(n.Keys[i+1:], n.Keys[i:])
	n.Keys[i] = k
}

// insertChild - insert child to children of Node on the i-position
func (n *Node[V]) insertChild(i int, child string) {
	n.Children = append(n.Children, child)
	copy(n.Children[i+1:], n.Children[i:])
	n.Children[i] = child
}

// newSplitNode - internal function. Create additional Node from nodeToSplit and returns it
func newSplitNode[V constraints.Ordered](t int, nodeToSplit *Node[V], name string) *Node[V] {
	n := NewNode[V](t, name)
	n.Leaf = nodeToSplit.Leaf
	n.Keys = append(n.Keys, nodeToSplit.Keys[t:]...)

	if !nodeToSplit.Leaf {
		n.Children = append(n.Children, nodeToSplit.Children[t:]...)
	}

	return n
}

// deleteMaxKey - delete max key in array of Node's keys
func (n *Node[V]) deleteMaxKey() V {
	maxKey := n.Keys[len(n.Keys)-1]
	n.deleteKeyByIndex(len(n.Keys) - 1)

	return maxKey
}

// deleteMinKey - delete min key in array of Node's keys
func (n *Node[V]) deleteMinKey() V {
	minKey := n.Keys[0]
	n.deleteKeyByIndex(0)

	return minKey
}

// deleteKeyByIndex - delete key by key in array of Node's keys
func (n *Node[V]) deleteKeyByIndex(i int) {
	n.Keys = append(n.Keys[:i], n.Keys[i+1:]...)
}
