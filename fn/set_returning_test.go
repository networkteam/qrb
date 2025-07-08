package fn_test

import (
	"testing"

	. "github.com/networkteam/qrb"
	"github.com/networkteam/qrb/fn"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestSetReturningFunctions(t *testing.T) {
	t.Run("generate_series with two arguments", func(t *testing.T) {
		// Example: SELECT * FROM generate_series(2,4);
		q := Select(N("*")).
			From(fn.GenerateSeries(Int(2), Int(4))).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT * FROM generate_series(2, 4)
		`, sql)
	})

	t.Run("generate_series with three arguments", func(t *testing.T) {
		// Example: SELECT * FROM generate_series(5,1,-2);
		q := Select(N("*")).
			From(fn.GenerateSeries(Int(5), Int(1), Int(-2))).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT * FROM generate_series(5, 1, -2)
		`, sql)
	})

	t.Run("generate_series with alias", func(t *testing.T) {
		// Example: SELECT current_date + s.a AS dates FROM generate_series(0,14,7) AS s(a);
		q := Select(Func("current_date").Plus(N("s.a"))).
			From(fn.GenerateSeries(Int(0), Int(14), Int(7))).As("s").ColumnAliases("a").
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT current_date() + s.a FROM generate_series(0, 14, 7) AS s (a)
		`, sql)
	})

	t.Run("generate_series with timestamp", func(t *testing.T) {
		// Example: SELECT * FROM generate_series('2008-03-01 00:00'::timestamp, '2008-03-04 12:00', '10 hours');
		q := Select(N("*")).
			From(fn.GenerateSeries(
				Arg("2008-03-01 00:00").Cast("timestamp"),
				String("2008-03-04 12:00"),
				String("10 hours"),
			)).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT * FROM generate_series($1::timestamp, '2008-03-04 12:00', '10 hours')
		`, sql)
	})

	t.Run("generate_subscripts basic", func(t *testing.T) {
		// Example: SELECT generate_subscripts('{NULL,1,NULL,2}'::int[], 1) AS s;
		q := Select(fn.GenerateSubscripts(
			Arg("{NULL,1,NULL,2}").Cast("int[]"),
			Int(1),
		).As("s")).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT generate_subscripts($1::int[], 1) AS s
		`, sql)
	})

	t.Run("generate_subscripts in FROM clause", func(t *testing.T) {
		// Simplified example from docs
		q := Select(N("a")).As("array").
			Select(N("s")).As("subscript").
			Select(N("a").Subscript(N("s"))).As("value").
			From(
				Select(
					fn.GenerateSubscripts(N("a"), Int(1)).As("s"),
					N("a"),
				).From(N("arrays")),
			).As("foo").SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT a AS array, s AS subscript, a[s] AS value 
		FROM (SELECT generate_subscripts(a, 1) AS s, a FROM arrays) AS foo
		`, sql)
	})

	t.Run("generate_subscripts with reverse", func(t *testing.T) {
		q := Select(N("*")).
			From(fn.GenerateSubscripts(
				N("some_array"),
				Int(1),
				Bool(true),
			)).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT * FROM generate_subscripts(some_array, 1, true)
		`, sql)
	})

	t.Run("multiple generate_subscripts in FROM", func(t *testing.T) {
		// Example from docs: unnest2 function implementation
		// Simplified test for multiple FROM items
		q := Select(
			N("i"),
			N("j"),
		).From(
			fn.GenerateSubscripts(N("some_array"), Int(1)),
		).As("g1").ColumnAliases("i").
			From(fn.GenerateSubscripts(N("some_array"), Int(2))).As("g2").ColumnAliases("j").
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT i, j FROM generate_subscripts(some_array, 1) AS g1 (i), generate_subscripts(some_array, 2) AS g2 (j)
		`, sql)
	})
}
