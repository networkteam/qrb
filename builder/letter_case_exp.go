package builder

type LetterCaseBuilder struct {
	ExpBase

	name    string
	identer Identer
}

func LetterCase(name string, identer Identer) LetterCaseBuilder {
	return LetterCaseBuilder{
		name:    name,
		identer: identer,
	}
}

var _ Exp = LetterCaseBuilder{}

func (l LetterCaseBuilder) IsExp() {}

func (l LetterCaseBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteString(l.name)
	sb.WriteRune('(')
	sb.WriteString(l.identer.Ident())
	sb.WriteRune(')')
}
