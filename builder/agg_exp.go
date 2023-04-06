package builder

type AggExpBuilder struct {
	ExpBase

	name               string
	distinct           bool
	exps               []Exp
	orderBys           []orderByClause
	filterConjunction  []Exp
	withinGroupOrderBy bool
}

type OrderByAggExpBuilder struct {
	AggExpBuilder
}

func Agg(name string, exps []Exp) AggExpBuilder {
	exp := AggExpBuilder{
		name: name,
		exps: exps,
	}
	exp.Exp = exp // self-reference for base methods
	return exp
}

// aggregate_name (expression [ , ... ] [ order_by_clause ] ) [ FILTER ( WHERE filter_clause ) ]
// aggregate_name (DISTINCT expression [ , ... ] [ order_by_clause ] ) [ FILTER ( WHERE filter_clause ) ]
// aggregate_name ( [ expression [ , ... ] ] ) WITHIN GROUP ( order_by_clause ) [ FILTER ( WHERE filter_clause ) ]

var _ Exp = AggExpBuilder{}

func (b AggExpBuilder) IsExp() {}

func (b AggExpBuilder) Distinct() AggExpBuilder {
	newBuilder := b

	newBuilder.distinct = true

	newBuilder.Exp = newBuilder // self-reference for base methods
	return newBuilder
}

// [ ORDER BY expression [ ASC | DESC | USING operator ] [ NULLS { FIRST | LAST } ] [, ...] ]

// OrderBy adds an ORDER BY clause to the aggregate function.
// If AggExpBuilder.WithinGroup is called, the ORDER BY clause is used after the aggregate function in WITHIN GROUP.
func (b AggExpBuilder) OrderBy(exp Exp) OrderByAggExpBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.orderBys, b.orderBys, 1)

	newBuilder.orderBys = append(newBuilder.orderBys, orderByClause{
		exp: exp,
	})

	newBuilder.Exp = newBuilder // self-reference for base methods
	return OrderByAggExpBuilder{
		AggExpBuilder: newBuilder,
	}
}

func (b OrderByAggExpBuilder) Asc() OrderByAggExpBuilder {
	return b.setOrder(sortOrderAsc)
}

func (b OrderByAggExpBuilder) Desc() OrderByAggExpBuilder {
	return b.setOrder(sortOrderDesc)
}

func (b OrderByAggExpBuilder) setOrder(order sortOrder) OrderByAggExpBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.orderBys, b.orderBys, 0)

	lastIdx := len(newBuilder.orderBys) - 1
	newBuilder.orderBys[lastIdx].order = order

	newBuilder.Exp = newBuilder // self-reference for base methods
	return newBuilder
}

func (b OrderByAggExpBuilder) NullsFirst() OrderByAggExpBuilder {
	return b.setNulls(sortNullsFirst)
}

func (b OrderByAggExpBuilder) NullsLast() OrderByAggExpBuilder {
	return b.setNulls(sortNullsLast)
}

func (b OrderByAggExpBuilder) setNulls(nulls sortNulls) OrderByAggExpBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.orderBys, b.orderBys, 0)

	newBuilder.orderBys[len(newBuilder.orderBys)-1].nulls = nulls

	newBuilder.Exp = newBuilder // self-reference for base methods
	return newBuilder
}

// Filter adds a filter to the aggregate function.
// Multiple calls to Filter are joined with AND.
func (b AggExpBuilder) Filter(cond Exp) AggExpBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.filterConjunction, b.filterConjunction, 1)

	newBuilder.filterConjunction = append(newBuilder.filterConjunction, cond)

	newBuilder.Exp = newBuilder // self-reference for base methods
	return newBuilder
}

// WithinGroup adds a WITHIN GROUP order by clause after the aggregate function.
// Sort arguments are added via AggExpBuilder.OrderBy.
func (b AggExpBuilder) WithinGroup() AggExpBuilder {
	newBuilder := b

	newBuilder.withinGroupOrderBy = true

	newBuilder.Exp = newBuilder // self-reference for base methods
	return newBuilder
}

func (b AggExpBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteString(b.name)
	sb.WriteRune('(')
	if b.distinct {
		sb.WriteString("DISTINCT ")
	}
	for i, exp := range b.exps {
		if i > 0 {
			sb.WriteRune(',')
		}
		exp.WriteSQL(sb)
	}
	if !b.withinGroupOrderBy && len(b.orderBys) > 0 {
		sb.WriteString(" ORDER BY ")
		for i, clause := range b.orderBys {
			if i > 0 {
				sb.WriteRune(',')
			}
			clause.WriteSQL(sb)
		}
	}
	sb.WriteRune(')')

	if b.withinGroupOrderBy {
		sb.WriteString(" WITHIN GROUP (ORDER BY ")
		for i, clause := range b.orderBys {
			if i > 0 {
				sb.WriteRune(',')
			}
			clause.WriteSQL(sb)
		}
		sb.WriteRune(')')
	}

	if len(b.filterConjunction) > 0 {
		sb.WriteString(" FILTER (WHERE ")
		And(b.filterConjunction...).WriteSQL(sb)
		sb.WriteRune(')')
	}
}
