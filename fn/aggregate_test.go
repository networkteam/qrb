package fn_test

import (
	"testing"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/fn"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestAggregateExpressions(t *testing.T) {
	// See https://www.postgresql.org/docs/15/sql-expressions.html#SYNTAX-AGGREGATES

	t.Run("example 1.1", func(t *testing.T) {
		q := qrb.Select(fn.ArrayAgg(qrb.N("a")).OrderBy(qrb.N("b")).Desc()).From(qrb.N("table"))
		sql, _, _ := qrb.Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, "SELECT array_agg(a ORDER BY b DESC) FROM table", sql)
	})

	t.Run("example 1.2", func(t *testing.T) {
		q := qrb.Select(fn.StringAgg(qrb.N("a"), qrb.String(",")).OrderBy(qrb.N("a"))).From(qrb.N("table"))
		sql, _, _ := qrb.Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, "SELECT string_agg(a,',' ORDER BY a) FROM table", sql)
	})

	t.Run("example 1.3", func(t *testing.T) {
		// SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY income) FROM households
		q := qrb.Select(fn.PercentileCont(qrb.Float(0.5)).WithinGroup().OrderBy(qrb.N("income"))).From(qrb.N("households"))
		sql, _, _ := qrb.Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, "SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY income) FROM households", sql)
	})

	// See https://www.postgresql.org/docs/15/functions-aggregate.html#FUNCTIONS-GROUPING-TABLE

	t.Run("example 2.1", func(t *testing.T) {
		q := qrb.Select(qrb.N("make"), qrb.N("model"), fn.Grouping(qrb.N("make"), qrb.N("model")), fn.Sum(qrb.N("sales"))).
			From(qrb.N("items_sold")).
			GroupBy().Rollup(qrb.Exps(qrb.N("make"), qrb.N("model")))
		sql, _, _ := qrb.Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, "SELECT make, model, GROUPING(make,model), sum(sales) FROM items_sold GROUP BY ROLLUP (make,model)", sql)
	})

	// See https://tapoueh.org/blog/2017/11/the-mode-ordered-set-aggregate-function/

	t.Run("3.1", func(t *testing.T) {
		// There might be the case that tables or columns are named with uppercase letters or use keywords.
		// The qrb.N function can be used with a quoted identifier.
		q := qrb.Select(qrb.N(`"Title"`)).As(`"Album"`).
			Select(fn.StringAgg(qrb.N(`"Genre"."Name"`), qrb.String(",")).Distinct().OrderBy(qrb.N(`"Genre"."Name"`))).As(`"Genres"`).
			From(qrb.N(`"Track"`)).
			Join(qrb.N(`"Genre"`)).Using(`"GenreId"`).
			Join(qrb.N(`"Album"`)).Using(`"AlbumId"`).
			GroupBy(qrb.N(`"Title"`)).
			Having(fn.Count(qrb.N(`"Genre"."Name"`)).Distinct().Gt(qrb.Int(1)))

		sql, _, _ := qrb.Build(q).ToSQL()

		testhelper.AssertSQLEquals(t, `
		SELECT "Title" AS "Album",
			   string_agg(
					   DISTINCT "Genre"."Name", ','
					   ORDER BY "Genre"."Name"
				   )
					 AS "Genres"
		FROM "Track"
				 JOIN "Genre" USING ("GenreId")
				 JOIN "Album" USING ("AlbumId")
		GROUP BY "Title"
		HAVING count(DISTINCT "Genre"."Name") > 1
		`, sql)
	})
}
