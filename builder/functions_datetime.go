package builder

// EXTRACT(field FROM source)

// Extract builds the EXTRACT(field FROM source) function.
func Extract(field string, from Exp) ExpBase {
	return ExpBase{
		Exp: extractExp{
			field: field,
			from:  from,
		},
	}
}

type extractExp struct {
	field string
	from  Exp
}

func (c extractExp) IsExp() {}

func (c extractExp) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("EXTRACT(")
	sb.WriteString(c.field)
	sb.WriteString(" FROM ")
	c.from.WriteSQL(sb)
	sb.WriteRune(')')
}
