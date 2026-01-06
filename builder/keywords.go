package builder

import "strings"

// pgReservedKeywords contains PostgreSQL reserved keywords that must be quoted
// when used as identifiers. Keywords are stored in uppercase for case-insensitive lookup.
// Source: https://www.postgresql.org/docs/current/sql-keywords-appendix.html
var pgReservedKeywords = map[string]struct{}{
	"ALL":          {},
	"ANALYSE":      {},
	"ANALYZE":      {},
	"AND":          {},
	"ANY":          {},
	"ARRAY":        {},
	"AS":           {},
	"ASYMMETRIC":   {},
	"BINARY":       {},
	"BOTH":         {},
	"CASE":         {},
	"CAST":         {},
	"CHECK":        {},
	"COLLATE":      {},
	"COLUMN":       {},
	"CONCURRENTLY": {},
	"CONSTRAINT":   {},
	"CREATE":       {},
	"CROSS":        {},
	"DEFAULT":      {},
	"DEFERRABLE":   {},
	"DESC":         {},
	"DISTINCT":     {},
	"DO":           {},
	"ELSE":         {},
	"END":          {},
	"EXCEPT":       {},
	"FALSE":        {},
	"FETCH":        {},
	"FOR":          {},
	"FOREIGN":      {},
	"FREEZE":       {},
	"FROM":         {},
	"FULL":         {},
	"GRANT":        {},
	"GROUP":        {},
	"HAVING":       {},
	"INNER":        {},
	"INTERSECT":    {},
	"INTO":         {},
	"IS":           {},
	"ISNULL":       {},
	"JOIN":         {},
	"LATERAL":      {},
	"LEADING":      {},
	"LEFT":         {},
	"LIKE":         {},
	"LIMIT":        {},
	"NOTNULL":      {},
	"NOT":          {},
	"NULL":         {},
	"OFFSET":       {},
	"ON":           {},
	"ONLY":         {},
	"OR":           {},
	"ORDER":        {},
	"OVERLAPS":     {},
	"PLACING":      {},
	"PRIMARY":      {},
	"REFERENCES":   {},
	"RETURNING":    {},
	"RIGHT":        {},
	"SELECT":       {},
	"SIMILAR":      {},
	"SOME":         {},
	"SYMMETRIC":    {},
	"TABLE":        {},
	"TABLESAMPLE":  {},
	"THEN":         {},
	"TO":           {},
	"TRAILING":     {},
	"TRUE":         {},
	"UNION":        {},
	"UNIQUE":       {},
	"USER":         {},
	"USING":        {},
	"VARIADIC":     {},
	"VERBOSE":      {},
	"WHEN":         {},
	"WHERE":        {},
	"WINDOW":       {},
	"WITH":         {},
}

// isReservedKeyword checks if the given identifier (case-insensitive) is a PostgreSQL reserved keyword.
func isReservedKeyword(s string) bool {
	_, ok := pgReservedKeywords[strings.ToUpper(s)]
	return ok
}

// quoteIdentifier wraps an identifier in double quotes, escaping any internal double quotes.
func quoteIdentifier(s string) string {
	escaped := strings.ReplaceAll(s, `"`, `""`)
	return `"` + escaped + `"`
}

// isAlreadyQuoted checks if an identifier part is already quoted (starts and ends with double quote).
func isAlreadyQuoted(s string) bool {
	return len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"'
}

// quoteIdentifierIfKeyword processes an identifier string and quotes parts that are reserved keywords.
// It handles dotted paths by processing each part individually.
// Already quoted parts are left unchanged.
func quoteIdentifierIfKeyword(ident string) string {
	if ident == "" || ident == "*" {
		return ident
	}

	// Check for Unicode identifier prefix (U&) - don't modify these
	if strings.HasPrefix(ident, "U&") || strings.HasPrefix(ident, "u&") {
		return ident
	}

	// Split by dots to handle dotted paths like schema.table.column
	parts := splitIdentifier(ident)

	for i, part := range parts {
		// Skip if already quoted
		if isAlreadyQuoted(part) {
			continue
		}

		// Skip asterisk
		if part == "*" {
			continue
		}

		// Quote if it's a reserved keyword
		if isReservedKeyword(part) {
			parts[i] = quoteIdentifier(part)
		}
	}

	return strings.Join(parts, ".")
}

// splitIdentifier splits an identifier by dots, but respects quoted parts.
// e.g., `schema."my.table".column` -> ["schema", `"my.table"`, "column"]
func splitIdentifier(ident string) []string {
	var parts []string
	var current strings.Builder
	inQuote := false

	for i := 0; i < len(ident); i++ {
		ch := ident[i]

		if ch == '"' {
			inQuote = !inQuote
			current.WriteByte(ch)
		} else if ch == '.' && !inQuote {
			parts = append(parts, current.String())
			current.Reset()
		} else {
			current.WriteByte(ch)
		}
	}

	// Add the last part
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}
