# Very Simple - Immutable Hash Map

- branch factor of 64, giving support for up to 68,719,476,736 entries
- Immutable. Each operation on the hash map returns a new version.
- Memory efficient we keep references to the old versions of the HMT, the majority of the trie can stay un-changed.

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