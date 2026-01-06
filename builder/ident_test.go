package builder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/networkteam/qrb"
)

func TestN(t *testing.T) {
	tests := []struct {
		input         string
		expectInvalid bool
		expected      string
	}{
		// Regular identifiers (unchanged)
		{"column_name1", false, "column_name1"},
		{"users", false, "users"},
		{"táblá_ñámé", false, "táblá_ñámé"},
		{"öäüß_column", false, "öäüß_column"},
		{"space_trimmed ", false, "space_trimmed"},

		// Dotted paths without keywords (unchanged)
		{"public.users", false, "public.users"},
		{"schema.mytable.mycolumn", false, "schema.mytable.mycolumn"},

		// Asterisks (unchanged)
		{"*", false, "*"},
		{"mytable.*", false, "mytable.*"},
		{"public.mytable.*", false, "public.mytable.*"},

		// Dotted paths with "table" and "column" keywords
		{"schema.table.column", false, `schema."table"."column"`},
		{"table.*", false, `"table".*`},
		{"public.table.*", false, `public."table".*`},

		// Already quoted identifiers (unchanged)
		{`"MyTable".name`, false, `"MyTable".name`},
		{`public."MyTable".*`, false, `public."MyTable".*`},
		{`"My"."Table".name`, false, `"My"."Table".name`},
		{`"My""Quoted""Table".*`, false, `"My""Quoted""Table".*`},

		// Unicode identifiers (unchanged)
		{`U&"d\0061t\+000061"`, false, `U&"d\0061t\+000061"`},
		{`U&"\0441\043B\043E\043D"`, false, `U&"\0441\043B\043E\043D"`},
		{`U&"d!0061t!+000061" UESCAPE '!'`, false, `U&"d!0061t!+000061" UESCAPE '!'`},

		// Keywords - should be auto-quoted
		{"from", false, `"from"`},
		{"select", false, `"select"`},
		{"where", false, `"where"`},
		{"order", false, `"order"`},
		{"group", false, `"group"`},
		{"user", false, `"user"`},
		{"table", false, `"table"`},
		{"to", false, `"to"`},
		{"all", false, `"all"`},
		{"and", false, `"and"`},
		{"or", false, `"or"`},
		{"not", false, `"not"`},
		{"null", false, `"null"`},
		{"true", false, `"true"`},
		{"false", false, `"false"`},
		{"in", false, "in"}, // "in" is not a reserved keyword

		// Keywords with different cases - should be auto-quoted
		{"FROM", false, `"FROM"`},
		{"Select", false, `"Select"`},
		{"WHERE", false, `"WHERE"`},
		{"User", false, `"User"`},

		// Keywords in dotted paths - only keyword parts should be quoted
		{"mytable.from.id", false, `mytable."from".id`},
		{"schema.select.mycolumn", false, `schema."select".mycolumn`},
		{"public.user.name", false, `public."user".name`},
		{"from.to.where", false, `"from"."to"."where"`},
		{"table.from.id", false, `"table"."from".id`},
		{"schema.select.column", false, `schema."select"."column"`},

		// Already quoted keywords (unchanged)
		{`"from"`, false, `"from"`},
		{`"select"`, false, `"select"`},
		{`mytable."from".id`, false, `mytable."from".id`},
		{`"table"."from"."id"`, false, `"table"."from"."id"`},

		// Quoted identifier with dot inside (unchanged)
		{`schema."my.table".mycolumn`, false, `schema."my.table".mycolumn`},
		// Keywords in quoted identifier context
		{`table."from".id`, false, `"table"."from".id`},
		{`schema."my.table".column`, false, `schema."my.table"."column"`},

		// Invalid identifiers
		{"1column_name", true, ""},
		{`"MyTable.name`, true, ""},
		{`My"Table.name`, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			q := qrb.N(tt.input)

			sql, _, err := qrb.Build(q).ToSQL()
			if tt.expectInvalid {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, sql)
			}
		})
	}
}
