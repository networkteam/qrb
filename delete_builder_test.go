package qrb_test

import (
	"testing"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestDeleteBuilder(t *testing.T) {
	t.Run("examples", func(t *testing.T) {
		// From https://www.postgresql.org/docs/15/sql-delete.html#id-1.9.3.100.9

		t.Run("example 1", func(t *testing.T) {
			q := qrb.
				DeleteFrom("films").
				Where(qrb.N("kind").Neq(qrb.String("Musical")))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				DELETE FROM films WHERE kind <> 'Musical'
				`,
				nil,
				q,
			)
		})

		t.Run("example 2", func(t *testing.T) {
			q := qrb.
				DeleteFrom("films")

			testhelper.AssertSQLWriterEquals(
				t,
				`
				DELETE FROM films
				`,
				nil,
				q,
			)
		})
	})
}
