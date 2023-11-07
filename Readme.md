btree
=======================
Library for work with B-trees.

You can create a B-tree and use a list of functions to work with it.

In this library you have disk storage realisation: tree's structure is saved in json files.


You can make your own storage realisation implementing this NodeStorage interface:
```
type NodeStorage[V constraints.Ordered] interface {
	Name() string
	Read(name string) (*Node[V], error)
	Write(n *Node[V]) error
	Delete(name string) error
}
```

## Tree functions (DiskStorage realisation)
- [Empty tree's creation example](#empty-trees-creation-example)
- [Insert key to tree ](#insert-key-to-tree)
- [Exists element in tree](#exists-element-in-tree)
- [Delete element by key from tree](#delete-element-by-key-from-tree)

### Empty tree's creation example

```
intStorage, _ := btree.NewDiskStorage[int]("myIntTree", 3)
intTree, _ := btree.NewTree[int](3, intStorage) // empty int tree

stringStorage, _ := btree.NewDiskStorage[string]("myStringTree", 3)
stringTree, _ := btree.NewTree[string](3, stringStorage) // empty int tree
```

### Insert key to tree
```
storage, _ := btree.NewDiskStorage[int]("myTree", 3)
t, _ := btree.NewTree[int](3, storage) // empty int tree
t.Insert(22)
t.Insert(8)
t.Insert(4)
```

### Exists element in tree

```
storage, _ := btree.NewDiskStorage[int]("myTree", 3)
t, _ := btree.NewTree[int](3, storage) // empty int tree
t.Insert(22)
t.Insert(8)
t.Insert(4)

resultNil := t.Exists(15) // false
result    := t.Exists(8)  // true
```

### Delete element by key from tree
```
storage, _ := btree.NewDiskStorage[int]("myTree", 3)
t, _ := btree.NewTree[int](3, storage) // empty int tree
t.Insert(22)
t.Insert(8)
t.Insert(4)

err := t.Delete(22) // without err
```