package hmt

import (
	"hash/maphash"
)

type HashArrayMappedTrie[T any] struct {
	s maphash.Seed
	// TODO implement collection of historical versions of this HashArrayMapped Trie
	// will contains an ordered reference of root tree pointers
	generations []*Trie[T]
	root        *Trie[T]
}

func New[T any]() *HashArrayMappedTrie[T] {
	return &HashArrayMappedTrie[T]{
		s:    maphash.MakeSeed(),
		root: newTrie[T](),
	}
}

func (h *HashArrayMappedTrie[T]) Get(k Key) (*Entry[T], error) {
	v, err := hashCode(h.s, k)
	if err != nil {
		return nil, err
	}
	return retrieve[T](h.root, v), nil
}
func (h *HashArrayMappedTrie[T]) Set(k Key, v T) (*HashArrayMappedTrie[T], error) {
	vKey, err := hashCode(h.s, k)
	if err != nil {
		return h, err
	}
	insertTrie[T](
		h.root,
		vKey,
		Entry[T]{
			Key:   k,
			Value: v,
		})
	return h, nil
}

func (h *HashArrayMappedTrie[T]) Entries() []*Entry[T] {
	return h.root.ChildEntries()
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

func insertTrie[T any](tr *Trie[T], valueKey uint64, v Entry[T]) *Trie[T] {

	if tr.value == nil {
		tr.key = valueKey
		tr.value = &v
		return tr
	} else if tr.key == valueKey {
		tr.value = &v
		return tr
	}
	// go deeper
	idx := lookupPath(valueKey)
	if tr.connections[idx] == nil {
		tr.connections[idx] = newTrie[T]()
	}
	return insertTrie(tr.connections[idx], incrementPath(valueKey), v)
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
