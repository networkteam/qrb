package builder

type Operator string

const (
	opEqual              Operator = "="
	opLessThan           Operator = "<"
	opLessThanOrEqual    Operator = "<="
	opGreaterThan        Operator = ">"
	opGreaterThanOrEqual Operator = ">="
	opNotEqual           Operator = "<>"

	opConcat Operator = "||"

	// opRegexpMatch matches a string with a POSIX regular expression pattern, case-sensitive.
	opRegexpMatch Operator = "~"
	// opRegexpIMatch matches a string with a POSIX regular expression pattern, case-insensitive.
	opRegexpIMatch Operator = "~*"
	// opRegexpNotMatch does not match a string with a POSIX regular expression pattern, case-sensitive.
	opRegexpNotMatch Operator = "!~"
	// opRegexpINotMatch does not match a string with a POSIX regular expression pattern, case-insensitive.
	opRegexpINotMatch Operator = "!~*"
)

// Op allows to use arbitrary operators.
//
// Example:
//
//	N("a").Op(Operator("^"), Int(5))
func (b ExpBase) Op(op Operator, rgt Exp) ExpBase {
	return ExpBase{
		Exp: opExp{
			lft: b.Exp,
			op:  op,
			rgt: rgt,
		},
	}
}

type opExp struct {
	lft Exp
	op  Operator
	rgt Exp
}

func (c opExp) IsExp() {}

func (c opExp) WriteSQL(sb *SQLBuilder) {
	c.lft.WriteSQL(sb)
	sb.WriteRune(' ')
	sb.WriteString(string(c.op))
	sb.WriteRune(' ')
	c.rgt.WriteSQL(sb)
}

// Common operators

func (b ExpBase) Eq(rgt Exp) Exp {
	return b.Op(opEqual, rgt)
}

func (b ExpBase) Neq(rgt Exp) Exp {
	return b.Op(opNotEqual, rgt)
}

func (b ExpBase) Lt(rgt Exp) Exp {
	return b.Op(opLessThan, rgt)
}

func (b ExpBase) Lte(rgt Exp) Exp {
	return b.Op(opLessThanOrEqual, rgt)
}

func (b ExpBase) Gt(rgt Exp) Exp {
	return b.Op(opGreaterThan, rgt)
}

func (b ExpBase) Gte(rgt Exp) Exp {
	return b.Op(opGreaterThanOrEqual, rgt)
}

// --- String operators

func (b ExpBase) Concat(rgt Exp) ExpBase {
	return b.Op(opConcat, rgt)
}

// --- Unary expressions

// Not builds a NOT expression.
func Not(e Exp) Exp {
	return unaryExp{
		prefix: "NOT",
		exp:    e,
	}
}

// IsNull builds an IS NULL expression.
func (b ExpBase) IsNull() Exp {
	return unaryExp{
		exp:    b.Exp,
		suffix: "IS NULL",
	}
}

// IsNotNull builds an IS NOT NULL expression.
func (b ExpBase) IsNotNull() Exp {
	return unaryExp{
		exp:    b.Exp,
		suffix: "IS NOT NULL",
	}
}

type unaryExp struct {
	prefix string
	exp    Exp
	suffix string
}

func (u unaryExp) IsExp() {}

func (u unaryExp) WriteSQL(sb *SQLBuilder) {
	if u.prefix != "" {
		sb.WriteString(u.prefix)
		sb.WriteRune(' ')
	}
	u.exp.WriteSQL(sb)
	if u.suffix != "" {
		sb.WriteRune(' ')
		sb.WriteString(u.suffix)
	}
}

// --- Junction expressions

func And(exps ...Exp) Exp {
	return junctionExp{
		exps: exps,
		op:   "AND",
	}
}

func Or(exps ...Exp) Exp {
	return junctionExp{
		exps: exps,
		op:   "OR",
	}
}

type junctionExp struct {
	exps []Exp
	op   string
}

func (c junctionExp) IsExp() {}

func (c junctionExp) WriteSQL(sb *SQLBuilder) {
	if len(c.exps) == 1 {
		c.exps[0].WriteSQL(sb)
		return
	}

	for i, exp := range c.exps {
		if i > 0 {
			sb.WriteRune(' ')
			sb.WriteString(c.op)
			sb.WriteRune(' ')
		}
		// Check if the expression is a junction expression and wrap it in parentheses.
		if _, ok := exp.(junctionExp); ok {
			sb.WriteRune('(')
			exp.WriteSQL(sb)
			sb.WriteRune(')')
		} else {
			exp.WriteSQL(sb)
		}
	}
}

func (b ExpBase) Cast(typ string) Exp {
	return castExp{
		exp: b.Exp,
		typ: typ,
	}
}

type castExp struct {
	exp Exp
	typ string
}

func (c castExp) IsExp() {}

func (c castExp) WriteSQL(sb *SQLBuilder) {
	_, needsNoParens := c.exp.(noParensExp)
	if !needsNoParens {
		sb.WriteRune('(')
	}
	c.exp.WriteSQL(sb)
	if !needsNoParens {
		sb.WriteRune(')')
	}
	sb.WriteString("::")
	sb.WriteString(c.typ)
}

func (b ExpBase) In(selectOrExpressions SelectOrExpressions) Exp {
	return inExp{
		lft: b.Exp,
		op:  "IN",
		rgt: selectOrExpressions,
	}
}

func (b ExpBase) NotIn(selectOrExpressions SelectOrExpressions) Exp {
	return inExp{
		lft: b.Exp,
		op:  "NOT IN",
		rgt: selectOrExpressions,
	}
}

type SelectOrExpressions interface {
	Exp
	isSelectOrExpressions()
}

type inExp struct {
	lft Exp
	op  string
	rgt SelectOrExpressions
}

func (c inExp) IsExp() {}

func (c inExp) WriteSQL(sb *SQLBuilder) {
	c.lft.WriteSQL(sb)
	sb.WriteRune(' ')
	sb.WriteString(c.op)
	sb.WriteRune(' ')
	c.rgt.WriteSQL(sb)
}

func Exists(subquery SelectExp) Exp {
	return existsExp{
		subquery: subquery,
	}
}

type existsExp struct {
	subquery SelectExp
}

func (c existsExp) IsExp() {}

func (c existsExp) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("EXISTS ")
	c.subquery.WriteSQL(sb)
}
