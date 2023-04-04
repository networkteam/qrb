package builder

// string LIKE pattern [ESCAPE escape-character]

func (b ExpBase) Like(rgt Exp) MatchingBuilder {
	return matchingExp{
		lft: b.Exp,
		rgt: rgt,
		op:  "LIKE",
	}
}

func (b ExpBase) ILike(rgt Exp) MatchingBuilder {
	return matchingExp{
		lft: b.Exp,
		rgt: rgt,
		op:  "ILIKE",
	}
}

// string NOT LIKE pattern [ESCAPE escape-character]

func (b ExpBase) NotLike(rgt Exp) MatchingBuilder {
	return matchingExp{
		lft: b.Exp,
		rgt: rgt,
		op:  "NOT LIKE",
	}
}

func (b ExpBase) NotILike(rgt Exp) MatchingBuilder {
	return matchingExp{
		lft: b.Exp,
		rgt: rgt,
		op:  "NOT ILIKE",
	}
}

// string SIMILAR TO pattern [ESCAPE escape-character]

func (b ExpBase) SimilarTo(rgt Exp) MatchingBuilder {
	return matchingExp{
		lft: b.Exp,
		rgt: rgt,
		op:  "SIMILAR TO",
	}
}

// string NOT SIMILAR TO pattern [ESCAPE escape-character]

func (b ExpBase) NotSimilarTo(rgt Exp) MatchingBuilder {
	return matchingExp{
		lft: b.Exp,
		rgt: rgt,
		op:  "NOT SIMILAR TO",
	}
}

// ---

type MatchingBuilder interface {
	Exp
	Escape(escapeCharacter rune) MatchingBuilder
}

type matchingExp struct {
	lft    Exp
	rgt    Exp
	op     string
	escape *rune
}

func (l matchingExp) Escape(escapeCharacter rune) MatchingBuilder {
	l.escape = &escapeCharacter
	return l
}

func (l matchingExp) IsExp() {}

func (l matchingExp) WriteSQL(sb *SQLBuilder) {
	l.lft.WriteSQL(sb)
	sb.WriteRune(' ')
	sb.WriteString(l.op)
	sb.WriteRune(' ')
	l.rgt.WriteSQL(sb)
	if l.escape != nil {
		sb.WriteString(" ESCAPE ")
		sb.WriteString(pqQuoteLiteral(string(*l.escape)))
	}
}
