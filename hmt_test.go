package hmt

import (
    "crypto/rand"
    "hash/maphash"
    "log"
    "math/big"
    "reflect"
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
    testVal := Value("Some Value of anything")
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

    expectedValue := Value("Hello Friend")
    key, err := hashCode(s, Value("message"))
    if err != nil {
        t.Fatal(err)
    }

    root := newTrie()
    insertTrie(root, key, expectedValue)

    found := retrieve(root, key)
    if found == nil {
        t.Fatal("Should have found the value but got nil")
    }
    if !reflect.DeepEqual(*found, expectedValue) {
        t.Fatal("should be the same but values are differrent \nfound:", found, "expect:", expectedValue)
    }
}

func TestPublicInterface(t *testing.T) {
    nt := New()

    doubleUp := make(map[string]Value)
    for i := 0; i < 1000; i++ {
        key := Key(makeUID(32))
        value := Value(makeUID(32))
        _, err := nt.Set(key, value)
        if err != nil {
            t.Fatal("Bad insert")
        }
        doubleUp[string(key)] = value
    }

    for k, v := range doubleUp {
        val, err := nt.Get(Key(k))
        if err != nil {
            return
        }
        if string(*val) != string(v) {
            t.Fatal("Expected to get the value we set in the duplicate pure go map..")
        }
    }

}