package qrb_test

import (
	"testing"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestInsertBuilder(t *testing.T) {
	t.Run("examples", func(t *testing.T) {
		// From https://www.postgresql.org/docs/15/sql-insert.html#id-1.9.3.152.9

		t.Run("example 1", func(t *testing.T) {
			q := qrb.
				InsertInto("films").
				Values(qrb.String("UA502"), qrb.String("Bananas"), qrb.Int(105), qrb.String("1971-07-13"), qrb.String("Comedy"), qrb.String("82 minutes"))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				INSERT INTO films VALUES
    				('UA502', 'Bananas', 105, '1971-07-13', 'Comedy', '82 minutes')
				`,
				nil,
				q,
			)
		})

		t.Run("examples", func(t *testing.T) {
			q := qrb.
				InsertInto("films").
				ColumnNames("code", "title", "did", "date_prod", "kind").
				Values(qrb.String("T_601"), qrb.String("Yojimbo"), qrb.Int(106), qrb.String("1961-06-16"), qrb.String("Drama"))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				INSERT INTO films (code, title, did, date_prod, kind)
					VALUES ('T_601', 'Yojimbo', 106, '1961-06-16', 'Drama')
				`,
				nil,
				q,
			)
		})
	})

	t.Run("set map", func(t *testing.T) {
		q := qrb.
			InsertInto("films").
			SetMap(map[string]any{
				"code":      "UA502",
				"title":     "Bananas",
				"did":       105,
				"date_prod": "1971-07-13",
				"kind":      "Comedy",
				"length":    "82 minutes",
			})

		testhelper.AssertSQLWriterEquals(
			t,
			`
				INSERT INTO films (code,date_prod,did,kind,length,title) VALUES
    				($1, $2, $3, $4, $5, $6)
				`,
			[]any{"UA502", "1971-07-13", 105, "Comedy", "82 minutes", "Bananas"},
			q,
		)
	})
}
