package lock

import (
	"runtime"
	"sync/atomic"

	"github.com/bitstonks/leftright/pkg/deque"
)

// Operation is the generic input accepted by the left-right data structure. Named for readability. Should be immutable.
type Operation interface{}

// OpResult is the generic output returned by the left-right data structure update method. Named for readability.
type OpResult interface{}

// LeftRightStructure represents the data structure that we are operating over in out left-right system.
type LeftRightStructure interface {
	// Update is the method that lets us mutate the data structure. It has to be deterministic, because we need to apply
	// it twice - once on the left and once on the right part of the structure, and we're relying on those being equal.
	Update(Operation) OpResult
}

// LeftRightLock provides the core of the left-right pattern.
type LeftRightLock struct {
	// data holds the left and right structures, which we'll be reading and updating.
	data [2]LeftRightStructure
	// opQ is a queue used to store operations that were only applied on a single side of the structure.
	opQ deque.Deque
	// numReaders is an array of 2 atomic ints, counting numbers of readers on the left/right instance.
	numReaders [2]*int64
	// sideToLock is an atomic var (0 or 1) and determines which side the reader locks before it starts to read.
	sideToLock *int32
	// sideToRead is an atomic var (0 or 1) and determines which side the reader should read.
	sideToRead *int32
}

// NewLeftRightLock creates a LeftRightLock. The two structures provided have to be equal.
func NewLeftRightLock(left, right LeftRightStructure) *LeftRightLock {
	m := &LeftRightLock{
		data:       [2]LeftRightStructure{left, right},
		opQ:        deque.NewDeque(),
		numReaders: [2]*int64{new(int64), new(int64)},
		// Start reading on Left as this is initialized to 0
		sideToLock: new(int32),
		sideToRead: new(int32),
	}
	return m
}

// RLock should be called by a go routine before it starts reading. This is a wait-free operation.
func (lr *LeftRightLock) RLock() (LeftRightStructure, int32) {
	lockIdx := atomic.LoadInt32(lr.sideToLock)
	atomic.AddInt64(lr.numReaders[lockIdx], 1)
	return lr.data[atomic.LoadInt32(lr.sideToRead)], lockIdx
}

// RUnlock should be called by a go routine after it stops reading.
func (lr *LeftRightLock) RUnlock(lockIdx int32) {
	atomic.AddInt64(lr.numReaders[lockIdx], -1)
}

// Publish swaps read and write sides, thus publishing all mutations made on the writing side. This function may have to
// wait up to one read operation for it finish the swap. After the swap it will also apply all outstanding update
// operations that were only applied on one side.
func (lr *LeftRightLock) Publish() []OpResult {
	lr.swap()
	return lr.reapplyOpHistory()
}

// Write runs the Update method ont the writeable side with the given operator.
func (lr *LeftRightLock) Write(op Operation) OpResult {
	lr.opQ.PushBack(op)
	sideToWrite := 1 - atomic.LoadInt32(lr.sideToRead)
	return lr.data[sideToWrite].Update(op)
}

// swap read and write sides, thus publishing all mutations made on the writing side. This function may have to
// wait for any outstanding reads before the swap is considered complete.
func (lr *LeftRightLock) swap() (newSideToWrite int32) {
	// At this point both sides are safe to read, so redirect reads to the new sides, but keep locks on the same side.
	newSideToRead := 1 - atomic.LoadInt32(lr.sideToRead)
	atomic.StoreInt32(lr.sideToRead, newSideToRead)

	lockIdx := 1 - atomic.LoadInt32(lr.sideToLock)
	// Wait for all readers from previous iteration to complete. We have to do this because there is a race condition
	// in which a reader can lock a different side than it reads from. This wait ensures that the side it reads from
	// is always correct even if the lock isn't. This is tested in TestRaceConditionWait.
	for atomic.LoadInt64(lr.numReaders[lockIdx]) != 0 {
		runtime.Gosched()
	}
	// Switch locking side and wait for all the readers to evacuate
	atomic.StoreInt32(lr.sideToLock, lockIdx)
	lockIdx = 1 - lockIdx
	for atomic.LoadInt64(lr.numReaders[lockIdx]) != 0 {
		runtime.Gosched()
	}
	return 1 - newSideToRead
}

// reapplyOpHistory will apply all outstanding update operations that were only applied on one side and return the
// result.
func (lr *LeftRightLock) reapplyOpHistory() (results []OpResult) {
	sideToWrite := 1 - atomic.LoadInt32(lr.sideToRead)
	for lr.opQ.Len() > 0 {
		op, _ := lr.opQ.PopFront()
		results = append(results, lr.data[sideToWrite].Update(op))
	}
	return
}
