package btree

import (
	"os"
	"reflect"
	"testing"

	"golang.org/x/exp/constraints"
)

func TestNewTreeStorage(t *testing.T) {
	type args struct {
		t int
	}
	type testCase[V constraints.Ordered] struct {
		name    string
		args    args
		want    *TreeStorage[V]
		wantErr bool
	}
	tests := []testCase[int]{
		{
			name:    "t_is_0",
			args:    args{t: 0},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative_t",
			args:    args{t: -100},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success_creating_empty_tree",
			args: args{t: 2},
			want: &TreeStorage[int]{
				t: 2,
				storage: &FileStorage[int]{
					folderName: "success_creating_empty_tree",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTreeStorage[int](tt.args.t, tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTreeStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTreeStorage() got = %v, want %v", got, tt.want)
			}
			os.RemoveAll("success_creating_empty_tree")
		})
	}
}

func TestTreeStorage_Exists(t1 *testing.T) {
	type args[V constraints.Ordered] struct {
		k V
	}
	type testCase[V constraints.Ordered] struct {
		name    string
		t       *TreeStorage[V]
		args    args[V]
		want    bool
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "empty_tree",
			t:       createTreeStorage(3, []string{}, "empty_tree"),
			args:    args[string]{k: "A"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "tree_with_one_element-not found",
			t:       createTreeStorage(3, []string{"A"}, "tree_with_one_element-not found"),
			args:    args[string]{k: "B"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "tree_with_one_element-found",
			t:       createTreeStorage(3, []string{"A"}, "tree_with_one_element-found"),
			args:    args[string]{k: "A"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "tree_with_several_elements_in_root-found",
			t:       createTreeStorage(3, []string{"A", "B", "C", "D"}, "tree_with_several_elements_in_root-found"),
			args:    args[string]{k: "C"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "tree_with_root_and_one_child-not_found",
			t:       createTreeStorage(3, []string{"A", "B", "D", "E", "F", "C"}, "tree_with_root_and_one_child-not_found"),
			args:    args[string]{k: "K"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "tree_with_root_and_one_child-found",
			t:       createTreeStorage(3, []string{"A", "B", "D", "E", "F", "C"}, "tree_with_root_and_one_child-found"),
			args:    args[string]{k: "F"},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			got, err := tt.t.Exists(tt.args.k)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t1.Errorf("Exists() got = %v, want %v", got, tt.want)
			}
			os.RemoveAll(tt.name)
		})
	}
}

func TestTreeStorage_Insert(t1 *testing.T) {
	type args[V constraints.Ordered] struct {
		k V
	}

	type testCase[V constraints.Ordered] struct {
		name         string
		t            *TreeStorage[V]
		args         args[V]
		wantErr      bool
		wantRootNode *nodeStorage[V]
	}

	tests := []testCase[string]{
		{
			name:    "insert_leaf_in_the_end",
			t:       createTreeStorage(3, []string{"A", "B", "C"}, "insert_leaf_in_the_end"),
			args:    args[string]{k: "D"},
			wantErr: false,
			wantRootNode: &nodeStorage[string]{
				Name:     "0",
				Keys:     []string{"A", "B", "C", "D"},
				Children: []string{},
				Leaf:     true,
			},
		},
		{
			name:    "insert_leaf_in_the_beginning",
			t:       createTreeStorage(3, []string{"B", "C", "D"}, "insert_leaf_in_the_beginning"),
			args:    args[string]{k: "A"},
			wantErr: false,
			wantRootNode: &nodeStorage[string]{
				Name:     "0",
				Keys:     []string{"A", "B", "C", "D"},
				Children: []string{},
				Leaf:     true,
			},
		},
		{
			name:    "insert_leaf_in_the_middle",
			t:       createTreeStorage(3, []string{"A", "B", "D"}, "insert_leaf_in_the_middle"),
			args:    args[string]{k: "C"},
			wantErr: false,
			wantRootNode: &nodeStorage[string]{
				Name:     "0",
				Keys:     []string{"A", "B", "C", "D"},
				Children: []string{},
				Leaf:     true,
			},
		},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if err := tt.t.Insert(tt.args.k); (err != nil) != tt.wantErr {
				t1.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
			v := []validNodeStorage[string]{
				{
					nodeName: "0",
					want:     tt.wantRootNode,
				},
			}
			checkTreeStructure(tt.t, t1, v)
			os.RemoveAll(tt.name)
		})
	}
}

func TestTreStorage_Insert_to_full_node(t1 *testing.T) {
	testFolder := "insert_to_full_node"
	defer os.RemoveAll(testFolder)

	t := createTreeStorage(3, []string{"A", "B", "D", "E", "F"}, testFolder)

	validTree := []validNodeStorage[string]{
		{
			nodeName: "0",
			want: &nodeStorage[string]{
				Name:     "0",
				Keys:     []string{"A", "B", "D", "E", "F"},
				Children: []string{},
				Leaf:     true,
			},
		},
	}

	// check tree structure before insert
	checkTreeStructure(t, t1, validTree)

	t.Insert("C")

	// check tree structure after insert
	validTreeAfterInsert := []validNodeStorage[string]{
		{
			nodeName: "0",
			want: &nodeStorage[string]{
				Name:     "0",
				Keys:     []string{"D"},
				Children: []string{"00", "01"},
				Leaf:     false,
			},
		},
		{
			nodeName: "00",
			want: &nodeStorage[string]{
				Name:     "00",
				Keys:     []string{"A", "B", "C"},
				Children: []string{},
				Leaf:     true,
			},
		},
		{
			nodeName: "01",
			want: &nodeStorage[string]{
				Name:     "01",
				Keys:     []string{"E", "F"},
				Children: []string{},
				Leaf:     true,
			},
		},
	}

	checkTreeStructure(t, t1, validTreeAfterInsert)
}

func TestTreStorage_Insert_split_leaf(t1 *testing.T) {
	testFolder := "insert_split_leaf"
	defer os.RemoveAll(testFolder)

	t := createTreeStorage(3, []string{"A", "B", "D", "E", "F", "C"}, testFolder)

	validTree := []validNodeStorage[string]{
		{
			nodeName: "0",
			want: &nodeStorage[string]{
				Name:     "0",
				Keys:     []string{"D"},
				Children: []string{"00", "01"},
				Leaf:     false,
			},
		},
		{
			nodeName: "00",
			want: &nodeStorage[string]{
				Name:     "00",
				Keys:     []string{"A", "B", "C"},
				Children: []string{},
				Leaf:     true,
			},
		},
		{
			nodeName: "01",
			want: &nodeStorage[string]{
				Name:     "01",
				Keys:     []string{"E", "F"},
				Children: []string{},
				Leaf:     true,
			},
		},
	}

	// check tree structure before insert
	checkTreeStructure(t, t1, validTree)

	t.Insert("G")
	t.Insert("K")
	t.Insert("M")

	// check tree structure after insert
	validTreeAfterInsert := []validNodeStorage[string]{
		{
			nodeName: "0",
			want: &nodeStorage[string]{
				Name:     "0",
				Keys:     []string{"D"},
				Children: []string{"00", "01"},
				Leaf:     false,
			},
		},
		{
			nodeName: "00",
			want: &nodeStorage[string]{
				Name:     "00",
				Keys:     []string{"A", "B", "C"},
				Children: []string{},
				Leaf:     true,
			},
		},
		{
			nodeName: "01",
			want: &nodeStorage[string]{
				Name:     "01",
				Keys:     []string{"E", "F", "G", "K", "M"},
				Children: []string{},
				Leaf:     true,
			},
		},
	}

	checkTreeStructure(t, t1, validTreeAfterInsert)

	t.Insert("S")

	validTreeAfterInsert2 := []validNodeStorage[string]{
		{
			nodeName: "0",
			want: &nodeStorage[string]{
				Name:     "0",
				Keys:     []string{"D", "G"},
				Children: []string{"00", "01", "02"},
				Leaf:     false,
			},
		},
		{
			nodeName: "00",
			want: &nodeStorage[string]{
				Name:     "00",
				Keys:     []string{"A", "B", "C"},
				Children: []string{},
				Leaf:     true,
			},
		},
		{
			nodeName: "01",
			want: &nodeStorage[string]{
				Name:     "01",
				Keys:     []string{"E", "F"},
				Children: []string{},
				Leaf:     true,
			},
		},
		{
			nodeName: "02",
			want: &nodeStorage[string]{
				Name:     "02",
				Keys:     []string{"K", "M", "S"},
				Children: []string{},
				Leaf:     true,
			},
		},
	}

	checkTreeStructure(t, t1, validTreeAfterInsert2)
}

func checkTreeStructure(t *TreeStorage[string], t1 *testing.T, treeData []validNodeStorage[string]) {
	for _, vn := range treeData {
		n, err := t.storage.readNode(vn.nodeName)
		if err != nil {
			t1.Errorf("Error reading node %s: %v", vn.nodeName, err)
			continue
		}

		if !reflect.DeepEqual(n, vn.want) {
			t1.Errorf("Node %s has unexpected content. Got: %+v, Want: %+v", vn.nodeName, n, vn.want)
		}
	}
}

type validNodeStorage[V constraints.Ordered] struct {
	nodeName string
	want     *nodeStorage[V]
}

func createTreeStorage(t int, elements []string, name string) *TreeStorage[string] {
	tree, _ := NewTreeStorage[string](t, name)
	for _, el := range elements {
		tree.Insert(el)
	}

	return tree
}
