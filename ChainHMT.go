package hmt

type ChainHMT[T any] struct {
	err error
	ht  *HMT[T]
}

func (ch *ChainHMT[T]) Get(k Key) (res *Entry[T]) {
	if ch.err != nil {
		return nil
	}
	res, ch.err = ch.ht.Get(k)
	return res
}
func (ch *ChainHMT[T]) Set(k Key, val T) *ChainHMT[T] {
	if ch.err != nil {
		return ch
	}
	ch.ht, ch.err = ch.ht.Set(k, val)
	return ch
}

func (ch *ChainHMT[T]) Del(k Key) *ChainHMT[T] {
	if ch.err != nil {
		return nil
	}
	ch.ht, ch.err = ch.ht.Del(k)
	return ch
}

func (ch *ChainHMT[T]) Entries() []*Entry[T] {
	return ch.ht.Entries()
}
func (ch *ChainHMT[T]) Error() error {
	return ch.err
}

func (ch *ChainHMT[T]) AsHMT() (*HMT[T], error) {
	return ch.ht, ch.err
}
