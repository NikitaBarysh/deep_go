package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type node struct {
	key, value  int
	left, right *node
}

type OrderedMap struct {
	root *node
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{root: nil, size: 0}
}

func (m *OrderedMap) Insert(key, value int) {
	m.root = insert(m.root, key, value)
	m.size++

}

func insert(n *node, key, value int) *node {
	if n == nil {
		return &node{key: key, value: value}
	}
	if key < n.key {
		n.left = insert(n.left, key, value)
		return n
	}
	if key > n.key {
		n.right = insert(n.right, key, value)
		return n
	}
	n.value = value
	return n
}

func (m *OrderedMap) Erase(key int) {
	m.root = erase(m.root, key)
	m.size--
}

func erase(n *node, key int) *node {
	if n == nil {
		return nil
	}
	if key < n.key {
		n.left = erase(n.left, key)
		return n
	}
	if key > n.key {
		n.right = erase(n.right, key)
		return n
	}
	if n.left == nil {
		return n.right
	}
	if n.right == nil {
		return n.left
	}
	minNode := n.right
	for minNode.left != nil {
		minNode = minNode.left
	}
	n.key, n.value = minNode.key, minNode.value
	n.right = erase(n.right, minNode.key)
	return n
}

func (m *OrderedMap) Contains(key int) bool {
	n := m.root
	for n != nil {
		if key < n.key {
			n = n.left
		} else if key > n.key {
			n = n.right
		} else {
			return true
		}
	}
	return false
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	var inOrder func(n *node)
	inOrder = func(n *node) {
		if n == nil {
			return
		}
		inOrder(n.left)
		action(n.key, n.value)
		inOrder(n.right)
	}
	inOrder(m.root)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
