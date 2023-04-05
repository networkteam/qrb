package builder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/networkteam/qrb"
)

func TestBind(t *testing.T) {
	t.Run("single named arg", func(t *testing.T) {
		q := qrb.Select(qrb.N("*")).From(qrb.N("employees")).Where(qrb.N("id").Eq(qrb.Bind("id")))

		sql, args, err := qrb.
			Build(q).
			WithNamedArgs(map[string]any{"id": 42}).
			ToSQL()

		require.NoError(t, err)
		assert.Equal(t, "SELECT * FROM employees WHERE id = $1", sql)
		assert.Equal(t, []any{42}, args)
	})

	t.Run("re-use named arg", func(t *testing.T) {
		q := qrb.
			Select(qrb.N("*")).
			From(qrb.N("employees")).
			Where(qrb.Or(
				qrb.N("firstname").ILike(qrb.Bind("search")),
				qrb.N("lastname").ILike(qrb.Bind("search")),
			))

		sql, args, err := qrb.
			Build(q).
			WithNamedArgs(map[string]any{"search": "Jo%"}).
			ToSQL()

		require.NoError(t, err)
		assert.Equal(t, "SELECT * FROM employees WHERE firstname ILIKE $1 OR lastname ILIKE $1", sql)
		assert.Equal(t, []any{"Jo%"}, args)
	})

	t.Run("multiple named args", func(t *testing.T) {
		q := qrb.
			Select(qrb.N("*")).
			From(qrb.N("employees")).
			Where(qrb.And(
				qrb.Or(
					qrb.N("firstname").ILike(qrb.Bind("search")),
					qrb.N("lastname").ILike(qrb.Bind("search")),
				),
				qrb.N("active").Eq(qrb.Bind("active")),
			))

		sql, args, err := qrb.
			Build(q).
			WithNamedArgs(map[string]any{"search": "Jo%", "active": true}).
			ToSQL()

		require.NoError(t, err)
		assert.Equal(t, "SELECT * FROM employees WHERE (firstname ILIKE $1 OR lastname ILIKE $1) AND active = $2", sql)
		assert.Equal(t, []any{"Jo%", true}, args)
	})

	t.Run("missing named arg", func(t *testing.T) {
		q := qrb.Select(qrb.N("*")).From(qrb.N("employees")).Where(qrb.N("id").Eq(qrb.Bind("id")))

		_, _, err := qrb.
			Build(q).
			ToSQL()

		assert.EqualError(t, err, "missing named argument \"id\"")
	})

	t.Run("mixed Arg and Bind", func(t *testing.T) {
		q := qrb.
			Select(qrb.N("*")).
			From(qrb.N("employees")).
			Where(qrb.And(
				qrb.Or(
					qrb.N("firstname").ILike(qrb.Bind("search")),
					qrb.N("lastname").ILike(qrb.Bind("search")),
				),
				qrb.N("active").Eq(qrb.Arg(true)),
			))

		sql, args, err := qrb.
			Build(q).
			WithNamedArgs(map[string]any{"search": "Jo%"}).
			ToSQL()

		require.NoError(t, err)
		assert.Equal(t, "SELECT * FROM employees WHERE (firstname ILIKE $1 OR lastname ILIKE $1) AND active = $2", sql)
		assert.Equal(t, []any{"Jo%", true}, args)
	})
}
