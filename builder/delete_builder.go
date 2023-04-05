package builder

// [ WITH [ RECURSIVE ] with_query [, ...] ]
// DELETE FROM [ ONLY ] table_name [ * ] [ [ AS ] alias ]
//     [ USING from_item [, ...] ]
//     [ WHERE condition | WHERE CURRENT OF cursor_name ]
//     [ RETURNING * | output_expression [ [ AS ] output_name ] [, ...] ]

func DeleteFrom(tableName string) DeleteBuilder {
	return DeleteBuilder{
		tableName: tableName,
	}
}

type DeleteBuilder struct {
	tableName        string
	alias            string
	whereConjunction []Exp
}

func (b DeleteBuilder) As(alias string) DeleteBuilder {
	newBuilder := b
	newBuilder.alias = alias
	return newBuilder
}

// Where adds a WHERE condition to the delete.
// Multiple calls to Where are joined with AND.
func (b DeleteBuilder) Where(cond Exp) DeleteBuilder {
	newBuilder := b

	newBuilder.whereConjunction = make([]Exp, len(b.whereConjunction), len(b.whereConjunction)+1)
	copy(newBuilder.whereConjunction, b.whereConjunction)

	newBuilder.whereConjunction = append(newBuilder.whereConjunction, cond)
	return newBuilder
}

func (b DeleteBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("DELETE FROM ")
	sb.WriteString(b.tableName)
	if b.alias != "" {
		sb.WriteString(" AS ")
		sb.WriteString(b.alias)
	}
	if len(b.whereConjunction) > 0 {
		sb.WriteString(" WHERE ")
		And(b.whereConjunction...).WriteSQL(sb)
	}
}
