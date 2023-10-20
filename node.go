package btree

import "golang.org/x/exp/constraints"

// node is the structure of tree's node.
// node's n is a count of keys in node so far
// node's keys is an array of ordered keys (each key has ordered type)
// node's children is an array of nodes (children of this node)
// leaf is a sign: node is leaf or not
type node[V constraints.Ordered] struct {
	n        int
	keys     []V
	children []*node[V]
	leaf     bool
}

// newNode - internal function for creating empty node
func newNode[V constraints.Ordered](t int) *node[V] {
	return &node[V]{
		n:        0,
		keys:     make([]V, 0),
		children: make([]*node[V], 0),
		leaf:     true,
	}
}

// newSplitNode - internal function. Create additional node from nodeToSplit and returns it
func newSplitNode[V constraints.Ordered](t int, nodeToSplit *node[V]) *node[V] {
	n := newNode[V](t)
	n.leaf = nodeToSplit.leaf
	n.n = t - 1
	n.keys = append(n.keys, nodeToSplit.keys[t:]...)

	if !nodeToSplit.leaf {
		n.children = append(n.children, nodeToSplit.children[t:]...)
	}

	return n
}

// insertIndex returns the index where a key should be inserted in a sorted key's slice
func insertIndex[V constraints.Ordered](n *node[V], k V) int {
	i := 0
	for i < n.n && k > n.keys[i] {
		i++
	}
	return i
}

// insertIndex - insert key to node on the i-position in key's array
func insertKeyToNode[V constraints.Ordered](n *node[V], i int, k V) {
	if len(n.keys) <= n.n {
		n.keys = append(n.keys, k)
	}

	copy(n.keys[i+1:], n.keys[i:])
	n.keys[i] = k
	n.n++
}

// insertNodeToNodes - insert node to array of nodes on the i-position
func insertNodeToNodes[V constraints.Ordered](nodes []*node[V], i int, n *node[V]) []*node[V] {
	return append(nodes[:i], append([]*node[V]{n}, nodes[i:]...)...)
}

// shouldSplitChildren - function for verifying should split node or not
func shouldSplitChildren[V constraints.Ordered](n *node[V], i, t int) bool {
	return len(n.children) > 0 && n.children[i].n == 2*t-1
}

// exists - verify node has key k in array of keys or not
func exists[V constraints.Ordered](n *node[V], k V) bool {
	s, _ := search(n, k)

	return s != nil
}

// search - search node by key and return it with founded index
func search[V constraints.Ordered](n *node[V], k V) (*node[V], int) {
	if n == nil || n.n == 0 {
		return nil, 0
	}

	i := insertIndex(n, k)
	if i < n.n && k == n.keys[i] {
		return n, i
	}

	if n.leaf {
		return nil, 0
	}

	return search(n.children[i], k)
}

// deleteMaxKey - delete max key in array of node's keys
func (n *node[V]) deleteMaxKey() V {
	maxKey := n.keys[len(n.keys)-1]
	n.deleteKeyByIndex(len(n.keys) - 1)

	return maxKey
}

// deleteMinKey - delete min key in array of node's keys
func (n *node[V]) deleteMinKey() V {
	minKey := n.keys[0]
	n.deleteKeyByIndex(0)

	return minKey
}

// deleteKeyByIndex - delete key by index in array of node's keys
func (n *node[V]) deleteKeyByIndex(i int) {
	n.keys = append(n.keys[:i], n.keys[i+1:]...)
	n.n--
}

// deleteKeyByKey - delete key by key in array of node's keys
func (n *node[V]) deleteKeyByKey(k V) {
	// Find the index 'i' where 'k' is located in the keys slice
	i := -1
	for index, key := range n.keys {
		if key == k {
			i = index
			break
		}
	}

	if i != -1 {
		n.deleteKeyByIndex(i)
	}
}

// merge function - merge two child nodes to one
func (n *node[V]) merge(index int) {
	leftChild := n.children[index]
	rightChild := n.children[index+1]

	leftChild.keys = append(leftChild.keys, n.keys[index])
	leftChild.keys = append(leftChild.keys, rightChild.keys...)
	leftChild.n = len(leftChild.keys)

	if !leftChild.leaf {
		leftChild.children = append(leftChild.children, rightChild.children...)
	}

	n.keys = append(n.keys[:index], n.keys[index+1:]...)
	n.n = len(n.keys)
	n.children = append(n.children[:index+1], n.children[index+2:]...)

	if len(n.keys) == 0 {
		n.keys = leftChild.keys
		n.children = leftChild.children
		n.n = leftChild.n
		if len(n.children) == 0 {
			n.leaf = true
		}
	}
}
