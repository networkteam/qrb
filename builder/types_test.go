package builder

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpType(t *testing.T) {
	validTypes := []string{
		`integer`,
		`text`,
		`boolean`,
		`varchar(255)`,
		`custom_type`,
		`"QuotedType"`,
		`integer[]`,
		`text [ ] [ ]`,
		`"QuotedArrayType"[]`,
		`integer[][]`,
		`text[16]`,
	}

	invalidIdentifiers := []string{
		"1int",
		`"MyTable.name`,
		`My"Table.name`,
	}

	for _, id := range validTypes {
		t.Run(id, func(t *testing.T) {
			q := expType(id)

			sql, _, err := Build(q).ToSQL()
			require.NoError(t, err)
			assert.Equal(t, strings.TrimSpace(id), sql)
			assert.Equal(t, strings.TrimSpace(id), string(q))
		})
	}

	for _, id := range invalidIdentifiers {
		t.Run(id, func(t *testing.T) {
			q := expType(id)

			_, _, err := Build(q).ToSQL()
			require.Error(t, err)
		})
	}
}
