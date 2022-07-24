# Very Simple - Immutable Hash Map

- branch factor of 64, giving support for up to 68,719,476,736 entries
- Immutable. Insert and delete operations result in a new version of the HMT.
- Memory efficient we keep references to the old versions of the Trie. The majority of keys and values will be unchanged.

Using clojure and writing my own small lisp interpreter has inspired me to look at this data structure.
Combined with a current sabbatical, now seems like the perfect time to write and publish something like this to keep my skills share with the new Go-1.18 generics.

**Usage**

```go
package main

import (
	"log"
	"github.com/ebuckley/hmt"
)

func main() {
    // There is a convienient wrapper around HMT's ChainHMT that gives you a fluent API (if you choose)
    points := hmt.New[int]().Chain()

    points.
        Set(hmt.Key("ersin"), 1).
        Set(hmt.Key("emily"), 1337).
        Set(hmt.Key("vishi"), 99).
        Set(hmt.Key("vladamir"), 69)

    if points.Error() != nil {
        log.Fatalln("something went wrong with the API")
    }
    entries := points.Entries()
    if len(entries) != 4 {
        log.Fatalln("should have 4 entries")
    }
}

```

## Read more

- Enjoyable blog about the basic approach of using a hash at different depths to find the https://worace.works/2016/05/24/hash-array-mapped-tries/
- The guy wo invented this (Phillip Bagwell) http://lampwww.epfl.ch/papers/idealhashtrees.pdf
- A cool project to follow along with if you are interested in immutable databases http://aosabook.org/en/500L/pages/dbdb-dog-bed-database.html