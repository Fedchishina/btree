package btree

import (
	"encoding/json"
	"errors"
	"os"

	"golang.org/x/exp/constraints"
)

// DiskStorage - is a storage for keeping files of Tree. Format of files in this realisation - json
// - param folderName is a name of folder where will be saved files of tree
type DiskStorage[V constraints.Ordered] struct {
	folderName string
}

// NewDiskStorage - function for creating of DiskStorage
// - param folderName is name of folder where will be saved files of tree
// - param t is a min degree of b-tree. It can't be less than 2
func NewDiskStorage[V constraints.Ordered](folderName string, t int) (*DiskStorage[V], error) {
	if t < 2 {
		return nil, errors.New("t can't be less than 2")
	}

	if err := os.Mkdir(folderName, os.ModePerm); err != nil {
		return nil, err
	}

	root := NewNode[V](t, RootName)
	s := &DiskStorage[V]{
		folderName: folderName,
	}

	if err := s.Write(root); err != nil {
		return nil, err
	}

	return s, nil
}

// Name - this function returns name of DiskStorage where we keep a Tree
func (fs *DiskStorage[V]) Name() string {
	return fs.folderName
}

// Read - function for reading Node by name from DiskStorage
// - param name - is name of Node file
func (fs *DiskStorage[V]) Read(name string) (*Node[V], error) {
	data, err := os.ReadFile(fs.filePath(name))
	if err != nil {
		return nil, err
	}

	var n Node[V]
	if err = json.Unmarshal(data, &n); err != nil {
		return nil, err
	}

	return &n, nil
}

// Write - function for writing Node to DiskStorage
func (fs *DiskStorage[V]) Write(n *Node[V]) error {
	jsonData, err := json.Marshal(n)
	if err != nil {
		return err
	}

	err = os.WriteFile(fs.filePath(n.Name), jsonData, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Delete - function for deleting Node from DiskStorage
// param name - is name of Node file
func (fs *DiskStorage[V]) Delete(name string) error {
	return os.Remove(fs.filePath(name))
}

// filePath - this function returns filePath of Node in DiskStorage
func (fs *DiskStorage[V]) filePath(name string) string {
	return fs.folderName + "/" + name + ".json"
}
