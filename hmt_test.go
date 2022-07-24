package hmt

import (
	"crypto/rand"
	"hash/maphash"
	"log"
	"math/big"
	"testing"
)

func makeUID(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		b, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			log.Fatalln("Fatal error making a secure unique ID:", err)
		}
		s[i] = letters[b.Int64()]
	}
	return string(s)
}

func TestHashAdvancing(t *testing.T) {

	var code uint64 = 17735010357689101119
	t.Log("Value is ", code)
	currIdx := lookupPath(code)
	if currIdx != 63 {
		t.Fatal("first idx should be 63, got:", currIdx)
	}
	nextCode := incrementPath(code)
	nextIdx := lookupPath(nextCode)
	if nextIdx != 44 {
		t.Fatal("next chunk of bytes should be 44, got: ", nextIdx)
	}
}

func TestHasher(t *testing.T) {
	s := maphash.MakeSeed()
	testVal := Key("Some Value of anything")
	code, err := hashCode(s, testVal)
	if err != nil {
		t.Fatal("Should not error:", err)
	}

	code2, err := hashCode(s, testVal)
	if err != nil {
		t.Fatal("Should not error:", err)
	}
	if code != code2 {
		t.Fatal("Hashcodes should be consistent for one maphash \ncode1 and code2 are:", code, code2)
	}
}

func TestBasicTrieOps(t *testing.T) {
	var s = maphash.MakeSeed()

	expectedValue := "Hello Friend"
	entry := Entry[string]{
		Key:   Key("message"),
		Value: expectedValue,
	}
	key, err := hashCode(s, entry.Key)
	if err != nil {
		t.Fatal(err)
	}

	root := newTrie[string]()
	nextGen := insertTrie[string](root, key, entry)

	found := retrieve[string](nextGen, key)
	if found == nil {
		t.Fatal("Should have found the value but got nil")
	}

	if found.Value != expectedValue {
		t.Fatal("should be the same but values are differrent \nfound:", found, "expect:", expectedValue)
	}
}
func BenchmarkPublicInterface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nt := New[string]()

		doubleUp := make(map[string]string)
		for i := 0; i < 1000; i++ {
			key := Key(makeUID(32))
			value := makeUID(32)
			var err error
			nt, err = nt.Set(key, value)
			if err != nil {
				b.Fatal("Bad insert")
			}
			doubleUp[string(key)] = value
		}

		for k, v := range doubleUp {
			val, err := nt.Get(Key(k))
			if err != nil {
				b.Fatal("should not error but got", err)
			}
			if val == nil {
				b.Fatal("Should return a value and not be nil")
			}
			if val.Value != v {
				b.Fatal("Expected to get the value we set in the duplicate pure go map..")
			}
		}
		generations := 0
		genHt := nt
		for genHt.previous != nil {
			genHt = genHt.previous
			generations++
		}
		if generations != 1000 {
			b.Fatal("there should be 1000 generations of the immutable hash table but found: ", generations)
		}
	}
}
func TestPublicInterface(t *testing.T) {
	nt := New[string]()

	doubleUp := make(map[string]string)
	for i := 0; i < 1000; i++ {
		key := Key(makeUID(32))
		value := makeUID(32)
		var err error
		nt, err = nt.Set(key, value)
		if err != nil {
			t.Fatal("Bad insert")
		}
		doubleUp[string(key)] = value
	}

	for k, v := range doubleUp {
		val, err := nt.Get(Key(k))
		if err != nil {
			log.Fatalln(err)
		}
		if val == nil {
			t.Fatal("Should return a value and not be nil")
		}
		if val.Value != v {
			t.Fatal("Expected to get the value we set in the duplicate pure go map..")
		}
	}
	generations := 0
	genHt := nt
	for genHt.previous != nil {
		genHt = genHt.previous
		generations++
	}
	if generations != 1000 {
		t.Fatal("there should be 1000 generations of the immutable hash table but found: ", generations)
	}
}
func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {

	}
}
func TestStoringComplexTypes(t *testing.T) {
	type Account struct {
		Type   int
		Amount int64
	}
	ht := New[Account]()

	var err error
	ht, err = ht.Set(Key("ebuckley"), Account{
		Type:   1,
		Amount: 89889,
	})
	if err != nil {
		log.Fatalln("should not error")
	}

	ht, err = ht.Set(Key("jbezos"), Account{
		Type:   2,
		Amount: 1337,
	})
	if err != nil {
		log.Fatalln("do not error either")
	}

	v, err := ht.Get(Key("ebuckley"))
	if v.Value.Type != 1 {
		log.Fatalln("should be type 1")
	}
	if v.Value.Amount != 89889 {
		log.Fatalln("should have balance 89889")
	}

	vals := ht.Entries()
	if len(vals) != 2 {
		log.Fatalln("Should have 2 values in the hmt")
	}
}

func TestDelete(t *testing.T) {

	var err error
	type Account struct {
		Type   int
		Amount int64
	}
	ht := New[Account]()

	ht, err = ht.Set(Key("ebuckley"), Account{
		Type:   1,
		Amount: 89889,
	})
	if err != nil {
		t.Fatal("should not error but got", err)
	}
	ht, err = ht.Del(Key("ebuckley"))
	if err != nil {
		t.Fatal("should not error but got", err)
	}

	ent, err := ht.Get(Key("ebuckley"))
	if err != nil {
		t.Fatal("should not error but got", err)
	}
	if ent != nil {
		t.Fatal("The entry should be nil after being deleted")
	}
}

func TestChainableInterface(t *testing.T) {
	points := New[int]().Chain()

	points.
		Set(Key("ersin"), 1).
		Set(Key("emily"), 1337).
		Set(Key("vishi"), 99).
		Set(Key("vladamir"), 69)

	if points.Error() != nil {
		t.Fatal("Should not error out")
	}
	entries := points.Entries()
	if len(entries) != 4 {
		t.Fatal("expected 4 entries")
	}
}
