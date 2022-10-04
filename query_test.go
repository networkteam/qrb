package jrm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/networkteam/jrm"
)

func TestQuery(t *testing.T) {
	q := jrm.Query(
		jrm.With(
			"author_json",
			jrm.
				Select(
					jrm.Ident("authors.author_id"),
				).
				SelectAs(
					jrm.JsonBuildObject().
						Prop("id", jrm.Ident("authors.author_id")).
						Prop("name", jrm.Ident("authors.name")),
					"json",
				).
				From(jrm.Ident("authors")),
		).
			Select(
				jrm.Ident("posts.post_id"),
				jrm.JsonBuildObject().
					Prop("title", jrm.Ident("posts.title")).
					Prop("author", jrm.Ident("author_json.json")),
			).
			From(jrm.Ident("posts")).
			LeftJoin(
				jrm.Ident("author_json"),
				jrm.Eq(jrm.Ident("posts.author_id"), jrm.Ident("author_json.author_id")),
			),
	)

	sql, args := q.ToSQL()

	assert.Equal(t, "WITH author_json AS (SELECT authors.author_id,JSON_BUILD_OBJECT('id',authors.author_id,'name',authors.name) AS json FROM authors) SELECT posts.post_id,JSON_BUILD_OBJECT('author',author_json.json,'title',posts.title) FROM posts LEFT JOIN author_json ON posts.author_id = author_json.author_id", sql)
	assert.Empty(t, args)
}
