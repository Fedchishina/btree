btree
=======================
Library for work with B-trees.
In this library we have two realisation: keeping b-tree in memory and in file storage.

You can create a B-tree and use a list of functions to work with it.
## Tree functions (in-memory saving)
- [Empty tree's creation example](#empty-trees-creation-example)
- [Insert key to tree](#insert-key-to-tree)
- [Exists element](#exists-element)
- [Delete element by key from tree](#delete-element-by-key-from-tree)

## Tree functions (storage saving)
- [Empty tree's creation in storage example](#empty-trees-creation-in-storage-example)
- [Insert key to tree in storage](#insert-key-to-tree-in-storage)
- [Exists element in tree](#exists-element-in-tree)
- [Delete element by key from tree in storage](#delete-element-by-key-from-tree-in-storage)

### Empty tree's creation example

```
t := btree.New[int]() // empty int tree
t := btree.New[string]() // empty string tree
```

### Insert key to tree
```
t := btree.New[int]() // empty int tree
t.Insert(22)
t.Insert(8)
t.Insert(4)
```

### Exists element

```
t := tree.New[int]()
t.Insert(22)
t.Insert(8)
t.Insert(4)

resultNil := t.Exists(15) // false
result    := t.Exists(8)  // true
```

### Delete element by key from tree
```
t := btree.New[int]()
t.Insert(22)
t.Insert(8)
t.Insert(4)

err := t.Delete(22) // without err
```

### Empty tree's creation in storage example

```
t := btree.NewTreeStorage[int](3, "myTree") // empty int tree
t := btree.NewTreeStorage[string](3, "myTree") // empty string tree
```

### Insert key to tree in storage
```
t := btree.NewTreeStorage[int](3, "myTree") // empty int tree
t.Insert(22)
t.Insert(8)
t.Insert(4)
```

### Exists element in tree

```
t := btree.NewTreeStorage[int](3, "myTree") 
t.Insert(22)
t.Insert(8)
t.Insert(4)

resultNil := t.Exists(15) // false
result    := t.Exists(8)  // true
```

### Delete element by key from tree in storage
```
t := btree.NewTreeStorage[int](3, "myTree") 
t.Insert(22)
t.Insert(8)
t.Insert(4)

err := t.Delete(22) // without err
```