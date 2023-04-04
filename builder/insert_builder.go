package builder

import "sort"

// [ WITH [ RECURSIVE ] with_query [, ...] ]
// INSERT INTO table_name [ AS alias ] [ ( column_name [, ...] ) ]
//     [ OVERRIDING { SYSTEM | USER } VALUE ]
//     { DEFAULT VALUES | VALUES ( { expression | DEFAULT } [, ...] ) [, ...] | query }
//     [ ON CONFLICT [ conflict_target ] conflict_action ]
//     [ RETURNING * | output_expression [ [ AS ] output_name ] [, ...] ]

func InsertInto(tableName string) InsertBuilder {
	return InsertBuilder{
		tableName: tableName,
	}
}

type InsertBuilder struct {
	tableName   string
	alias       string
	columnNames []string
	valueLists  [][]Exp
}

func (b InsertBuilder) As(alias string) InsertBuilder {
	newBuilder := b
	newBuilder.alias = alias
	return newBuilder
}

func (b InsertBuilder) ColumnNames(columnName string, rest ...string) InsertBuilder {
	newBuilder := b
	newBuilder.columnNames = append([]string{columnName}, rest...)
	return newBuilder
}

// Values appends the given values to insert.
// It can be called multiple times to insert multiple rows.
func (b InsertBuilder) Values(values ...Exp) InsertBuilder {
	newBuilder := b

	newBuilder.valueLists = make([][]Exp, len(newBuilder.valueLists), len(newBuilder.valueLists)+1)
	copy(newBuilder.valueLists, b.valueLists)

	newBuilder.valueLists = append(newBuilder.valueLists, values)
	return newBuilder
}

// SetMap sets the column names and values to insert from the given map.
// It overwrites any previous column names and values.
func (b InsertBuilder) SetMap(m map[string]any) InsertBuilder {
	newBuilder := b

	columnNames := make([]string, 0, len(m))
	for columnName := range m {
		columnNames = append(columnNames, columnName)
	}

	// Make sure the order of column names is stable.
	sort.Strings(columnNames)

	values := make([]Exp, len(m))
	for i, columnName := range columnNames {
		values[i] = Arg(m[columnName])
	}

	newBuilder.columnNames = columnNames
	newBuilder.valueLists = [][]Exp{values}
	return newBuilder
}

func (b InsertBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("INSERT INTO ")
	sb.WriteString(b.tableName)
	if b.alias != "" {
		sb.WriteString(" AS ")
		sb.WriteString(b.alias)
	}
	if b.columnNames != nil {
		sb.WriteString(" (")
		for i, columnName := range b.columnNames {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(columnName)
		}
		sb.WriteString(")")
	}
	if b.valueLists != nil {
		sb.WriteString(" VALUES ")
		for i, valueList := range b.valueLists {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString("(")
			for j, value := range valueList {
				if j > 0 {
					sb.WriteString(",")
				}
				value.WriteSQL(sb)
			}
			sb.WriteString(")")
		}
	}

}
