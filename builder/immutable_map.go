package builder

type immutableMap[K comparable, V any] struct {
	k K
	v V

	parent *immutableMap[K, V]
}

func newImmutableMap[K comparable, V any]() immutableMap[K, V] {
	return immutableMap[K, V]{}
}

func (m immutableMap[K, V]) Get(key K) (v V, ok bool) {
	current := &m
	for {
		if current.parent == nil {
			return v, false
		}
		if current.k == key {
			return current.v, true
		}
		current = current.parent
	}
}

func (m immutableMap[K, V]) Has(key K) bool {
	current := &m
	for {
		if current.parent == nil {
			return false
		}
		if current.k == key {
			return true
		}
		current = current.parent
	}
}

func (m immutableMap[K, V]) Set(key K, value V) immutableMap[K, V] {
	newMap := m
	// Check if key already exists and rebuild map without it
	if m.Has(key) {
		newMap = m.Delete(key)
	}

	return immutableMap[K, V]{
		k:      key,
		v:      value,
		parent: &newMap,
	}
}

func (m immutableMap[K, V]) Delete(key K) immutableMap[K, V] {
	// Rebuild map without key
	newMap := newImmutableMap[K, V]()

	current := &m
	for {
		if current.parent == nil {
			return newMap
		}
		if current.k != key {
			// New variable to take fresh pointer
			newParent := newMap
			newMap = immutableMap[K, V]{
				k:      current.k,
				v:      current.v,
				parent: &newParent,
			}
		}
		current = current.parent
	}
}

func (m immutableMap[K, V]) Each(do func(k K, v V)) {
	current := &m
	for {
		if current.parent == nil {
			return
		}
		do(current.k, current.v)
		current = current.parent
	}
}
