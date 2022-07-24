package hmt

import (
	"hash/maphash"
)

type HMT[T any] struct {
	s        maphash.Seed
	previous *HMT[T]
	root     *Trie[T]
}

func New[T any]() *HMT[T] {
	return &HMT[T]{
		s:    maphash.MakeSeed(),
		root: newTrie[T](),
	}
}

func (h *HMT[T]) Get(k Key) (*Entry[T], error) {
	v, err := hashCode(h.s, k)
	if err != nil {
		return nil, err
	}
	return retrieve[T](h.root, v), nil
}
func (h *HMT[T]) Set(k Key, v T) (*HMT[T], error) {
	vKey, err := hashCode(h.s, k)
	if err != nil {
		return h, err
	}
	newRoot := insertTrie[T](
		h.root,
		vKey,
		Entry[T]{
			Key:   k,
			Value: v,
		})
	newHmt := &HMT[T]{
		s:        h.s,
		previous: h,
		root:     newRoot,
	}
	return newHmt, nil
}

func (h *HMT[T]) Del(k Key) (*HMT[T], error) {
	hashkey, err := hashCode(h.s, k)
	if err != nil {
		return nil, err
	}
	newRoot := delete(h.root, hashkey)
	if newRoot == nil {
		return h, nil
	}

	return &HMT[T]{
		s:        h.s,
		previous: h,
		root:     newRoot,
	}, nil
}

func (h *HMT[T]) Entries() []*Entry[T] {
	return h.root.ChildEntries()
}

func (h *HMT[T]) Generations() (res []*HMT[T]) {
	res = append(res, h)
	if h.previous != nil {
		res = append(res, h.previous.Generations()...)
	}
	return
}

func (h *HMT[T]) Chain() *ChainHMT[T] {
	return &ChainHMT[T]{
		err: nil,
		ht:  h,
	}
}

type Key []byte

type Entry[T any] struct {
	Key   Key
	Value T
}

type Trie[T any] struct {
	key         uint64
	value       *Entry[T]
	connections [64]*Trie[T]
}

func (t Trie[T]) ChildEntries() []*Entry[T] {
	response := make([]*Entry[T], 0)
	if t.value == nil {
		return response
	} else {
		response = append(response, t.value)
	}
	for _, childTrie := range t.connections {
		if childTrie != nil {
			response = append(response, childTrie.ChildEntries()...)
		}
	}
	return response
}

func newTrie[T any]() *Trie[T] {
	itr := &Trie[T]{
		key:         0,
		value:       nil,
		connections: [64]*Trie[T]{},
	}
	return itr
}

func hashCode(s maphash.Seed, bs Key) (uint64, error) {
	var h maphash.Hash
	h.SetSeed(s)
	_, err := h.Write(bs)
	if err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}

// insertTrie will mutate the path of visited nodes, creating copies on the way to finding the right path
// returns the new and improved root *tr
func insertTrie[T any](tr *Trie[T], valueKey uint64, v Entry[T]) *Trie[T] {
	ret := copyTrie(tr)
	if ret.value == nil {
		ret.key = valueKey
		ret.value = &v
		return ret
	}
	if tr.key == valueKey {
		ret.value = &v
		return ret
	}
	// go deeper
	idx := lookupPath(valueKey)
	if ret.connections[idx] == nil {
		ret.connections[idx] = newTrie[T]()
	}
	ret.connections[idx] = insertTrie(ret.connections[idx], incrementPath(valueKey), v)
	return ret
}
func copyTrie[T any](tr *Trie[T]) *Trie[T] {
	return &Trie[T]{
		key:         tr.key,
		value:       tr.value,
		connections: tr.connections, // connections are still the same shared subtree
	}
}

// should always return a value between 0 & 63 to give us the index to lookup in the trie.
func lookupPath(val uint64) uint64 {
	return val & 63
}
func incrementPath(val uint64) uint64 {
	return val >> 6
}

func retrieve[T any](tr *Trie[T], valueKey uint64) *Entry[T] {
	if valueKey == tr.key {
		return tr.value
	}
	idx := lookupPath(valueKey)
	child := tr.connections[idx]
	if child == nil {
		return nil
	}
	return retrieve(child, incrementPath(valueKey))
}

// will immutably delete the child, creating copies of the nodes that it visits. if it's never found, return nil..
func delete[T any](tr *Trie[T], valueKey uint64) *Trie[T] {
	ntr := copyTrie(tr)
	if valueKey == tr.key {
		ntr.value = nil
		return ntr
	}
	idx := lookupPath(valueKey)
	newValueKey := incrementPath(valueKey)
	child := ntr.connections[idx]
	// if the child does not exist, then return nil. We should not copy nodes if the tree is not mutated (deleting a non-existing key)
	if child == nil {
		return nil
	}
	result := delete(child, newValueKey)
	if result == nil {
		return tr // return the original trie if no child for deletion is found..
	}
	// update the subtree
	ntr.connections[idx] = result
	return ntr
}
