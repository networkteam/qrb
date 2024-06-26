package builder

type subqueryExp struct {
	op  string
	exp Exp
}

func (s subqueryExp) IsExp() {}

func (s subqueryExp) WriteSQL(sb *SQLBuilder) {
	sb.WriteString(s.op)
	sb.WriteRune(' ')

	_, isSelect := s.exp.(SelectExp)
	if !isSelect {
		sb.WriteRune('(')
	}
	s.exp.WriteSQL(sb)
	if !isSelect {
		sb.WriteRune(')')
	}
}

func Any(exp Exp) Exp {
	return subqueryExp{
		op:  "ANY",
		exp: exp,
	}
}

func All(exp Exp) Exp {
	return subqueryExp{
		op:  "ALL",
		exp: exp,
	}
}
