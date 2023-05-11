package fn_test

import (
	"testing"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/builder"
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

	t.Run("functions", func(t *testing.T) {
		tt := []struct {
			name        string
			fn          builder.Exp
			expectedSQL string
		}{
			{
				name:        "array_agg",
				fn:          fn.ArrayAgg(qrb.N("title")).OrderBy(qrb.N("title")).Desc(),
				expectedSQL: "array_agg(title ORDER BY title DESC)",
			},
			{
				name:        "avg",
				fn:          fn.Avg(qrb.N("price")),
				expectedSQL: "avg(price)",
			},
			{
				name:        "bit_and",
				fn:          fn.BitAnd(qrb.N("flags")),
				expectedSQL: "bit_and(flags)",
			},
			{
				name:        "bit_or",
				fn:          fn.BitOr(qrb.N("flags")),
				expectedSQL: "bit_or(flags)",
			},
			{
				name:        "bit_xor",
				fn:          fn.BitXor(qrb.N("flags")),
				expectedSQL: "bit_xor(flags)",
			},
			{
				name:        "bool_and",
				fn:          fn.BoolAnd(qrb.N("active")),
				expectedSQL: "bool_and(active)",
			},
			{
				name:        "bool_or",
				fn:          fn.BoolOr(qrb.N("active")),
				expectedSQL: "bool_or(active)",
			},
			{
				name:        "count",
				fn:          fn.Count(qrb.N("*")),
				expectedSQL: "count(*)",
			},
			{
				name:        "json_agg",
				fn:          fn.JsonAgg(qrb.N("title")),
				expectedSQL: `json_agg(title)`,
			},
			{
				name:        "jsonb_agg",
				fn:          fn.JsonbAgg(qrb.N("title")),
				expectedSQL: `jsonb_agg(title)`,
			},
			{
				name:        "json_object_agg",
				fn:          fn.JsonObjectAgg(qrb.N("title"), qrb.N("price")),
				expectedSQL: `json_object_agg(title, price)`,
			},
			{
				name:        "string_agg",
				fn:          fn.StringAgg(qrb.N("title"), qrb.String(",")).OrderBy(qrb.N("title")).Desc(),
				expectedSQL: "string_agg(title, ',' ORDER BY title DESC)",
			},
			{
				name:        "max",
				fn:          fn.Max(qrb.N("price")),
				expectedSQL: "max(price)",
			},
			{
				name:        "min",
				fn:          fn.Min(qrb.N("price")),
				expectedSQL: "min(price)",
			},
			{
				name:        "range_agg",
				fn:          fn.RangeAgg(qrb.N("price")),
				expectedSQL: "range_agg(price)",
			},
			{
				name:        "range_intersect_agg",
				fn:          fn.RangeIntersectAgg(qrb.N("price")),
				expectedSQL: "range_intersect_agg(price)",
			},
			{
				name:        "sum",
				fn:          fn.Sum(qrb.N("price")),
				expectedSQL: "sum(price)",
			},
			{
				name:        "xmlagg",
				fn:          fn.Xmlagg(qrb.N("title")),
				expectedSQL: "xmlagg(title)",
			},
			{
				name:        "mode",
				fn:          fn.Mode().WithinGroup().OrderBy(qrb.N("price")).Asc(),
				expectedSQL: "mode() WITHIN GROUP (ORDER BY price ASC)",
			},
			{
				name:        "percentile_cont",
				fn:          fn.PercentileCont(qrb.Float(0.5)).WithinGroup().OrderBy(qrb.N("price")).Asc(),
				expectedSQL: "percentile_cont(0.5) WITHIN GROUP (ORDER BY price ASC)",
			},
			{
				name:        "percentile_disc",
				fn:          fn.PercentileDisc(qrb.Float(0.5)).WithinGroup().OrderBy(qrb.N("price")).Asc(),
				expectedSQL: "percentile_disc(0.5) WITHIN GROUP (ORDER BY price ASC)",
			},
			{
				name:        "rank",
				fn:          fn.Rank().WithinGroup().OrderBy(qrb.N("price")).Asc(),
				expectedSQL: "rank() WITHIN GROUP (ORDER BY price ASC)",
			},
			{
				name:        "dense_rank",
				fn:          fn.DenseRank().WithinGroup().OrderBy(qrb.N("price")).Asc(),
				expectedSQL: "dense_rank() WITHIN GROUP (ORDER BY price ASC)",
			},
			{
				name:        "percent_rank",
				fn:          fn.PercentRank().WithinGroup().OrderBy(qrb.N("price")).Asc(),
				expectedSQL: "percent_rank() WITHIN GROUP (ORDER BY price ASC)",
			},
			{
				name:        "cume_dist",
				fn:          fn.CumeDist().WithinGroup().OrderBy(qrb.N("price")).Asc(),
				expectedSQL: "cume_dist() WITHIN GROUP (ORDER BY price ASC)",
			},
			{
				name:        "grouping",
				fn:          fn.Grouping(qrb.N("price"), qrb.N("title")),
				expectedSQL: "GROUPING(price,title)",
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				sql, _, _ := qrb.Build(tc.fn).ToSQL()
				testhelper.AssertSQLEquals(t, tc.expectedSQL, sql)
			})
		}
	})
}
