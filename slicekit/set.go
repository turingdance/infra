package slicekit

import "sync"

// Set 定义了一个set，其中元素存储在map中
type Set[T comparable] struct {
	locker   *sync.RWMutex
	elements map[T]bool
}

// NewSet 创建一个新的set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		locker:   &sync.RWMutex{},
		elements: make(map[T]bool)}
}

// Add 向set中添加元素
func (s *Set[T]) Add(value T) {
	s.locker.Lock()
	s.elements[value] = true
	s.locker.Unlock()
}

// Remove 从set中移除元素
func (s *Set[T]) Remove(value T) {
	s.locker.Lock()
	delete(s.elements, value)
	s.locker.Unlock()
}

// Contains 检查元素是否在set中
func (s *Set[T]) Contains(value T) bool {
	s.locker.Lock()
	_, ok := s.elements[value]
	s.locker.Unlock()
	return ok
}

// Size 返回set的大小
func (s *Set[T]) Size() int {
	return len(s.elements)
}

// Size 返回set的大小
func (s *Set[T]) Range(fun func(value T) bool) {
	s.locker.Lock()
	for k, _ := range s.elements {
		if !fun(k) {
			break
		}
	}
	s.locker.Unlock()
}

// Clear 清空set
func (s *Set[T]) Clear() {
	s.locker.Lock()
	for k := range s.elements {
		delete(s.elements, k)
	}
	s.locker.Unlock()
}
