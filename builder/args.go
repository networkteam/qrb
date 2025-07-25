package builder

// Arg creates an expression that represents an argument that will be bound to a placeholder with the given value.
// Each call to Arg will create a new placeholder and emit the argument when writing the query.
func Arg(argument any) ExpBase {
	return ExpBase{
		Exp: argExp{
			argument: argument,
		},
	}
}

type argExp struct {
	argument any
}

func (a argExp) IsExp() {}

func (a argExp) WriteSQL(sb *SQLBuilder) {
	p := sb.CreatePlaceholder(a.argument)
	sb.WriteString(p)
}

type Expressions struct {
	exps []Exp
	ExpBase
}

func ToExpressions(exps ...Exp) Expressions {
	return Expressions{
		exps: exps,
		ExpBase: ExpBase{
			Exp: Expressions{exps: exps},
		},
	}
}

func unwrapExpressions(exps []Expressions) [][]Exp {
	result := make([][]Exp, len(exps))
	for i, e := range exps {
		result[i] = e.exps
	}
	return result
}

func (e Expressions) IsExp() {}

func (e Expressions) WriteSQL(sb *SQLBuilder) {
	sb.WriteRune('(')
	for i, exp := range e.exps {
		if i > 0 {
			sb.WriteRune(',')
		}
		exp.WriteSQL(sb)
	}
	sb.WriteRune(')')
}

func (e Expressions) isSelectOrExpressions() {}

// Args creates argument expressions for the given arguments.
func Args[T any](arguments ...T) Expressions {
	exps := make([]Exp, len(arguments))
	for i, arg := range arguments {
		exps[i] = Arg(arg)
	}
	return ToExpressions(exps...)
}

// Bind creates an expression that represents an argument that will be bound to a placeholder with the given value.
// If Bind is called again with the same name, the same placeholder will be used.
func Bind(argName string) ExpBase {
	return ExpBase{
		Exp: bindExp{
			name: argName,
		},
	}
}

type bindExp struct {
	name string
}

func (b bindExp) IsExp() {}

func (b bindExp) WriteSQL(sb *SQLBuilder) {
	p := sb.BindPlaceholder(b.name)
	sb.WriteString(p)
}
