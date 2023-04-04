package builder_test

import (
	"testing"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestArray(t *testing.T) {
	t.Run("with ints", func(t *testing.T) {
		b := qrb.Array(
			qrb.Int(1),
			qrb.Int(2),
			qrb.Int(3),
		)

		testhelper.AssertSQLWriterEquals(t, "ARRAY[1,2,3]", nil, b)
	})

	t.Run("with ident", func(t *testing.T) {
		b := qrb.Array(
			qrb.Int(1),
			qrb.Int(2),
			qrb.N("bar"),
		)

		testhelper.AssertSQLWriterEquals(t, "ARRAY[1,2,bar]", nil, b)
	})

	t.Run("with placeholder", func(t *testing.T) {
		b := qrb.Array(
			qrb.Int(1),
			qrb.Int(2),
			qrb.Arg(3),
		)

		testhelper.AssertSQLWriterEquals(t, "ARRAY[1,2,$1]", []any{3}, b)
	})
}

func TestString(t *testing.T) {
	tt := []struct {
		input    string
		expected string
	}{
		{
			input:    "foo",
			expected: "'foo'",
		},
		{
			input:    "foo'bar",
			expected: "'foo''bar'",
		},
		{
			input:    `with some \n escapes`,
			expected: `E'with some \\n escapes'`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.input, func(t *testing.T) {
			b := qrb.String(tc.input)

			testhelper.AssertSQLWriterEquals(t, tc.expected, nil, b)
		})
	}
}
