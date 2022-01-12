package lock

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type testData map[string]string

func newTestData() *testData {
	s := make(testData)
	return &s
}

type input struct {
	key string
	val string
}

func (s *testData) Update(op Operation) OpResult {
	inp := op.(input)
	(*s)[inp.key] = inp.val
	return inp.key
}

func TestBasic(t *testing.T) {
	lock := NewLeftRightLock(newTestData(), newTestData())

	// Able to read, but nothing's inside
	m, art := lock.RLock()
	td := m.(*testData)
	_, ok := (*td)["test"]
	lock.RUnlock(art)
	assert.False(t, ok)

	val := lock.Write(input{"test", "123"})
	assert.Equal(t, "test", val)

	// Still nothing because changes have not been published.
	m, art = lock.RLock()
	td = m.(*testData)
	_, ok = (*td)["test"]
	lock.RUnlock(art)
	assert.False(t, ok)

	results := lock.Publish()
	assert.Len(t, results, 1)
	assert.Equal(t, "test", results[0])

	m, art = lock.RLock()
	td = m.(*testData)
	v, ok := (*td)["test"]
	lock.RUnlock(art)
	assert.True(t, ok)
	assert.Equal(t, "123", v)
}

func TestSimpleWait(t *testing.T) {
	write := func(wg *sync.WaitGroup, completed *int32, lock *LeftRightLock) {
		lock.Write(input{"test", "123"})
		lock.Publish()
		atomic.AddInt32(completed, 1)
		wg.Done()
	}

	// Set up concurrency primitives for testing
	lock := NewLeftRightLock(newTestData(), newTestData())
	var wg sync.WaitGroup
	wg.Add(1)
	writeCompleted := new(int32)
	*writeCompleted = 0

	// Start a read lock and verify that writing is blocked.
	m, art := lock.RLock()
	go write(&wg, writeCompleted, lock)
	time.Sleep(time.Millisecond * 20)

	// Verify that the data was not published.
	// NOTE: If we started a new read here, we would probably be able to see changes as it would read from the new side.
	// But writer is still blocked because of this longstanding read.
	v, ok := (*m.(*testData))["test"]
	assert.False(t, ok, v)
	assert.Equal(t, int32(0), atomic.LoadInt32(writeCompleted))
	lock.RUnlock(art)

	// Wait for write to finish
	wg.Wait()
	assert.Equal(t, int32(1), atomic.LoadInt32(writeCompleted))

	m, art = lock.RLock()
	v, ok = (*m.(*testData))["test"]
	assert.True(t, ok)
	assert.Equal(t, "123", v)
	lock.RUnlock(art)
}

func TestRaceConditionWait(t *testing.T) {
	write := func(wg *sync.WaitGroup, completed *int32, lock *LeftRightLock) {
		lock.Write(input{"test", "123"})
		lock.Publish()
		atomic.AddInt32(completed, 1)
		wg.Done()
	}

	// Set up concurrency primitives for testing
	lock := NewLeftRightLock(newTestData(), newTestData())
	var wg sync.WaitGroup
	wg.Add(1)
	writeCompleted := new(int32)
	*writeCompleted = 0

	// Simulate a race condition in which some previous reader was able to read Left, but lock Right
	art := int32(1)
	atomic.StoreInt64(lock.numReaders[art], 1)

	go write(&wg, writeCompleted, lock)
	time.Sleep(time.Millisecond * 20)
	// Verify that writer is blocked
	assert.Equal(t, int32(0), atomic.LoadInt32(writeCompleted))

	// Unlock read to allow writer to finish.
	lock.RUnlock(art)
	wg.Wait()
	assert.Equal(t, int32(1), atomic.LoadInt32(writeCompleted))
}
