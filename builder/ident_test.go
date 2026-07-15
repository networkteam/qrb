package builder_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/builder"
	"github.com/networkteam/qrb/internal/testhelper"
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
		// The U& prefix may use an upper or lower case U
		{`u&"d\0061t"`, false, `u&"d\0061t"`},
		// The UESCAPE keyword is case-insensitive
		{`U&"d!0061t" uescape '!'`, false, `U&"d!0061t" uescape '!'`},
		// A doubled backslash is a literal escape character
		{`U&"a\\b"`, false, `U&"a\\b"`},

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
		// The plus sign form of a Unicode escape sequence requires exactly six hexadecimal digits
		{`U&"a\+0041"`, true, ""},
		// A single backslash is not a valid Unicode escape sequence
		{`"a\b"`, true, ""},
		// The escape character must not be a hexadecimal digit, plus sign, quote or whitespace
		{`U&"x!0041" UESCAPE 'A'`, true, ""},
		{`U&"x!0041" UESCAPE '+'`, true, ""},
		{`U&"x!0041" UESCAPE '"'`, true, ""},
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

func TestIdentExp_Unqualified(t *testing.T) {
	tests := []struct {
		input         string
		expectInvalid bool
		expected      string
	}{
		// Identifiers without qualification (unchanged)
		{"mycolumn", false, "mycolumn"},
		{"_col$1", false, "_col$1"},
		{"táblá_ñámé", false, "táblá_ñámé"},
		{"*", false, "*"},
		{`"MyColumn"`, false, `"MyColumn"`},
		{`"My""Quoted""Column"`, false, `"My""Quoted""Column"`},

		// A dot inside a quoted identifier is no qualification (unchanged)
		{`"my.column"`, false, `"my.column"`},
		{`"col""with.dot"`, false, `"col""with.dot"`},
		{`U&"a.b"`, false, `U&"a.b"`},

		// Reserved keywords are quoted when writing
		{"from", false, `"from"`},
		{`"from"`, false, `"from"`},

		// Unicode identifiers without qualification (unchanged, UESCAPE clause is kept)
		{`U&"d\0061t\+000061"`, false, `U&"d\0061t\+000061"`},
		{`U&"\0441\043B\043E\043D"`, false, `U&"\0441\043B\043E\043D"`},
		{`u&"d\0061t"`, false, `u&"d\0061t"`},
		{`U&"a\\b"`, false, `U&"a\\b"`},
		{`U&"d!0061t!+000061" UESCAPE '!'`, false, `U&"d!0061t!+000061" UESCAPE '!'`},
		// The escape character may be a dot, which must not be treated as qualification
		{`U&"d.0061t" UESCAPE '.'`, false, `U&"d.0061t" UESCAPE '.'`},

		// Qualified identifiers are reduced to the last segment
		{"pipeline_instances.status", false, "status"},
		{"mytable.mycolumn", false, "mycolumn"},
		{"myschema.mytable.mycolumn", false, "mycolumn"},
		{"mycatalog.myschema.mytable.mycolumn", false, "mycolumn"},
		{"mytable.*", false, "*"},
		{"public.mytable.*", false, "*"},
		{"táblá.ñámé", false, "ñámé"},

		// Reserved keywords as last segment are quoted when writing
		{"mytable.from", false, `"from"`},
		{"from.to.where", false, `"where"`},
		{"table.column", false, `"column"`},
		{"FROM.SELECT", false, `"SELECT"`},

		// Quoted segments
		{`"MyTable".name`, false, "name"},
		{`"My"."Table".name`, false, "name"},
		{`mytable."MyColumn"`, false, `"MyColumn"`},
		{`"My""Quoted""Table".mycolumn`, false, "mycolumn"},
		{`mytable."My""Quoted""Column"`, false, `"My""Quoted""Column"`},
		{`"ends""".mycolumn`, false, "mycolumn"},
		{`mytable."a"""`, false, `"a"""`},

		// Dots inside quoted segments are no separators
		{`myschema."my.table".mycolumn`, false, "mycolumn"},
		{`"a.b.c".mycolumn`, false, "mycolumn"},
		{`"my.schema"."my.table"."my.column"`, false, `"my.column"`},

		// Unicode identifiers with qualification
		{`U&"tbl".mycolumn`, false, "mycolumn"},
		{`U&"t\0062l".mycolumn`, false, "mycolumn"},
		{`u&"tbl".mycolumn`, false, "mycolumn"},

		// A trailing UESCAPE clause belongs to the U& prefix of the first segment,
		// so it is dropped together with the qualification
		{`U&"t!0062l".mycolumn UESCAPE '!'`, false, "mycolumn"},
		{`U&"t!0062l".mycolumn uescape '!'`, false, "mycolumn"},
		// The escape character may be a dot, which must not be treated as qualification
		{`U&"t.0062l".mycolumn UESCAPE '.'`, false, "mycolumn"},

		// Invalid identifiers stay invalid (returned unchanged, so building reports an error)
		{"", true, ""},
		{"1column.mycolumn", true, ""},
		{"my table.mycolumn", true, ""},
		{`"MyTable.mycolumn`, true, ""},
		{`My"Table.mycolumn`, true, ""},
		{"mytable.", true, ""},
		{".mycolumn", true, ""},
		{"mytable..mycolumn", true, ""},
		// Asterisks are only allowed as last segment
		{"*.mycolumn", true, ""},
		{"mytable.*.mycolumn", true, ""},
		// Invalid escape characters in the UESCAPE clause
		{`U&"x!0041".mycolumn UESCAPE 'A'`, true, ""},
		{`U&"x!0041".mycolumn UESCAPE '+'`, true, ""},
		{`U&"x!0041".mycolumn UESCAPE '"'`, true, ""},
		// The plus sign form of a Unicode escape sequence requires exactly six hexadecimal digits
		{`U&"a\+0041".mycolumn`, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			q := qrb.N(tt.input).Unqualified()

			sql, _, err := qrb.Build(q).ToSQL()
			if tt.expectInvalid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, sql)

			// Unqualified is idempotent
			sql, _, err = qrb.Build(q.Unqualified()).ToSQL()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, sql)
		})
	}
}

func TestIdentExp_String(t *testing.T) {
	var _ fmt.Stringer = builder.IdentExp{}

	assert.Equal(t, "pipeline_instances.status", qrb.N("pipeline_instances.status").String())
	assert.Equal(t, "mycolumn", qrb.N("mycolumn").String())
	// The raw identifier is returned without keyword quoting
	assert.Equal(t, "mytable.from", qrb.N("mytable.from").String())
	assert.Equal(t, `myschema."my.table"`, qrb.N(`myschema."my.table"`).String())
}

func TestIdentExp_UnqualifiedString(t *testing.T) {
	assert.Equal(t, "status", qrb.N("pipeline_instances.status").UnqualifiedString())
	assert.Equal(t, "status", qrb.N("status").UnqualifiedString())
	// Quoting is kept as given
	assert.Equal(t, `"my.column"`, qrb.N(`myschema."my.table"."my.column"`).UnqualifiedString())
	// The raw identifier is returned without keyword quoting (it is quoted when written, e.g. by Set)
	assert.Equal(t, "from", qrb.N("mytable.from").UnqualifiedString())
	assert.Equal(t, "mycolumn", qrb.N(`U&"t!0062l".mycolumn UESCAPE '!'`).UnqualifiedString())
	// An invalid identifier is returned unchanged
	assert.Equal(t, "my table.mycolumn", qrb.N("my table.mycolumn").UnqualifiedString())
}

func TestIdentExp_Unqualified_Usage(t *testing.T) {
	// Table-qualified identifiers as generated by tools like construct
	distributors := struct {
		builder.Identer
		DID   builder.IdentExp
		DName builder.IdentExp
	}{
		Identer: qrb.N("distributors"),
		DID:     qrb.N("distributors.did"),
		DName:   qrb.N("distributors.dname"),
	}

	t.Run("on conflict target", func(t *testing.T) {
		q := qrb.InsertInto(distributors).
			ColumnNames("did", "dname").
			Values(qrb.Int(7), qrb.String("Redline GmbH")).
			OnConflict(distributors.DID.Unqualified()).DoNothing()

		testhelper.AssertSQLWriterEquals(
			t,
			`
			INSERT INTO distributors (did, dname) VALUES (7, 'Redline GmbH')
			ON CONFLICT (did) DO NOTHING
			`,
			nil,
			q,
		)
	})

	t.Run("update set", func(t *testing.T) {
		q := qrb.Update(distributors).
			Set(distributors.DName.UnqualifiedString(), qrb.Arg("ACME")).
			Where(distributors.DID.Eq(qrb.Arg(7)))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			UPDATE distributors SET dname = $1 WHERE distributors.did = $2
			`,
			[]any{"ACME", 7},
			q,
		)
	})
}
