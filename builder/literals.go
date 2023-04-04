package builder

import "strconv"

func String(s string) Exp {
	return expStr(s)
}

type expStr string

func (e expStr) IsExp() {}

func (e expStr) WriteSQL(sb *SQLBuilder) {
	sb.WriteString(pqQuoteLiteral(string(e)))
}

func Float(f float64) Exp {
	return expFloat(f)
}

type expFloat float64

func (e expFloat) IsExp() {}

func (e expFloat) WriteSQL(sb *SQLBuilder) {
	sb.WriteString(strconv.FormatFloat(float64(e), 'f', -1, 64))
}

func Int(s int) Exp {
	return expInt(s)
}

type expInt int

func (e expInt) IsExp() {}

func (e expInt) WriteSQL(sb *SQLBuilder) {
	sb.WriteString(strconv.Itoa(int(e)))
}

func Bool(b bool) Exp {
	return expBool(b)
}

type expBool bool

func (e expBool) IsExp() {}

func (e expBool) WriteSQL(sb *SQLBuilder) {
	sb.WriteString(strconv.FormatBool(bool(e)))
}

// Array builds an array literal.
//
// Make sure that all elements are of the same type.
func Array(elems ...Exp) Exp {
	return expArray(elems)
}

type expArray []Exp

func (e expArray) IsExp() {}

func (e expArray) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("ARRAY[")
	for i, elem := range e {
		if i > 0 {
			sb.WriteRune(',')
		}
		elem.WriteSQL(sb)
	}
	sb.WriteString("]")
}

// Null builds the null literal.
func Null() Exp {
	return expNull{}
}

type expNull struct{}

func (e expNull) IsExp() {}

func (e expNull) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("NULL")
}

// Interval builds an interval constant.
func Interval(spec string) Exp {
	return expInterval{
		spec: spec,
	}
}

type expInterval struct {
	spec string
}

func (e expInterval) IsExp() {}

func (e expInterval) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("INTERVAL ")
	sb.WriteString(pqQuoteLiteral(e.spec))
}
