package backend

import (
	"sort"

	"github.com/digitalrebar/digitalrebar/go/common/store"
)

// Index declares a struct field that can be indexed for a given
// Model, along with the function that should be used to sort things
// in key order.  Unless otherwise indicated, any methods that act on an Index
// return a new Index with its own reference to indexed items.
type Index struct {
	Key    string
	sorted bool
	less   func(store.KeySaver, store.KeySaver) bool
	gtRef  func(store.KeySaver) bool
	gteRef func(store.KeySaver) bool
	objs   []store.KeySaver
}

func NewIndex(key string, less func(store.KeySaver, store.KeySaver) bool) *Index {
	return &Index{Key: key, less: less}
}

// Sort sorts the values the Index references in their native order in
// a stable fashion. It operates in-place and returns the same index.
func (i *Index) Sort() *Index {
	less := func(j, k int) bool { return i.less(i.objs[j], i.objs[k]) }
	if !i.sorted {
		sort.SliceStable(i.objs, less)
		i.sorted = true
	}
	return i
}

// cp makes a copy of the current Index with a copy of the passed-in
// slice of objects that the index should reference.
func (i *Index) cp(newObjs []store.KeySaver) *Index {
	objs := make([]store.KeySaver, len(newObjs))
	copy(objs, newObjs)
	return &Index{
		Key:    i.Key,
		sorted: i.sorted,
		less:   i.less,
		gtRef:  i.gtRef,
		gteRef: i.gteRef,
		objs:   objs,
	}
}

// Subset causes the index to discard all elements that fall outside
// the first index for which lower returns true and the first index
// for which upper returns true.  The index must be sorted first, or
// Subset will panic.
//
// lower and upper should be thunks that examine the object passed to
// determine where the subset should start and end at, and must choose
// items based on what the index is currently sorted by.
//
// Subset returns a new *Index.
func (i *Index) Subset(lower, upper func(store.KeySaver) bool) *Index {
	if !i.sorted {
		panic("Cannot take subset of unsorted Index")
	}
	totalCount := len(i.objs)
	start := sort.Search(totalCount, func(j int) bool { return lower(i.objs[j]) })
	var objs []store.KeySaver
	if start == totalCount {
		objs = []store.KeySaver{}
	} else {
		objs = i.objs[start:sort.Search(totalCount, func(j int) bool { return upper(i.objs[j]) })]
	}
	return i.cp(objs)
}

// A couple of helper functions for using Subset to implement the
// usual comparison operators
func alwaysFalse(store.KeySaver) bool { return false }
func alwaysTrue(store.KeySaver) bool  { return true }

// SetComparators allows you to add static comparators that test for
// greater-than and greater-than-or-equal against values stored in the
// Index.  The comaprators must test the field that the index is
// sorted by, or you will get nonsensical results for the comparison
// operators.
//
// Once comparators are set for an index, the usual comparison methods
// (Lt, Eq, Gt, etc.)  will work without panicing, and will be
// inherited by any Indexes created from this one.
//
// SetComparators retuns a new index.
func (i *Index) SetComparators(gt, gte func(store.KeySaver) bool) *Index {
	res := i.cp(i.objs)
	res.gtRef = gt
	res.gteRef = gte
	return res
}

// Lt returns an index with all items that are less than the current
// comparators
func (i *Index) Lt() *Index {
	return i.Subset(alwaysFalse, i.gteRef)
}

// Lte returns an index with all the items that are less than or equal
// to the current comparators.
func (i *Index) Lte() *Index {
	return i.Subset(alwaysFalse, i.gtRef)
}

// Eq returns an index with all the items that are equal to the
// current comparators.
func (i *Index) Eq() *Index {
	return i.Subset(i.gteRef, i.gtRef)
}

// Gte returns an index with all the items that are greater than or
// equal to the current comparators
func (i *Index) Gte() *Index {
	return i.Subset(i.gteRef, alwaysTrue)
}

// Gt returns an index with all the items that are greater-than the
// current comparators
func (i *Index) Gt() *Index {
	return i.Subset(i.gtRef, alwaysTrue)
}

// Ne returns an index with all the items that are not equal to the current
// comparators.
func (i *Index) Ne() *Index {
	lt := i.Lt()
	gt := i.Gt()
	objs := lt.objs
	objs = append(objs, gt.objs...)
	return i.cp(objs)
}

// Filter returns an index with all the items that do not match the
// passed filter stripped out.  Unlike the comparators based on
// Subset, it is indifferent to the sort order of the index and always
// looks at all the objects in the index.
//
// Filter returns a new *Index
func (i *Index) Filter(f func(store.KeySaver) bool) *Index {
	objs := []store.KeySaver{}
	for _, obj := range i.objs {
		if f(obj) {
			objs = append(objs, obj)
		}
	}
	return i.cp(objs)
}

// Offset returns a new *Index with the first n things removed.
func (i *Index) Offset(num int) *Index {
	return i.cp(i.objs[num:])
}

// Limit returns a new *Index with all but the first n things removed.
func (i *Index) Limit(num int) *Index {
	return i.cp(i.objs[:num])
}

// Chain returns a new index based on other but with our items.
func (i *Index) Chain(other *Index) *Index {
	res := other.cp(i.objs)
	res.sorted = false
	return res
}

// Reverse performs an in-place reversal of the objects contained in
// the index.  It returns the current Index.
func (i *Index) Reverse() *Index {
	for lower, upper := 0, len(i.objs)-1; lower < upper; lower, upper = lower+1, upper-1 {
		i.objs[lower], i.objs[upper] = i.objs[upper], i.objs[lower]
	}
	i.sorted = false
	return i
}

// Items returns the current items the index has filtered.
func (i *Index) Items() []store.KeySaver {
	res := make([]store.KeySaver, len(i.objs))
	copy(res, i.objs)
	return res
}

// SetItems returns a new *Index with its items set to s.
// The new Index is not sorted.
func (i *Index) SetItems(s []store.KeySaver) *Index {
	res := i.cp(s)
	res.sorted = false
	return res
}
