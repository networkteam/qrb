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

func (b FuncBuilder) IsExp()            {}
func (b FuncBuilder) isFromExp()        {}
func (b FuncBuilder) isFromLateralExp() {}
func (b FuncBuilder) NoParensExp()      {}

func (b FuncBuilder) WithOrdinality() FuncBuilder {
	newBuilder := b
	newBuilder.withOrdinality = true

	newBuilder.Exp = newBuilder // self-reference for base methods

	return newBuilder
}

func (b FuncBuilder) As(alias string) FuncBuilder {
	newBuilder := b
	newBuilder.alias = alias

	newBuilder.Exp = newBuilder // self-reference for base methods

	return newBuilder
}

// ColumnDefinition adds a column definition to the function call.
// To add multiple column definitions, call this method multiple times.
func (b FuncBuilder) ColumnDefinition(name, typ string) FuncBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.columnDefs, b.columnDefs, 1)

	newBuilder.columnDefs = append(newBuilder.columnDefs, funcColumnDefinition{
		name: name,
		typ:  typ,
	})

	newBuilder.Exp = newBuilder // self-reference for base methods

	return newBuilder
}

var errWithOrdinalityAndColumnDefinitions = errors.New("func: WITH ORDINALITY is not supported with column definitions, use ROWS FROM instead")

func (b FuncBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteString(b.name)
	sb.WriteRune('(')
	for i, arg := range b.args {
		if i > 0 {
			sb.WriteRune(',')
		}
		arg.WriteSQL(sb)
	}
	sb.WriteRune(')')
	if b.withOrdinality {
		sb.WriteString(" WITH ORDINALITY")
	}
	if b.alias != "" {
		sb.WriteString(" AS ")
		sb.WriteString(b.alias)
	}
	if len(b.columnDefs) > 0 {
		if b.withOrdinality {
			sb.AddError(errWithOrdinalityAndColumnDefinitions)
			return
		}
		if b.alias == "" {
			sb.WriteString(" AS")
		}
		sb.WriteString(" (")
		for i, def := range b.columnDefs {
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
