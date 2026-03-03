package ddl

import "github.com/networkteam/qrb/builder"

// CreateTable starts building a CREATE TABLE statement.
func CreateTable(tableName builder.Identer) builder.CreateTableBuilder {
	return builder.CreateTable(tableName)
}

// CreateSchema starts building a CREATE SCHEMA statement.
func CreateSchema(schemaName builder.Identer) builder.CreateSchemaBuilder {
	return builder.CreateSchema(schemaName)
}

// CreateIndex starts building a CREATE INDEX statement.
func CreateIndex(indexName string) builder.CreateIndexBuilder {
	return builder.CreateIndex(indexName)
}

// DropTable starts building a DROP TABLE statement.
func DropTable(tableName builder.Identer, rest ...builder.Identer) builder.DropTableBuilder {
	return builder.DropTable(tableName, rest...)
}

// DropSchema starts building a DROP SCHEMA statement.
func DropSchema(schemaName builder.Identer, rest ...builder.Identer) builder.DropSchemaBuilder {
	return builder.DropSchema(schemaName, rest...)
}

// AlterTable starts building an ALTER TABLE statement.
func AlterTable(tableName builder.Identer) builder.AlterTableBuilder {
	return builder.AlterTable(tableName)
}
