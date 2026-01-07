package builder

import "sort"

// [ WITH [ RECURSIVE ] with_query [, ...] ]
// UPDATE [ ONLY ] table_name [ * ] [ [ AS ] alias ]
//     SET { column_name = { expression | DEFAULT } |
//           ( column_name [, ...] ) = [ ROW ] ( { expression | DEFAULT } [, ...] ) |
//           ( column_name [, ...] ) = ( sub-SELECT )
//         } [, ...]
//     [ FROM from_item [, ...] ]
//     [ WHERE condition | WHERE CURRENT OF cursor_name ]
//     [ RETURNING * | output_expression [ [ AS ] output_name ] [, ...] ]

func Update(tableName Identer) UpdateBuilder {
	return UpdateBuilder{
		tableName: tableName,
	}
}

type UpdateBuilder struct {
	withQueries      withQueries
	tableName        Identer
	alias            string
	setItems         []updateSetItem
	from             []fromItem
	whereConjunction []Exp
	returningItems   returningItems
}

func (b UpdateBuilder) isWithQuery() {}

type updateSetItem struct {
	columnName string
	value      Exp
}

func (b UpdateBuilder) As(alias string) UpdateBuilder {
	newBuilder := b
	newBuilder.alias = alias
	return newBuilder
}

func (b UpdateBuilder) Set(columnName string, value Exp) UpdateBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.setItems, b.setItems, 1)

	newBuilder.setItems = append(newBuilder.setItems, updateSetItem{
		columnName: columnName,
		value:      value,
	})
	return newBuilder
}

// SetMap sets the items in the set clause to the given map.
// It overwrites any previous set clause items.
func (b UpdateBuilder) SetMap(m map[string]any) UpdateBuilder {
	newBuilder := b

	columnNames := make([]string, 0, len(m))
	for columnName := range m {
		columnNames = append(columnNames, columnName)
	}

	// Make sure the order of column names is stable.
	sort.Strings(columnNames)

	setItems := make([]updateSetItem, len(m))
	for i, columnName := range columnNames {
		setItems[i] = updateSetItem{
			columnName: columnName,
			value:      Arg(m[columnName]),
		}
	}

	newBuilder.setItems = setItems
	return newBuilder
}

func (b UpdateBuilder) From(from FromExp) FromUpdateBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.from, b.from, 1)

	newBuilder.from = append(newBuilder.from, fromItem{
		from: from,
	})

	return FromUpdateBuilder{
		UpdateBuilder: newBuilder,
	}
}

type FromUpdateBuilder struct {
	UpdateBuilder
}

// As sets the alias for the last added from item.
func (b FromUpdateBuilder) As(alias string) FromUpdateBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.from, b.from, 0)

	lastIdx := len(newBuilder.from) - 1
	newBuilder.from[lastIdx].alias = alias

	return newBuilder
}

// ColumnAliases sets the column aliases for the last added from item.
func (b FromUpdateBuilder) ColumnAliases(aliases ...string) FromUpdateBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.from, b.from, 0)

	lastIdx := len(newBuilder.from) - 1
	newBuilder.from[lastIdx].columnAliases = aliases

	return newBuilder
}

// Where adds a WHERE condition to the update.
// Multiple calls to Where are joined with AND.
func (b UpdateBuilder) Where(cond Exp) UpdateBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.whereConjunction, b.whereConjunction, 1)

	newBuilder.whereConjunction = append(newBuilder.whereConjunction, cond)
	return newBuilder
}

func (b UpdateBuilder) Returning(outputExpression Exp) ReturningUpdateBuilder {
	newBuilder := b
	newBuilder.returningItems = b.returningItems.cloneSlice(1)

	newBuilder.returningItems = append(newBuilder.returningItems, returningItem{
		outputExpression: outputExpression,
	})

	return ReturningUpdateBuilder{newBuilder}
}

type ReturningUpdateBuilder struct {
	UpdateBuilder
}

// As sets the output name for the last output expression.
func (b ReturningUpdateBuilder) As(outputName string) UpdateBuilder {
	newBuilder := b.UpdateBuilder
	newBuilder.returningItems = b.returningItems.cloneSlice(0)

	lastIdx := len(newBuilder.returningItems) - 1
	newBuilder.returningItems[lastIdx].outputName = outputName

	return newBuilder
}

// WriteSQL writes the update as an expression.
func (b UpdateBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteRune('(')
	b.innerWriteSQL(sb)
	sb.WriteRune(')')
}

func (b UpdateBuilder) innerWriteSQL(sb *SQLBuilder) {
	if len(b.withQueries) > 0 {
		b.withQueries.WriteSQL(sb)
	}

	sb.WriteString("UPDATE ")
	b.tableName.WriteSQL(sb)
	if b.alias != "" {
		sb.WriteString(" AS ")
		sb.WriteString(b.alias)
	}
	sb.WriteString(" SET ")
	for i, setItem := range b.setItems {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(quoteIdentifierIfKeyword(setItem.columnName))
		sb.WriteString(" = ")
		setItem.value.WriteSQL(sb)
	}
	if len(b.from) > 0 {
		sb.WriteString(" FROM ")
		for i, f := range b.from {
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

// ApplyIf applies the given function to the builder if the condition is true.
// It returns the builder itself if the condition is false, otherwise it returns the result of the function.
// It's especially helpful for building a query conditionally.
func (b UpdateBuilder) ApplyIf(cond bool, apply func(q UpdateBuilder) UpdateBuilder) UpdateBuilder {
	if cond && apply != nil {
		return apply(b)
	}
	return b
}
