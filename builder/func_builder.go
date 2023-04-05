package builder

import (
	"errors"
)

// function_name ( [ argument [, ...] ] ) [ WITH ORDINALITY ]

func Func(name string, args ...Exp) FuncBuilder {
	exp := FuncBuilder{
		name: name,
		args: args,
	}
	exp.Exp = exp // self-reference for base methods
	return exp
}

type FuncBuilder struct {
	ExpBase
	name           string
	args           []Exp
	withOrdinality bool
	alias          string
	columnDefs     []funcColumnDefinition
}

type funcColumnDefinition struct {
	name string
	typ  string
}

func (f FuncBuilder) IsExp()            {}
func (f FuncBuilder) isFromExp()        {}
func (f FuncBuilder) isFromLateralExp() {}

func (f FuncBuilder) WithOrdinality() FuncBuilder {
	newBuilder := f
	newBuilder.withOrdinality = true

	newBuilder.Exp = newBuilder // self-reference for base methods

	return newBuilder
}

func (f FuncBuilder) As(alias string) FuncBuilder {
	newBuilder := f
	newBuilder.alias = alias

	newBuilder.Exp = newBuilder // self-reference for base methods

	return newBuilder
}

// ColumnDefinition adds a column definition to the function call.
// To add multiple column definitions, call this method multiple times.
func (f FuncBuilder) ColumnDefinition(name, typ string) FuncBuilder {
	newBuilder := f

	newBuilder.columnDefs = make([]funcColumnDefinition, len(f.columnDefs), len(f.columnDefs)+1)
	copy(newBuilder.columnDefs, f.columnDefs)

	newBuilder.columnDefs = append(newBuilder.columnDefs, funcColumnDefinition{
		name: name,
		typ:  typ,
	})

	newBuilder.Exp = newBuilder // self-reference for base methods

	return newBuilder
}

var errWithOrdinalityAndColumnDefinitions = errors.New("func: WITH ORDINALITY is not supported with column definitions, use ROWS FROM instead")

func (f FuncBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteString(f.name)
	sb.WriteRune('(')
	for i, arg := range f.args {
		if i > 0 {
			sb.WriteRune(',')
		}
		arg.WriteSQL(sb)
	}
	sb.WriteRune(')')
	if f.withOrdinality {
		sb.WriteString(" WITH ORDINALITY")
	}
	if f.alias != "" {
		sb.WriteString(" AS ")
		sb.WriteString(f.alias)
	}
	if len(f.columnDefs) > 0 {
		if f.withOrdinality {
			sb.AddError(errWithOrdinalityAndColumnDefinitions)
			return
		}
		if f.alias == "" {
			sb.WriteString(" AS")
		}
		sb.WriteString(" (")
		for i, def := range f.columnDefs {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(def.name)
			sb.WriteRune(' ')
			sb.WriteString(def.typ)
		}
		sb.WriteString(")")
	}
}
