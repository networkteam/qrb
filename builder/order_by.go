package builder

type sortOrder string

const (
	sortOrderAsc  sortOrder = "ASC"
	sortOrderDesc sortOrder = "DESC"
)

type sortNulls string

const (
	sortNullsFirst sortNulls = "NULLS FIRST"
	sortNullsLast  sortNulls = "NULLS LAST"
)

type orderByClause struct {
	exp   Exp
	order sortOrder
	nulls sortNulls
}

func (s orderByClause) WriteSQL(sb *SQLBuilder) {
	s.exp.WriteSQL(sb)
	if s.order != "" {
		sb.WriteRune(' ')
		sb.WriteString(string(s.order))
	}
	if s.nulls != "" {
		sb.WriteRune(' ')
		sb.WriteString(string(s.nulls))
	}
}
