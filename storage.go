package btree

import (
	"golang.org/x/exp/constraints"
)

// RootName - is name of root Node
const RootName = "0"

type NodeStorage[V constraints.Ordered] interface {
	Name() string
	Read(name string) (*Node[V], error)
	Write(n *Node[V]) error
	Delete(name string) error
}
