package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type myStruct struct {
	name    string
	country string
}

func TestImmutableMap_Set(t *testing.T) {
	m := newImmutableMap[string, myStruct]()
	m2 := m.Set("foo", myStruct{name: "foo", country: "DE"})
	m3 := m2.Set("bar", myStruct{name: "bar", country: "UK"})

	{
		_, ok := m.Get("foo")
		assert.False(t, ok)
	}

	{
		foo, ok := m2.Get("foo")
		assert.True(t, ok)
		assert.Equal(t, "foo", foo.name)

		_, ok = m.Get("bar")
		assert.False(t, ok)
	}

	{
		foo, ok := m3.Get("foo")
		assert.True(t, ok)
		assert.Equal(t, "foo", foo.name)

		bar, ok := m3.Get("bar")
		assert.True(t, ok)
		assert.Equal(t, "bar", bar.name)
	}

	// Time to update an existing value

	m4 := m3.Set("foo", myStruct{name: "foo", country: "US"})

	{
		foo, ok := m4.Get("foo")
		assert.True(t, ok)
		assert.Equal(t, "US", foo.country)
	}

	// Now iterate all items
	{
		var names []string
		m4.Each(func(k string, v myStruct) {
			names = append(names, v.name)
		})
		assert.Equal(t, []string{"foo", "bar"}, names)
	}
}
