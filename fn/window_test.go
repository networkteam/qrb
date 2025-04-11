package fn_test

import (
	"testing"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/fn"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestWindowFunctions(t *testing.T) {
	t.Run("row_number", func(t *testing.T) {
		q := qrb.Select(fn.RowNumber().Over()).From(qrb.N("table"))
		sql, _, _ := qrb.Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, "SELECT row_number() OVER () FROM table", sql)
	})
}
