package algo

import "sync"

type Map[K comparable, V any] struct {
	sync.RWMutex
	Map map[K]V
}
