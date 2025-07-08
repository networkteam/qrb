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

// Mapping of operators to their precedence, higher number means higher precedence.
// We also include other operators (not using Op) to have a complete mapping.
// See https://www.postgresql.org/docs/15/sql-syntax-lexical.html#SQL-PRECEDENCE.
var opPrecedence = map[Operator]int{
	Operator("."):  7,
	Operator("::"): 6,
	// 5: [ ] Array element selection
	// 4: unary plus / minus
	opPow:    3,
	opMult:   2,
	opDivide: 2,
	opMod:    2,
	opPlus:   1,
	opMinus:  1,
	// 0: any other operator (clever use of zero value)
	Operator("BETWEEN"):  -1,
	Operator("IN"):       -1,
	Operator("LIKE"):     -1,
	Operator("ILIKE"):    -1,
	Operator("SIMILAR"):  -1,
	opLessThan:           -2,
	opGreaterThan:        -2,
	opEqual:              -2,
	opLessThanOrEqual:    -2,
	opGreaterThanOrEqual: -2,
	opNotEqual:           -2,
	Operator("IS"):       -3,
	Operator("ISNULL"):   -3,
	Operator("NOTNULL"):  -3,
	Operator("NOT"):      -4,
	Operator("AND"):      -5,
	Operator("OR"):       -5,
}

type Precedencer interface {
	Precedence() int
}

// Op allows to use arbitrary operators.
//
// Example:
//
//	N("a").Op(Operator("^"), Int(5))
func (b ExpBase) Op(op Operator, rgt Exp) ExpBase {
	// Unwrap any ExpBase.
	if rgtIsExpBase, ok := rgt.(ExpBase); ok {
		rgt = rgtIsExpBase.Exp
	}

	return ExpBase{
		Exp: opExp{
			lft: b.Exp,
			op:  op,
			rgt: rgt,
		},
	}
}

type opExp struct {
	lft      Exp
	op       Operator
	rgt      Exp
	unspaced bool
}

func (c opExp) IsExp() {}

func (c opExp) WriteSQL(sb *SQLBuilder) {
	lftNeedsParens := false
	if lftPrecedence, ok := c.lft.(Precedencer); ok {
		lftNeedsParens = lftPrecedence.Precedence() < c.Precedence()
	}

	if lftNeedsParens {
		sb.WriteRune('(')
	}
	c.lft.WriteSQL(sb)
	if lftNeedsParens {
		sb.WriteRune(')')
	}

	if !c.unspaced {
		sb.WriteRune(' ')
	}
	sb.WriteString(string(c.op))
	if !c.unspaced {
		sb.WriteRune(' ')
	}

	rgtNeedsParens := false
	if rgtPrecedence, ok := c.rgt.(Precedencer); ok {
		rgtNeedsParens = rgtPrecedence.Precedence() < c.Precedence()
		// Special case: if the right expression is an opExp with a different operator and the same precedence (e.g. + / -), we need parentheses.
		if rgtOpExp, ok := c.rgt.(opExp); ok && rgtOpExp.op != c.op && rgtPrecedence.Precedence() == c.Precedence() {
			rgtNeedsParens = true
		}
	}

	if rgtNeedsParens {
		sb.WriteRune('(')
	}
	c.rgt.WriteSQL(sb)
	if rgtNeedsParens {
		sb.WriteRune(')')
	}
}

func (c opExp) Precedence() int {
	return opPrecedence[c.op]
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

// --- Subscript

type subscriptExp struct {
	base       Exp
	subscript  Exp
	upperBound Exp
}

func (s subscriptExp) IsExp() {}

func (s subscriptExp) Precedence() int {
	return 5 // Array element selection has precedence 5
}

func (s subscriptExp) WriteSQL(sb *SQLBuilder) {
	// According to PostgreSQL docs, parentheses can be omitted only for 
	// column references and positional parameters
	needsParens := true
	
	// Check if base is a column reference (IdentExp), positional parameter (argExp),
	// or another subscript expression (for chaining like arr[1][2])
	switch s.base.(type) {
	case IdentExp:
		needsParens = false
	case argExp:
		needsParens = false
	case subscriptExp:
		needsParens = false
	case ExpBase:
		// Check if ExpBase wraps an argExp, IdentExp, or subscriptExp
		if expBase, ok := s.base.(ExpBase); ok {
			switch expBase.Exp.(type) {
			case argExp:
				needsParens = false
			case IdentExp:
				needsParens = false
			case subscriptExp:
				needsParens = false
			}
		}
	}

	if needsParens {
		sb.WriteRune('(')
	}
	s.base.WriteSQL(sb)
	if needsParens {
		sb.WriteRune(')')
	}

	sb.WriteRune('[')
	s.subscript.WriteSQL(sb)
	if s.upperBound != nil {
		sb.WriteRune(':')
		s.upperBound.WriteSQL(sb)
	}
	sb.WriteRune(']')
}

// Subscript allows to access an array element by index or array slice by lower and upper bounds.
// When called with one argument: expression[index]
// When called with two arguments: expression[lower:upper]
func (b ExpBase) Subscript(index Exp, upperBound ...Exp) ExpBase {
	var upper Exp
	if len(upperBound) > 0 {
		upper = upperBound[0]
	}
	return ExpBase{
		Exp: subscriptExp{
			base:       b.Exp,
			subscript:  index,
			upperBound: upper,
		},
	}
}

// --- Unary expressions

// Not builds a NOT expression.
func Not(e Exp) Exp {
	return unaryExp{
		prefix:     "NOT",
		exp:        e,
		precedence: opPrecedence["NOT"],
	}
}

// IsNull builds an IS NULL expression.
func (b ExpBase) IsNull() Exp {
	return unaryExp{
		exp:        b.Exp,
		suffix:     "IS NULL",
		precedence: opPrecedence["IS"],
	}
}

// IsNotNull builds an IS NOT NULL expression.
func (b ExpBase) IsNotNull() Exp {
	return unaryExp{
		exp:        b.Exp,
		suffix:     "IS NOT NULL",
		precedence: opPrecedence["IS"],
	}
}

type unaryExp struct {
	prefix     string
	exp        Exp
	suffix     string
	precedence int
}

func (u unaryExp) IsExp() {}

func (u unaryExp) Precedence() int {
	return u.precedence
}

func (u unaryExp) WriteSQL(sb *SQLBuilder) {
	if u.prefix != "" {
		sb.WriteString(u.prefix)
		sb.WriteRune(' ')
	}

	needsParens := false
	if expPrecedence, ok := u.exp.(Precedencer); ok {
		needsParens = expPrecedence.Precedence() < u.precedence
	}

	if needsParens {
		sb.WriteRune('(')
	}
	u.exp.WriteSQL(sb)
	if needsParens {
		sb.WriteRune(')')
	}

	if u.suffix != "" {
		sb.WriteRune(' ')
		sb.WriteString(u.suffix)
	}
}

// --- Junction expressions

// And builds an AND expression of non-nil expressions.
func And(exps ...Exp) Exp {
	return junctionExp{
		exps: nonNil(exps),
		op:   "AND",
	}
}

// Or builds an OR expression of non-nil expressions.
func Or(exps ...Exp) Exp {
	return junctionExp{
		exps: nonNil(exps),
		op:   "OR",
	}
}

type junctionExp struct {
	exps []Exp
	op   string
}

func (c junctionExp) IsExp() {}

func (c junctionExp) Precedence() int {
	return opPrecedence[Operator(c.op)]
}

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

func (b ExpBase) Cast(typ string) ExpBase {
	exp := opExp{
		lft:      b.Exp,
		op:       Operator("::"),
		rgt:      expType(typ),
		unspaced: true,
	}
	return ExpBase{Exp: exp}
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

// --- Regexp

func (b ExpBase) RegexpMatch(pattern Exp) Exp {
	return b.Op(opRegexpMatch, pattern)
}

func (b ExpBase) RegexpIMatch(pattern Exp) Exp {
	return b.Op(opRegexpIMatch, pattern)
}

func (b ExpBase) RegexpNotMatch(pattern Exp) Exp {
	return b.Op(opRegexpNotMatch, pattern)
}

func (b ExpBase) RegexpINotMatch(pattern Exp) Exp {
	return b.Op(opRegexpINotMatch, pattern)
}
