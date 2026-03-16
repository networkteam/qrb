package builder

// columnDef represents a column definition in a CREATE TABLE or ALTER TABLE ADD COLUMN statement.
type columnDef struct {
	name              string
	typeName          string // raw SQL type string
	notNull           bool
	defaultExp        Exp
	primaryKey        bool
	unique            bool
	check             Exp
	references        *columnReference
	generatedIdentity string // "ALWAYS" or "BY DEFAULT" for GENERATED ... AS IDENTITY
	generatedAs       Exp    // expression for GENERATED ALWAYS AS (expr) STORED
	generatedStored   bool   // true = STORED, false = VIRTUAL
}

type columnReference struct {
	table             Identer
	columns           []string
	onDelete          string
	onUpdate          string
	deferrable        *bool
	initiallyDeferred bool
}

func (c columnDef) writeSQL(sb *SQLBuilder) {
	sb.WriteString(quoteIdentifierIfKeyword(c.name))
	sb.WriteRune(' ')
	sb.WriteString(c.typeName)
	if c.notNull {
		sb.WriteString(" NOT NULL")
	}
	if c.defaultExp != nil {
		sb.WriteString(" DEFAULT ")
		c.defaultExp.WriteSQL(sb)
	}
	if c.primaryKey {
		sb.WriteString(" PRIMARY KEY")
	}
	if c.unique {
		sb.WriteString(" UNIQUE")
	}
	if c.generatedIdentity != "" {
		sb.WriteString(" GENERATED ")
		sb.WriteString(c.generatedIdentity)
		sb.WriteString(" AS IDENTITY")
	}
	if c.generatedAs != nil {
		sb.WriteString(" GENERATED ALWAYS AS (")
		c.generatedAs.WriteSQL(sb)
		sb.WriteRune(')')
		if c.generatedStored {
			sb.WriteString(" STORED")
		} else {
			sb.WriteString(" VIRTUAL")
		}
	}
	if c.check != nil {
		sb.WriteString(" CHECK (")
		c.check.WriteSQL(sb)
		sb.WriteRune(')')
	}
	if c.references != nil {
		c.references.writeSQL(sb)
	}
}

func (r columnReference) writeSQL(sb *SQLBuilder) {
	sb.WriteString(" REFERENCES ")
	r.table.WriteSQL(sb)
	if len(r.columns) > 0 {
		sb.WriteString(" (")
		writeColumnList(sb, r.columns)
		sb.WriteRune(')')
	}
	if r.onDelete != "" {
		sb.WriteString(" ON DELETE ")
		sb.WriteString(r.onDelete)
	}
	if r.onUpdate != "" {
		sb.WriteString(" ON UPDATE ")
		sb.WriteString(r.onUpdate)
	}
	writeDeferrable(sb, r.deferrable, r.initiallyDeferred)
}

func writeDeferrable(sb *SQLBuilder, deferrable *bool, initiallyDeferred bool) {
	if deferrable != nil {
		if *deferrable {
			sb.WriteString(" DEFERRABLE")
		} else {
			sb.WriteString(" NOT DEFERRABLE")
		}
		if initiallyDeferred {
			sb.WriteString(" INITIALLY DEFERRED")
		}
	}
}

type tableConstraintKind int

const (
	tableConstraintPrimaryKey tableConstraintKind = iota
	tableConstraintUnique
	tableConstraintForeignKey
	tableConstraintCheck
)

type tableConstraint struct {
	constraintName    string
	kind              tableConstraintKind
	columns           []string
	refTable          Identer
	refColumns        []string
	onDelete          string
	onUpdate          string
	checkExp          Exp
	deferrable        *bool
	initiallyDeferred bool
}

func (c tableConstraint) writeSQL(sb *SQLBuilder) {
	if c.constraintName != "" {
		sb.WriteString("CONSTRAINT ")
		sb.WriteString(quoteIdentifierIfKeyword(c.constraintName))
		sb.WriteRune(' ')
	}
	switch c.kind {
	case tableConstraintPrimaryKey:
		sb.WriteString("PRIMARY KEY (")
		writeColumnList(sb, c.columns)
		sb.WriteRune(')')
	case tableConstraintUnique:
		sb.WriteString("UNIQUE (")
		writeColumnList(sb, c.columns)
		sb.WriteRune(')')
	case tableConstraintForeignKey:
		sb.WriteString("FOREIGN KEY (")
		writeColumnList(sb, c.columns)
		sb.WriteString(") REFERENCES ")
		c.refTable.WriteSQL(sb)
		if len(c.refColumns) > 0 {
			sb.WriteString(" (")
			writeColumnList(sb, c.refColumns)
			sb.WriteRune(')')
		}
		if c.onDelete != "" {
			sb.WriteString(" ON DELETE ")
			sb.WriteString(c.onDelete)
		}
		if c.onUpdate != "" {
			sb.WriteString(" ON UPDATE ")
			sb.WriteString(c.onUpdate)
		}
	case tableConstraintCheck:
		sb.WriteString("CHECK (")
		c.checkExp.WriteSQL(sb)
		sb.WriteRune(')')
	}
	writeDeferrable(sb, c.deferrable, c.initiallyDeferred)
}

func writeColumnList(sb *SQLBuilder, columns []string) {
	for i, col := range columns {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(quoteIdentifierIfKeyword(col))
	}
}
