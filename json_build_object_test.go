package jrm_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/networkteam/jrm"
)

func TestJsonBuildObject(t *testing.T) {
	t.Run("nested", func(t *testing.T) {
		b := jrm.JsonBuildObject().
			Prop("name", jrm.String("Henry")).
			Prop("address", jrm.JsonBuildObject().
				Prop("street", jrm.String("Main Street")),
			)

		sql, args := expToSQL(b)

		assert.Equal(t, "JSON_BUILD_OBJECT('address',JSON_BUILD_OBJECT('street','Main Street'),'name','Henry')", sql)
		assert.Empty(t, args)
	})

	t.Run("immutability", func(t *testing.T) {
		b1 := jrm.JsonBuildObject().
			Prop("name", jrm.String("Henry"))

		b2 := b1.Prop("age", jrm.Int(42))

		{
			sql, args := expToSQL(b1)

			assert.Equal(t, "JSON_BUILD_OBJECT('name','Henry')", sql)
			assert.Empty(t, args)
		}

		{
			sql, args := expToSQL(b2)

			assert.Equal(t, "JSON_BUILD_OBJECT('age',42,'name','Henry')", sql)
			assert.Empty(t, args)
		}
	})
}

func BenchmarkJsonBuildObject_Build(b *testing.B) {
	var j *jrm.JsonBuildObjectBuilder
	b.Run("nested", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			j = jrm.JsonBuildObject().
				Prop("firstname", jrm.String("Henry")).
				Prop("lastname", jrm.String("Ford")).
				Prop("birthyear", jrm.Int(1863)).
				Prop("address", jrm.JsonBuildObject().
					Prop("street", jrm.String("Main Street")).
					Prop("city", jrm.String("Dearborn")),
				)
		}
	})

	b.Run("nested init", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			j = jrm.JsonBuildObject().
				Init().
				Prop("firstname", jrm.String("Henry")).
				Prop("lastname", jrm.String("Ford")).
				Prop("birthyear", jrm.Int(1863)).
				Prop("address", jrm.JsonBuildObject().
					Init().
					Prop("street", jrm.String("Main Street")).
					Prop("city", jrm.String("Dearborn")).
					Done(),
				).
				Done()
		}
	})

	_ = j
}

func expToSQL(exp jrm.Exp) (string, []any) {
	sb := new(strings.Builder)
	args := exp.WriteSQL(sb)
	sql := sb.String()

	return sql, args
}
