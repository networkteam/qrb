package builder

import (
	"errors"
	"fmt"
	"regexp"
)

type expType string

func (e expType) IsExp() {}

var ErrInvalidType = errors.New("type: invalid")

func (e expType) WriteSQL(sb *SQLBuilder) {
	if sb.Validating() {
		if !isValidType(string(e)) {
			sb.AddError(fmt.Errorf("%w: %s", ErrInvalidType, string(e)))
			return
		}
	}

	sb.WriteString(string(e))
}

var validTypeRegex = regexp.MustCompile(`(?ms)\A(` +
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
	`)` + // End of identifier match
	`(?:\(\d+\))?` + // Optional parentheses with numeric value
	`(?:\s+UESCAPE\s+'[^0-9A-Fa-f"+''"[:space:]]')?` + // Optional UESCAPE clause with single character not in the excluded set
	`(\s*\[\s*\d*\s*\])*` + // Match array notation: zero or more `[]`, optionally with spaces and a number
	`\z` + // End of string
	`)`,
)

func isValidType(s string) bool {
	return validTypeRegex.MatchString(s)
}
