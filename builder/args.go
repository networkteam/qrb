package builder

// Arg creates an expression that represents an argument that will be bound to a placeholder with the given value.
// Each call to Arg will create a new placeholder and emit the argument when writing the query.
func Arg(argument any) Exp {
	return argExp{
		argument: argument,
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

// Args creates argument expressions for the given arguments.
func Args(argument any, rest ...any) []Exp {
	exps := make([]Exp, 1+len(rest))
	exps[0] = Arg(argument)
	for i, arg := range rest {
		exps[i+1] = Arg(arg)
	}
	return exps
}

// Bind creates an expression that represents an argument that will be bound to a placeholder with the given value.
// If Bind is called again with the same name, the same placeholder will be used.
func Bind(argName string) Exp {
	return bindExp{
		name: argName,
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
