package builder

// CreateTable starts building a CREATE TABLE statement.
func CreateTable(tableName Identer) CreateTableBuilder {
	return CreateTableBuilder{
		tableName: tableName,
	}
}

// CreateTableBuilder builds a CREATE TABLE statement.
type CreateTableBuilder struct {
	tableName   Identer
	ifNotExists bool
	unlogged    bool
	columns     []columnDef
	constraints []tableConstraint
}

// IfNotExists adds IF NOT EXISTS to the CREATE TABLE statement.
func (b CreateTableBuilder) IfNotExists() CreateTableBuilder {
	newBuilder := b
	newBuilder.ifNotExists = true
	return newBuilder
}

// Unlogged adds UNLOGGED to the CREATE TABLE statement.
func (b CreateTableBuilder) Unlogged() CreateTableBuilder {
	newBuilder := b
	newBuilder.unlogged = true
	return newBuilder
}

// Column adds a column definition to the CREATE TABLE statement.
func (b CreateTableBuilder) Column(name string, typeName string) ColumnCreateTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.columns, b.columns, 1)
	newBuilder.columns = append(newBuilder.columns, columnDef{
		name:     name,
		typeName: typeName,
	})
	return ColumnCreateTableBuilder{CreateTableBuilder: newBuilder}
}

// Constraint adds a named table-level constraint.
func (b CreateTableBuilder) Constraint(name string) ConstraintCreateTableBuilder {
	return ConstraintCreateTableBuilder{
		CreateTableBuilder: b,
		constraintName:     name,
	}
}

// PrimaryKey adds an anonymous table-level PRIMARY KEY constraint.
func (b CreateTableBuilder) PrimaryKey(columns ...string) CreateTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.constraints, b.constraints, 1)
	newBuilder.constraints = append(newBuilder.constraints, tableConstraint{
		kind:    tableConstraintPrimaryKey,
		columns: columns,
	})
	return newBuilder
}

// Unique adds an anonymous table-level UNIQUE constraint.
func (b CreateTableBuilder) Unique(columns ...string) CreateTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.constraints, b.constraints, 1)
	newBuilder.constraints = append(newBuilder.constraints, tableConstraint{
		kind:    tableConstraintUnique,
		columns: columns,
	})
	return newBuilder
}

// ForeignKey starts an anonymous table-level FOREIGN KEY constraint.
func (b CreateTableBuilder) ForeignKey(columns ...string) ForeignKeyCreateTableBuilder {
	return ForeignKeyCreateTableBuilder{
		CreateTableBuilder: b,
		columns:            columns,
	}
}

// Check adds an anonymous table-level CHECK constraint.
func (b CreateTableBuilder) Check(exp Exp) CreateTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.constraints, b.constraints, 1)
	newBuilder.constraints = append(newBuilder.constraints, tableConstraint{
		kind:     tableConstraintCheck,
		checkExp: exp,
	})
	return newBuilder
}

// WriteSQL writes the CREATE TABLE statement.
func (b CreateTableBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("CREATE ")
	if b.unlogged {
		sb.WriteString("UNLOGGED ")
	}
	sb.WriteString("TABLE ")
	if b.ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}
	b.tableName.WriteSQL(sb)
	sb.WriteString(" (")
	idx := 0
	for _, col := range b.columns {
		if idx > 0 {
			sb.WriteString(",")
		}
		col.writeSQL(sb)
		idx++
	}
	for _, c := range b.constraints {
		if idx > 0 {
			sb.WriteString(",")
		}
		c.writeSQL(sb)
		idx++
	}
	sb.WriteRune(')')
}

// --- ColumnCreateTableBuilder ---

// ColumnCreateTableBuilder is returned after adding a column, providing column constraint methods.
type ColumnCreateTableBuilder struct {
	CreateTableBuilder
}

// NotNull adds a NOT NULL constraint to the last column.
func (b ColumnCreateTableBuilder) NotNull() ColumnCreateTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.columns, b.columns, 0)
	lastIdx := len(newBuilder.columns) - 1
	newBuilder.columns[lastIdx].notNull = true
	return newBuilder
}

// Default adds a DEFAULT expression to the last column.
func (b ColumnCreateTableBuilder) Default(exp Exp) ColumnCreateTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.columns, b.columns, 0)
	lastIdx := len(newBuilder.columns) - 1
	newBuilder.columns[lastIdx].defaultExp = exp
	return newBuilder
}

// PrimaryKey adds a PRIMARY KEY constraint to the last column.
func (b ColumnCreateTableBuilder) PrimaryKey() ColumnCreateTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.columns, b.columns, 0)
	lastIdx := len(newBuilder.columns) - 1
	newBuilder.columns[lastIdx].primaryKey = true
	return newBuilder
}

// Unique adds a UNIQUE constraint to the last column.
func (b ColumnCreateTableBuilder) Unique() ColumnCreateTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.columns, b.columns, 0)
	lastIdx := len(newBuilder.columns) - 1
	newBuilder.columns[lastIdx].unique = true
	return newBuilder
}

// Check adds a CHECK constraint to the last column.
func (b ColumnCreateTableBuilder) Check(exp Exp) ColumnCreateTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.columns, b.columns, 0)
	lastIdx := len(newBuilder.columns) - 1
	newBuilder.columns[lastIdx].check = exp
	return newBuilder
}

// References adds a REFERENCES constraint to the last column.
func (b ColumnCreateTableBuilder) References(table Identer, columns ...string) ReferencesCreateTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.columns, b.columns, 0)
	lastIdx := len(newBuilder.columns) - 1
	newBuilder.columns[lastIdx].references = &columnReference{
		table:   table,
		columns: columns,
	}
	return ReferencesCreateTableBuilder(newBuilder)
}

// --- ReferencesCreateTableBuilder ---

// ReferencesCreateTableBuilder is returned after adding REFERENCES, providing ON DELETE/ON UPDATE methods.
type ReferencesCreateTableBuilder struct {
	CreateTableBuilder
}

// OnDelete adds ON DELETE action to the last column's REFERENCES constraint.
func (b ReferencesCreateTableBuilder) OnDelete(action string) ReferencesCreateTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.columns, b.columns, 0)
	lastIdx := len(newBuilder.columns) - 1
	ref := *newBuilder.columns[lastIdx].references
	ref.onDelete = action
	newBuilder.columns[lastIdx].references = &ref
	return newBuilder
}

// OnUpdate adds ON UPDATE action to the last column's REFERENCES constraint.
func (b ReferencesCreateTableBuilder) OnUpdate(action string) ReferencesCreateTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.columns, b.columns, 0)
	lastIdx := len(newBuilder.columns) - 1
	ref := *newBuilder.columns[lastIdx].references
	ref.onUpdate = action
	newBuilder.columns[lastIdx].references = &ref
	return newBuilder
}

// --- ConstraintCreateTableBuilder ---

// ConstraintCreateTableBuilder holds a constraint name and provides methods to define the constraint kind.
type ConstraintCreateTableBuilder struct {
	CreateTableBuilder
	constraintName string
}

// PrimaryKey creates a named PRIMARY KEY constraint.
func (b ConstraintCreateTableBuilder) PrimaryKey(columns ...string) CreateTableBuilder {
	newBuilder := b.CreateTableBuilder
	cloneSlice(&newBuilder.constraints, b.constraints, 1)
	newBuilder.constraints = append(newBuilder.constraints, tableConstraint{
		constraintName: b.constraintName,
		kind:           tableConstraintPrimaryKey,
		columns:        columns,
	})
	return newBuilder
}

// Unique creates a named UNIQUE constraint.
func (b ConstraintCreateTableBuilder) Unique(columns ...string) CreateTableBuilder {
	newBuilder := b.CreateTableBuilder
	cloneSlice(&newBuilder.constraints, b.constraints, 1)
	newBuilder.constraints = append(newBuilder.constraints, tableConstraint{
		constraintName: b.constraintName,
		kind:           tableConstraintUnique,
		columns:        columns,
	})
	return newBuilder
}

// ForeignKey creates a named FOREIGN KEY constraint.
func (b ConstraintCreateTableBuilder) ForeignKey(columns ...string) ForeignKeyCreateTableBuilder {
	return ForeignKeyCreateTableBuilder{
		CreateTableBuilder: b.CreateTableBuilder,
		constraintName:     b.constraintName,
		columns:            columns,
	}
}

// Check creates a named CHECK constraint.
func (b ConstraintCreateTableBuilder) Check(exp Exp) CreateTableBuilder {
	newBuilder := b.CreateTableBuilder
	cloneSlice(&newBuilder.constraints, b.constraints, 1)
	newBuilder.constraints = append(newBuilder.constraints, tableConstraint{
		constraintName: b.constraintName,
		kind:           tableConstraintCheck,
		checkExp:       exp,
	})
	return newBuilder
}

// --- ForeignKeyCreateTableBuilder ---

// ForeignKeyCreateTableBuilder holds foreign key columns and provides References method.
type ForeignKeyCreateTableBuilder struct {
	CreateTableBuilder
	constraintName string
	columns        []string
}

// References specifies the referenced table and columns for a FOREIGN KEY constraint.
func (b ForeignKeyCreateTableBuilder) References(table Identer, columns ...string) ReferencesConstraintCreateTableBuilder {
	newBuilder := b.CreateTableBuilder
	cloneSlice(&newBuilder.constraints, b.constraints, 1)
	newBuilder.constraints = append(newBuilder.constraints, tableConstraint{
		constraintName: b.constraintName,
		kind:           tableConstraintForeignKey,
		columns:        b.columns,
		refTable:       table,
		refColumns:     columns,
	})
	return ReferencesConstraintCreateTableBuilder{CreateTableBuilder: newBuilder}
}

// --- ReferencesConstraintCreateTableBuilder ---

// ReferencesConstraintCreateTableBuilder provides ON DELETE/ON UPDATE for table-level FOREIGN KEY constraints.
type ReferencesConstraintCreateTableBuilder struct {
	CreateTableBuilder
}

// OnDelete adds ON DELETE action to the last table-level FOREIGN KEY constraint.
func (b ReferencesConstraintCreateTableBuilder) OnDelete(action string) ReferencesConstraintCreateTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.constraints, b.constraints, 0)
	lastIdx := len(newBuilder.constraints) - 1
	newBuilder.constraints[lastIdx].onDelete = action
	return newBuilder
}

// OnUpdate adds ON UPDATE action to the last table-level FOREIGN KEY constraint.
func (b ReferencesConstraintCreateTableBuilder) OnUpdate(action string) ReferencesConstraintCreateTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.constraints, b.constraints, 0)
	lastIdx := len(newBuilder.constraints) - 1
	newBuilder.constraints[lastIdx].onUpdate = action
	return newBuilder
}
