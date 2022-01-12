package lock_test

import (
	"fmt"
	"github.com/bitstonks/leftright/pkg/lock"
)

type store map[string]string
func newStore() *store {
	s := make(store)
	return &s
}

type keyVal struct {
	key string
	val string
}

func (s *store) Update(op lock.Operation) lock.OpResult {
	inp := op.(keyVal)
	(*s)[inp.key] = inp.val
	return inp.key
}

func Example_simple() {
	// Create a left right lock with two empty maps
	lock := lock.NewLeftRightLock(newStore(), newStore())

	// Write a value to one side, which cannot be read yet.
	lock.Write(keyVal{"test", "value"})
	// Toggle left/right parts to publish written data.
	lock.Publish()

	// Call RLock to read the data. Save the locking artefact, because you need it to call RUnlock later.
	data, artefact := lock.RLock()
	s := *data.(*store)  // No generics so we have to manually cast :(
	// Read from `s` in any way you want.
	fmt.Println(s["test"])
	// Output: value
	lock.RUnlock(artefact)
}
