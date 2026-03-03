package builder

// CreateIndex starts building a CREATE INDEX statement.
func CreateIndex(indexName string) CreateIndexBuilder {
	return CreateIndexBuilder{
		indexName: indexName,
	}
}

// CreateIndexBuilder builds a CREATE INDEX statement.
type CreateIndexBuilder struct {
	indexName    string
	unique       bool
	concurrently bool
	ifNotExists  bool
	tableName    Identer
	using        string
	columns      []Exp
	include      []string
	where        []Exp
}

// Unique adds UNIQUE to the CREATE INDEX statement.
func (b CreateIndexBuilder) Unique() CreateIndexBuilder {
	newBuilder := b
	newBuilder.unique = true
	return newBuilder
}

// Concurrently adds CONCURRENTLY to the CREATE INDEX statement.
func (b CreateIndexBuilder) Concurrently() CreateIndexBuilder {
	newBuilder := b
	newBuilder.concurrently = true
	return newBuilder
}

// IfNotExists adds IF NOT EXISTS to the CREATE INDEX statement.
func (b CreateIndexBuilder) IfNotExists() CreateIndexBuilder {
	newBuilder := b
	newBuilder.ifNotExists = true
	return newBuilder
}

// On sets the table for the index.
func (b CreateIndexBuilder) On(tableName Identer) CreateIndexBuilder {
	newBuilder := b
	newBuilder.tableName = tableName
	return newBuilder
}

// Using sets the index method (e.g. btree, hash, gin, gist).
func (b CreateIndexBuilder) Using(method string) CreateIndexBuilder {
	newBuilder := b
	newBuilder.using = method
	return newBuilder
}

// Columns sets the indexed columns or expressions.
func (b CreateIndexBuilder) Columns(columns ...Exp) CreateIndexBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.columns, b.columns, len(columns))
	newBuilder.columns = append(newBuilder.columns, columns...)
	return newBuilder
}

// Include adds columns to the INCLUDE clause.
func (b CreateIndexBuilder) Include(columns ...string) CreateIndexBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.include, b.include, len(columns))
	newBuilder.include = append(newBuilder.include, columns...)
	return newBuilder
}

// Where adds a WHERE condition for a partial index.
// Multiple calls to Where are joined with AND.
func (b CreateIndexBuilder) Where(cond Exp) CreateIndexBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.where, b.where, 1)
	newBuilder.where = append(newBuilder.where, cond)
	return newBuilder
}

// WriteSQL writes the CREATE INDEX statement.
func (b CreateIndexBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("CREATE ")
	if b.unique {
		sb.WriteString("UNIQUE ")
	}
	sb.WriteString("INDEX ")
	if b.concurrently {
		sb.WriteString("CONCURRENTLY ")
	}
	if b.ifNotExists {
		sb.WriteString("IF NOT EXISTS ")
	}
	sb.WriteString(quoteIdentifierIfKeyword(b.indexName))
	if b.tableName != nil {
		sb.WriteString(" ON ")
		b.tableName.WriteSQL(sb)
	}
	if b.using != "" {
		sb.WriteString(" USING ")
		sb.WriteString(b.using)
	}
	if len(b.columns) > 0 {
		sb.WriteString(" (")
		for i, col := range b.columns {
			if i > 0 {
				sb.WriteString(",")
			}
			col.WriteSQL(sb)
		}
		sb.WriteRune(')')
	}
	if len(b.include) > 0 {
		sb.WriteString(" INCLUDE (")
		writeColumnList(sb, b.include)
		sb.WriteRune(')')
	}
	if len(b.where) > 0 {
		sb.WriteString(" WHERE ")
		And(b.where...).WriteSQL(sb)
	}
}
