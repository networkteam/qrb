package builder_test

import (
	"testing"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/builder"
	"github.com/networkteam/qrb/fn"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestJsonbQuery(t *testing.T) {
	t.Run("select json object", func(t *testing.T) {
		b := qrb.SelectJsonb(
			fn.JsonbBuildObject().
				Prop("id", qrb.N("authors.author_id")).
				Prop("name", qrb.N("authors.name")),
		).
			From(qrb.N("authors")).
			Where(qrb.N("authors.author_id").Eq(qrb.Arg(123)))

		testhelper.AssertSQLWriterEquals(
			t,
			"SELECT jsonb_build_object('id',authors.author_id,'name',authors.name) FROM authors WHERE authors.author_id = $1",
			[]any{123},
			b,
		)

		// We can now modify an existing JSON selection!
		// Each SelectBuilder acts as a kind of query blueprint that can be used to modify later.
		withPostCount := b.ApplySelectJsonb(func(obj builder.JsonbBuildObjectBuilder) builder.JsonbBuildObjectBuilder {
			return obj.Prop("postCount", fn.Count(qrb.N("posts")))
		})

		testhelper.AssertSQLWriterEquals(
			t,
			"SELECT jsonb_build_object('id',authors.author_id,'name',authors.name,'postCount',count(posts)) FROM authors WHERE authors.author_id = $1",
			[]any{123},
			withPostCount,
		)
	})
}
