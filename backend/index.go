package backend

import (
	"sort"

	"github.com/digitalrebar/digitalrebar/go/common/store"
)

// Index declares a struct field that can be indexed for a given
// Model, along with the function that should be used to sort things
// in key order.
type Index struct {
	Key  string
	objs []store.KeySaver
	less func(store.KeySaver, store.KeySaver) bool
}

// Sort sorts the values the Index references in their native order.
// This implements a stable sort.
func (i *Index) Sort() *Index {
	less := func(j, k int) bool { return i.less(i.objs[j], i.objs[k]) }
	sort.SliceStable(i.objs, less)
	return i
}

// Subset causes the index to discard all elements that fall
// outside the first index for which lower returns true and the
// first index for which upper returns true.
//
// lower and upper should be thunks that examine the object passed
// to determine where the subset should start and end at.
//
func (i *Index) Subset(lower, upper func(store.KeySaver) bool) *Index {
	totalCount := len(i.objs)
	start := sort.Search(totalCount, func(j int) bool { return lower(i.objs[j]) })
	if start == totalCount {
		i.objs = []store.KeySaver{}
	} else {
		i.objs = i.objs[start:sort.Search(totalCount, func(j int) bool { return upper(i.objs[j]) })]
	}
	return i
}

// Offset removes the first n objects from the index
func (i *Index) Offset(num int) *Index {
	i.objs = i.objs[num:]
	return i
}

// Limit discards all but the first n objects from the index
func (i *Index) Limit(num int) *Index {
	i.objs = i.objs[:num]
	return i
}

// Chain passes the index's current items to the other index
// The other index will sort the items and return them.
func (i *Index) Chain(other *Index) *Index {
	other.objs = i.objs
	return other.Sort()
}

// Items returns the current items the index has filtered.
func (i *Index) Items() []store.KeySaver {
	res := make([]store.KeySaver, len(i.objs))
	copy(res, i.objs)
	return res
}
