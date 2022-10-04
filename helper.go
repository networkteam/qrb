package jrm

import "strings"

func pqQuoteString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
