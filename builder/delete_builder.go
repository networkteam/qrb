package builder

// [ WITH [ RECURSIVE ] with_query [, ...] ]
// DELETE FROM [ ONLY ] table_name [ * ] [ [ AS ] alias ]
//     [ USING from_item [, ...] ]
//     [ WHERE condition | WHERE CURRENT OF cursor_name ]
//     [ RETURNING * | output_expression [ [ AS ] output_name ] [, ...] ]

func DeleteFrom(tableName IdentExp) DeleteBuilder {
	return DeleteBuilder{
		tableName: tableName,
	}
}

type DeleteBuilder struct {
	withQueries      withQueries
	tableName        IdentExp
	alias            string
	using            []fromItem
	whereConjunction []Exp
	returningItems   returningItems
}

func (b DeleteBuilder) As(alias string) DeleteBuilder {
	newBuilder := b
	newBuilder.alias = alias
	return newBuilder
}

// Using adds a USING clause to the delete.
func (b DeleteBuilder) Using(from FromExp) FromDeleteBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.using, b.using, 1)

	newBuilder.using = append(newBuilder.using, fromItem{
		from: from,
	})

	return FromDeleteBuilder{
		DeleteBuilder: newBuilder,
	}
}

type FromDeleteBuilder struct {
	DeleteBuilder
}

// As sets the alias for the last added from item.
func (b FromDeleteBuilder) As(alias string) FromDeleteBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.using, b.using, 0)

	lastIdx := len(newBuilder.using) - 1
	newBuilder.using[lastIdx].alias = alias

	return newBuilder
}

// ColumnAliases sets the column aliases for the last added from item.
func (b FromDeleteBuilder) ColumnAliases(aliases ...string) FromDeleteBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.using, b.using, 0)

	lastIdx := len(newBuilder.using) - 1
	newBuilder.using[lastIdx].columnAliases = aliases

	return newBuilder
}

// Where adds a WHERE condition to the delete.
// Multiple calls to Where are joined with AND.
func (b DeleteBuilder) Where(cond Exp) DeleteBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.whereConjunction, b.whereConjunction, 1)

	newBuilder.whereConjunction = append(newBuilder.whereConjunction, cond)
	return newBuilder
}

func (b DeleteBuilder) Returning(outputExpression Exp) ReturningDeleteBuilder {
	newBuilder := b
	newBuilder.returningItems = b.returningItems.cloneSlice(1)

	newBuilder.returningItems = append(newBuilder.returningItems, returningItem{
		outputExpression: outputExpression,
	})

	return ReturningDeleteBuilder{newBuilder}
}

type ReturningDeleteBuilder struct {
	DeleteBuilder
}

// As sets the output name for the last output expression.
func (b ReturningDeleteBuilder) As(outputName string) DeleteBuilder {
	newBuilder := b.DeleteBuilder
	newBuilder.returningItems = b.returningItems.cloneSlice(0)

	lastIdx := len(newBuilder.returningItems) - 1
	newBuilder.returningItems[lastIdx].outputName = outputName

	return newBuilder
}

func (b DeleteBuilder) WriteSQL(sb *SQLBuilder) {
	if len(b.withQueries) > 0 {
		b.withQueries.WriteSQL(sb)
	}

	sb.WriteString("DELETE FROM ")
	b.tableName.WriteSQL(sb)
	if b.alias != "" {
		sb.WriteString(" AS ")
		sb.WriteString(b.alias)
	}
	if len(b.using) > 0 {
		sb.WriteString(" USING ")
		for i, f := range b.using {
			if i > 0 {
				sb.WriteString(",")
			}
			f.WriteSQL(sb)
		}
	}
	if len(b.whereConjunction) > 0 {
		sb.WriteString(" WHERE ")
		And(b.whereConjunction...).WriteSQL(sb)
	}
	if len(b.returningItems) > 0 {
		b.returningItems.WriteSQL(sb)
	}
}
