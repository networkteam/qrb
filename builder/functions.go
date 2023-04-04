package builder

func funcExp(name string, args []Exp) FuncExp {
	e := FuncExp{
		name: name,
		args: args,
	}
	e.Exp = e // self-reference for base methods
	return e
}

type FuncExp struct {
	ExpBase
	name string
	args []Exp
}

func (c FuncExp) IsExp() {}

func (c FuncExp) WriteSQL(sb *SQLBuilder) {
	sb.WriteString(c.name)
	sb.WriteRune('(')
	for i, exp := range c.args {
		if i > 0 {
			sb.WriteRune(',')
		}
		exp.WriteSQL(sb)
	}
	sb.WriteRune(')')
}
