package builder

// A lightweight immutable map implementation that uses a slice as the underlying data structure.
// It offers only O(N) lookup performance and performs full copies on set.
// But it offers better performance for typical usage scenarios when building JSON objects
// than more sophisticated immutable map implementations using HMATs like
// https://github.com/arr-ai/frozen or https://github.com/raviqqe/hamt.
type immutableSliceMap[K comparable, V any] []mapEntry[K, V]

type mapEntry[K comparable, V any] struct {
	k K
	v V
}

func newImmutableSliceMap[K comparable, V any]() immutableSliceMap[K, V] {
	return make(immutableSliceMap[K, V], 0)
}

func (m immutableSliceMap[K, V]) Get(key K) (v V, ok bool) {
	for _, entry := range m {
		if entry.k == key {
			return entry.v, true
		}
	}
	return v, false
}

func (m immutableSliceMap[K, V]) Has(key K) bool {
	for _, entry := range m {
		if entry.k == key {
			return true
		}
	}
	return false
}

func (m immutableSliceMap[K, V]) Set(key K, value V) immutableSliceMap[K, V] {
	newMap := make(immutableSliceMap[K, V], len(m), len(m)+1)
	copy(newMap, m)

	found := false
	for i, entry := range newMap {
		if entry.k == key {
			newMap[i].v = value
			found = true
		}
	}
	if !found {
		newMap = append(newMap, mapEntry[K, V]{k: key, v: value})
	}
	return newMap
}

func (m immutableSliceMap[K, V]) Delete(key K) immutableSliceMap[K, V] {
	// Find index of key (if it exits) and copy everything before and after it to a new slice via `copy`

	idx := -1
	for i, entry := range m {
		if entry.k == key {
			idx = i
			break
		}
	}
	if idx == -1 {
		return m
	}

	newMap := make(immutableSliceMap[K, V], len(m)-1)

	copy(newMap, m[:idx])
	copy(newMap[idx:], m[idx+1:])

	return newMap
}

func (m immutableSliceMap[K, V]) Each(do func(k K, v V)) {
	for _, entry := range m {
		do(entry.k, entry.v)
	}
}

func (m immutableSliceMap[K, V]) clone() immutableSliceMap[K, V] {
	const extraCapacity = 16
	newMap := make(immutableSliceMap[K, V], len(m), len(m)+extraCapacity)
	copy(newMap, m)
	return newMap
}

func (m *immutableSliceMap[K, V]) mutatingSet(key K, value V) {
	found := false
	for i, entry := range *m {
		if entry.k == key {
			(*m)[i].v = value
			found = true
		}
	}
	if !found {
		*m = append(*m, mapEntry[K, V]{k: key, v: value})
	}
}
