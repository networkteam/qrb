package builder

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// N writes the given name / identifier.
//
// It will validate the identifier when writing the query,
// but it will not detect all invalid identifiers that are invalid in PostgreSQL (especially considering reserved keywords).
func N(s string) IdentExp {
	exp := IdentExp{ident: strings.TrimSpace(s)}
	exp.Exp = exp // self-reference for base methods
	return exp
}

type IdentExp struct {
	ExpBase
	ident string
}

func (i IdentExp) IsExp()       {}
func (i IdentExp) isFromExp()   {}
func (i IdentExp) NoParensExp() {}

func (i IdentExp) Ident() string {
	return i.ident
}

type Identer interface {
	Exp
	Ident() string
	isFromExp()
}

var ErrInvalidIdentifier = errors.New("identifier: invalid")

func (i IdentExp) WriteSQL(sb *SQLBuilder) {
	if sb.Validating() {
		if !isValidIdentifier(i.ident) {
			sb.AddError(fmt.Errorf("%w: %s", ErrInvalidIdentifier, i.ident))
			return
		}
	}

	sb.WriteString(i.ident)
}

var validIdentifierRegex = regexp.MustCompile(`(?ms)\A(` +
	`(?:U&)?` + // Optional U& prefix for Unicode escape sequences
	`(?:` +
	`(?:[_\p{L}][_\p{L}\p{Nd}$]{0,62}` + // Unquoted identifier
	`|"` + // Quoted identifier
	`(?:` +
	`[^"\\]|""` + // Any character except double quotes or backslashes; two double quotes are allowed
	`|\\(?:\+?[0-9A-Fa-f]{4}|\+?[0-9A-Fa-f]{6})` + // Unicode escape sequence: \+? followed by four or six hexadecimal digits
	`)+"` +
	`)\.)*` + // Allow for dotted paths
	`(?:` +
	`[_\p{L}][_\p{L}\p{Nd}$]{0,62}` + // Unquoted identifier
	`|"(([^"\\]|"")` + // Quoted identifier (same as above)
	`|\\(?:\+?[0-9A-Fa-f]{4}|\+?[0-9A-Fa-f]{6})` +
	`)+"` +
	`|\*` + // Allow for asterisks
	`)(?:\s+UESCAPE\s+'[^0-9A-Fa-f"+''"[:space:]]')?` + // Optional UESCAPE clause with single character not in the excluded set
	`\z` + // End of string
	`)`,
)

func isValidIdentifier(s string) bool {
	return validIdentifierRegex.MatchString(s)
}
