package builder

type Exp interface {
	IsExp()
	SQLWriter
}

// ExpBase is a base type for expressions to allow embedding of various default operators.
type ExpBase struct {
	Exp
}

func (ExpBase) IsExp() {}

// noParensExp is a marker interface for expressions that do not need to be wrapped in parentheses (again) e.g. when combined via other operators (e.g. cast).
type noParensExp interface {
	Exp
	NoParensExp()
}
