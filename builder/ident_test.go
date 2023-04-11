package builder_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/networkteam/qrb"
)

func TestN(t *testing.T) {
	validIdentifiers := []string{
		"column_name1",
		`"MyTable".name`,
		"public.\"MyTable\".*",
		`"My"."Table".name`,
		`"My""Quoted""Table".*`,
		"space_trimmed ",
		"táblá_ñámé",
		"öäüß_column",
		`U&"d\0061t\+000061"`,
		`U&"\0441\043B\043E\043D"`,
		`U&"d!0061t!+000061" UESCAPE '!'`,
	}

	invalidIdentifiers := []string{
		"1column_name",
		`"MyTable.name`,
		`My"Table.name`,
	}

	for _, id := range validIdentifiers {
		t.Run(id, func(t *testing.T) {
			q := qrb.N(id)

			sql, _, err := qrb.Build(q).ToSQL()
			require.NoError(t, err)
			assert.Equal(t, strings.TrimSpace(id), sql)
			assert.Equal(t, strings.TrimSpace(id), q.Ident())
		})
	}

	for _, id := range invalidIdentifiers {
		t.Run(id, func(t *testing.T) {
			q := qrb.N(id)

			_, _, err := qrb.Build(q).ToSQL()
			require.Error(t, err)
		})
	}
}
