package builder

// DropTable starts building a DROP TABLE statement.
func DropTable(tableName Identer, rest ...Identer) DropTableBuilder {
	names := make([]Identer, 0, 1+len(rest))
	names = append(names, tableName)
	names = append(names, rest...)
	return DropTableBuilder{
		tableNames: names,
	}
}

// DropTableBuilder builds a DROP TABLE statement.
type DropTableBuilder struct {
	tableNames []Identer
	ifExists   bool
	cascade    bool
	restrict   bool
}

// IfExists adds IF EXISTS to the DROP TABLE statement.
func (b DropTableBuilder) IfExists() DropTableBuilder {
	newBuilder := b
	newBuilder.ifExists = true
	return newBuilder
}

// Cascade adds CASCADE to the DROP TABLE statement.
func (b DropTableBuilder) Cascade() DropTableBuilder {
	newBuilder := b
	newBuilder.cascade = true
	newBuilder.restrict = false
	return newBuilder
}

// Restrict adds RESTRICT to the DROP TABLE statement.
func (b DropTableBuilder) Restrict() DropTableBuilder {
	newBuilder := b
	newBuilder.restrict = true
	newBuilder.cascade = false
	return newBuilder
}

// WriteSQL writes the DROP TABLE statement.
func (b DropTableBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("DROP TABLE ")
	if b.ifExists {
		sb.WriteString("IF EXISTS ")
	}
	for i, name := range b.tableNames {
		if i > 0 {
			sb.WriteString(",")
		}
		name.WriteSQL(sb)
	}
	if b.cascade {
		sb.WriteString(" CASCADE")
	}
	if b.restrict {
		sb.WriteString(" RESTRICT")
	}
}

// DropSchema starts building a DROP SCHEMA statement.
func DropSchema(schemaName Identer, rest ...Identer) DropSchemaBuilder {
	names := make([]Identer, 0, 1+len(rest))
	names = append(names, schemaName)
	names = append(names, rest...)
	return DropSchemaBuilder{
		schemaNames: names,
	}
}

// DropSchemaBuilder builds a DROP SCHEMA statement.
type DropSchemaBuilder struct {
	schemaNames []Identer
	ifExists    bool
	cascade     bool
	restrict    bool
}

// IfExists adds IF EXISTS to the DROP SCHEMA statement.
func (b DropSchemaBuilder) IfExists() DropSchemaBuilder {
	newBuilder := b
	newBuilder.ifExists = true
	return newBuilder
}

// Cascade adds CASCADE to the DROP SCHEMA statement.
func (b DropSchemaBuilder) Cascade() DropSchemaBuilder {
	newBuilder := b
	newBuilder.cascade = true
	newBuilder.restrict = false
	return newBuilder
}

// Restrict adds RESTRICT to the DROP SCHEMA statement.
func (b DropSchemaBuilder) Restrict() DropSchemaBuilder {
	newBuilder := b
	newBuilder.restrict = true
	newBuilder.cascade = false
	return newBuilder
}

// WriteSQL writes the DROP SCHEMA statement.
func (b DropSchemaBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("DROP SCHEMA ")
	if b.ifExists {
		sb.WriteString("IF EXISTS ")
	}
	for i, name := range b.schemaNames {
		if i > 0 {
			sb.WriteString(",")
		}
		name.WriteSQL(sb)
	}
	if b.cascade {
		sb.WriteString(" CASCADE")
	}
	if b.restrict {
		sb.WriteString(" RESTRICT")
	}
}
