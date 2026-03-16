package builder

// CreateTable starts building a CREATE TABLE statement.
func CreateTable(tableName Identer) CreateTableBuilder {
	return CreateTableBuilder{
		tableName: tableName,
	}
}

// CreateTableBuilder builds a CREATE TABLE statement.
type CreateTableBuilder struct {
	tableName      Identer
	ifNotExists    bool
	unlogged       bool
	temporary      bool
	columns        []columnDef
	constraints    []tableConstraint
	likeSource     Identer
	likeOptions    []string
	partitionBy    string // "RANGE", "LIST", "HASH"
	partitionExprs []Exp
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

// Temporary adds TEMPORARY to the CREATE TABLE statement.
func (b CreateTableBuilder) Temporary() CreateTableBuilder {
	newBuilder := b
	newBuilder.temporary = true
	return newBuilder
}

// PartitionByRange adds PARTITION BY RANGE to the CREATE TABLE statement.
func (b CreateTableBuilder) PartitionByRange(exprs ...Exp) CreateTableBuilder {
	newBuilder := b
	newBuilder.partitionBy = "RANGE"
	newBuilder.partitionExprs = exprs
	return newBuilder
}

// PartitionByList adds PARTITION BY LIST to the CREATE TABLE statement.
func (b CreateTableBuilder) PartitionByList(exprs ...Exp) CreateTableBuilder {
	newBuilder := b
	newBuilder.partitionBy = "LIST"
	newBuilder.partitionExprs = exprs
	return newBuilder
}

// PartitionByHash adds PARTITION BY HASH to the CREATE TABLE statement.
func (b CreateTableBuilder) PartitionByHash(exprs ...Exp) CreateTableBuilder {
	newBuilder := b
	newBuilder.partitionBy = "HASH"
	newBuilder.partitionExprs = exprs
	return newBuilder
}

// Like adds a LIKE source_table clause to the CREATE TABLE statement.
func (b CreateTableBuilder) Like(source Identer) LikeCreateTableBuilder {
	newBuilder := b
	newBuilder.likeSource = source
	return LikeCreateTableBuilder{CreateTableBuilder: newBuilder}
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
	if b.temporary {
		sb.WriteString("TEMPORARY ")
	}
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
	if b.likeSource != nil {
		sb.WriteString("LIKE ")
		b.likeSource.WriteSQL(sb)
		for _, opt := range b.likeOptions {
			sb.WriteString(" ")
			sb.WriteString(opt)
		}
		idx++
	}
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
	if b.partitionBy != "" {
		sb.WriteString(" PARTITION BY ")
		sb.WriteString(b.partitionBy)
		sb.WriteString(" (")
		for i, expr := range b.partitionExprs {
			if i > 0 {
				sb.WriteString(",")
			}
			expr.WriteSQL(sb)
		}
		sb.WriteRune(')')
	}
}

// --- LikeCreateTableBuilder ---

// LikeCreateTableBuilder is returned after adding a LIKE clause, providing INCLUDING/EXCLUDING methods.
type LikeCreateTableBuilder struct {
	CreateTableBuilder
}

func (b LikeCreateTableBuilder) appendLikeOption(option string) LikeCreateTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.likeOptions, b.likeOptions, 1)
	newBuilder.likeOptions = append(newBuilder.likeOptions, option)
	return newBuilder
}

// IncludingAll adds INCLUDING ALL to the LIKE clause.
func (b LikeCreateTableBuilder) IncludingAll() LikeCreateTableBuilder {
	return b.appendLikeOption("INCLUDING ALL")
}

// IncludingComments adds INCLUDING COMMENTS to the LIKE clause.
func (b LikeCreateTableBuilder) IncludingComments() LikeCreateTableBuilder {
	return b.appendLikeOption("INCLUDING COMMENTS")
}

// IncludingCompression adds INCLUDING COMPRESSION to the LIKE clause.
func (b LikeCreateTableBuilder) IncludingCompression() LikeCreateTableBuilder {
	return b.appendLikeOption("INCLUDING COMPRESSION")
}

// IncludingConstraints adds INCLUDING CONSTRAINTS to the LIKE clause.
func (b LikeCreateTableBuilder) IncludingConstraints() LikeCreateTableBuilder {
	return b.appendLikeOption("INCLUDING CONSTRAINTS")
}

// IncludingDefaults adds INCLUDING DEFAULTS to the LIKE clause.
func (b LikeCreateTableBuilder) IncludingDefaults() LikeCreateTableBuilder {
	return b.appendLikeOption("INCLUDING DEFAULTS")
}

// IncludingGenerated adds INCLUDING GENERATED to the LIKE clause.
func (b LikeCreateTableBuilder) IncludingGenerated() LikeCreateTableBuilder {
	return b.appendLikeOption("INCLUDING GENERATED")
}

// IncludingIdentity adds INCLUDING IDENTITY to the LIKE clause.
func (b LikeCreateTableBuilder) IncludingIdentity() LikeCreateTableBuilder {
	return b.appendLikeOption("INCLUDING IDENTITY")
}

// IncludingIndexes adds INCLUDING INDEXES to the LIKE clause.
func (b LikeCreateTableBuilder) IncludingIndexes() LikeCreateTableBuilder {
	return b.appendLikeOption("INCLUDING INDEXES")
}

// IncludingStatistics adds INCLUDING STATISTICS to the LIKE clause.
func (b LikeCreateTableBuilder) IncludingStatistics() LikeCreateTableBuilder {
	return b.appendLikeOption("INCLUDING STATISTICS")
}

// IncludingStorage adds INCLUDING STORAGE to the LIKE clause.
func (b LikeCreateTableBuilder) IncludingStorage() LikeCreateTableBuilder {
	return b.appendLikeOption("INCLUDING STORAGE")
}

// ExcludingComments adds EXCLUDING COMMENTS to the LIKE clause.
func (b LikeCreateTableBuilder) ExcludingComments() LikeCreateTableBuilder {
	return b.appendLikeOption("EXCLUDING COMMENTS")
}

// ExcludingCompression adds EXCLUDING COMPRESSION to the LIKE clause.
func (b LikeCreateTableBuilder) ExcludingCompression() LikeCreateTableBuilder {
	return b.appendLikeOption("EXCLUDING COMPRESSION")
}

// ExcludingConstraints adds EXCLUDING CONSTRAINTS to the LIKE clause.
func (b LikeCreateTableBuilder) ExcludingConstraints() LikeCreateTableBuilder {
	return b.appendLikeOption("EXCLUDING CONSTRAINTS")
}

// ExcludingDefaults adds EXCLUDING DEFAULTS to the LIKE clause.
func (b LikeCreateTableBuilder) ExcludingDefaults() LikeCreateTableBuilder {
	return b.appendLikeOption("EXCLUDING DEFAULTS")
}

// ExcludingGenerated adds EXCLUDING GENERATED to the LIKE clause.
func (b LikeCreateTableBuilder) ExcludingGenerated() LikeCreateTableBuilder {
	return b.appendLikeOption("EXCLUDING GENERATED")
}

// ExcludingIdentity adds EXCLUDING IDENTITY to the LIKE clause.
func (b LikeCreateTableBuilder) ExcludingIdentity() LikeCreateTableBuilder {
	return b.appendLikeOption("EXCLUDING IDENTITY")
}

// ExcludingIndexes adds EXCLUDING INDEXES to the LIKE clause.
func (b LikeCreateTableBuilder) ExcludingIndexes() LikeCreateTableBuilder {
	return b.appendLikeOption("EXCLUDING INDEXES")
}

// ExcludingStatistics adds EXCLUDING STATISTICS to the LIKE clause.
func (b LikeCreateTableBuilder) ExcludingStatistics() LikeCreateTableBuilder {
	return b.appendLikeOption("EXCLUDING STATISTICS")
}

// ExcludingStorage adds EXCLUDING STORAGE to the LIKE clause.
func (b LikeCreateTableBuilder) ExcludingStorage() LikeCreateTableBuilder {
	return b.appendLikeOption("EXCLUDING STORAGE")
}

// --- ColumnCreateTableBuilder ---

// ColumnCreateTableBuilder is returned after adding a column, providing column constraint methods.
type ColumnCreateTableBuilder struct {
	CreateTableBuilder
}

func (b ColumnCreateTableBuilder) cloneLastColumn() (ColumnCreateTableBuilder, *columnDef) {
	newBuilder := b
	cloneSlice(&newBuilder.columns, b.columns, 0)
	return newBuilder, &newBuilder.columns[len(newBuilder.columns)-1]
}

// NotNull adds a NOT NULL constraint to the last column.
func (b ColumnCreateTableBuilder) NotNull() ColumnCreateTableBuilder {
	newBuilder, col := b.cloneLastColumn()
	col.notNull = true
	return newBuilder
}

// Default adds a DEFAULT expression to the last column.
func (b ColumnCreateTableBuilder) Default(exp Exp) ColumnCreateTableBuilder {
	newBuilder, col := b.cloneLastColumn()
	col.defaultExp = exp
	return newBuilder
}

// PrimaryKey adds a PRIMARY KEY constraint to the last column.
func (b ColumnCreateTableBuilder) PrimaryKey() ColumnCreateTableBuilder {
	newBuilder, col := b.cloneLastColumn()
	col.primaryKey = true
	return newBuilder
}

// Unique adds a UNIQUE constraint to the last column.
func (b ColumnCreateTableBuilder) Unique() ColumnCreateTableBuilder {
	newBuilder, col := b.cloneLastColumn()
	col.unique = true
	return newBuilder
}

// Check adds a CHECK constraint to the last column.
func (b ColumnCreateTableBuilder) Check(exp Exp) ColumnCreateTableBuilder {
	newBuilder, col := b.cloneLastColumn()
	col.check = exp
	return newBuilder
}

// GeneratedAlwaysAsIdentity adds GENERATED ALWAYS AS IDENTITY to the last column.
func (b ColumnCreateTableBuilder) GeneratedAlwaysAsIdentity() ColumnCreateTableBuilder {
	newBuilder, col := b.cloneLastColumn()
	col.generatedIdentity = "ALWAYS"
	return newBuilder
}

// GeneratedByDefaultAsIdentity adds GENERATED BY DEFAULT AS IDENTITY to the last column.
func (b ColumnCreateTableBuilder) GeneratedByDefaultAsIdentity() ColumnCreateTableBuilder {
	newBuilder, col := b.cloneLastColumn()
	col.generatedIdentity = "BY DEFAULT"
	return newBuilder
}

// GeneratedAlwaysAs adds GENERATED ALWAYS AS (expression) to the last column.
// Call Stored() or Virtual() on the returned builder to specify the storage type.
func (b ColumnCreateTableBuilder) GeneratedAlwaysAs(exp Exp) GeneratedColumnCreateTableBuilder {
	newBuilder, col := b.cloneLastColumn()
	col.generatedAs = exp
	return GeneratedColumnCreateTableBuilder{CreateTableBuilder: newBuilder.CreateTableBuilder}
}

// References adds a REFERENCES constraint to the last column.
func (b ColumnCreateTableBuilder) References(table Identer, columns ...string) ReferencesCreateTableBuilder {
	newBuilder, col := b.cloneLastColumn()
	col.references = &columnReference{
		table:   table,
		columns: columns,
	}
	return ReferencesCreateTableBuilder(newBuilder)
}

// --- GeneratedColumnCreateTableBuilder ---

// GeneratedColumnCreateTableBuilder is returned after GeneratedAlwaysAs, providing Stored/Virtual methods.
type GeneratedColumnCreateTableBuilder struct {
	CreateTableBuilder
}

// Stored sets the generated column to STORED.
func (b GeneratedColumnCreateTableBuilder) Stored() ColumnCreateTableBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.columns, b.columns, 0)
	newBuilder.columns[len(newBuilder.columns)-1].generatedStored = true
	return ColumnCreateTableBuilder{CreateTableBuilder: newBuilder.CreateTableBuilder}
}

// Virtual sets the generated column to VIRTUAL.
func (b GeneratedColumnCreateTableBuilder) Virtual() ColumnCreateTableBuilder {
	return ColumnCreateTableBuilder{CreateTableBuilder: b.CreateTableBuilder}
}

// --- ReferencesCreateTableBuilder ---

// ReferencesCreateTableBuilder is returned after adding REFERENCES, providing ON DELETE/ON UPDATE methods.
type ReferencesCreateTableBuilder struct {
	CreateTableBuilder
}

func (b ReferencesCreateTableBuilder) cloneLastRef() (ReferencesCreateTableBuilder, *columnReference) {
	newBuilder := b
	cloneSlice(&newBuilder.columns, b.columns, 0)
	lastIdx := len(newBuilder.columns) - 1
	ref := *newBuilder.columns[lastIdx].references
	newBuilder.columns[lastIdx].references = &ref
	return newBuilder, &ref
}

// OnDelete starts an ON DELETE referential action for the last column's REFERENCES constraint.
func (b ReferencesCreateTableBuilder) OnDelete() ReferentialActionBuilder[ReferencesCreateTableBuilder] {
	return ReferentialActionBuilder[ReferencesCreateTableBuilder]{
		parent: b,
		setter: func(b ReferencesCreateTableBuilder, action string) ReferencesCreateTableBuilder {
			newBuilder, ref := b.cloneLastRef()
			ref.onDelete = action
			return newBuilder
		},
	}
}

// OnUpdate starts an ON UPDATE referential action for the last column's REFERENCES constraint.
func (b ReferencesCreateTableBuilder) OnUpdate() ReferentialActionBuilder[ReferencesCreateTableBuilder] {
	return ReferentialActionBuilder[ReferencesCreateTableBuilder]{
		parent: b,
		setter: func(b ReferencesCreateTableBuilder, action string) ReferencesCreateTableBuilder {
			newBuilder, ref := b.cloneLastRef()
			ref.onUpdate = action
			return newBuilder
		},
	}
}

// Deferrable adds DEFERRABLE to the last column's REFERENCES constraint.
func (b ReferencesCreateTableBuilder) Deferrable() ReferencesCreateTableBuilder {
	newBuilder, ref := b.cloneLastRef()
	v := true
	ref.deferrable = &v
	return newBuilder
}

// NotDeferrable adds NOT DEFERRABLE to the last column's REFERENCES constraint.
func (b ReferencesCreateTableBuilder) NotDeferrable() ReferencesCreateTableBuilder {
	newBuilder, ref := b.cloneLastRef()
	v := false
	ref.deferrable = &v
	return newBuilder
}

// InitiallyDeferred adds INITIALLY DEFERRED to the last column's REFERENCES constraint.
func (b ReferencesCreateTableBuilder) InitiallyDeferred() ReferencesCreateTableBuilder {
	newBuilder, ref := b.cloneLastRef()
	ref.initiallyDeferred = true
	return newBuilder
}

// InitiallyImmediate adds INITIALLY IMMEDIATE to the last column's REFERENCES constraint.
func (b ReferencesCreateTableBuilder) InitiallyImmediate() ReferencesCreateTableBuilder {
	newBuilder, ref := b.cloneLastRef()
	ref.initiallyDeferred = false
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

func (b ReferencesConstraintCreateTableBuilder) cloneLastConstraint() (ReferencesConstraintCreateTableBuilder, *tableConstraint) {
	newBuilder := b
	cloneSlice(&newBuilder.constraints, b.constraints, 0)
	return newBuilder, &newBuilder.constraints[len(newBuilder.constraints)-1]
}

// OnDelete starts an ON DELETE referential action for the last table-level FOREIGN KEY constraint.
func (b ReferencesConstraintCreateTableBuilder) OnDelete() ReferentialActionBuilder[ReferencesConstraintCreateTableBuilder] {
	return ReferentialActionBuilder[ReferencesConstraintCreateTableBuilder]{
		parent: b,
		setter: func(b ReferencesConstraintCreateTableBuilder, action string) ReferencesConstraintCreateTableBuilder {
			newBuilder, c := b.cloneLastConstraint()
			c.onDelete = action
			return newBuilder
		},
	}
}

// OnUpdate starts an ON UPDATE referential action for the last table-level FOREIGN KEY constraint.
func (b ReferencesConstraintCreateTableBuilder) OnUpdate() ReferentialActionBuilder[ReferencesConstraintCreateTableBuilder] {
	return ReferentialActionBuilder[ReferencesConstraintCreateTableBuilder]{
		parent: b,
		setter: func(b ReferencesConstraintCreateTableBuilder, action string) ReferencesConstraintCreateTableBuilder {
			newBuilder, c := b.cloneLastConstraint()
			c.onUpdate = action
			return newBuilder
		},
	}
}

// Deferrable adds DEFERRABLE to the last table-level FOREIGN KEY constraint.
func (b ReferencesConstraintCreateTableBuilder) Deferrable() ReferencesConstraintCreateTableBuilder {
	newBuilder, c := b.cloneLastConstraint()
	v := true
	c.deferrable = &v
	return newBuilder
}

// NotDeferrable adds NOT DEFERRABLE to the last table-level FOREIGN KEY constraint.
func (b ReferencesConstraintCreateTableBuilder) NotDeferrable() ReferencesConstraintCreateTableBuilder {
	newBuilder, c := b.cloneLastConstraint()
	v := false
	c.deferrable = &v
	return newBuilder
}

// InitiallyDeferred adds INITIALLY DEFERRED to the last table-level FOREIGN KEY constraint.
func (b ReferencesConstraintCreateTableBuilder) InitiallyDeferred() ReferencesConstraintCreateTableBuilder {
	newBuilder, c := b.cloneLastConstraint()
	c.initiallyDeferred = true
	return newBuilder
}

// InitiallyImmediate adds INITIALLY IMMEDIATE to the last table-level FOREIGN KEY constraint.
func (b ReferencesConstraintCreateTableBuilder) InitiallyImmediate() ReferencesConstraintCreateTableBuilder {
	newBuilder, c := b.cloneLastConstraint()
	c.initiallyDeferred = false
	return newBuilder
}
