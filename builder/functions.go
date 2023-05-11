package builder

func FuncExp(name string, args []Exp) ExpBase {
	return ExpBase{
		Exp: funcExp{
			name: name,
			args: args,
		},
	}
}

type funcExp struct {
	ExpBase
	name string
	args []Exp
}

func (c funcExp) IsExp() {}

func (c funcExp) WriteSQL(sb *SQLBuilder) {
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
