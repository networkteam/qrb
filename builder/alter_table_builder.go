package builder

type alterActionKind int

const (
	alterActionAddColumn alterActionKind = iota
	alterActionDropColumn
	alterActionAddConstraint
	alterActionDropConstraint
	alterActionRenameColumn
	alterActionRenameTo
	alterActionAlterColumnType
	alterActionAlterColumnSetDefault
	alterActionAlterColumnDropDefault
	alterActionAlterColumnSetNotNull
	alterActionAlterColumnDropNotNull
)

type alterAction struct {
	kind        alterActionKind
	column      columnDef
	ifNotExists bool // for ADD COLUMN IF NOT EXISTS
	ifExists    bool // for DROP COLUMN IF EXISTS, DROP CONSTRAINT IF EXISTS
	constraint  tableConstraint
	oldName     string
	newName     string
	columnName  string
	typeName    string
	defaultExp  Exp
}

// AlterTable starts building an ALTER TABLE statement.
func AlterTable(tableName Identer) AlterTableBuilder {
	return AlterTableBuilder{
		tableName: tableName,
	}
}

// AlterTableBuilder builds an ALTER TABLE statement.
type AlterTableBuilder struct {
	tableName Identer
	ifExists  bool
	actions   []alterAction
}

// IfExists adds IF EXISTS to the ALTER TABLE statement.
func (b AlterTableBuilder) IfExists() AlterTableBuilder {
	newBuilder := b
	newBuilder.ifExists = true
	return newBuilder
}

// AddColumn adds an ADD COLUMN action.
func (b AlterTableBuilder) AddColumn(name string, typeName string) AddColumnAlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 1)
	newBuilder.actions = append(newBuilder.actions, alterAction{
		kind:   alterActionAddColumn,
		column: columnDef{name: name, typeName: typeName},
	})
	return AddColumnAlterTableBuilder{AlterTableBuilder: newBuilder}
}

// AddColumnIfNotExists adds an ADD COLUMN IF NOT EXISTS action.
func (b AlterTableBuilder) AddColumnIfNotExists(name string, typeName string) AddColumnAlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 1)
	newBuilder.actions = append(newBuilder.actions, alterAction{
		kind:        alterActionAddColumn,
		column:      columnDef{name: name, typeName: typeName},
		ifNotExists: true,
	})
	return AddColumnAlterTableBuilder{AlterTableBuilder: newBuilder}
}

// DropColumn adds a DROP COLUMN action.
func (b AlterTableBuilder) DropColumn(name string) AlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 1)
	newBuilder.actions = append(newBuilder.actions, alterAction{
		kind:       alterActionDropColumn,
		columnName: name,
	})
	return newBuilder
}

// DropColumnIfExists adds a DROP COLUMN IF EXISTS action.
func (b AlterTableBuilder) DropColumnIfExists(name string) AlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 1)
	newBuilder.actions = append(newBuilder.actions, alterAction{
		kind:       alterActionDropColumn,
		columnName: name,
		ifExists:   true,
	})
	return newBuilder
}

// AddConstraint starts adding a named constraint.
func (b AlterTableBuilder) AddConstraint(name string) AddConstraintAlterTableBuilder {
	return AddConstraintAlterTableBuilder{
		AlterTableBuilder: b,
		constraintName:    name,
	}
}

// DropConstraint adds a DROP CONSTRAINT action.
func (b AlterTableBuilder) DropConstraint(name string) AlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 1)
	newBuilder.actions = append(newBuilder.actions, alterAction{
		kind: alterActionDropConstraint,
		constraint: tableConstraint{
			constraintName: name,
		},
	})
	return newBuilder
}

// DropConstraintIfExists adds a DROP CONSTRAINT IF EXISTS action.
func (b AlterTableBuilder) DropConstraintIfExists(name string) AlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 1)
	newBuilder.actions = append(newBuilder.actions, alterAction{
		kind:     alterActionDropConstraint,
		ifExists: true,
		constraint: tableConstraint{
			constraintName: name,
		},
	})
	return newBuilder
}

// RenameColumn adds a RENAME COLUMN action.
func (b AlterTableBuilder) RenameColumn(oldName, newName string) AlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 1)
	newBuilder.actions = append(newBuilder.actions, alterAction{
		kind:    alterActionRenameColumn,
		oldName: oldName,
		newName: newName,
	})
	return newBuilder
}

// RenameTo adds a RENAME TO action.
func (b AlterTableBuilder) RenameTo(newName string) AlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 1)
	newBuilder.actions = append(newBuilder.actions, alterAction{
		kind:    alterActionRenameTo,
		newName: newName,
	})
	return newBuilder
}

// AlterColumn starts an ALTER COLUMN sub-builder.
func (b AlterTableBuilder) AlterColumn(name string) AlterColumnBuilder {
	return AlterColumnBuilder{
		AlterTableBuilder: b,
		columnName:        name,
	}
}

// WriteSQL writes the ALTER TABLE statement.
func (b AlterTableBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("ALTER TABLE ")
	if b.ifExists {
		sb.WriteString("IF EXISTS ")
	}
	b.tableName.WriteSQL(sb)
	for i, action := range b.actions {
		if i > 0 {
			sb.WriteString(",")
		} else {
			sb.WriteRune(' ')
		}
		action.writeSQL(sb)
	}
}

func (a alterAction) writeSQL(sb *SQLBuilder) {
	switch a.kind {
	case alterActionAddColumn:
		sb.WriteString("ADD COLUMN ")
		if a.ifNotExists {
			sb.WriteString("IF NOT EXISTS ")
		}
		a.column.writeSQL(sb)
	case alterActionDropColumn:
		sb.WriteString("DROP COLUMN ")
		if a.ifExists {
			sb.WriteString("IF EXISTS ")
		}
		sb.WriteString(quoteIdentifierIfKeyword(a.columnName))
	case alterActionAddConstraint:
		sb.WriteString("ADD ")
		a.constraint.writeSQL(sb)
	case alterActionDropConstraint:
		sb.WriteString("DROP CONSTRAINT ")
		if a.ifExists {
			sb.WriteString("IF EXISTS ")
		}
		sb.WriteString(quoteIdentifierIfKeyword(a.constraint.constraintName))
	case alterActionRenameColumn:
		sb.WriteString("RENAME COLUMN ")
		sb.WriteString(quoteIdentifierIfKeyword(a.oldName))
		sb.WriteString(" TO ")
		sb.WriteString(quoteIdentifierIfKeyword(a.newName))
	case alterActionRenameTo:
		sb.WriteString("RENAME TO ")
		sb.WriteString(quoteIdentifierIfKeyword(a.newName))
	case alterActionAlterColumnType:
		sb.WriteString("ALTER COLUMN ")
		sb.WriteString(quoteIdentifierIfKeyword(a.columnName))
		sb.WriteString(" TYPE ")
		sb.WriteString(a.typeName)
	case alterActionAlterColumnSetDefault:
		sb.WriteString("ALTER COLUMN ")
		sb.WriteString(quoteIdentifierIfKeyword(a.columnName))
		sb.WriteString(" SET DEFAULT ")
		a.defaultExp.WriteSQL(sb)
	case alterActionAlterColumnDropDefault:
		sb.WriteString("ALTER COLUMN ")
		sb.WriteString(quoteIdentifierIfKeyword(a.columnName))
		sb.WriteString(" DROP DEFAULT")
	case alterActionAlterColumnSetNotNull:
		sb.WriteString("ALTER COLUMN ")
		sb.WriteString(quoteIdentifierIfKeyword(a.columnName))
		sb.WriteString(" SET NOT NULL")
	case alterActionAlterColumnDropNotNull:
		sb.WriteString("ALTER COLUMN ")
		sb.WriteString(quoteIdentifierIfKeyword(a.columnName))
		sb.WriteString(" DROP NOT NULL")
	}
}

// --- AddColumnAlterTableBuilder ---

// AddColumnAlterTableBuilder provides column constraint methods after AddColumn.
type AddColumnAlterTableBuilder struct {
	AlterTableBuilder
}

// NotNull adds NOT NULL to the last added column.
func (b AddColumnAlterTableBuilder) NotNull() AddColumnAlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 0)
	lastIdx := len(newBuilder.actions) - 1
	newBuilder.actions[lastIdx].column.notNull = true
	return newBuilder
}

// Default adds a DEFAULT expression to the last added column.
func (b AddColumnAlterTableBuilder) Default(exp Exp) AddColumnAlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 0)
	lastIdx := len(newBuilder.actions) - 1
	newBuilder.actions[lastIdx].column.defaultExp = exp
	return newBuilder
}

// Unique adds a UNIQUE constraint to the last added column.
func (b AddColumnAlterTableBuilder) Unique() AddColumnAlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 0)
	lastIdx := len(newBuilder.actions) - 1
	newBuilder.actions[lastIdx].column.unique = true
	return newBuilder
}

// PrimaryKey adds a PRIMARY KEY constraint to the last added column.
func (b AddColumnAlterTableBuilder) PrimaryKey() AddColumnAlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 0)
	lastIdx := len(newBuilder.actions) - 1
	newBuilder.actions[lastIdx].column.primaryKey = true
	return newBuilder
}

// Check adds a CHECK constraint to the last added column.
func (b AddColumnAlterTableBuilder) Check(exp Exp) AddColumnAlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 0)
	lastIdx := len(newBuilder.actions) - 1
	newBuilder.actions[lastIdx].column.check = exp
	return newBuilder
}

// References adds a REFERENCES constraint to the last added column.
func (b AddColumnAlterTableBuilder) References(table Identer, columns ...string) ReferencesAlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 0)
	lastIdx := len(newBuilder.actions) - 1
	newBuilder.actions[lastIdx].column.references = &columnReference{
		table:   table,
		columns: columns,
	}
	return ReferencesAlterTableBuilder{AlterTableBuilder: newBuilder.AlterTableBuilder}
}

// --- ReferencesAlterTableBuilder ---

// ReferencesAlterTableBuilder provides ON DELETE/ON UPDATE for ALTER TABLE ADD COLUMN REFERENCES.
type ReferencesAlterTableBuilder struct {
	AlterTableBuilder
}

// OnDelete adds ON DELETE action.
func (b ReferencesAlterTableBuilder) OnDelete(action string) ReferencesAlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 0)
	lastIdx := len(newBuilder.actions) - 1
	ref := *newBuilder.actions[lastIdx].column.references
	ref.onDelete = action
	newBuilder.actions[lastIdx].column.references = &ref
	return newBuilder
}

// OnUpdate adds ON UPDATE action.
func (b ReferencesAlterTableBuilder) OnUpdate(action string) ReferencesAlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 0)
	lastIdx := len(newBuilder.actions) - 1
	ref := *newBuilder.actions[lastIdx].column.references
	ref.onUpdate = action
	newBuilder.actions[lastIdx].column.references = &ref
	return newBuilder
}

// --- AddConstraintAlterTableBuilder ---

// AddConstraintAlterTableBuilder holds a constraint name and provides methods to define the constraint kind.
type AddConstraintAlterTableBuilder struct {
	AlterTableBuilder
	constraintName string
}

// PrimaryKey adds a named PRIMARY KEY constraint.
func (b AddConstraintAlterTableBuilder) PrimaryKey(columns ...string) AlterTableBuilder {
	newBuilder := b.AlterTableBuilder
	cloneSlice(&newBuilder.actions, b.actions, 1)
	newBuilder.actions = append(newBuilder.actions, alterAction{
		kind: alterActionAddConstraint,
		constraint: tableConstraint{
			constraintName: b.constraintName,
			kind:           tableConstraintPrimaryKey,
			columns:        columns,
		},
	})
	return newBuilder
}

// Unique adds a named UNIQUE constraint.
func (b AddConstraintAlterTableBuilder) Unique(columns ...string) AlterTableBuilder {
	newBuilder := b.AlterTableBuilder
	cloneSlice(&newBuilder.actions, b.actions, 1)
	newBuilder.actions = append(newBuilder.actions, alterAction{
		kind: alterActionAddConstraint,
		constraint: tableConstraint{
			constraintName: b.constraintName,
			kind:           tableConstraintUnique,
			columns:        columns,
		},
	})
	return newBuilder
}

// ForeignKey starts a named FOREIGN KEY constraint.
func (b AddConstraintAlterTableBuilder) ForeignKey(columns ...string) AddForeignKeyAlterTableBuilder {
	return AddForeignKeyAlterTableBuilder{
		AlterTableBuilder: b.AlterTableBuilder,
		constraintName:    b.constraintName,
		columns:           columns,
	}
}

// Check adds a named CHECK constraint.
func (b AddConstraintAlterTableBuilder) Check(exp Exp) AlterTableBuilder {
	newBuilder := b.AlterTableBuilder
	cloneSlice(&newBuilder.actions, b.actions, 1)
	newBuilder.actions = append(newBuilder.actions, alterAction{
		kind: alterActionAddConstraint,
		constraint: tableConstraint{
			constraintName: b.constraintName,
			kind:           tableConstraintCheck,
			checkExp:       exp,
		},
	})
	return newBuilder
}

// --- AddForeignKeyAlterTableBuilder ---

// AddForeignKeyAlterTableBuilder provides References method for ALTER TABLE ADD CONSTRAINT ... FOREIGN KEY.
type AddForeignKeyAlterTableBuilder struct {
	AlterTableBuilder
	constraintName string
	columns        []string
}

// References specifies the referenced table and columns.
func (b AddForeignKeyAlterTableBuilder) References(table Identer, columns ...string) ReferencesConstraintAlterTableBuilder {
	newBuilder := b.AlterTableBuilder
	cloneSlice(&newBuilder.actions, b.actions, 1)
	newBuilder.actions = append(newBuilder.actions, alterAction{
		kind: alterActionAddConstraint,
		constraint: tableConstraint{
			constraintName: b.constraintName,
			kind:           tableConstraintForeignKey,
			columns:        b.columns,
			refTable:       table,
			refColumns:     columns,
		},
	})
	return ReferencesConstraintAlterTableBuilder{AlterTableBuilder: newBuilder}
}

// --- ReferencesConstraintAlterTableBuilder ---

// ReferencesConstraintAlterTableBuilder provides ON DELETE/ON UPDATE for ALTER TABLE ADD CONSTRAINT ... FOREIGN KEY ... REFERENCES.
type ReferencesConstraintAlterTableBuilder struct {
	AlterTableBuilder
}

// OnDelete adds ON DELETE action.
func (b ReferencesConstraintAlterTableBuilder) OnDelete(action string) ReferencesConstraintAlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 0)
	lastIdx := len(newBuilder.actions) - 1
	newBuilder.actions[lastIdx].constraint.onDelete = action
	return newBuilder
}

// OnUpdate adds ON UPDATE action.
func (b ReferencesConstraintAlterTableBuilder) OnUpdate(action string) ReferencesConstraintAlterTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.actions, b.actions, 0)
	lastIdx := len(newBuilder.actions) - 1
	newBuilder.actions[lastIdx].constraint.onUpdate = action
	return newBuilder
}

// --- AlterColumnBuilder ---

// AlterColumnBuilder provides ALTER COLUMN sub-actions.
type AlterColumnBuilder struct {
	AlterTableBuilder
	columnName string
}

// Type adds an ALTER COLUMN ... TYPE action.
func (b AlterColumnBuilder) Type(typeName string) AlterTableBuilder {
	newBuilder := b.AlterTableBuilder
	cloneSlice(&newBuilder.actions, b.actions, 1)
	newBuilder.actions = append(newBuilder.actions, alterAction{
		kind:       alterActionAlterColumnType,
		columnName: b.columnName,
		typeName:   typeName,
	})
	return newBuilder
}

// SetDefault adds an ALTER COLUMN ... SET DEFAULT action.
func (b AlterColumnBuilder) SetDefault(exp Exp) AlterTableBuilder {
	newBuilder := b.AlterTableBuilder
	cloneSlice(&newBuilder.actions, b.actions, 1)
	newBuilder.actions = append(newBuilder.actions, alterAction{
		kind:       alterActionAlterColumnSetDefault,
		columnName: b.columnName,
		defaultExp: exp,
	})
	return newBuilder
}

// DropDefault adds an ALTER COLUMN ... DROP DEFAULT action.
func (b AlterColumnBuilder) DropDefault() AlterTableBuilder {
	newBuilder := b.AlterTableBuilder
	cloneSlice(&newBuilder.actions, b.actions, 1)
	newBuilder.actions = append(newBuilder.actions, alterAction{
		kind:       alterActionAlterColumnDropDefault,
		columnName: b.columnName,
	})
	return newBuilder
}

// SetNotNull adds an ALTER COLUMN ... SET NOT NULL action.
func (b AlterColumnBuilder) SetNotNull() AlterTableBuilder {
	newBuilder := b.AlterTableBuilder
	cloneSlice(&newBuilder.actions, b.actions, 1)
	newBuilder.actions = append(newBuilder.actions, alterAction{
		kind:       alterActionAlterColumnSetNotNull,
		columnName: b.columnName,
	})
	return newBuilder
}

// DropNotNull adds an ALTER COLUMN ... DROP NOT NULL action.
func (b AlterColumnBuilder) DropNotNull() AlterTableBuilder {
	newBuilder := b.AlterTableBuilder
	cloneSlice(&newBuilder.actions, b.actions, 1)
	newBuilder.actions = append(newBuilder.actions, alterAction{
		kind:       alterActionAlterColumnDropNotNull,
		columnName: b.columnName,
	})
	return newBuilder
}
