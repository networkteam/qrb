package qrb_test

import (
	"testing"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/builder"
	"github.com/networkteam/qrb/fn"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestDeleteBuilder(t *testing.T) {
	t.Run("examples", func(t *testing.T) {
		// From https://www.postgresql.org/docs/15/sql-delete.html#id-1.9.3.100.8

		t.Run("example 0.1", func(t *testing.T) {
			q := qrb.
				DeleteFrom(qrb.N("films")).
				Using(qrb.N("producers")).
				Where(qrb.And(
					qrb.N("producer_id").Eq(qrb.N("producers.id")),
					qrb.N("producers.name").Eq(qrb.String("foo")),
				))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				DELETE FROM films USING producers
				  WHERE producer_id = producers.id AND producers.name = 'foo'
				`,
				nil,
				q,
			)
		})

		t.Run("example 0.2", func(t *testing.T) {
			q := qrb.
				DeleteFrom(qrb.N("films")).
				Where(qrb.N("producer_id").In(
					qrb.Select(qrb.N("id")).From(qrb.N("producers")).Where(qrb.N("name").Eq(qrb.String("foo"))),
				))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				DELETE FROM films
				  WHERE producer_id IN (SELECT id FROM producers WHERE name = 'foo')
				`,
				nil,
				q,
			)
		})

		// From https://www.postgresql.org/docs/15/sql-delete.html#id-1.9.3.100.9

		t.Run("example 1.1", func(t *testing.T) {
			q := qrb.
				DeleteFrom(qrb.N("films")).
				Where(qrb.N("kind").Neq(qrb.String("Musical")))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				DELETE FROM films WHERE kind <> 'Musical'
				`,
				nil,
				q,
			)
		})

		t.Run("example 1.2", func(t *testing.T) {
			q := qrb.
				DeleteFrom(qrb.N("films"))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				DELETE FROM films
				`,
				nil,
				q,
			)
		})

		t.Run("example 1.3", func(t *testing.T) {
			q := qrb.
				DeleteFrom(qrb.N("tasks")).
				Where(qrb.N("status").Eq(qrb.String("DONE"))).
				Returning(qrb.N("*"))

			testhelper.AssertSQLWriterEquals(
				t,
				`
				DELETE FROM tasks WHERE status = 'DONE' RETURNING *
				`,
				nil,
				q,
			)
		})
	})

	t.Run("with", func(t *testing.T) {
		// Example borrowed from https://stackoverflow.com/a/37225172/749191

		var listens builder.Identer = qrb.N("listens")

		q := qrb.
			With("max_table").As(
			qrb.Select(qrb.N("uid")).Select(fn.Max(qrb.N("ts")).Minus(qrb.Int(10000))).As("mx").From(qrb.N("listens")).GroupBy(qrb.N("uid")),
		).
			DeleteFrom(listens).
			Where(qrb.N("ts").Lt(qrb.Select(qrb.N("mx")).From(qrb.N("max_table")).Where(qrb.N("max_table.uid").Eq(qrb.N("listens.uid")))))

		testhelper.AssertSQLWriterEquals(
			t,
			`
				WITH max_table AS (
					SELECT uid, max(ts) - 10000 AS mx
					FROM listens 
					GROUP BY uid
				) 
				DELETE FROM listens 
				WHERE ts < (SELECT mx
								   FROM max_table 
								   WHERE max_table.uid = listens.uid)
				`,
			nil,
			q,
		)
	})

	t.Run("", func(t *testing.T) {
		q := qrb.DeleteFrom(qrb.N("list_line_items")).
			Where(qrb.And(
				qrb.N("shopping_cart_id").Eq(qrb.Arg(99)),
				qrb.Exps(qrb.N("supplier_id"), qrb.N("article_id"), qrb.N("unit_id")).In(
					qrb.Select(
						fn.Unnest(qrb.Arg([]int{1, 2, 3}).Cast("integer[]")),
						fn.Unnest(qrb.Arg([]int{4, 5, 6}).Cast("integer[]")),
						fn.Unnest(qrb.Arg([]int{7, 8, 9}).Cast("integer[]")),
					),
				),
			))

		testhelper.AssertSQLWriterEquals(
			t,
			`
				DELETE FROM list_line_items
				WHERE shopping_cart_id = $1
				AND (supplier_id, article_id, unit_id) IN (
					SELECT unnest($2::integer[]), unnest($3::integer[]), unnest($4::integer[])
				)
				`,
			[]any{
				99,
				[]int{1, 2, 3},
				[]int{4, 5, 6},
				[]int{7, 8, 9},
			},
			q,
		)
	})
}
