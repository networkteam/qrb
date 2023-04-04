package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImmutableSliceMap_Get(t *testing.T) {
	m1 := newImmutableSliceMap[string, int]()
	m2 := m1.Set("a", 1)
	m3 := m2.Set("b", 2)
	m4 := m3.Delete("a")

	tests := []struct {
		name   string
		m      immutableSliceMap[string, int]
		key    string
		wantV  int
		wantOk bool
	}{
		{"EmptyMap", m1, "a", 0, false},
		{"OneElement", m2, "a", 1, true},
		{"TwoElements", m3, "a", 1, true},
		{"TwoElementsNoKey", m3, "c", 0, false},
		{"DeletedKey", m4, "a", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, ok := tt.m.Get(tt.key)
			if v != tt.wantV || ok != tt.wantOk {
				t.Errorf("ImmutableMap.Get(%s) = (%v, %v), want (%v, %v)", tt.key, v, ok, tt.wantV, tt.wantOk)
			}
		})
	}
}

func TestImmutableSliceMap_Set(t *testing.T) {
	m1 := newImmutableSliceMap[string, int]()
	m2 := m1.Set("a", 1)
	m3a := m2.Set("b", 2)
	m3b := m2.Set("c", 4)
	m4 := m3a.Set("a", 3)

	m5a := m3a.Set("d", 5).Set("e", 6)
	m5b := m3a.Set("f", 7).Set("g", 8)

	v, _ := m1.Get("a")
	if v != 0 {
		t.Errorf("ImmutableMap.Set() mutation detected in m1")
	}

	v, _ = m2.Get("a")
	if v != 1 {
		t.Errorf("ImmutableMap.Set() mutation detected in m2")
	}

	v, _ = m3a.Get("a")
	if v != 1 {
		t.Errorf("ImmutableMap.Set() mutation detected in m3a")
	}

	v, _ = m4.Get("a")
	if v != 3 {
		t.Errorf("ImmutableMap.Set() mutation detected in m4")
	}

	v, _ = m2.Get("b")
	if v != 0 {
		t.Errorf("ImmutableMap.Set() mutation detected in m2")
	}

	v, _ = m3a.Get("b")
	if v != 2 {
		t.Errorf("ImmutableMap.Set() mutation detected in m3a")
	}

	v, _ = m3b.Get("c")
	if v != 4 {
		t.Errorf("ImmutableMap.Set() mutation detected in m3b")
	}

	v, _ = m4.Get("b")
	if v != 2 {
		t.Errorf("ImmutableMap.Set() mutation detected in m4")
	}

	v, _ = m5a.Get("d")
	if v != 5 {
		t.Errorf("ImmutableMap.Set() mutation detected in m5a")
	}

	v, _ = m5b.Get("f")
	if v != 7 {
		t.Errorf("ImmutableMap.Set() mutation detected in m5b")
	}

	assertMapKeys(t, m1, nil, "m1 map should keep keys")
	assertMapKeys(t, m2, []string{"a"}, "m2 map should keep keys")
	assertMapKeys(t, m3a, []string{"a", "b"}, "m3a map should keep keys")
	assertMapKeys(t, m3b, []string{"a", "c"}, "m3b map should keep keys")
	assertMapKeys(t, m4, []string{"a", "b"}, "m4 order should be kept after override")
	assertMapKeys(t, m5a, []string{"a", "b", "d", "e"}, "m5 order")
	assertMapKeys(t, m5b, []string{"a", "b", "f", "g"}, "m5 order")
}

func TestImmutableSliceMap_Delete(t *testing.T) {
	m1 := newImmutableSliceMap[string, int]()
	m2 := m1.
		Set("a", 1).
		Set("b", 2)
	m3 := m2.Delete("a")

	if m1.Has("a") {
		t.Errorf("ImmutableMap.Delete() mutation detected in m1")
	}

	if m2.Has("a") == false {
		t.Errorf("ImmutableMap.Delete() mutation detected in m2")
	}

	if m3.Has("a") {
		t.Errorf("ImmutableMap.Delete() mutation detected in m3")
	}

	if m3.Has("b") == false {
		t.Errorf("ImmutableMap.Delete() mutation detected in m3")
	}
}

func TestImmutableSliceMap_Clone(t *testing.T) {
	m1 := newImmutableSliceMap[string, int]()
	m2 := m1.
		Set("a", 1).
		Set("b", 2)

	m3a := m2.clone()
	m3b := m2.clone()

	m3a.mutatingSet("c", 3)
	m3a.mutatingSet("a", -1)

	m3b.mutatingSet("c", 4)
	m3b.mutatingSet("d", 5)

	if m1.Has("a") {
		t.Errorf("ImmutableMap.Delete() mutation detected in m1")
	}

	if v, _ := m2.Get("a"); v != 1 {
		t.Errorf("ImmutableMap.Delete() mutation detected in m2")
	}
	if v, _ := m2.Get("b"); v != 2 {
		t.Errorf("ImmutableMap.Delete() mutation detected in m2")
	}
	if m2.Has("c") {
		t.Errorf("ImmutableMap.Delete() mutation detected in m2")
	}

	if v, _ := m3a.Get("a"); v != -1 {
		t.Errorf("ImmutableMap.Delete() mutation detected in m3a")
	}
	if !m3a.Has("c") {
		t.Errorf("ImmutableMap.Delete() mutation detected in m3a")
	}
	if m3a.Has("d") {
		t.Errorf("ImmutableMap.Delete() mutation detected in m3a")
	}

	if v, _ := m3b.Get("a"); v != 1 {
		t.Errorf("ImmutableMap.Delete() mutation detected in m3b")
	}
	if !m3b.Has("c") {
		t.Errorf("ImmutableMap.Delete() mutation detected in m3b")
	}
	if !m3b.Has("d") {
		t.Errorf("ImmutableMap.Delete() mutation detected in m3b")
	}
}

func assertMapKeys[K comparable, V any](t *testing.T, m immutableSliceMap[K, V], expectedKeys []K, msg string) {
	t.Helper()

	var orderedKeys []K
	m.Each(func(k K, v V) {
		orderedKeys = append(orderedKeys, k)
	})

	assert.Equal(t, expectedKeys, orderedKeys, msg)
}
