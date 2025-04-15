package fn_test

import (
	"testing"

	. "github.com/networkteam/qrb"
	"github.com/networkteam/qrb/fn"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestArrayFunctions(t *testing.T) {
	t.Run("unnest", func(t *testing.T) {
		q := Select(N("*")).
			From(
				fn.Unnest(
					Array(Int(1), Int(2)),
					Array(String("foo"), String("bar"), String("baz")),
				),
			).As("x").ColumnAliases("a", "b").
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT * FROM unnest(ARRAY[1,2], ARRAY['foo','bar','baz']) AS x (a,b)
		`, sql)
	})
}
