package builder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/builder"
	"github.com/networkteam/qrb/fn"
)

func TestJsonBuildObject(t *testing.T) {
	t.Run("nested", func(t *testing.T) {
		b := fn.JsonBuildObject().
			Prop("name", qrb.String("Henry")).
			Prop("address", fn.JsonBuildObject().
				Prop("street", qrb.String("Main Street")),
			)

		sql, args, _ := qrb.Build(b).ToSQL()

		assert.Equal(t, "json_build_object('name','Henry','address',json_build_object('street','Main Street'))", sql)
		assert.Empty(t, args)
	})

	t.Run("immutability", func(t *testing.T) {
		b1 := fn.JsonBuildObject().
			Prop("name", qrb.String("Henry"))

		b2 := b1.Prop("age", qrb.Int(42))

		b3 := b2.Prop("name", qrb.String("John J."))

		b4 := b3.Unset("age")

		{
			sql, args, _ := qrb.Build(b1).ToSQL()

			assert.Equal(t, "json_build_object('name','Henry')", sql)
			assert.Empty(t, args)
		}

		{
			sql, args, _ := qrb.Build(b2).ToSQL()

			assert.Equal(t, "json_build_object('name','Henry','age',42)", sql)
			assert.Empty(t, args)
		}

		{
			sql, args, _ := qrb.Build(b3).ToSQL()

			assert.Equal(t, "json_build_object('name','John J.','age',42)", sql)
			assert.Empty(t, args)
		}

		{
			sql, args, _ := qrb.Build(b4).ToSQL()

			assert.Equal(t, "json_build_object('name','John J.')", sql)
			assert.Empty(t, args)
		}
	})

	t.Run("optimized via Start()", func(t *testing.T) {
		b1 := fn.JsonBuildObject().
			Prop("name", qrb.String("Henry"))

		b2 := b1.Start().
			Prop("age", qrb.Int(42)).
			Prop("height", qrb.Int(123)).
			End()

		b3 := b1.Start().
			Prop("age", qrb.Int(21)).
			Prop("height", qrb.Int(134)).
			Prop("name", qrb.String("John J.")).
			End()

		{
			sql, args, _ := qrb.Build(b1).ToSQL()

			assert.Equal(t, "json_build_object('name','Henry')", sql)
			assert.Empty(t, args)
		}

		{
			sql, args, _ := qrb.Build(b2).ToSQL()

			assert.Equal(t, "json_build_object('name','Henry','age',42,'height',123)", sql)
			assert.Empty(t, args)
		}

		{
			sql, args, _ := qrb.Build(b3).ToSQL()

			assert.Equal(t, "json_build_object('name','John J.','age',21,'height',134)", sql)
			assert.Empty(t, args)
		}
	})
}

// A list of random words to use for property names
var words = []string{
	"abruptly",
	"absurd",
	"abyss",
	"affix",
	"askew",
	"avenue",
	"awkward",
	"axiom",
	"azure",
	"bagpipes",
	"bandwagon",
	"banjo",
	"bayou",
	"beekeeper",
	"bikini",
	"blizzard",
	"boggle",
	"bookworm",
	"boxcar",
	"boxful",
	"buckaroo",
	"buffalo",
	"buffoon",
	"buxom",
	"buzzard",
	"buzzing",
	"buzzwords",
	"caliph",
	"cobweb",
	"cockiness",
	"croquet",
	"crypt",
	"curacao",
	"cycle",
	"daiquiri",
	"dirndl",
	"disavow",
	"dizzying",
	"duplex",
	"dwarves",
	"embezzle",
	"equip",
	"espionage",
}

func BenchmarkJsonBuildObject_Build(b *testing.B) {
	var j builder.JsonBuildObjectBuilder

	b.Run("nested", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			j = fn.JsonBuildObject().
				Prop("firstname", qrb.String("Henry")).
				Prop("lastname", qrb.String("Ford")).
				Prop("birthyear", qrb.Int(1863)).
				Prop("address", fn.JsonBuildObject().
					Prop("street", qrb.String("Main Street")).
					Prop("city", qrb.String("Dearborn")),
				)
		}
	})

	b.Run("long", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			j = fn.JsonBuildObject()
			for _, word := range words {
				j = j.Prop(word, qrb.String(word))
			}
		}
	})

	b.Run("long with build", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			jb := fn.JsonBuildObject().Start()
			for _, word := range words {
				jb = jb.Prop(word, qrb.String(word))
			}
			j = jb.End()
		}
	})

	_ = j
}

func BenchmarkJsonBuildObject_Write(b *testing.B) {
	nested := fn.JsonBuildObject().
		Prop("firstname", qrb.String("Henry")).
		Prop("lastname", qrb.String("Ford")).
		Prop("birthyear", qrb.Int(1863)).
		Prop("address", fn.JsonBuildObject().
			Prop("street", qrb.String("Main Street")).
			Prop("city", qrb.String("Dearborn")),
		)

	long := fn.JsonBuildObject()
	for _, word := range words {
		long = long.Prop(word, qrb.String(word))
	}

	b.ResetTimer()

	var sql string
	var args []any

	b.Run("nested", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sql, args, _ = qrb.Build(nested).ToSQL()
		}
	})

	b.Run("long", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sql, args, _ = qrb.Build(long).ToSQL()
		}
	})

	_ = sql
	_ = args
}
