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
	withQueries                      withQueries
	tableName                        string
	alias                            string
	columnNames                      []string
	defaultValues                    bool
	valueLists                       [][]Exp
	query                            SelectExp
	conflictTargets                  []conflictTarget
	conflictTargetWhereConjunction   []Exp
	conflictConstraintName           string
	conflictAction                   string
	conflictDoUpdateSetItems         []updateSetItem
	conflictDoUpdateWhereConjunction []Exp
	returningItems                   returningItems
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
	cloneSlice(&newBuilder.valueLists, b.valueLists, 1)

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

// OnConflict sets the ON CONFLICT clause with a conflict target expression to the insert.
// Multiple conflict targets or none can be specified (e.g. for index column names).
// Specify no conflict target for later addition of ON CONSTRAINT or ON CONFLICT DO NOTHING.
func (b InsertBuilder) OnConflict(conflictTargets ...Exp) OnConflictInsertBuilder {
	newBuilder := b

	cloneSlice(&newBuilder.conflictTargets, b.conflictTargets, len(conflictTargets))

	for _, target := range conflictTargets {
		newBuilder.conflictTargets = append(newBuilder.conflictTargets, conflictTarget{
			exp: target,
		})
	}

	return OnConflictInsertBuilder{
		builder: newBuilder,
	}
}

type OnConflictInsertBuilder struct {
	builder InsertBuilder
}

func (b OnConflictInsertBuilder) OnConstraint(constraintName string) OnConflictInsertBuilder {
	newBuilder := b.builder

	newBuilder.conflictConstraintName = constraintName

	return OnConflictInsertBuilder{newBuilder}

}

func (b OnConflictInsertBuilder) DoUpdate() OnConflictDoUpdateInsertBuilder {
	newBuilder := b.builder

	newBuilder.conflictAction = "DO UPDATE"

	return OnConflictDoUpdateInsertBuilder{newBuilder}
}

func (b OnConflictInsertBuilder) DoNothing() InsertBuilder {
	newBuilder := b.builder

	newBuilder.conflictAction = "DO NOTHING"

	return newBuilder
}

// Where adds a WHERE condition as the index predicate to the conflict target.
// Multiple calls to Where are joined with AND.
func (b OnConflictInsertBuilder) Where(cond Exp) OnConflictInsertBuilder {
	newBuilder := b.builder
	cloneSlice(&newBuilder.conflictTargetWhereConjunction, b.builder.conflictTargetWhereConjunction, 1)

	newBuilder.conflictTargetWhereConjunction = append(newBuilder.conflictTargetWhereConjunction, cond)

	return OnConflictInsertBuilder{
		builder: newBuilder,
	}
}

type OnConflictDoUpdateInsertBuilder struct {
	InsertBuilder
}

// Set adds a SET column = value to the DO UPDATE conflict action.
func (b OnConflictDoUpdateInsertBuilder) Set(columnName string, value Exp) OnConflictDoUpdateInsertBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.conflictDoUpdateSetItems, b.conflictDoUpdateSetItems, 1)

	newBuilder.conflictDoUpdateSetItems = append(newBuilder.conflictDoUpdateSetItems, updateSetItem{
		columnName: columnName,
		value:      value,
	})
	return newBuilder
}

// Where adds a WHERE condition to the DO UPDATE conflict action.
// Multiple calls to Where are joined with AND.
func (b OnConflictDoUpdateInsertBuilder) Where(cond Exp) OnConflictDoUpdateInsertBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.conflictDoUpdateWhereConjunction, b.conflictDoUpdateWhereConjunction, 1)

	newBuilder.conflictDoUpdateWhereConjunction = append(newBuilder.conflictDoUpdateWhereConjunction, cond)

	return newBuilder
}

type conflictTarget struct {
	exp Exp
}

func (b InsertBuilder) Returning(outputExpression Exp) ReturningInsertBuilder {
	newBuilder := b
	newBuilder.returningItems = b.returningItems.cloneSlice(1)

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
	newBuilder.returningItems = b.returningItems.cloneSlice(0)

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

func (i returningItems) cloneSlice(additionalCapacity int) returningItems {
	newSlice := make(returningItems, len(i), len(i)+additionalCapacity)
	copy(newSlice, i)
	return newSlice
}

var ErrInsertValuesAndQuery = errors.New("insert: cannot set both values and query")

// WriteSQL writes the insert as an expression.
func (b InsertBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteRune('(')
	b.innerWriteSQL(sb)
	sb.WriteRune(')')
}

var ErrInsertConflictConstraintAndTarget = errors.New("insert: cannot set both conflict constraint name and targets")

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

	if b.conflictAction != "" {
		sb.WriteString(" ON CONFLICT")
		if b.conflictConstraintName != "" && len(b.conflictTargets) > 0 {
			sb.AddError(ErrInsertConflictConstraintAndTarget)
			return
		}
		if b.conflictConstraintName != "" {
			sb.WriteString(" ON CONSTRAINT ")
			sb.WriteString(b.conflictConstraintName)
		}
		if len(b.conflictTargets) > 0 {
			sb.WriteString(" (")
			for i, target := range b.conflictTargets {
				if i > 0 {
					sb.WriteString(",")
				}
				target.exp.WriteSQL(sb)
			}
			sb.WriteString(")")
		}
		if len(b.conflictTargetWhereConjunction) > 0 {
			sb.WriteString(" WHERE ")
			And(b.conflictTargetWhereConjunction...).WriteSQL(sb)
		}
		sb.WriteString(" ")
		sb.WriteString(b.conflictAction)
		if b.conflictAction == "DO UPDATE" {
			if len(b.conflictDoUpdateSetItems) > 0 {
				sb.WriteString(" SET ")
				for i, item := range b.conflictDoUpdateSetItems {
					if i > 0 {
						sb.WriteString(",")
					}
					sb.WriteString(item.columnName)
					sb.WriteString(" = ")
					item.value.WriteSQL(sb)
				}
			}
			if len(b.conflictDoUpdateWhereConjunction) > 0 {
				sb.WriteString(" WHERE ")
				And(b.conflictDoUpdateWhereConjunction...).WriteSQL(sb)
			}
		}
	}

	if len(b.returningItems) > 0 {
		b.returningItems.WriteSQL(sb)
	}
}
