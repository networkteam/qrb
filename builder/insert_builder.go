package builder

import (
	"errors"
	"sort"
)

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
	withQueries    withQueries
	tableName      string
	alias          string
	columnNames    []string
	defaultValues  bool
	valueLists     [][]Exp
	query          SelectExp
	returningItems returningItems
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

// DefaultValues sets the DEFAULT VALUES clause to insert a row with default values.
// If InsertBuilder.Values is called after this method, it will overrule the DEFAULT VALUES clause.
func (b InsertBuilder) DefaultValues() InsertBuilder {
	newBuilder := b
	newBuilder.defaultValues = true
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

// Query sets a select query as the values to insert.
func (b InsertBuilder) Query(query SelectExp) InsertBuilder {
	newBuilder := b
	newBuilder.query = query
	return newBuilder
}

func (b InsertBuilder) Returning(outputExpression Exp) ReturningInsertBuilder {
	newBuilder := b

	newBuilder.returningItems = make(returningItems, len(b.returningItems), len(b.returningItems)+1)
	copy(newBuilder.returningItems, b.returningItems)

	newBuilder.returningItems = append(newBuilder.returningItems, returningItem{
		outputExpression: outputExpression,
	})

	return ReturningInsertBuilder{newBuilder}
}

type ReturningInsertBuilder struct {
	InsertBuilder
}

// As sets the output name for the last output expression.
func (b ReturningInsertBuilder) As(outputName string) InsertBuilder {
	newBuilder := b.InsertBuilder

	newBuilder.returningItems = make(returningItems, len(b.returningItems), len(b.returningItems)+1)
	copy(newBuilder.returningItems, b.returningItems)

	lastIdx := len(newBuilder.returningItems) - 1
	newBuilder.returningItems[lastIdx].outputName = outputName

	return newBuilder
}

type returningItem struct {
	outputExpression Exp
	outputName       string
}

type returningItems []returningItem

func (i returningItems) WriteSQL(sb *SQLBuilder) {
	sb.WriteString(" RETURNING ")
	for j, item := range i {
		if j > 0 {
			sb.WriteString(",")
		}
		item.outputExpression.WriteSQL(sb)
		if item.outputName != "" {
			sb.WriteString(" AS ")
			sb.WriteString(item.outputName)
		}
	}
}

var ErrInsertValuesAndQuery = errors.New("insert: cannot set both values and query")

// WriteSQL writes the insert as an expression.
func (b InsertBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteRune('(')
	b.innerWriteSQL(sb)
	sb.WriteRune(')')
}

func (b InsertBuilder) innerWriteSQL(sb *SQLBuilder) {
	if len(b.withQueries) > 0 {
		b.withQueries.WriteSQL(sb)
	}

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
	if b.valueLists != nil && b.query != nil {
		sb.AddError(ErrInsertValuesAndQuery)
		return
	}
	if b.query != nil {
		sb.WriteString(" ")
		b.query.innerWriteSQL(sb)
	} else if b.valueLists != nil {
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
	} else if b.defaultValues {
		sb.WriteString(" DEFAULT VALUES")
	}
	if len(b.returningItems) > 0 {
		b.returningItems.WriteSQL(sb)
	}
}
