package builder

// CreateSchema starts building a CREATE SCHEMA statement.
func CreateSchema(schemaName Identer) CreateSchemaBuilder {
	return CreateSchemaBuilder{
		schemaName: schemaName,
	}
}

// CreateSchemaBuilder builds a CREATE SCHEMA statement.
type CreateSchemaBuilder struct {
	schemaName    Identer
	ifNotExists   bool
	authorization string
}

// IfNotExists adds IF NOT EXISTS to the CREATE SCHEMA statement.
func (b CreateSchemaBuilder) IfNotExists() CreateSchemaBuilder {
	newBuilder := b
	newBuilder.ifNotExists = true
	return newBuilder
}

// Authorization adds AUTHORIZATION to the CREATE SCHEMA statement.
func (b CreateSchemaBuilder) Authorization(role string) CreateSchemaBuilder {
	newBuilder := b
	newBuilder.authorization = role
	return newBuilder
}

// WriteSQL writes the CREATE SCHEMA statement.
func (b CreateSchemaBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("CREATE SCHEMA ")
	if b.ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}
	b.schemaName.WriteSQL(sb)
	if b.authorization != "" {
		sb.WriteString(" AUTHORIZATION ")
		sb.WriteString(quoteIdentifierIfKeyword(b.authorization))
	}
}
