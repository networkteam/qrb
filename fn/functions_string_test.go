package fn_test

import (
	"testing"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/fn"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestLetterCases(t *testing.T) {
	t.Run("lower", func(t *testing.T) {
		q := qrb.Select(fn.Lower(qrb.N("a"))).From(qrb.N("table"))
		sql, _, _ := qrb.Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, "SELECT lower(a) FROM table", sql)
	})

	t.Run("lower with arg", func(t *testing.T) {
		q := qrb.Select(qrb.N("id")).From(qrb.N("table")).Where(fn.Lower(qrb.N("name")).Eq(fn.Lower(qrb.Arg("foo"))))
		sql, _, _ := qrb.Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, "SELECT id FROM table WHERE lower(name) = lower($1)", sql)
	})

	t.Run("upper", func(t *testing.T) {
		q := qrb.Select(fn.Upper(qrb.N("a")).Eq(qrb.Arg("foo"))).From(qrb.N("table"))
		sql, _, _ := qrb.Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, "SELECT upper(a) = $1 FROM table", sql)
	})

	t.Run("init cap", func(t *testing.T) {
		q := qrb.Select(fn.Initcap(qrb.N("a"))).From(qrb.N("table"))
		sql, _, _ := qrb.Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, "SELECT initcap(a) FROM table", sql)
	})
}
