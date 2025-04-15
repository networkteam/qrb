package builder_test

import (
	"testing"

	"github.com/networkteam/qrb/builder"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestAggBuilder(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		b := builder.Agg("my_agg", []builder.Exp{builder.N("foo")})
		testhelper.AssertSQLWriterEquals(t, `my_agg(foo)`, nil, b)
	})

	t.Run("basic", func(t *testing.T) {
		b := builder.Agg("my_agg", []builder.Exp{builder.N("foo"), builder.N("bar")}).
			Distinct().
			OrderBy(builder.N("foo")).NullsFirst().Asc().
			OrderBy(builder.N("bar")).NullsLast().Desc().
			Filter(builder.N("foo").Gt(builder.Int(1)))
		testhelper.AssertSQLWriterEquals(t, `my_agg(DISTINCT foo,bar ORDER BY foo ASC NULLS FIRST,bar DESC NULLS LAST) FILTER (WHERE foo > 1)`, nil, b)
	})
}
