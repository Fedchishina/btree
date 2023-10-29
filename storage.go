package btree

import (
	"encoding/json"
	"errors"
	"os"

	"golang.org/x/exp/constraints"
)

type storage[V constraints.Ordered] interface {
	name() string
	readNode(nodeName string) (*nodeStorage[V], error)
	writeNode(n *nodeStorage[V]) error
}

// FileStorage - is a storage for keeping files of tree. Format of files in this realisation - json
// param folderName - is name of folder where will be saved files of tree
type FileStorage[V constraints.Ordered] struct {
	folderName string
}

// NewFileStorage - function for creating of file storage
// param folderName - is name of folder where will be saved files of tree
// - param t is a min degree of b-tree. It can't be less than 2
func NewFileStorage[V constraints.Ordered](folderName string, t int) (*FileStorage[V], error) {
	if t < 2 {
		return nil, errors.New("t can't be less than 2")
	}

	if err := os.Mkdir(folderName, os.ModePerm); err != nil {
		return nil, err
	}

	root := newNodeStorage[V](t, RootName)
	s := &FileStorage[V]{
		folderName: folderName,
	}

	if err := s.writeNode(root); err != nil {
		return nil, err
	}

	return s, nil
}

// readNode - function for reading node by name
// param nodeName - is name of nodeStorage file
func (fs *FileStorage[V]) readNode(nodeName string) (*nodeStorage[V], error) {
	data, err := os.ReadFile(fs.filePath(nodeName))
	if err != nil {
		return nil, err
	}

	var n nodeStorage[V]
	if err = json.Unmarshal(data, &n); err != nil {
		return nil, err
	}

	return &n, nil
}

// writeNode - function for writing node to storage
func (fs *FileStorage[V]) writeNode(n *nodeStorage[V]) error {
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

// filePath - this function returns filePath of node in storage
func (fs *FileStorage[V]) filePath(nodeName string) string {
	return fs.folderName + "/" + nodeName + ".json"
}

// name - this function returns name of storage where we keep a treeStorage
func (fs *FileStorage[V]) name() string {
	return fs.folderName
}
