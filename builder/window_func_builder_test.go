package builder_test

import (
	"testing"

	. "github.com/networkteam/qrb"
	"github.com/networkteam/qrb/fn"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestWindowFuncBuilder(t *testing.T) {
	t.Run("Example 1", func(t *testing.T) {
		b := Select(
			N("depname"),
			N("empno"),
			N("salary"),
			fn.Avg(N("salary")).Over().PartitionBy(N("depname")),
		).From(N("empsalary"))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT depname, empno, salary, avg(salary) OVER (PARTITION BY depname) FROM empsalary
			`,
			nil,
			b,
		)
	})

	t.Run("Example 2", func(t *testing.T) {
		b := Select(
			N("depname"),
			N("empno"),
			N("salary"),
			fn.Rank().Over().PartitionBy(N("depname")).OrderBy(N("salary")).Desc(),
		).From(N("empsalary"))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT depname, empno, salary,
				   rank() OVER (PARTITION BY depname ORDER BY salary DESC)
			FROM empsalary
			`,
			nil,
			b,
		)
	})

	t.Run("Example 3", func(t *testing.T) {
		b := Select(
			N("salary"),
			fn.Sum(N("salary")).Over(),
		).From(N("empsalary"))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT salary, sum(salary) OVER () FROM empsalary
			`,
			nil,
			b,
		)
	})

	t.Run("Example 4", func(t *testing.T) {
		b := Select(
			N("salary"),
			fn.Sum(N("salary")).Over().OrderBy(N("salary")),
		).From(N("empsalary"))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT salary, sum(salary) OVER (ORDER BY salary) FROM empsalary
			`,
			nil,
			b,
		)
	})

	t.Run("Example 5", func(t *testing.T) {
		b := Select(
			N("depname"),
			N("empno"),
			N("salary"),
			N("enroll_date"),
		).
			From(
				Select(
					N("depname"),
					N("empno"),
					N("salary"),
					N("enroll_date"),
				).
					Select(
						fn.Rank().Over().PartitionBy(N("depname")).OrderBy(N("salary")).Desc().OrderBy(N("empno")),
					).As("pos").
					From(N("empsalary")),
			).As("salaries").
			Where(N("pos").Lt(Int(3)))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT depname, empno, salary, enroll_date
			FROM
			  (SELECT depname, empno, salary, enroll_date,
					  rank() OVER (PARTITION BY depname ORDER BY salary DESC, empno) AS pos
				 FROM empsalary
			  ) AS salaries
			WHERE pos < 3
			`,
			nil,
			b,
		)
	})

	t.Run("Example 6", func(t *testing.T) {
		b := Select(
			fn.Sum(N("salary")).Over("w"),
			fn.Avg(N("salary")).Over("w"),
		).
			From(N("empsalary")).
			Window("w").As().PartitionBy(N("depname")).OrderBy(N("salary")).Desc().
			SelectBuilder

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT sum(salary) OVER w, avg(salary) OVER w
			  FROM empsalary
			  WINDOW w AS (PARTITION BY depname ORDER BY salary DESC)
			`,
			nil,
			b,
		)
	})

	t.Run("Example 6 variant 1", func(t *testing.T) {
		b := Select(
			fn.Sum(N("salary")).Over("w"),
			fn.Avg(N("salary")).Over("w"),
			fn.RowNumber().Over("w"),
		).
			From(N("empsalary")).
			Window("w").As().PartitionBy(N("depname")).OrderBy(N("salary")).Desc().
			SelectBuilder

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT sum(salary) OVER w, avg(salary) OVER w, row_number() OVER w
			  FROM empsalary
			  WINDOW w AS (PARTITION BY depname ORDER BY salary DESC)
			`,
			nil,
			b,
		)
	})
}
