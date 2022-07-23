package hmt

import (
	"hash/maphash"
)

type HashArrayMappedTrie struct {
	s           maphash.Seed
	generations []*Trie // TODO implement collection of historical versions of this HashArrayMapped Trie
	root        *Trie
}

func New() *HashArrayMappedTrie {
	return &HashArrayMappedTrie{
		s:    maphash.MakeSeed(),
		root: newTrie(),
	}
}

func (h *HashArrayMappedTrie) Get(k Key) (*Value, error) {
	v, err := hashCode(h.s, Value(k))
	if err != nil {
		return nil, err
	}
	return retrieve(h.root, v), nil
}
func (h *HashArrayMappedTrie) Set(k Key, v Value) (*HashArrayMappedTrie, error) {
	vKey, err := hashCode(h.s, Value(k))
	if err != nil {
		return h, err
	}
	insertTrie(h.root, vKey, v)
	return h, nil
}

type Key []byte
type Value []byte

type Trie struct {
	key         uint64
	value       *Value
	connections [64]*Trie
}

func newTrie() *Trie {
	itr := &Trie{
		key:         0,
		value:       nil,
		connections: [64]*Trie{},
	}
	return itr
}

func hashCode(s maphash.Seed, bs Value) (uint64, error) {
	var h maphash.Hash
	h.SetSeed(s)
	_, err := h.Write(bs)
	if err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}

func insertTrie(tr *Trie, valueKey uint64, v Value) *Trie {

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
		tr.connections[idx] = newTrie()
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

func retrieve(tr *Trie, valueKey uint64) *Value {
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
