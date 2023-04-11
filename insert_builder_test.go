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
				InsertInto(qrb.N("films")).
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

		t.Run("example 2", func(t *testing.T) {
			q := qrb.
				InsertInto(qrb.N("films")).
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

		t.Run("example 3a", func(t *testing.T) {
			q := qrb.
				InsertInto(qrb.N("films")).
				Values(qrb.String("UA502"), qrb.String("Bananas"), qrb.Int(105), qrb.Default(), qrb.String("Comedy"), qrb.String("82 minutes"))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				INSERT INTO films VALUES
    				('UA502', 'Bananas', 105, DEFAULT, 'Comedy', '82 minutes')
				`,
				nil,
				q,
			)
		})

		t.Run("example 3b", func(t *testing.T) {
			q := qrb.
				InsertInto(qrb.N("films")).
				ColumnNames("code", "title", "did", "date_prod", "kind").
				Values(qrb.String("T_601"), qrb.String("Yojimbo"), qrb.Int(106), qrb.Default(), qrb.String("Drama"))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				INSERT INTO films (code, title, did, date_prod, kind)
    				VALUES ('T_601', 'Yojimbo', 106, DEFAULT, 'Drama')
				`,
				nil,
				q,
			)
		})

		t.Run("example 4", func(t *testing.T) {
			q := qrb.
				InsertInto(qrb.N("films")).
				DefaultValues()

			testhelper.AssertSQLWriterEquals(
				t,
				`
				INSERT INTO films DEFAULT VALUES
				`,
				nil,
				q,
			)
		})

		t.Run("example 5", func(t *testing.T) {
			q := qrb.
				InsertInto(qrb.N("films")).
				ColumnNames("code", "title", "did", "date_prod", "kind").
				Values(qrb.String("B6717"), qrb.String("Tampopo"), qrb.Int(110), qrb.String("1985-02-10"), qrb.String("Comedy")).
				Values(qrb.String("HG120"), qrb.String("The Dinner Game"), qrb.Int(140), qrb.Default(), qrb.String("Comedy"))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				INSERT INTO films (code, title, did, date_prod, kind) VALUES
					('B6717', 'Tampopo', 110, '1985-02-10', 'Comedy'),
					('HG120', 'The Dinner Game', 140, DEFAULT, 'Comedy')
				`,
				nil,
				q,
			)
		})

		t.Run("example 6", func(t *testing.T) {
			q := qrb.
				InsertInto(qrb.N("films")).
				Query(qrb.Select(qrb.N("*")).From(qrb.N("tmp_films")).Where(qrb.N("date_prod").Lt(qrb.String("2004-05-07"))))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				INSERT INTO films SELECT * FROM tmp_films WHERE date_prod < '2004-05-07'
				`,
				nil,
				q,
			)
		})

		t.Run("example 7a", func(t *testing.T) {
			q := qrb.
				InsertInto(qrb.N("tictactoe")).
				ColumnNames("game", "board[1:3][1:3]").
				Values(qrb.Int(1), qrb.String(`{{" "," "," "},{" "," "," "},{" "," "," "}}`))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				INSERT INTO tictactoe (game, board[1:3][1:3])
    				VALUES (1, '{{" "," "," "},{" "," "," "},{" "," "," "}}')
				`,
				nil,
				q,
			)
		})

		t.Run("example 7b", func(t *testing.T) {
			q := qrb.
				InsertInto(qrb.N("tictactoe")).
				ColumnNames("game", "board").
				Values(qrb.Int(2), qrb.String(`{{X," "," "},{" ",O," "},{" ",X," "}}`))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				INSERT INTO tictactoe (game, board)
    				VALUES (2, '{{X," "," "},{" ",O," "},{" ",X," "}}')
				`,
				nil,
				q,
			)
		})

		t.Run("example 8", func(t *testing.T) {
			q := qrb.
				InsertInto(qrb.N("distributors")).
				ColumnNames("did", "dname").
				Values(qrb.Default(), qrb.String("XYZ Widgets")).
				Returning(qrb.N("did"))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				INSERT INTO distributors (did, dname) VALUES (DEFAULT, 'XYZ Widgets')
   					RETURNING did
				`,
				nil,
				q,
			)
		})

		t.Run("example 9", func(t *testing.T) {
			q := qrb.
				With("upd").As(
				qrb.Update(qrb.N("employees")).
					Set("sales_count", qrb.N("sales_count").Plus(qrb.Int(1))).
					Where(qrb.N("id").Eq(
						qrb.Select(qrb.N("sales_person")).From(qrb.N("accounts")).Where(qrb.N("name").Eq(qrb.String("Acme Corporation"))),
					)).
					Returning(qrb.N("*")),
			).
				InsertInto(qrb.N("employees_log")).
				Query(qrb.Select(qrb.N("*"), qrb.N("current_timestamp")).From(qrb.N("upd")))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				WITH upd AS (
				  UPDATE employees SET sales_count = sales_count + 1 WHERE id =
					(SELECT sales_person FROM accounts WHERE name = 'Acme Corporation')
					RETURNING *
				)
				INSERT INTO employees_log SELECT *, current_timestamp FROM upd
				`,
				nil,
				q,
			)
		})

		t.Run("example 10", func(t *testing.T) {
			q := qrb.InsertInto(qrb.N("distributors")).ColumnNames("did", "dname").
				Values(qrb.Int(5), qrb.String("Gizmo Transglobal")).
				Values(qrb.Int(6), qrb.String("Associated Computing,Inc")).
				OnConflict(qrb.N("did")).DoUpdate().Set("dname", qrb.N("EXCLUDED.dname"))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				INSERT INTO distributors (did, dname)
				VALUES (5, 'Gizmo Transglobal'), (6, 'Associated Computing,Inc')
				ON CONFLICT (did) DO UPDATE SET dname = EXCLUDED.dname
				`,
				nil,
				q,
			)
		})

		t.Run("example 11", func(t *testing.T) {
			q := qrb.InsertInto(qrb.N("distributors")).ColumnNames("did", "dname").
				Values(qrb.Int(7), qrb.String("Redline GmbH")).
				OnConflict(qrb.N("did")).DoNothing()

			testhelper.AssertSQLWriterEquals(
				t,
				`
				INSERT INTO distributors (did, dname) VALUES (7, 'Redline GmbH')
    			ON CONFLICT (did) DO NOTHING
				`,
				nil,
				q,
			)
		})

		t.Run("example 12a", func(t *testing.T) {
			q := qrb.InsertInto(qrb.N("distributors")).As("d").ColumnNames("did", "dname").
				Values(qrb.Int(8), qrb.String("Anvil Distribution")).
				OnConflict(qrb.N("did")).DoUpdate().
				Set("dname", qrb.N("EXCLUDED.dname").Concat(qrb.String(" (formerly ")).Concat(qrb.N("d.dname")).Concat(qrb.String(")"))).
				Where(qrb.N("d.zipcode").Neq(qrb.String("21201")))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				INSERT INTO distributors AS d (did, dname) VALUES (8, 'Anvil Distribution')
				ON CONFLICT (did) DO UPDATE
				SET dname = EXCLUDED.dname || ' (formerly ' || d.dname || ')'
				WHERE d.zipcode <> '21201'
				`,
				nil,
				q,
			)
		})

		t.Run("example 12b", func(t *testing.T) {
			q := qrb.InsertInto(qrb.N("distributors")).ColumnNames("did", "dname").
				Values(qrb.Int(9), qrb.String("Antwerp Design")).
				OnConflict().OnConstraint("distributors_pkey").DoNothing()

			testhelper.AssertSQLWriterEquals(
				t,
				`
				INSERT INTO distributors (did, dname) VALUES (9, 'Antwerp Design')
    			ON CONFLICT ON CONSTRAINT distributors_pkey DO NOTHING
				`,
				nil,
				q,
			)
		})

		t.Run("example 13", func(t *testing.T) {
			q := qrb.InsertInto(qrb.N("distributors")).ColumnNames("did", "dname").
				Values(qrb.Int(10), qrb.String("Conrad International")).
				OnConflict(qrb.N("did")).Where(qrb.N("is_active")).DoNothing()

			testhelper.AssertSQLWriterEquals(
				t,
				`
				INSERT INTO distributors (did, dname) VALUES (10, 'Conrad International')
    				ON CONFLICT (did) WHERE is_active DO NOTHING
				`,
				nil,
				q,
			)
		})
	})

	t.Run("set map", func(t *testing.T) {
		q := qrb.
			InsertInto(qrb.N("films")).
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

	t.Run("values with args", func(t *testing.T) {
		q := qrb.
			InsertInto(qrb.N("films")).
			ColumnNames("code", "date_prod", "did", "kind", "length", "title").
			Values(qrb.Args("UA502", "1971-07-13", 105, "Comedy", "82 minutes", "Bananas")...)

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
