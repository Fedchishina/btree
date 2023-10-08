package btree

import (
	"reflect"
	"testing"

	"golang.org/x/exp/constraints"
)

func TestNew(t *testing.T) {
	type args struct {
		t int
	}
	type testCase[V constraints.Ordered] struct {
		name    string
		args    args
		want    *Tree[V]
		wantErr bool
	}
	tests := []testCase[int]{
		{
			name:    "t is 0",
			args:    args{t: 0},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative t",
			args:    args{t: -100},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success creating empty tree",
			args: args{t: 2},
			want: &Tree[int]{
				t: 2,
				root: &node[int]{
					n:        0,
					keys:     make([]int, 0),
					children: make([]*node[int], 0),
					leaf:     true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New[int](tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTree_Insert(t1 *testing.T) {
	type args[V constraints.Ordered] struct {
		k V
	}
	type testCase[V constraints.Ordered] struct {
		name string
		t    Tree[V]
		args args[V]
		want Tree[V]
	}
	tests := []testCase[string]{
		{
			name: "insert leaf in the end",
			t: Tree[string]{
				t: 3,
				root: &node[string]{
					n:        3,
					keys:     []string{"A", "B", "C", "", ""},
					children: []*node[string]{},
					leaf:     true,
				},
			},
			args: args[string]{k: "D"},
			want: Tree[string]{
				t: 3,
				root: &node[string]{
					n:        4,
					keys:     []string{"A", "B", "C", "D", ""},
					children: []*node[string]{},
					leaf:     true,
				},
			},
		},
		{
			name: "insert leaf in the beginning",
			t: Tree[string]{
				t: 3,
				root: &node[string]{
					n:        3,
					keys:     []string{"B", "C", "D", "", ""},
					children: []*node[string]{},
					leaf:     true,
				},
			},
			args: args[string]{k: "A"},
			want: Tree[string]{
				t: 3,
				root: &node[string]{
					n:        4,
					keys:     []string{"A", "B", "C", "D", ""},
					children: []*node[string]{},
					leaf:     true,
				},
			},
		},
		{
			name: "insert leaf in the middle",
			t: Tree[string]{
				t: 3,
				root: &node[string]{
					n:        3,
					keys:     []string{"A", "C", "D", "", ""},
					children: []*node[string]{},
					leaf:     true,
				},
			},
			args: args[string]{k: "B"},
			want: Tree[string]{
				t: 3,
				root: &node[string]{
					n:        4,
					keys:     []string{"A", "B", "C", "D", ""},
					children: []*node[string]{},
					leaf:     true,
				},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			tt.t.Insert(tt.args.k)
			if !reflect.DeepEqual(tt.t, tt.want) {
				t1.Errorf("Insert() = %#+v, want %#+v", tt.t, tt.want)
			}
		})
	}
}

func TestTree_Insert_to_full_node(t1 *testing.T) {
	t := getTree([]string{"A", "B", "D", "E", "F"}, 3)
	validTree := []validNode[string]{
		{node: t.root, keys: []string{"A", "B", "D", "E", "F"}, leaf: true, n: 5},
	}

	// check tree's structure before insert
	for _, n := range validTree {
		checkNode(t1, &n)
	}

	t.Insert("C")

	// check tree's structure after insert
	validTreeAfterInsert := []validNode[string]{
		{node: t.root, keys: []string{"D"}, leaf: false, n: 1},
		{node: t.root.children[0], keys: []string{"A", "B", "C"}, leaf: true, n: 3},
		{node: t.root.children[1], keys: []string{"E", "F"}, leaf: true, n: 2},
	}
	for _, n := range validTreeAfterInsert {
		checkNode(t1, &n)
	}
}

func TestTree_Insert_split_leaf(t1 *testing.T) {
	t := getTree([]string{"A", "B", "D", "E", "F", "C"}, 3)
	validTree := []validNode[string]{
		{node: t.root, keys: []string{"D"}, leaf: false, n: 1},
		{node: t.root.children[0], keys: []string{"A", "B", "C"}, leaf: true, n: 3},
		{node: t.root.children[1], keys: []string{"E", "F"}, leaf: true, n: 2},
	}

	// check tree's structure before insert
	for _, n := range validTree {
		checkNode(t1, &n)
	}

	t.Insert("G")
	t.Insert("K")
	t.Insert("M")

	// check tree's structure after insert
	validTreeAfterInsert := []validNode[string]{
		{node: t.root, keys: []string{"D"}, leaf: false, n: 1},
		{node: t.root.children[0], keys: []string{"A", "B", "C"}, leaf: true, n: 3},
		{node: t.root.children[1], keys: []string{"E", "F", "G", "K", "M"}, leaf: true, n: 5},
	}
	for _, n := range validTreeAfterInsert {
		checkNode(t1, &n)
	}

	t.Insert("S")
	// check tree's structure after insert
	validTreeAfterInsert2 := []validNode[string]{
		{node: t.root, keys: []string{"D", "G"}, leaf: false, n: 2},
		{node: t.root.children[0], keys: []string{"A", "B", "C"}, leaf: true, n: 3},
		{node: t.root.children[1], keys: []string{"E", "F"}, leaf: true, n: 2},
		{node: t.root.children[2], keys: []string{"K", "M", "S"}, leaf: true, n: 3},
	}
	for _, n := range validTreeAfterInsert2 {
		checkNode(t1, &n)
	}
}

func TestTree_Exists(t1 *testing.T) {
	type args[V constraints.Ordered] struct {
		key V
	}
	type testCase[V constraints.Ordered] struct {
		name string
		t    *Tree[V]
		args args[V]
		want bool
	}
	tests := []testCase[string]{
		{
			name: "empty tree",
			t:    getTree([]string{}, 3),
			args: args[string]{key: "A"},
			want: false,
		},
		{
			name: "tree with one element - not found",
			t:    getTree([]string{"A"}, 3),
			args: args[string]{key: "B"},
			want: false,
		},
		{
			name: "tree with one element - found",
			t:    getTree([]string{"A"}, 3),
			args: args[string]{key: "A"},
			want: true,
		},
		{
			name: "tree with several elements in root - found",
			t:    getTree([]string{"A", "B", "C", "D"}, 3),
			args: args[string]{key: "C"},
			want: true,
		},
		{
			name: "tree with root and one child - not found",
			t:    getTree([]string{"A", "B", "D", "E", "F", "C"}, 3),
			args: args[string]{key: "K"},
			want: false,
		},
		{
			name: "tree with root and one child - found",
			t:    getTree([]string{"A", "B", "D", "E", "F", "C"}, 3),
			args: args[string]{key: "F"},
			want: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if got := tt.t.Exists(tt.args.key); got != tt.want {
				t1.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTree_Delete(t1 *testing.T) {
	type args[V constraints.Ordered] struct {
		k V
	}
	type testCase[V constraints.Ordered] struct {
		name           string
		t              *Tree[V]
		args           args[V]
		wantErr        bool
		tAfterDeleting *Tree[V]
	}
	tests := []testCase[string]{
		{
			name:    "empty tree",
			t:       getTree([]string{}, 3),
			args:    args[string]{k: "A"},
			wantErr: true,
		},
		{
			name:    "tree with one element - not found",
			t:       getTree([]string{"A", "B", "C"}, 3),
			args:    args[string]{k: "D"},
			wantErr: true,
		},
		{
			name:           "delete first leaf",
			t:              getTree([]string{"A", "B", "C"}, 3),
			args:           args[string]{k: "A"},
			wantErr:        false,
			tAfterDeleting: getTree([]string{"B", "C"}, 3),
		},
		{
			name:           "delete last leaf",
			t:              getTree([]string{"A", "B", "C"}, 3),
			args:           args[string]{k: "C"},
			wantErr:        false,
			tAfterDeleting: getTree([]string{"A", "B"}, 3),
		},
		{
			name:           "delete middle leaf",
			t:              getTree([]string{"A", "B", "C"}, 3),
			args:           args[string]{k: "B"},
			wantErr:        false,
			tAfterDeleting: getTree([]string{"A", "C"}, 3),
		},
		{
			name:           "delete key node (left child)",
			t:              getTree([]string{"A", "B", "D", "E", "F", "C", "G", "K", "M"}, 3),
			args:           args[string]{k: "D"},
			wantErr:        false,
			tAfterDeleting: getTree([]string{"A", "B", "E", "F", "C", "G", "K", "M"}, 3),
		},
		{
			name:           "delete key node (right child)",
			t:              getTree([]string{"A", "B", "D", "E", "F", "C", "G", "K", "M", "N", "O"}, 3),
			args:           args[string]{k: "G"},
			wantErr:        false,
			tAfterDeleting: getTree([]string{"A", "B", "D", "E", "F", "C", "K", "M", "N", "O"}, 3),
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if err := tt.t.Delete(tt.args.k); (err != nil) != tt.wantErr {
				t1.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.t, tt.tAfterDeleting) {
				t1.Errorf("Delete() got = %v, want %v", tt.t, tt.tAfterDeleting)
			}
		})
	}
}

func TestTree_Delete_get_key_from_left_child(t1 *testing.T) {
	t := getTree([]string{"A", "B", "D", "E", "F", "C", "G", "K", "M"}, 3)
	validTree := []validNode[string]{
		{node: t.root, keys: []string{"D"}, leaf: false, n: 1},
		{node: t.root.children[0], keys: []string{"A", "B", "C"}, leaf: true, n: 3},
		{node: t.root.children[1], keys: []string{"E", "F", "G", "K", "M"}, leaf: true, n: 5},
	}

	// check tree's structure before insert
	for _, n := range validTree {
		checkNode(t1, &n)
	}

	t.Delete("D")

	// check tree's structure after delete
	validTreeAfterInsert := []validNode[string]{
		{node: t.root, keys: []string{"C"}, leaf: false, n: 1},
		{node: t.root.children[0], keys: []string{"A", "B"}, leaf: true, n: 2},
		{node: t.root.children[1], keys: []string{"E", "F", "G", "K", "M"}, leaf: true, n: 5},
	}
	for _, n := range validTreeAfterInsert {
		checkNode(t1, &n)
	}
}

func TestTree_Delete_get_key_from_right_child(t1 *testing.T) {
	t := getTree([]string{"A", "B", "D", "E", "F", "C", "G", "K", "M", "N", "O"}, 3)
	validTree := []validNode[string]{
		{node: t.root, keys: []string{"D", "G"}, leaf: false, n: 2},
		{node: t.root.children[0], keys: []string{"A", "B", "C"}, leaf: true, n: 3},
		{node: t.root.children[1], keys: []string{"E", "F"}, leaf: true, n: 2},
		{node: t.root.children[2], keys: []string{"K", "M", "N", "O"}, leaf: true, n: 4},
	}

	// check tree's structure before insert
	for _, n := range validTree {
		checkNode(t1, &n)
	}

	t.Delete("G")

	// check tree's structure after delete
	validTreeAfterInsert := []validNode[string]{
		{node: t.root, keys: []string{"D", "K"}, leaf: false, n: 2},
		{node: t.root.children[0], keys: []string{"A", "B", "C"}, leaf: true, n: 3},
		{node: t.root.children[1], keys: []string{"E", "F"}, leaf: true, n: 2},
		{node: t.root.children[2], keys: []string{"M", "N", "O"}, leaf: true, n: 3},
	}
	for _, n := range validTreeAfterInsert {
		checkNode(t1, &n)
	}
}

func TestTree_Delete_merge_node(t1 *testing.T) {
	t := getTree([]string{"A", "B", "D", "E", "F", "C"}, 3)
	validTree := []validNode[string]{
		{node: t.root, keys: []string{"D"}, leaf: false, n: 1},
		{node: t.root.children[0], keys: []string{"A", "B", "C"}, leaf: true, n: 3},
		{node: t.root.children[1], keys: []string{"E", "F"}, leaf: true, n: 2},
	}

	// check tree's structure before insert
	for _, n := range validTree {
		checkNode(t1, &n)
	}

	t.Delete("D")

	// check tree's structure after delete
	validTreeAfterInsert := []validNode[string]{
		{node: t.root, keys: []string{"C"}, leaf: false, n: 1},
		{node: t.root.children[0], keys: []string{"A", "B"}, leaf: true, n: 2},
		{node: t.root.children[1], keys: []string{"E", "F"}, leaf: true, n: 2},
	}
	for _, n := range validTreeAfterInsert {
		checkNode(t1, &n)
	}

	t.Delete("C")

	// check tree's structure after delete
	validTreeAfterSecondInsert := []validNode[string]{
		{node: t.root, keys: []string{"A", "B", "E", "F"}, leaf: true, n: 4},
	}
	for _, n := range validTreeAfterSecondInsert {
		checkNode(t1, &n)
	}
}

type validNode[V constraints.Ordered] struct {
	node *node[V]
	n    int
	keys []V
	leaf bool
}

func checkNode[V constraints.Ordered](t *testing.T, vn *validNode[V]) {
	if vn == nil {
		return
	}

	if len(vn.keys) != len(vn.node.keys) {
		t.Errorf("Error in len - Want keys: %v, have keys %v",
			vn.keys,
			vn.node.keys,
		)
	}

	for i := 0; i < len(vn.keys); i++ {
		if vn.keys[i] != vn.node.keys[i] {
			t.Errorf("Error - Want keys: %v, have keys %v",
				vn.keys,
				vn.node.keys,
			)
		}
	}

	if vn.node.leaf != vn.leaf {
		t.Errorf("Error - Want leaf: %v, have leaf: %v",
			vn.leaf,
			vn.node.leaf,
		)
	}

	if vn.node.n != vn.n {
		t.Errorf("Error - Want n: %v, have n: %v",
			vn.n,
			vn.node.n,
		)
	}
}

func getTree(elements []string, t int) *Tree[string] {
	tree, _ := New[string](t)
	for _, el := range elements {
		tree.Insert(el)
	}

	return tree
}
