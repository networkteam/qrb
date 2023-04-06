package builder

// [ WITH [ RECURSIVE ] with_query [, ...] ]
// with_query: with_query_name [ ( column_name [, ...] ) ] AS [ [ NOT ] MATERIALIZED ] ( select | values | insert | update | delete )
//             [ SEARCH { BREADTH | DEPTH } FIRST BY column_name [, ...] SET search_seq_col_name ]
// TODO:       [ CYCLE column_name [, ...] SET cycle_mark_col_name [ TO cycle_mark_value DEFAULT cycle_mark_default ] USING cycle_path_col_name ]

// With starts a new WITH clause.
func With(queryName string) WithWithBuilder {
	return WithWithBuilder{
		builder: WithBuilder{
			withQueries: []withQuery{
				{
					queryName: queryName,
				},
			},
		},
	}
}

// WithRecursive starts a new WITH RECURSIVE clause.
func WithRecursive(queryName string) WithWithBuilder {
	return WithWithBuilder{
		builder: WithBuilder{
			withQueries: []withQuery{
				{
					recursive: true,
					queryName: queryName,
				},
			},
		},
	}
}

// With adds a WITH query to the with clause.
// The actual query must be supplied via As.
func (b WithBuilder) With(queryName string) WithWithBuilder {
	return b.startWithQuery(queryName, false)
}

// WithRecursive adds a WITH RECURSIVE query to the select builder.
// The actual query must be supplied via As.
func (b WithBuilder) WithRecursive(queryName string) WithWithBuilder {
	return b.startWithQuery(queryName, true)
}

func (b WithBuilder) startWithQuery(queryName string, recursive bool) WithWithBuilder {
	newBuilder := b
	newBuilder.withQueries = b.withQueries.cloneSlice(1)
	newBuilder.withQueries = append(newBuilder.withQueries, withQuery{
		recursive: recursive,
		queryName: queryName,
	})

	return WithWithBuilder{
		builder: newBuilder,
	}
}

// ColumnNames sets the column names for the currently started WITH query.
func (b WithWithBuilder) ColumnNames(names ...string) WithWithBuilder {
	newBuilder := b.builder
	newBuilder.withQueries = b.builder.withQueries.cloneSlice(0)

	lastIdx := len(newBuilder.withQueries) - 1
	newBuilder.withQueries[lastIdx].columnNames = names

	return WithWithBuilder{
		builder: newBuilder,
	}
}

func (b WithWithBuilder) As(query WithQuery) WithBuilder {
	return b.asWithMaterialized(query, nil)
}

func (b WithWithBuilder) AsNotMaterialized(builder WithQuery) WithBuilder {
	materialized := false
	return b.asWithMaterialized(builder, &materialized)
}

func (b WithWithBuilder) AsMaterialized(builder WithQuery) WithBuilder {
	materialized := true
	return b.asWithMaterialized(builder, &materialized)
}

func (b WithWithBuilder) asWithMaterialized(query WithQuery, materialized *bool) WithBuilder {
	newBuilder := b.builder
	newBuilder.withQueries = b.builder.withQueries.cloneSlice(0)

	lastIdx := len(newBuilder.withQueries) - 1
	newBuilder.withQueries[lastIdx].query = query
	newBuilder.withQueries[lastIdx].materialized = materialized

	return newBuilder
}

const withSearchTypeDepth = "DEPTH"
const withSearchTypeBreadth = "BREADTH"

func (b WithBuilder) SearchDepthFirst() WithSearchBuilder {
	return WithSearchBuilder{
		builder:    b,
		searchType: withSearchTypeDepth,
	}
}

func (b WithBuilder) SearchBreadthFirst() WithSearchBuilder {
	return WithSearchBuilder{
		builder:    b,
		searchType: withSearchTypeBreadth,
	}
}

func (b WithSearchBuilder) By(columnName Exp, columnNames ...Exp) WithSearchByBuilder {
	return WithSearchByBuilder{
		builder:       b.builder,
		searchType:    b.searchType,
		byColumnNames: append([]Exp{columnName}, columnNames...),
	}
}

func (b WithSearchByBuilder) Set(searchColumnName string) WithBuilder {
	newBuilder := b.builder
	newBuilder.withQueries = b.builder.withQueries.cloneSlice(0)

	lastIdx := len(newBuilder.withQueries) - 1
	newBuilder.withQueries[lastIdx].search = &withQuerySearch{
		searchType:    b.searchType,
		byColumnNames: b.byColumnNames,
		setColumnName: searchColumnName,
	}

	return newBuilder
}

// Select starts a new SelectBuilder following the with clause.
func (b WithBuilder) Select(exps ...Exp) SelectSelectBuilder {
	selectBuilder := SelectBuilder{
		withQueries: b.withQueries,
	}
	return selectBuilder.Select(exps...)
}

// InsertInto starts a new InsertBuilder following the with clause.
func (b WithBuilder) InsertInto(tableName string) InsertBuilder {
	return InsertBuilder{
		withQueries: b.withQueries,
		tableName:   tableName,
	}
}

// Update starts a new UpdateBuilder following the with clause.
func (b WithBuilder) Update(tableName string) UpdateBuilder {
	return UpdateBuilder{
		withQueries: b.withQueries,
		tableName:   tableName,
	}
}

type WithQuery interface {
	SQLWriter
	// isWithQuery is a marker method to ensure that multiple builder types can be used as WITH queries.
	isWithQuery()
}

// WithBuilder builds the WITH clause.
type WithBuilder struct {
	withQueries withQueries
}

type WithWithBuilder struct {
	builder WithBuilder
}

type WithSearchBuilder struct {
	builder    WithBuilder
	searchType string
}

type WithSearchByBuilder struct {
	builder       WithBuilder
	searchType    string
	byColumnNames []Exp
}

type withQuery struct {
	recursive    bool
	queryName    string
	columnNames  []string
	materialized *bool
	query        WithQuery
	search       *withQuerySearch
}

type withQuerySearch struct {
	searchType    string
	byColumnNames []Exp
	setColumnName string
}

type withQueries []withQuery

func (q withQueries) hasRecursiveWith() bool {
	for _, w := range q {
		if w.recursive {
			return true
		}
	}
	return false
}

func (q withQueries) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("WITH ")
	if q.hasRecursiveWith() {
		// from the docs: When there are multiple queries in the WITH clause, RECURSIVE should be written only once, immediately after WITH. It applies to all queries in the WITH clause, though it has no effect on queries that do not use recursion or forward references.
		sb.WriteString("RECURSIVE ")
	}
	for i, w := range q {
		if i > 0 {
			sb.WriteString(",")
		}
		w.writeSQL(sb)
	}
	sb.WriteRune(' ')
}

func (w withQuery) writeSQL(sb *SQLBuilder) {
	sb.WriteString(w.queryName)
	if len(w.columnNames) > 0 {
		sb.WriteRune('(')
		for i, c := range w.columnNames {
			if i > 0 {
				sb.WriteRune(',')
			}
			sb.WriteString(c)
		}
		sb.WriteRune(')')
	}
	sb.WriteString(" AS ")
	if w.materialized != nil {
		if *w.materialized == false {
			sb.WriteString("NOT ")
		}
		sb.WriteString("MATERIALIZED ")
	}
	w.query.WriteSQL(sb)
	if w.search != nil {
		sb.WriteString(" SEARCH ")
		sb.WriteString(w.search.searchType)
		sb.WriteString(" FIRST BY ")
		for i, exp := range w.search.byColumnNames {
			if i > 0 {
				sb.WriteRune(',')
			}
			exp.WriteSQL(sb)
		}
		sb.WriteString(" SET ")
		sb.WriteString(w.search.setColumnName)
	}
}

func (q withQueries) cloneSlice(additionalCapacity int) withQueries {
	newSlice := make(withQueries, len(q), len(q)+additionalCapacity)
	copy(newSlice, q)
	return newSlice
}
