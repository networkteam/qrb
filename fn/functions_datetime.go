package fn

import "github.com/networkteam/qrb/builder"

// EXTRACT(field FROM source)

// Extract builds the EXTRACT(field FROM source) function.
func Extract(field string, from builder.Exp) builder.ExpBase {
	return builder.ExpBase{
		Exp: extractExp{
			field: field,
			from:  from,
		},
	}
}

type extractExp struct {
	field string
	from  builder.Exp
}

func (c extractExp) IsExp() {}

func (c extractExp) WriteSQL(sb *builder.SQLBuilder) {
	sb.WriteString("EXTRACT(")
	sb.WriteString(c.field)
	sb.WriteString(" FROM ")
	c.from.WriteSQL(sb)
	sb.WriteRune(')')
}
