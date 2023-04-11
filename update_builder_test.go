package qrb_test

import (
	"testing"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestUpdateBuilder(t *testing.T) {
	t.Run("examples", func(t *testing.T) {
		// From https://www.postgresql.org/docs/15/sql-update.html#id-1.9.3.183.9

		t.Run("example 1", func(t *testing.T) {
			q := qrb.
				Update(qrb.N("films")).
				Set("kind", qrb.String("Dramatic")).
				Where(qrb.N("kind").Eq(qrb.String("Drama")))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				UPDATE films SET kind = 'Dramatic' WHERE kind = 'Drama'
				`,
				nil,
				q,
			)
		})

		t.Run("example 2", func(t *testing.T) {
			q := qrb.
				Update(qrb.N("weather")).
				Set("temp_lo", qrb.N("temp_lo").Plus(qrb.Int(1))).
				Set("temp_hi", qrb.N("temp_lo").Plus(qrb.Int(15))).
				Set("prcp", qrb.Default()).
				Where(qrb.And(
					qrb.N("city").Eq(qrb.String("San Francisco")),
					qrb.N("date").Eq(qrb.String("2003-07-03")),
				))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				UPDATE weather SET temp_lo = temp_lo + 1, temp_hi = temp_lo + 15, prcp = DEFAULT
				  WHERE city = 'San Francisco' AND date = '2003-07-03'
				`,
				nil,
				q,
			)
		})

		t.Run("example 2", func(t *testing.T) {
			q := qrb.
				Update(qrb.N("weather")).
				Set("temp_lo", qrb.N("temp_lo").Plus(qrb.Int(1))).
				Set("temp_hi", qrb.N("temp_lo").Plus(qrb.Int(15))).
				Set("prcp", qrb.Default()).
				Where(qrb.And(
					qrb.N("city").Eq(qrb.String("San Francisco")),
					qrb.N("date").Eq(qrb.String("2003-07-03")),
				)).
				Returning(qrb.N("temp_lo")).
				Returning(qrb.N("temp_hi")).
				Returning(qrb.N("prcp"))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				UPDATE weather SET temp_lo = temp_lo + 1, temp_hi = temp_lo + 15, prcp = DEFAULT
				  WHERE city = 'San Francisco' AND date = '2003-07-03'
				  RETURNING temp_lo, temp_hi, prcp
				`,
				nil,
				q,
			)
		})
	})

	t.Run("with", func(t *testing.T) {
		// Example borrowed from https://medium.com/@mnu/update-a-postgresql-table-using-a-with-query-648eefaae2a6

		q := qrb.With("line_journey_pattern").As(
			qrb.Select(qrb.N("jp.id")).As("journey_pattern_id").
				Select(qrb.N("l.name")).As("line_name").
				From(qrb.N("journey_patterns")).As("jp").
				Join(qrb.N("routes")).As("r").On(qrb.N("jp.route_id").Eq(qrb.N("r.id"))).
				Join(qrb.N("lines")).As("l").On(qrb.N("r.line_id").Eq(qrb.N("l.id"))).
				Where(qrb.And(
					qrb.N("l.name").IsNotNull(),
					qrb.N("l.name").Neq(qrb.String("")),
				)),
		).Update(qrb.N("journey_patterns")).As("jp").
			Set("name", qrb.N("ljp.line_name").Concat(qrb.String(" - ")).Concat(qrb.N("jp.name"))).
			From(qrb.N("line_journey_pattern")).As("ljp").
			Where(qrb.N("ljp.journey_pattern_id").Eq(qrb.N("jp.id")))

		testhelper.AssertSQLWriterEquals(
			t,
			`
				WITH line_journey_pattern AS (
					SELECT jp.id AS journey_pattern_id, l.name AS line_name
					FROM journey_patterns AS jp
							 JOIN routes AS r ON jp.route_id = r.id
							 JOIN lines AS l ON r.line_id = l.id
					WHERE l.name IS NOT NULL
					  AND l.name <> ''
				)
				UPDATE journey_patterns AS jp
				SET name = ljp.line_name || ' - ' || jp.name
				FROM line_journey_pattern AS ljp
				WHERE ljp.journey_pattern_id = jp.id
				`,
			nil,
			q,
		)
	})

	t.Run("set map", func(t *testing.T) {
		q := qrb.
			Update(qrb.N("films")).
			SetMap(map[string]any{
				"code": "UA502",
				"kind": "Comedy",
			}).
			Where(qrb.N("kind").Eq(qrb.String("Drama")))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			UPDATE films SET code = $1, kind = $2 WHERE kind = 'Drama'	
			`,
			[]any{"UA502", "Comedy"},
			q,
		)
	})
}
