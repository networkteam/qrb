package jrm

import (
	"sort"
	"strconv"
	"strings"
)

type SQLWriter interface {
	WriteSQL(sb *strings.Builder) (args []any)
}

type Exp interface {
	isExp()
	SQLWriter
}

type Json interface {
}

func String(s string) Exp {
	return expStr(s)
}

type expStr string

func (e expStr) isExp() {}

func (e expStr) WriteSQL(sb *strings.Builder) (args []any) {
	sb.WriteString(pqQuoteString(string(e)))
	return nil
}

func Int(s int) Exp {
	return expInt(s)
}

type expInt int

func (e expInt) isExp() {}

func (e expInt) WriteSQL(sb *strings.Builder) (args []any) {
	sb.WriteString(strconv.Itoa(int(e)))
	return nil
}

func JsonBuildObject() *JsonBuildObjectBuilder {
	return &JsonBuildObjectBuilder{
		props: newImmutableMap[string, Exp](),
	}
}

type JsonBuildObjectBuilder struct {
	props *immutableMap[string, Exp]
}

func (b *JsonBuildObjectBuilder) isExp() {}

func (b *JsonBuildObjectBuilder) WriteSQL(sb *strings.Builder) (args []any) {
	sb.WriteString("JSON_BUILD_OBJECT(")

	keys := b.props.Keys()
	sort.Strings(keys)

	i := 0
	for _, key := range keys {
		value, _ := b.props.Get(key)
		if i > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(pqQuoteString(key))
		sb.WriteRune(',')
		valueArgs := value.WriteSQL(sb)
		args = append(args, valueArgs...)

		i++
	}
	sb.WriteRune(')')
	return args
}

func (b *JsonBuildObjectBuilder) Prop(key string, value Exp) *JsonBuildObjectBuilder {
	newProps := b.props.Set(key, value)
	return &JsonBuildObjectBuilder{
		props: newProps,
	}
}

type JsonBuildObjectBuilderInit struct {
	builder *JsonBuildObjectBuilder
}

func (b *JsonBuildObjectBuilder) Init() *JsonBuildObjectBuilderInit {
	return &JsonBuildObjectBuilderInit{
		builder: b,
	}
}

func (i *JsonBuildObjectBuilderInit) Prop(key string, value Exp) *JsonBuildObjectBuilderInit {
	// Yeah, directly set it!
	i.builder.props.m[key] = value
	return i
}

func (i *JsonBuildObjectBuilderInit) Done() *JsonBuildObjectBuilder {
	return i.builder
}

// ---

type immutableMap[K comparable, V any] struct {
	m map[K]V
}

func newImmutableMap[K comparable, V any]() *immutableMap[K, V] {
	return &immutableMap[K, V]{
		m: make(map[K]V),
	}
}

func (m *immutableMap[K, V]) Get(key K) (V, bool) {
	v, ok := m.m[key]
	return v, ok
}

func (m *immutableMap[K, V]) Set(key K, value V) *immutableMap[K, V] {
	newMap := make(map[K]V, len(m.m)+1)
	for k, v := range m.m {
		newMap[k] = v
	}
	newMap[key] = value
	return &immutableMap[K, V]{
		m: newMap,
	}
}

func (m *immutableMap[K, V]) Delete(key K) *immutableMap[K, V] {
	newMap := make(map[K]V, len(m.m)-1)
	for k, v := range m.m {
		newMap[k] = v
	}
	delete(newMap, key)
	return &immutableMap[K, V]{
		m: newMap,
	}
}

func (m *immutableMap[K, V]) Len() int {
	return len(m.m)
}

func (m *immutableMap[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.m))
	for key := range m.m {
		keys = append(keys, key)
	}
	return keys
}
