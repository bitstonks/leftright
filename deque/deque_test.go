package deque

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFrontQueue(t *testing.T) {
	q := NewDeque()
	q.PushBack(1)
	q.PushBack(2)
	q.PushBack(3)
	assert.Equal(t, 3, q.Len())

	v, err := q.PopFront()
	assert.Nil(t, err)
	assert.Equal(t, 1, v)

	v, err = q.PopFront()
	assert.Nil(t, err)
	assert.Equal(t, 2, v)

	v, err = q.PopFront()
	assert.Nil(t, err)
	assert.Equal(t, 3, v)

	assert.Equal(t, 0, q.Len())

	v, err = q.PopFront()
	assert.NotNil(t, err)

	assert.Equal(t, 0, q.Len())
}

func TestBackQueue(t *testing.T) {
	q := NewDeque()
	q.PushFront(1)
	q.PushFront(2)
	q.PushFront(3)
	assert.Equal(t, 3, q.Len())

	v, err := q.PopBack()
	assert.Nil(t, err)
	assert.Equal(t, 1, v)

	v, err = q.PopBack()
	assert.Nil(t, err)
	assert.Equal(t, 2, v)

	v, err = q.PopBack()
	assert.Nil(t, err)
	assert.Equal(t, 3, v)

	assert.Equal(t, 0, q.Len())

	v, err = q.PopBack()
	assert.NotNil(t, err)

	assert.Equal(t, 0, q.Len())
}

func TestBackStack(t *testing.T) {
	q := NewDeque()
	q.PushBack(1)
	q.PushBack(2)
	q.PushBack(3)
	assert.Equal(t, 3, q.Len())

	v, err := q.PopBack()
	assert.Nil(t, err)
	assert.Equal(t, 3, v)

	v, err = q.PopBack()
	assert.Nil(t, err)
	assert.Equal(t, 2, v)

	v, err = q.PopBack()
	assert.Nil(t, err)
	assert.Equal(t, 1, v)

	assert.Equal(t, 0, q.Len())
	v, err = q.PopBack()
	assert.NotNil(t, err)

	assert.Equal(t, 0, q.Len())
}

func TestFrontStack(t *testing.T) {
	q := NewDeque()
	q.PushFront(1)
	q.PushFront(2)
	q.PushFront(3)
	assert.Equal(t, 3, q.Len())

	v, err := q.PopFront()
	assert.Nil(t, err)
	assert.Equal(t, 3, v)

	v, err = q.PopFront()
	assert.Nil(t, err)
	assert.Equal(t, 2, v)

	v, err = q.PopFront()
	assert.Nil(t, err)
	assert.Equal(t, 1, v)

	assert.Equal(t, 0, q.Len())
	v, err = q.PopFront()
	assert.NotNil(t, err)

	assert.Equal(t, 0, q.Len())
}

func TestResize(t *testing.T) {
	q := NewDequeWithCapacity(0)
	assert.Equal(t, 0, len(q.data))

	q.PushBack(1)
	assert.Equal(t, 1, len(q.data))
	q.PushBack(2)
	assert.Equal(t, 2, len(q.data))
	q.PushBack(3)
	assert.Equal(t, 4, len(q.data))
	q.PushBack(4)
	assert.Equal(t, 4, len(q.data))

	v, err := q.PopFront()
	assert.Equal(t, 4, len(q.data))
	assert.Nil(t, err)
	assert.Equal(t, 1, v)

	q.PushBack(5)
	assert.Equal(t, 4, len(q.data))
	q.PushBack(6)
	assert.Equal(t, 8, len(q.data))

	assert.Equal(t, 5, q.Len())

	for _, want := range []int{2, 3, 4, 5, 6} {
		got, err := q.PopFront()
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
	assert.Equal(t, 0, q.Len())
}
