package qrb_test

import (
	"testing"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/builder"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestUpdateBuilder(t *testing.T) {
	t.Run("examples", func(t *testing.T) {
		// From https://www.postgresql.org/docs/15/sql-update.html#id-1.9.3.183.9

		t.Run("example 1", func(t *testing.T) {
			q := qrb.
				Update("films").
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
				Update("weather").
				Set("temp_lo", qrb.N("temp_lo").Op(builder.OpAdd, qrb.Int(1))).
				Set("temp_hi", qrb.N("temp_lo").Op(builder.OpAdd, qrb.Int(15))).
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
				Update("weather").
				Set("temp_lo", qrb.N("temp_lo").Op(builder.OpAdd, qrb.Int(1))).
				Set("temp_hi", qrb.N("temp_lo").Op(builder.OpAdd, qrb.Int(15))).
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

	t.Run("set map", func(t *testing.T) {
		q := qrb.
			Update("films").
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
