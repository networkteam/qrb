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

func Update(tableName string) UpdateBuilder {
	return UpdateBuilder{
		tableName: tableName,
	}
}

type UpdateBuilder struct {
	tableName        string
	alias            string
	setItems         []updateSetItem
	whereConjunction []Exp
}

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

	newBuilder.setItems = make([]updateSetItem, len(newBuilder.setItems), len(newBuilder.setItems)+1)
	copy(newBuilder.setItems, b.setItems)

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

// Where adds a WHERE condition to the update.
// Multiple calls to Where are joined with AND.
func (b UpdateBuilder) Where(cond Exp) UpdateBuilder {
	newBuilder := b

	newBuilder.whereConjunction = make([]Exp, len(b.whereConjunction), len(b.whereConjunction)+1)
	copy(newBuilder.whereConjunction, b.whereConjunction)

	newBuilder.whereConjunction = append(newBuilder.whereConjunction, cond)
	return newBuilder
}

func (b UpdateBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("UPDATE ")
	sb.WriteString(b.tableName)
	if b.alias != "" {
		sb.WriteString(" AS ")
		sb.WriteString(b.alias)
	}
	sb.WriteString(" SET ")
	for i, setItem := range b.setItems {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(setItem.columnName)
		sb.WriteString(" = ")
		setItem.value.WriteSQL(sb)
	}
	if len(b.whereConjunction) > 0 {
		sb.WriteString(" WHERE ")
		And(b.whereConjunction...).WriteSQL(sb)
	}
}
