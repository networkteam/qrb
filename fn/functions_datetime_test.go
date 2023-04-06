package fn_test

import (
	"testing"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/fn"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestDatetime(t *testing.T) {
	t.Run("extract", func(t *testing.T) {
		b := fn.Extract("YEARS", qrb.Func("age", qrb.Arg("2023-03-30").Cast("date"), qrb.N("u.birthday")))

		testhelper.AssertSQLWriterEquals(
			t,
			"EXTRACT(YEARS FROM age($1::date,u.birthday))",
			[]any{"2023-03-30"},
			b,
		)
	})
}
