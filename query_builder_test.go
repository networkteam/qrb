package qrb_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/networkteam/qrb"
)

func TestQueryBuilder_WithoutValidation(t *testing.T) {
	q := qrb.Select(qrb.Int(1)).From(qrb.N("1foo"))

	sql, _, err := qrb.Build(q).ToSQL()
	require.Error(t, err)

	sql, _, err = qrb.Build(q).WithoutValidation().ToSQL()
	require.NoError(t, err)

	require.Equal(t, "SELECT 1 FROM 1foo", sql)
}
