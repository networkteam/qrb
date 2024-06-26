package builder_test

import (
	"testing"

	"github.com/networkteam/qrb"
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
}
