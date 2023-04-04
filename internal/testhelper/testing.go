package testhelper

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/builder"
)

func AssertSQLWriterEquals(t *testing.T, expectedSQL string, expectedArgs []any, writer builder.SQLWriter) {
	t.Helper()

	sql, args, err := qrb.Build(writer).ToSQL()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	AssertSQLEquals(t, expectedSQL, sql)
	if len(expectedArgs) == 0 {
		assert.Empty(t, args)
	} else {
		assert.Equal(t, expectedArgs, args)
	}
}

func AssertSQLEquals(t *testing.T, expected string, actual string) {
	t.Helper()

	// Compare actual and expected and ignore whitespace differences (e.g. a space could also be multiple spaces / tab or newlines).
	// But make sure minimal whitespace is preserved (e.g. a space between two words should not be removed).
	// This is to make the tests more readable.
	// Replace all newlines with spaces to make the replacement regexp easier.
	expected = strings.ReplaceAll(expected, "\n", " ")
	// Normalize whitespace (remove multiple spaces, remove space after opening brackets or comma)
	pattern := regexp.MustCompile(`(\s|\(|,)\s+`)
	expected = pattern.ReplaceAllString(expected, "$1")
	// Remove space before closing brackets, ideally this would be handled by the pattern above
	expected = strings.ReplaceAll(expected, " )", ")")

	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(actual))
}
