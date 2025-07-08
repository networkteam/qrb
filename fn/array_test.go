package fn_test

import (
	"testing"

	. "github.com/networkteam/qrb"
	"github.com/networkteam/qrb/fn"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestArrayFunctions(t *testing.T) {
	t.Run("unnest", func(t *testing.T) {
		q := Select(N("*")).
			From(
				fn.Unnest(
					Array(Int(1), Int(2)),
					Array(String("foo"), String("bar"), String("baz")),
				),
			).As("x").ColumnAliases("a", "b").
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT * FROM unnest(ARRAY[1,2], ARRAY['foo','bar','baz']) AS x (a,b)
		`, sql)
	})

	t.Run("array_append", func(t *testing.T) {
		q := Select(fn.ArrayAppend(Array(Int(1), Int(2)), Int(3)))

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT array_append(ARRAY[1,2], 3)
		`, sql)
	})

	t.Run("array_prepend", func(t *testing.T) {
		q := Select(fn.ArrayPrepend(Int(1), Array(Int(2), Int(3))))

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT array_prepend(1, ARRAY[2,3])
		`, sql)
	})

	t.Run("array_cat", func(t *testing.T) {
		q := Select(fn.ArrayCat(Array(Int(1), Int(2)), Array(Int(3), Int(4))))

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT array_cat(ARRAY[1,2], ARRAY[3,4])
		`, sql)
	})

	t.Run("array_dims", func(t *testing.T) {
		q := Select(fn.ArrayDims(Array(Int(1), Int(2), Int(3))))

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT array_dims(ARRAY[1,2,3])
		`, sql)
	})

	t.Run("array_ndims", func(t *testing.T) {
		q := Select(fn.ArrayNdims(Array(Int(1), Int(2), Int(3))))

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT array_ndims(ARRAY[1,2,3])
		`, sql)
	})

	t.Run("array_length", func(t *testing.T) {
		q := Select(fn.ArrayLength(Array(Int(1), Int(2), Int(3)), Int(1)))

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT array_length(ARRAY[1,2,3], 1)
		`, sql)
	})

	t.Run("array_lower", func(t *testing.T) {
		q := Select(fn.ArrayLower(Array(Int(1), Int(2), Int(3)), Int(1)))

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT array_lower(ARRAY[1,2,3], 1)
		`, sql)
	})

	t.Run("array_upper", func(t *testing.T) {
		q := Select(fn.ArrayUpper(Array(Int(1), Int(2), Int(3)), Int(1)))

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT array_upper(ARRAY[1,2,3], 1)
		`, sql)
	})

	t.Run("array_remove", func(t *testing.T) {
		q := Select(fn.ArrayRemove(Array(Int(1), Int(2), Int(3), Int(2)), Int(2)))

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT array_remove(ARRAY[1,2,3,2], 2)
		`, sql)
	})

	t.Run("array_replace", func(t *testing.T) {
		q := Select(fn.ArrayReplace(Array(Int(1), Int(2), Int(3)), Int(2), Int(99)))

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT array_replace(ARRAY[1,2,3], 2, 99)
		`, sql)
	})

	t.Run("array_position", func(t *testing.T) {
		t.Run("without start", func(t *testing.T) {
			q := Select(fn.ArrayPosition(Array(String("a"), String("b"), String("c")), String("b")))

			sql, _, _ := Build(q).ToSQL()
			testhelper.AssertSQLEquals(t, `
			SELECT array_position(ARRAY['a','b','c'], 'b')
			`, sql)
		})

		t.Run("with start", func(t *testing.T) {
			q := Select(fn.ArrayPosition(Array(String("a"), String("b"), String("c"), String("b")), String("b"), Int(3)))

			sql, _, _ := Build(q).ToSQL()
			testhelper.AssertSQLEquals(t, `
			SELECT array_position(ARRAY['a','b','c','b'], 'b', 3)
			`, sql)
		})
	})

	t.Run("array_positions", func(t *testing.T) {
		q := Select(fn.ArrayPositions(Array(String("a"), String("b"), String("c"), String("b")), String("b")))

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT array_positions(ARRAY['a','b','c','b'], 'b')
		`, sql)
	})

	t.Run("array_to_string", func(t *testing.T) {
		t.Run("without null string", func(t *testing.T) {
			q := Select(fn.ArrayToString(Array(String("a"), String("b"), String("c")), String(",")))

			sql, _, _ := Build(q).ToSQL()
			testhelper.AssertSQLEquals(t, `
			SELECT array_to_string(ARRAY['a','b','c'], ',')
			`, sql)
		})

		t.Run("with null string", func(t *testing.T) {
			q := Select(fn.ArrayToString(Array(String("a"), String("b"), Null()), String(","), String("*")))

			sql, _, _ := Build(q).ToSQL()
			testhelper.AssertSQLEquals(t, `
			SELECT array_to_string(ARRAY['a','b',NULL], ',', '*')
			`, sql)
		})
	})

	t.Run("string_to_array", func(t *testing.T) {
		t.Run("without null string", func(t *testing.T) {
			q := Select(fn.StringToArray(String("a,b,c"), String(",")))

			sql, _, _ := Build(q).ToSQL()
			testhelper.AssertSQLEquals(t, `
			SELECT string_to_array('a,b,c', ',')
			`, sql)
		})

		t.Run("with null string", func(t *testing.T) {
			q := Select(fn.StringToArray(String("a,b,*"), String(","), String("*")))

			sql, _, _ := Build(q).ToSQL()
			testhelper.AssertSQLEquals(t, `
			SELECT string_to_array('a,b,*', ',', '*')
			`, sql)
		})
	})

	t.Run("array_fill", func(t *testing.T) {
		t.Run("without lower bounds", func(t *testing.T) {
			q := Select(fn.ArrayFill(String("x"), Array(Int(3))))

			sql, _, _ := Build(q).ToSQL()
			testhelper.AssertSQLEquals(t, `
			SELECT array_fill('x', ARRAY[3])
			`, sql)
		})

		t.Run("with lower bounds", func(t *testing.T) {
			q := Select(fn.ArrayFill(String("x"), Array(Int(3), Int(2)), Array(Int(2), Int(5))))

			sql, _, _ := Build(q).ToSQL()
			testhelper.AssertSQLEquals(t, `
			SELECT array_fill('x', ARRAY[3,2], ARRAY[2,5])
			`, sql)
		})
	})
}
