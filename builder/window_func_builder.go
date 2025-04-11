package builder

// function_name ([expression [, expression ... ]]) [ FILTER ( WHERE filter_clause ) ] OVER window_name
// function_name ([expression [, expression ... ]]) [ FILTER ( WHERE filter_clause ) ] OVER ( window_definition )
// function_name ( * ) [ FILTER ( WHERE filter_clause ) ] OVER window_name
// function_name ( * ) [ FILTER ( WHERE filter_clause ) ] OVER ( window_definition )

// Where `window_definition` has the syntax:
//
// [ existing_window_name ]
// [ PARTITION BY expression [, ...] ]
// [ ORDER BY expression [ ASC | DESC | USING operator ] [ NULLS { FIRST | LAST } ] [, ...] ]
// [ frame_clause ]

// The optional `frame_clause` can be one of:
//
// { RANGE | ROWS | GROUPS } frame_start [ frame_exclusion ]
// { RANGE | ROWS | GROUPS } BETWEEN frame_start AND frame_end [ frame_exclusion ]

// Where `frame_start` and `frame_end` can be one of
//
// UNBOUNDED PRECEDING
// offset PRECEDING
// CURRENT ROW
// offset FOLLOWING
// UNBOUNDED FOLLOWING

// And `frame_exclusion` can be one of:
//
// EXCLUDE CURRENT ROW
// EXCLUDE GROUP
// EXCLUDE TIES
// EXCLUDE NO OTHERS

// TODO Add frame_clause
// TODO Re-use windowDefinition to build the SQL

type WindowFuncCallBuilder struct {
	// FuncCall is the base function / aggregate call
	FuncCall           Exp
	existingWindowName string
	partitionBy        []Exp
	orderBys           []orderByClause
}

func (b WindowFuncCallBuilder) IsExp() {}

func (b WindowFuncCallBuilder) PartitionBy(exp Exp, exps ...Exp) WindowFuncCallBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.partitionBy, b.partitionBy, 1+len(exps))

	newBuilder.partitionBy = append(newBuilder.partitionBy, exp)
	newBuilder.partitionBy = append(newBuilder.partitionBy, exps...)
	return newBuilder
}

func (b WindowFuncCallBuilder) OrderBy(exp Exp) OrderByWindowFuncCallBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.orderBys, b.orderBys, 1)

	newBuilder.orderBys = append(newBuilder.orderBys, orderByClause{
		exp: exp,
	})

	return OrderByWindowFuncCallBuilder{
		WindowFuncCallBuilder: newBuilder,
	}
}

type OrderByWindowFuncCallBuilder struct {
	WindowFuncCallBuilder
}

func (b OrderByWindowFuncCallBuilder) Asc() OrderByWindowFuncCallBuilder {
	return b.setOrder(sortOrderAsc)
}

func (b OrderByWindowFuncCallBuilder) Desc() OrderByWindowFuncCallBuilder {
	return b.setOrder(sortOrderDesc)
}

func (b OrderByWindowFuncCallBuilder) setOrder(order sortOrder) OrderByWindowFuncCallBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.orderBys, b.orderBys, 0)

	lastIdx := len(newBuilder.orderBys) - 1
	newBuilder.orderBys[lastIdx].order = order

	return newBuilder
}

func (b OrderByWindowFuncCallBuilder) NullsFirst() OrderByWindowFuncCallBuilder {
	return b.setNulls(sortNullsFirst)
}

func (b OrderByWindowFuncCallBuilder) NullsLast() OrderByWindowFuncCallBuilder {
	return b.setNulls(sortNullsLast)
}

func (b OrderByWindowFuncCallBuilder) setNulls(nulls sortNulls) OrderByWindowFuncCallBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.orderBys, b.orderBys, 0)

	newBuilder.orderBys[len(newBuilder.orderBys)-1].nulls = nulls

	return newBuilder
}

func (b WindowFuncCallBuilder) WriteSQL(sb *SQLBuilder) {
	b.FuncCall.WriteSQL(sb)
	sb.WriteString(" OVER ")
	if b.existingWindowName != "" && len(b.partitionBy) == 0 && len(b.orderBys) == 0 {
		sb.WriteString(b.existingWindowName)
	} else {
		sb.WriteString("(")
		hasContent := false
		if b.existingWindowName != "" {
			sb.WriteString(b.existingWindowName)
			hasContent = true
		}
		if len(b.partitionBy) > 0 {
			if hasContent {
				sb.WriteRune(' ')
			}
			sb.WriteString("PARTITION BY ")
			for i, exp := range b.partitionBy {
				if i > 0 {
					sb.WriteRune(',')
				}
				exp.WriteSQL(sb)
			}
			hasContent = true
		}
		if len(b.orderBys) > 0 {
			if hasContent {
				sb.WriteRune(' ')
			}
			sb.WriteString("ORDER BY ")
			for i, clause := range b.orderBys {
				if i > 0 {
					sb.WriteRune(',')
				}
				clause.WriteSQL(sb)
			}
		}
		sb.WriteString(")")
	}
}
