package builder

// Build starts a new query builder based on the given SQLWriter.
// For executing the query, use qrbpgx.Build or qrbsql.Build which can set an executor specific to a driver.
func Build(builder SQLWriter) *QueryBuilder {
	opts := sqlBuilderOpts{
		validating: true,
	}
	return &QueryBuilder{
		builder: builder,
		opts:    opts,
	}
}

type QueryBuilder struct {
	builder   SQLWriter
	namedArgs map[string]any
	opts      sqlBuilderOpts
}

func (b *QueryBuilder) ToSQL() (sql string, args []any, err error) {
	return writeToSQLString(b.builder, b.namedArgs, b.opts)
}

func (b *QueryBuilder) WithNamedArgs(args map[string]any) *QueryBuilder {
	b.namedArgs = args
	return b
}

// WithoutValidation disables validation of the query while building.
//
// Errors might still occur when building the query, but no additional validation (like validating identifiers) will be performed.
func (b *QueryBuilder) WithoutValidation() *QueryBuilder {
	b.opts.validating = false
	return b
}
