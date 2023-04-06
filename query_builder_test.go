package qrb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/networkteam/qrb"
)

func TestQueryBuilder_WithoutValidation(t *testing.T) {
	q := qrb.Select(qrb.Int(1)).From(qrb.N("1foo"))

	sql, _, err := qrb.Build(q).ToSQL()
	require.Error(t, err)

	assert.NotEqual(t, "SELECT 1 FROM 1foo", sql)

	sql, _, err = qrb.Build(q).WithoutValidation().ToSQL()
	require.NoError(t, err)

	assert.Equal(t, "SELECT 1 FROM 1foo", sql)
}
