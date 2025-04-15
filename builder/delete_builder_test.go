package builder_test

import (
	"testing"

	"github.com/networkteam/qrb/builder"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestDeleteBuilder(t *testing.T) {
	t.Run("", func(t *testing.T) {
		q := builder.DeleteFrom(builder.N("employees")).As("e").
			Using(builder.N("companies")).As("c").
			Where(builder.And(
				builder.N("e.company_id").Eq(builder.N("c.id")),
				builder.N("c.deleted").Eq(builder.Bool(true)),
			))
		testhelper.AssertSQLWriterEquals(t, `DELETE FROM employees AS e USING companies AS c WHERE e.company_id = c.id AND c.deleted = true`, nil, q)
	})
}
