# Very Simple - Immutable Hash Array Mapped Tree

This package implements an easy to use immutable Hash Array Mapped Tree. Using clojure and writing my own small lisp interpreter has inspired me to look at this data structure.
Combined with a current sabbatical, now seems like the perfect time to write and publish something like this to improve my skills with the new Go-1.18 generics.

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
        Set(Key("ersin"), 1).
        Set(Key("emily"), 1337).
        Set(Key("vishi"), 99).
        Set(Key("vladamir"), 69)

    if points.Error() != nil {
        log.Fatalln("something went wrong with the API")
    }
    entries := points.Entries()
    if len(entries) != 4 {
        log.Fatalln("should have 4 entries")
    }
}
```