package builder

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// N writes the given name / identifier.
//
// It will validate the identifier when writing the query.
// Reserved PostgreSQL keywords are automatically quoted (e.g. "from" becomes `"from"`).
func N(s string) IdentExp {
	ident := strings.TrimSpace(s)
	exp := IdentExp{
		ident:       ident,
		quotedIdent: quoteIdentifierIfKeyword(ident),
	}
	exp.Exp = exp // self-reference for base methods
	return exp
}

type IdentExp struct {
	ExpBase
	ident       string
	quotedIdent string // pre-computed quoted identifier for performance
}

func (i IdentExp) IsExp()     {}
func (i IdentExp) isFromExp() {}

func (i IdentExp) Ident() string {
	return i.ident
}

// String returns the identifier as it was given to N, without any keyword quoting applied.
// It implements fmt.Stringer.
func (i IdentExp) String() string {
	return i.ident
}

// UnqualifiedString returns the identifier without its qualification as a raw string (see Unqualified),
// e.g. to use a table-qualified column as column name in Set of an update or insert builder.
func (i IdentExp) UnqualifiedString() string {
	return i.Unqualified().ident
}

// Unqualified returns the identifier without its qualification, i.e. only the last segment of a dotted path
// (e.g. the column of a table-qualified column):
//
//	N("pipeline_instances.status").Unqualified() // writes: status
//
// This is useful where PostgreSQL requires a bare column name, e.g. as conflict target in ON CONFLICT.
//
// Dots inside quoted segments are not separators, so N(`schema."my.table"`).Unqualified() writes "my.table".
// A trailing UESCAPE clause is dropped together with the qualification, since it can only belong to the
// U& prefix of the first segment. An unqualified or invalid identifier is returned unchanged
// (an invalid identifier still reports an error when the query is built).
func (i IdentExp) Unqualified() IdentExp {
	if !isValidIdentifier(i.ident) {
		return i
	}

	// Cut off a trailing UESCAPE clause before splitting, its escape character may be any
	// non-quote character, including a dot.
	base := trailingUescapeRegex.ReplaceAllString(i.ident, "")

	parts := splitIdentifier(base)
	if len(parts) < 2 {
		return i
	}
	return N(parts[len(parts)-1])
}

type Identer interface {
	Exp
	Ident() string
	isFromExp()
}

var ErrInvalidIdentifier = errors.New("identifier: invalid")

// trailingUescapeRegex matches a trailing UESCAPE clause as accepted by validIdentifierRegex.
var trailingUescapeRegex = regexp.MustCompile(`\s+(?i:UESCAPE)\s+'[^0-9A-Fa-f"+'[:space:]]'\z`)

func (i IdentExp) WriteSQL(sb *SQLBuilder) {
	if sb.Validating() {
		if !isValidIdentifier(i.ident) {
			sb.AddError(fmt.Errorf("%w: %s", ErrInvalidIdentifier, i.ident))
			return
		}
	}

	sb.WriteString(i.quotedIdent)
}

var validIdentifierRegex = regexp.MustCompile(`(?ms)\A(` +
	`(?:[Uu]&)?` + // Optional U& prefix for Unicode escape sequences (upper or lower case U)
	`(?:` +
	`(?:[_\p{L}][_\p{L}\p{Nd}$]{0,62}` + // Unquoted identifier
	`|"` + // Quoted identifier
	`(?:` +
	`[^"\\]|""` + // Any character except double quotes or backslashes; two double quotes are allowed
	`|\\(?:[0-9A-Fa-f]{4}|\+[0-9A-Fa-f]{6}|\\)` + // Unicode escape sequence: four hexadecimal digits, a plus sign followed by six hexadecimal digits, or a doubled backslash (literal escape character)
	`)+"` +
	`)\.)*` + // Allow for dotted paths
	`(?:` +
	`[_\p{L}][_\p{L}\p{Nd}$]{0,62}` + // Unquoted identifier
	`|"(([^"\\]|"")` + // Quoted identifier (same as above)
	`|\\(?:[0-9A-Fa-f]{4}|\+[0-9A-Fa-f]{6}|\\)` +
	`)+"` +
	`|\*` + // Allow for asterisks
	`)(?:\s+(?i:UESCAPE)\s+'[^0-9A-Fa-f"+''"[:space:]]')?` + // Optional UESCAPE clause (keyword is case-insensitive) with single character not in the excluded set
	`\z` + // End of string
	`)`,
)

func isValidIdentifier(s string) bool {
	return validIdentifierRegex.MatchString(s)
}
