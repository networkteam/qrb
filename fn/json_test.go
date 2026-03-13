package fn_test

import (
	"testing"

	. "github.com/networkteam/qrb"
	"github.com/networkteam/qrb/fn"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestJsonSetReturningFunctions(t *testing.T) {
	t.Run("jsonb_array_elements in FROM", func(t *testing.T) {
		q := Select(N("value")).
			From(fn.JsonbArrayElements(N("my_column"))).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT value FROM jsonb_array_elements(my_column)
		`, sql)
	})

	t.Run("json_array_elements in FROM", func(t *testing.T) {
		q := Select(N("value")).
			From(fn.JsonArrayElements(N("data"))).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT value FROM json_array_elements(data)
		`, sql)
	})

	t.Run("jsonb_array_elements_text in FROM", func(t *testing.T) {
		q := Select(N("value")).
			From(fn.JsonbArrayElementsText(N("my_column"))).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT value FROM jsonb_array_elements_text(my_column)
		`, sql)
	})

	t.Run("json_array_elements_text in FROM", func(t *testing.T) {
		q := Select(N("value")).
			From(fn.JsonArrayElementsText(N("data"))).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT value FROM json_array_elements_text(data)
		`, sql)
	})

	t.Run("jsonb_array_elements with alias", func(t *testing.T) {
		q := Select(N("elem")).
			From(fn.JsonbArrayElements(N("my_column"))).As("elem").
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT elem FROM jsonb_array_elements(my_column) AS elem
		`, sql)
	})

	t.Run("jsonb_each in FROM", func(t *testing.T) {
		q := Select(N("key"), N("value")).
			From(fn.JsonbEach(N("data"))).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT key, value FROM jsonb_each(data)
		`, sql)
	})

	t.Run("json_each in FROM", func(t *testing.T) {
		q := Select(N("key"), N("value")).
			From(fn.JsonEach(N("data"))).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT key, value FROM json_each(data)
		`, sql)
	})

	t.Run("jsonb_each_text in FROM", func(t *testing.T) {
		q := Select(N("key"), N("value")).
			From(fn.JsonbEachText(N("data"))).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT key, value FROM jsonb_each_text(data)
		`, sql)
	})

	t.Run("json_each_text in FROM", func(t *testing.T) {
		q := Select(N("key"), N("value")).
			From(fn.JsonEachText(N("data"))).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT key, value FROM json_each_text(data)
		`, sql)
	})

	t.Run("jsonb_object_keys in FROM", func(t *testing.T) {
		q := Select(N("jsonb_object_keys")).
			From(fn.JsonbObjectKeys(N("data"))).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT jsonb_object_keys FROM jsonb_object_keys(data)
		`, sql)
	})

	t.Run("json_object_keys in FROM", func(t *testing.T) {
		q := Select(N("json_object_keys")).
			From(fn.JsonObjectKeys(N("data"))).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT json_object_keys FROM json_object_keys(data)
		`, sql)
	})

	t.Run("jsonb_populate_recordset in FROM", func(t *testing.T) {
		q := Select(N("*")).
			From(fn.JsonbPopulateRecordset(Null(), N("data"))).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT * FROM jsonb_populate_recordset(NULL, data)
		`, sql)
	})

	t.Run("json_populate_recordset in FROM", func(t *testing.T) {
		q := Select(N("*")).
			From(fn.JsonPopulateRecordset(Null(), N("data"))).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT * FROM json_populate_recordset(NULL, data)
		`, sql)
	})

	t.Run("jsonb_path_query in FROM", func(t *testing.T) {
		q := Select(N("*")).
			From(fn.JsonbPathQuery(N("data"), String("$.items[*]"))).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT * FROM jsonb_path_query(data, '$.items[*]')
		`, sql)
	})

	t.Run("jsonb_path_query_tz in FROM", func(t *testing.T) {
		q := Select(N("*")).
			From(fn.JsonbPathQueryTZ(N("data"), String("$.items[*]"))).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT * FROM jsonb_path_query_tz(data, '$.items[*]')
		`, sql)
	})

	t.Run("jsonb_array_elements as expression", func(t *testing.T) {
		// Verify backward compatibility: can still be used as an expression
		q := Select(fn.JsonbArrayElements(N("data"))).
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT jsonb_array_elements(data)
		`, sql)
	})

	t.Run("jsonb_array_elements with cross join", func(t *testing.T) {
		q := Select(N("t.id"), N("elem")).
			From(N("my_table")).As("t").
			From(fn.JsonbArrayElements(N("t.data"))).As("elem").
			SelectBuilder

		sql, _, _ := Build(q).ToSQL()
		testhelper.AssertSQLEquals(t, `
		SELECT t.id, elem FROM my_table AS t, jsonb_array_elements(t.data) AS elem
		`, sql)
	})
}
