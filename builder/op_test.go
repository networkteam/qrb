package builder_test

import (
	"testing"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/builder"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestOp(t *testing.T) {
	t.Run("cast", func(t *testing.T) {
		t.Run("single exp", func(t *testing.T) {
			b := qrb.Arg("foo").Cast("text")

			testhelper.AssertSQLWriterEquals(
				t,
				"$1::text",
				[]any{"foo"},
				b,
			)
		})

		t.Run("combined exp", func(t *testing.T) {
			b := qrb.N("json_column").JsonExtractText(qrb.Arg("my_field")).Cast("int")

			testhelper.AssertSQLWriterEquals(
				t,
				"(json_column ->> $1)::int",
				[]any{"my_field"},
				b,
			)
		})

		t.Run("cast and like", func(t *testing.T) {
			b := qrb.N("articles.content").Cast("text").ILike(qrb.Arg("%foo%"))

			testhelper.AssertSQLWriterEquals(
				t,
				"articles.content::text ILIKE $1",
				[]any{"%foo%"},
				b,
			)
		})

		t.Run("array type", func(t *testing.T) {
			b := qrb.Array(qrb.Arg("foo"), qrb.Arg("bar")).Cast("uuid[]")

			testhelper.AssertSQLWriterEquals(
				t,
				"ARRAY[$1, $2]::uuid[]",
				[]any{"foo", "bar"},
				b,
			)
		})
	})

	t.Run("precedence", func(t *testing.T) {
		t.Run("plus and mult", func(t *testing.T) {
			b := qrb.N("a").Plus(qrb.N("b")).Mult(qrb.N("c"))

			testhelper.AssertSQLWriterEquals(
				t,
				"(a + b) * c",
				nil,
				b,
			)
		})

		t.Run("plus, plus and grouped minus", func(t *testing.T) {
			b := qrb.N("a").Plus(qrb.N("b")).Plus(qrb.N("c").Minus(qrb.N("d")))

			testhelper.AssertSQLWriterEquals(
				t,
				"a + b + (c - d)",
				nil,
				b,
			)
		})

		t.Run("plus, plus and minus", func(t *testing.T) {
			b := qrb.N("a").Plus(qrb.N("b")).Plus(qrb.N("c")).Minus(qrb.N("d"))

			testhelper.AssertSQLWriterEquals(
				t,
				"a + b + c - d",
				nil,
				b,
			)
		})

		t.Run("plus, minus and plus", func(t *testing.T) {
			b := qrb.N("a").Plus(qrb.N("b")).Minus(qrb.N("c").Plus(qrb.N("d")))

			testhelper.AssertSQLWriterEquals(
				t,
				"a + b - (c + d)",
				nil,
				b,
			)
		})

		t.Run("plus, plus and plus", func(t *testing.T) {
			b := qrb.N("a").Plus(qrb.N("b")).Plus(qrb.N("c").Plus(qrb.N("d")))

			testhelper.AssertSQLWriterEquals(
				t,
				"a + b + c + d",
				nil,
				b,
			)
		})
		t.Run("plus times plus", func(t *testing.T) {
			e1 := qrb.N("a").Plus(qrb.N("b"))
			e2 := qrb.N("c").Plus(qrb.N("d"))
			b := e1.Mult(e2)

			testhelper.AssertSQLWriterEquals(
				t,
				"(a + b) * (c + d)",
				nil,
				b,
			)
		})
	})

	t.Run("subscript", func(t *testing.T) {
		t.Run("simple column subscript", func(t *testing.T) {
			// mytable.arraycolumn[4]
			b := qrb.N("mytable.arraycolumn").Subscript(qrb.Int(4))

			testhelper.AssertSQLWriterEquals(
				t,
				"mytable.arraycolumn[4]",
				nil,
				b,
			)
		})

		t.Run("parameter subscript", func(t *testing.T) {
			// $1[10]
			b := qrb.Arg(1).Subscript(qrb.Int(10))

			testhelper.AssertSQLWriterEquals(
				t,
				"$1[10]",
				[]any{1},
				b,
			)
		})

		t.Run("function call subscript with parentheses", func(t *testing.T) {
			// (arrayfunction(a,b))[42]
			fn := qrb.Func("arrayfunction", qrb.N("a"), qrb.N("b"))
			b := fn.Subscript(qrb.Int(42))

			testhelper.AssertSQLWriterEquals(
				t,
				"(arrayfunction(a,b))[42]",
				nil,
				b,
			)
		})

		t.Run("array slice with parameter", func(t *testing.T) {
			// $1[10:42]
			b := qrb.Arg(1).Subscript(qrb.Int(10), qrb.Int(42))

			testhelper.AssertSQLWriterEquals(
				t,
				"$1[10:42]",
				[]any{1},
				b,
			)
		})

		t.Run("column array slice", func(t *testing.T) {
			// mytable.arraycolumn[1:5]
			b := qrb.N("mytable.arraycolumn").Subscript(qrb.Int(1), qrb.Int(5))

			testhelper.AssertSQLWriterEquals(
				t,
				"mytable.arraycolumn[1:5]",
				nil,
				b,
			)
		})

		t.Run("multidimensional array subscript", func(t *testing.T) {
			// mytable.two_d_column[17][34]
			b := qrb.N("mytable.two_d_column").Subscript(qrb.Int(17)).Subscript(qrb.Int(34))

			testhelper.AssertSQLWriterEquals(
				t,
				"mytable.two_d_column[17][34]",
				nil,
				b,
			)
		})

		t.Run("subscript with arithmetic expression", func(t *testing.T) {
			// (a + b)[1]
			expr := qrb.N("a").Plus(qrb.N("b"))
			b := expr.Subscript(qrb.Int(1))

			testhelper.AssertSQLWriterEquals(
				t,
				"(a + b)[1]",
				nil,
				b,
			)
		})

		t.Run("subscript without parentheses for high precedence", func(t *testing.T) {
			// a.b[1] - no parentheses needed as dot has higher precedence
			b := qrb.N("a.b").Subscript(qrb.Int(1))

			testhelper.AssertSQLWriterEquals(
				t,
				"a.b[1]",
				nil,
				b,
			)
		})

		t.Run("complex expression with subscript", func(t *testing.T) {
			// (a * b)[1:3]
			expr := qrb.N("a").Mult(qrb.N("b"))
			b := expr.Subscript(qrb.Int(1), qrb.Int(3))

			testhelper.AssertSQLWriterEquals(
				t,
				"(a * b)[1:3]",
				nil,
				b,
			)
		})
	})

	t.Run("comparison operators", func(t *testing.T) {
		t.Run("is distinct from", func(t *testing.T) {
			// Example: a IS DISTINCT FROM b
			b := qrb.N("a").Plus(qrb.Int(1)).IsDistinctFrom(qrb.N("b"))

			testhelper.AssertSQLWriterEquals(
				t,
				"a + 1 IS DISTINCT FROM b",
				nil,
				b,
			)
		})

		t.Run("is not distinct from", func(t *testing.T) {
			// Example: a IS NOT DISTINCT FROM b
			b := builder.ExpBase{Exp: qrb.Not(qrb.N("a"))}.IsNotDistinctFrom(qrb.N("b"))

			testhelper.AssertSQLWriterEquals(
				t,
				"(NOT a) IS NOT DISTINCT FROM b",
				nil,
				b,
			)
		})
	})
}
