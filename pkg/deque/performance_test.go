package deque

import "testing"

func BenchmarkChan_FillDrain(b *testing.B) {
	c := make(chan int, b.N)
	for i := 0; i < b.N; i++ {
		c <- i
	}

	for i := 0; i < b.N; i++ {
		<-c
	}
}

func BenchmarkDeque_FillDrain(b *testing.B) {
	d := NewDequeWithCapacity(b.N)
	for i := 0; i < b.N; i++ {
		d.PushBack(i)
	}

	for i := 0; i < b.N; i++ {
		d.PopFront()
	}
}

func BenchmarkChan_Push(b *testing.B) {
	c := make(chan int, b.N)
	for i := 0; i < b.N; i++ {
		c <- i
	}
}

func BenchmarkDeque_Push(b *testing.B) {
	d := NewDequeWithCapacity(b.N)
	for i := 0; i < b.N; i++ {
		d.PushBack(i)
	}
}

func BenchmarkChan_Queue(b *testing.B) {
	c := make(chan int, 1)
	for i := 0; i < b.N; i++ {
		c <- i
		<-c
	}
}

func BenchmarkDeque_Queue(b *testing.B) {
	d := NewDequeWithCapacity(1)
	for i := 0; i < b.N; i++ {
		d.PushBack(i)
		d.PopFront()
	}
}

func BenchmarkChan_Queue2(b *testing.B) {
	c := make(chan int, 2)
	for i := 0; i < b.N/2; i++ {
		c <- i
		c <- i
		<-c
		<-c
	}
}

func BenchmarkDeque_Queue2(b *testing.B) {
	d := NewDequeWithCapacity(2)
	for i := 0; i < b.N/2; i++ {
		d.PushBack(i)
		d.PushBack(i)
		d.PopFront()
		d.PopFront()
	}
}
