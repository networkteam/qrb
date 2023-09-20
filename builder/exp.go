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
