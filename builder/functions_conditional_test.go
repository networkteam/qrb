package builder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestCase(t *testing.T) {
	t.Run("no expression", func(t *testing.T) {
		b := qrb.Case().
			When(qrb.N("a").Eq(qrb.Int(1))).Then(qrb.String("one")).
			End()

		testhelper.AssertSQLWriterEquals(
			t,
			`
			CASE WHEN a = 1 THEN 'one'
			END`,
			nil,
			b,
		)
	})

	t.Run("with expression", func(t *testing.T) {
		b := qrb.Case(qrb.N("a")).
			When(qrb.Int(1)).Then(qrb.String("one")).
			When(qrb.Int(2)).Then(qrb.String("two")).
			Else(qrb.String("other")).
			End()

		testhelper.AssertSQLWriterEquals(
			t,
			`
			CASE a
				WHEN 1 THEN 'one'
				WHEN 2 THEN 'two'
				ELSE 'other'
			END`,
			nil,
			b,
		)
	})

	t.Run("add op", func(t *testing.T) {
		b := qrb.Case().
			When(qrb.N("a").Eq(qrb.Int(1))).Then(qrb.String("one")).
			End().Concat(qrb.String("-dings"))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			CASE WHEN a = 1 THEN 'one' END || '-dings'`,
			nil,
			b,
		)
	})

	t.Run("no when then", func(t *testing.T) {
		b := qrb.Case().End()

		_, _, err := qrb.Build(b).ToSQL()
		assert.EqualError(t, err, "case: no conditions given")
	})
}
