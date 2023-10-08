btree
=======================

Library for work with B-trees.

You can create a B-tree and use a list of functions to work with it.
## Tree functions
- [Empty tree's creation example](#empty-trees-creation-example)
- [Insert key to tree](#insert-key-to-tree)
- [Exists element](#exists-element)
- [Delete element by key from tree](#delete-element-by-key-from-tree)

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