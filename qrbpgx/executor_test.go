package qrbpgx_test

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/qrbpgx"
)

func TestExecutiveQueryBuilder(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	selectQuery := qrb.Select(qrb.N("id"), qrb.N("name")).
		From(qrb.N("users")).
		Where(qrb.N("name").Like(qrb.Bind("matchName"))).
		Limit(qrb.Arg(5))
	eqb := qrbpgx.Build(selectQuery).WithExecutor(mock).WithNamedArgs(map[string]interface{}{"matchName": "A%"})

	columns := []string{"id", "name"}

	mock.ExpectQuery("SELECT id,name FROM users WHERE name LIKE $1 LIMIT $2").WithArgs("A%", 5).WillReturnRows(pgxmock.NewRows(columns).AddRow(1, "Alice").AddRow(2, "Bob"))

	ctx := context.Background()
	_, err = eqb.Query(ctx)
	require.NoError(t, err)

	mock.ExpectQuery("SELECT id,name FROM users WHERE name LIKE $1 LIMIT $2").WithArgs("A%", 5).WillReturnRows(pgxmock.NewRows(columns).AddRow(1, "Alice"))

	row, err := eqb.QueryRow(ctx)
	require.NoError(t, err)

	var id int
	var name string
	err = row.Scan(&id, &name)
	require.NoError(t, err)

	assert.Equal(t, 1, id)
	assert.Equal(t, "Alice", name)

	eqb = qrbpgx.Build(qrb.InsertInto("users").Values(qrb.Default(), qrb.Arg("Robert"))).WithExecutor(mock)

	mock.ExpectExec("INSERT INTO users VALUES (DEFAULT,$1)").WithArgs("Robert").WillReturnResult(pgxmock.NewResult("INSERT", 1))

	_, err = eqb.Exec(ctx)
	require.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
