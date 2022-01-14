package lock

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
)

const N = 100

// Read `key` from a map in LeftRightLock `n` times.
func lrRead(lr LeftRightLock, key string, n int, wg *sync.WaitGroup) {
	for i := 0; i < n; i++ {
		data, art := lr.RLock()
		_, _ = (*data.(*testData))[fmt.Sprintf("%s%d", key, rand.Intn(N))]
		lr.RUnlock(art)
	}
	if wg != nil {
		wg.Done()
	}
}

// Read `key` from a sync.Map `n` times.
func smRead(m *sync.Map, key string, n int, wg *sync.WaitGroup) {
	for i := 0; i < n; i++ {
		_, _ = m.Load(fmt.Sprintf("%s%d", key, rand.Intn(N)))
	}
	if wg != nil {
		wg.Done()
	}
}

// Read `key` from a map protected by a RWMutex `n` times.
func rwRead(m *map[string]string, mu *sync.RWMutex, key string, n int, wg *sync.WaitGroup) {
	for i := 0; i < n; i++ {
		mu.RLock()
		_, _ = (*m)[fmt.Sprintf("%s%d", key, rand.Intn(N))]
		mu.RUnlock()
	}
	if wg != nil {
		wg.Done()
	}
}

func runConcurrent(n int, fn func(*sync.WaitGroup)) *sync.WaitGroup {
	if n == 0 {
		fn(nil)
		return nil
	}
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go fn(&wg)
	}
	return &wg
}

func conLrRead(nReads, nThreads int) *sync.WaitGroup {
	l, r := newTestData(), newTestData()
	l.Update(input{"test", "123"})
	r.Update(input{"test", "123"})
	lr := *NewLeftRightLock(l, r)
	return runConcurrent(nThreads, func(wg *sync.WaitGroup){lrRead(lr, "test", nReads, wg)})
}

func conSmRead(nReads, nThreads int) *sync.WaitGroup {
	var m sync.Map
	m.Store("test", "123")
	return runConcurrent(nThreads, func(wg *sync.WaitGroup){smRead(&m, "test", nReads, wg)})
}

func conRwRead(nReads, nThreads int) *sync.WaitGroup {
	mu := sync.RWMutex{}
	m := map[string]string{"test": "123"}
	return runConcurrent(nThreads, func(wg *sync.WaitGroup){rwRead(&m, &mu, "test", nReads, wg)})
}

func BenchmarkAtomic(b *testing.B) {
	x := new(int64)
	for i := 0; i < b.N; i++ {
		atomic.LoadInt64(x)
		atomic.LoadInt64(x)
		atomic.AddInt64(x, 1)
		atomic.AddInt64(x, -1)
	}
}

func BenchmarkAtomic10(b *testing.B) {
	x := new(int64)
	runConcurrent(2, func(wg *sync.WaitGroup){
		for i := 0; i < b.N; i++ {
			atomic.LoadInt64(x)
			atomic.LoadInt64(x)
			atomic.AddInt64(x, 1)
			atomic.AddInt64(x, -1)
		}
		wg.Done()
	}).Wait()
}

func BenchmarkLeftRight_Read(b *testing.B) {
	conLrRead(b.N, 0)
}

func BenchmarkSyncMap_Read(b *testing.B) {
	conSmRead(b.N, 0)
}

func BenchmarkRWLock_Read(b *testing.B) {
	conRwRead(b.N, 0)
}

func BenchmarkLeftRight_Read2(b *testing.B) {
	conLrRead(b.N, 2).Wait()
}

func BenchmarkSyncMap_Read2(b *testing.B) {
	conSmRead(b.N, 2).Wait()
}

func BenchmarkRWLock_Read2(b *testing.B) {
	conRwRead(b.N, 2).Wait()
}

func BenchmarkLeftRight_Read10(b *testing.B) {
	conLrRead(b.N, 10).Wait()
}

func BenchmarkSyncMap_Read10(b *testing.B) {
	conSmRead(b.N, 10).Wait()
}

func BenchmarkRWLock_Read10(b *testing.B) {
	conRwRead(b.N, 10).Wait()
}

func BenchmarkLeftRight_Read100(b *testing.B) {
	conLrRead(b.N, 100).Wait()
}

func BenchmarkSyncMap_Read100(b *testing.B) {
	conSmRead(b.N, 100).Wait()
}

func BenchmarkRWLock_Read100(b *testing.B) {
	conRwRead(b.N, 100).Wait()
}

func BenchmarkLeftRight_Write(b *testing.B) {
	l, r := newTestData(), newTestData()
	l.Update(input{"test", "123"})
	r.Update(input{"test", "123"})
	lr := *NewLeftRightLock(l, r)
	for i := 0; i < b.N; i++ {
		lr.Write(input{fmt.Sprintf("test%d", rand.Intn(N)), "123"})
		lr.Publish()
	}
}

func BenchmarkSyncMap_Write(b *testing.B) {
	var m sync.Map
	m.Store("test", "123")
	for i := 0; i < b.N; i++ {
		m.Store(fmt.Sprintf("test%d", rand.Intn(N)), "123")
	}
}

func BenchmarkRWLock_Write(b *testing.B) {
	mu := sync.RWMutex{}
	m := map[string]string{"test": "123"}
	for i := 0; i < b.N; i++ {
		mu.Lock()
		m[fmt.Sprintf("test%d", rand.Intn(N))] = "123"
		mu.Unlock()
	}
}

func BenchmarkLeftRight_WriteRead1(b *testing.B) {
	l, r := newTestData(), newTestData()
	l.Update(input{"test", "123"})
	r.Update(input{"test", "123"})
	lr := *NewLeftRightLock(l, r)
	wg := runConcurrent(1, func(wg *sync.WaitGroup){lrRead(lr, "test", b.N, wg)})
	for i := 0; i < b.N; i++ {
		lr.Write(input{fmt.Sprintf("test%d", rand.Intn(N)), "123"})
		lr.Publish()
	}
	wg.Wait()
}

func BenchmarkSyncMap_WriteRead1(b *testing.B) {
	var m sync.Map
	m.Store("test", "123")
	wg := runConcurrent(1, func(wg *sync.WaitGroup){smRead(&m, "test", b.N, wg)})
	for i := 0; i < b.N; i++ {
		m.Store(fmt.Sprintf("test%d", rand.Intn(N)), "123")
	}
	wg.Wait()
}

func BenchmarkRWLock_WriteRead1(b *testing.B) {
	var mu sync.RWMutex
	m := map[string]string{"test": "123"}
	wg := runConcurrent(1, func(wg *sync.WaitGroup){rwRead(&m, &mu, "test", b.N, wg)})
	for i := 0; i < b.N; i++ {
		mu.Lock()
		m[fmt.Sprintf("test%d", rand.Intn(N))] = "123"
		mu.Unlock()
	}
	wg.Wait()
}

func BenchmarkLeftRight_WriteRead4(b *testing.B) {
	l, r := newTestData(), newTestData()
	l.Update(input{"test", "123"})
	r.Update(input{"test", "123"})
	lr := *NewLeftRightLock(l, r)
	wg := runConcurrent(4, func(wg *sync.WaitGroup){lrRead(lr, "test", b.N, wg)})
	for i := 0; i < b.N; i++ {
		lr.Write(input{fmt.Sprintf("test%d", rand.Intn(N)), "123"})
		lr.Publish()
	}
	wg.Wait()
}

func BenchmarkSyncMap_WriteRead4(b *testing.B) {
	var m sync.Map
	m.Store("test", "123")
	wg := runConcurrent(4, func(wg *sync.WaitGroup){smRead(&m, "test", b.N, wg)})
	for i := 0; i < b.N; i++ {
		m.Store(fmt.Sprintf("test%d", rand.Intn(N)), "123")
	}
	wg.Wait()
}

func BenchmarkRWLock_WriteRead4(b *testing.B) {
	var mu sync.RWMutex
	m := map[string]string{"test": "123"}
	wg := runConcurrent(4, func(wg *sync.WaitGroup){rwRead(&m, &mu, "test", b.N, wg)})
	for i := 0; i < b.N; i++ {
		mu.Lock()
		m[fmt.Sprintf("test%d", rand.Intn(N))] = "123"
		mu.Unlock()
	}
	wg.Wait()
}
