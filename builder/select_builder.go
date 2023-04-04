package builder

import (
	"errors"
)

type SelectBuilder struct {
	withQueries []withQuery
	// combinations holds possible previous selects and combination via UNION, INTERSECT or EXCEPT with the current select.
	combinations []selectCombination
	// parts holds the parts of the current select.
	parts selectQueryParts
}

func (b SelectBuilder) IsExp()            {}
func (b SelectBuilder) isFromExp()        {}
func (b SelectBuilder) isFromLateralExp() {}

type selectQueryParts struct {
	distinct          bool
	distinctOn        []Exp
	selectJson        *JsonBuildObjectBuilder
	selectJsonAlias   string
	selectList        []outputExp
	from              []fromItem
	whereConjunction  []Exp
	groupByDistinct   bool
	groupBys          []groupingElement
	havingConjunction []Exp
	orderBys          []orderByClause
	limit             Exp
	offset            Exp
}

type combinationType string

const (
	combinationTypeUnion     combinationType = "UNION"
	combinationTypeIntersect combinationType = "INTERSECT"
	combinationTypeExcept    combinationType = "EXCEPT"
)

type selectCombination struct {
	parts           selectQueryParts
	combinationType combinationType
	all             bool
}

// [ WITH [ RECURSIVE ] with_query [, ...] ]
// with_query: with_query_name [ ( column_name [, ...] ) ] AS [ [ NOT ] MATERIALIZED ] ( select | values | insert | update | delete )
//             [ SEARCH { BREADTH | DEPTH } FIRST BY column_name [, ...] SET search_seq_col_name ]
// TODO:       [ CYCLE column_name [, ...] SET cycle_mark_col_name [ TO cycle_mark_value DEFAULT cycle_mark_default ] USING cycle_path_col_name ]

type WithQuery interface {
	SQLWriter
	isWithQuery()
}

// isWithQuery is a marker method to ensure that multiple builder types can be used as WITH queries.
func (b SelectBuilder) isWithQuery() {}

var _ WithQuery = SelectBuilder{}

// With adds a WITH query to the select builder.
// The actual query must be supplied via WithBuilder.As.
func (b SelectBuilder) With(queryName string) WithBuilder {
	newBuilder := b
	newBuilder.withQueries = make([]withQuery, len(b.withQueries), len(b.withQueries)+1)
	copy(newBuilder.withQueries, b.withQueries)

	newBuilder.withQueries = append(newBuilder.withQueries, withQuery{
		queryName: queryName,
	})

	return WithBuilder{
		builder: newBuilder,
	}
}

// WithRecursive adds a WITH RECURSIVE query to the select builder.
// The actual query must be supplied via WithBuilder.As.
func (b SelectBuilder) WithRecursive(queryName string) WithBuilder {
	newBuilder := b
	newBuilder.withQueries = make([]withQuery, len(b.withQueries), len(b.withQueries)+1)
	copy(newBuilder.withQueries, b.withQueries)

	newBuilder.withQueries = append(newBuilder.withQueries, withQuery{
		recursive: true,
		queryName: queryName,
	})

	return WithBuilder{
		builder: newBuilder,
	}
}

// ColumnNames sets the column names for the current WITH query.
func (b WithBuilder) ColumnNames(names ...string) WithBuilder {
	newBuilder := b
	newBuilder.builder.withQueries = make([]withQuery, len(newBuilder.builder.withQueries), len(newBuilder.builder.withQueries)+1)
	copy(newBuilder.builder.withQueries, b.builder.withQueries)

	lastIdx := len(newBuilder.builder.withQueries) - 1
	newBuilder.builder.withQueries[lastIdx].columnNames = names

	return newBuilder
}

func (b WithBuilder) As(builder WithQuery) WithSelectBuilder {
	return b.asWithMaterialized(builder, nil)
}

func (b WithBuilder) AsNotMaterialized(builder WithQuery) WithSelectBuilder {
	materialized := false
	return b.asWithMaterialized(builder, &materialized)
}

func (b WithBuilder) AsMaterialized(builder WithQuery) WithSelectBuilder {
	materialized := true
	return b.asWithMaterialized(builder, &materialized)
}

func (b WithBuilder) asWithMaterialized(builder WithQuery, materialized *bool) WithSelectBuilder {
	newBuilder := WithSelectBuilder{b.builder}
	newBuilder.withQueries = make([]withQuery, len(newBuilder.withQueries), len(newBuilder.withQueries)+1)
	copy(newBuilder.withQueries, b.builder.withQueries)

	lastIdx := len(newBuilder.withQueries) - 1
	newBuilder.withQueries[lastIdx].builder = builder
	newBuilder.withQueries[lastIdx].materialized = materialized

	return newBuilder
}

const withSearchTypeDepth = "DEPTH"
const withSearchTypeBreadth = "BREADTH"

func (b WithSelectBuilder) SearchDepthFirst() WithSearchBuilder {
	return WithSearchBuilder{
		builder:    b.SelectBuilder,
		searchType: withSearchTypeDepth,
	}
}

func (b WithSelectBuilder) SearchBreadthFirst() WithSearchBuilder {
	return WithSearchBuilder{
		builder:    b.SelectBuilder,
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

func (b WithSearchByBuilder) Set(searchColumnName string) WithSelectBuilder {
	newBuilder := WithSelectBuilder{b.builder}
	newBuilder.withQueries = make([]withQuery, len(newBuilder.withQueries), len(newBuilder.withQueries)+1)
	copy(newBuilder.withQueries, b.builder.withQueries)

	lastIdx := len(newBuilder.withQueries) - 1
	newBuilder.withQueries[lastIdx].search = &withQuerySearch{
		searchType:    b.searchType,
		byColumnNames: b.byColumnNames,
		setColumnName: searchColumnName,
	}

	return newBuilder
}

// WithBuilder starts building a WITH query.
type WithBuilder struct {
	builder SelectBuilder
}

// WithSelectBuilder is a SelectBuilder and can refer to the latest WITH query.
type WithSelectBuilder struct {
	SelectBuilder
}

type WithSearchBuilder struct {
	builder    SelectBuilder
	searchType string
}

type WithSearchByBuilder struct {
	builder       SelectBuilder
	searchType    string
	byColumnNames []Exp
}

type withQuery struct {
	recursive    bool
	queryName    string
	columnNames  []string
	materialized *bool
	builder      WithQuery
	search       *withQuerySearch
}

type withQuerySearch struct {
	searchType    string
	byColumnNames []Exp
	setColumnName string
}

// SELECT [ ALL | DISTINCT [ ON ( expression [, ...] ) ] ]

type SelectSelectBuilder struct {
	SelectBuilder
}

func (b SelectSelectBuilder) Distinct() SelectDistinctBuilder {
	newBuilder := b.SelectBuilder

	newBuilder.parts.distinct = true

	return SelectDistinctBuilder{
		SelectBuilder: newBuilder,
	}
}

type SelectDistinctBuilder struct {
	SelectBuilder
}

func (b SelectDistinctBuilder) On(exp Exp, exps ...Exp) SelectBuilder {
	newBuilder := b.SelectBuilder

	newBuilder.parts.distinctOn = append([]Exp{exp}, exps...)

	return newBuilder
}

func (b SelectSelectBuilder) As(alias string) SelectSelectBuilder {
	newBuilder := b.SelectBuilder

	newBuilder.parts.selectList = make([]outputExp, len(b.parts.selectList), len(b.parts.selectList)+1)
	copy(newBuilder.parts.selectList, b.parts.selectList)

	lastIdx := len(newBuilder.parts.selectList) - 1
	newBuilder.parts.selectList[lastIdx].alias = alias

	return SelectSelectBuilder{
		SelectBuilder: newBuilder,
	}
}

type outputExp struct {
	exp   Exp
	alias string
}

// Select adds the given expressions to the select list.
func (b SelectBuilder) Select(exps ...Exp) SelectSelectBuilder {
	newBuilder := b

	newBuilder.parts.selectList = make([]outputExp, len(b.parts.selectList), len(b.parts.selectList)+len(exps))
	copy(newBuilder.parts.selectList, b.parts.selectList)

	for _, exp := range exps {
		newBuilder.parts.selectList = append(newBuilder.parts.selectList, outputExp{
			exp: exp,
		})
	}
	return SelectSelectBuilder{
		SelectBuilder: newBuilder,
	}
}

// SelectJson sets the JSON selection for this builder.
//
// Any additional selections are added after the JSON selection.
// The JSON selection can be modified by SelectBuilder.SelectJson.
func (b SelectBuilder) SelectJson(apply func(obj JsonBuildObjectBuilder) JsonBuildObjectBuilder) SelectJsonSelectBuilder {
	newBuilder := b

	var obj JsonBuildObjectBuilder
	if newBuilder.parts.selectJson != nil {
		obj = *newBuilder.parts.selectJson
	}
	obj = apply(obj)
	newBuilder.parts.selectJson = &obj

	return SelectJsonSelectBuilder{
		SelectBuilder: newBuilder,
	}
}

type SelectJsonSelectBuilder struct {
	SelectBuilder
}

func (b SelectJsonSelectBuilder) As(alias string) SelectJsonSelectBuilder {
	newBuilder := b.SelectBuilder

	newBuilder.parts.selectJsonAlias = alias

	return SelectJsonSelectBuilder{
		SelectBuilder: newBuilder,
	}
}

type FromExp interface {
	SQLWriter // We do not actually use Exp here, since this cannot appear anywhere outside the FROM clause.
	isFromExp()
}

func (b SelectBuilder) From(from FromExp) FromSelectBuilder {
	newBuilder := b

	newBuilder.parts.from = make([]fromItem, len(b.parts.from), len(b.parts.from)+1)
	copy(newBuilder.parts.from, b.parts.from)

	newBuilder.parts.from = append(newBuilder.parts.from, fromItem{
		from: from,
	})
	return FromSelectBuilder{
		SelectBuilder: newBuilder,
	}
}

type FromLateralExp interface {
	FromExp
	isFromLateralExp()
}

func (b SelectBuilder) FromLateral(from FromLateralExp) FromSelectBuilder {
	newBuilder := b

	newBuilder.parts.from = make([]fromItem, len(b.parts.from), len(b.parts.from)+1)
	copy(newBuilder.parts.from, b.parts.from)

	newBuilder.parts.from = append(newBuilder.parts.from, fromItem{
		lateral: true,
		from:    from,
	})
	return FromSelectBuilder{
		SelectBuilder: newBuilder,
	}
}

func (b SelectBuilder) FromOnly(from FromExp) FromSelectBuilder {
	newBuilder := b

	newBuilder.parts.from = make([]fromItem, len(b.parts.from), len(b.parts.from)+1)
	copy(newBuilder.parts.from, b.parts.from)

	newBuilder.parts.from = append(newBuilder.parts.from, fromItem{
		only: true,
		from: from,
	})
	return FromSelectBuilder{
		SelectBuilder: newBuilder,
	}
}

// [ ONLY ] table_name [ * ] [ [ AS ] alias [ ( column_alias [, ...] ) ] ]
// TODO [ TABLESAMPLE sampling_method ( argument [, ...] ) [ REPEATABLE ( seed ) ] ]
// [ LATERAL ] ROWS FROM( function_name ( [ argument [, ...] ] ) [ AS ( column_definition [, ...] ) ] [, ...] ) [ WITH ORDINALITY ] [ [ AS ] alias [ ( column_alias [, ...] ) ] ]

type fromItem struct {
	lateral        bool
	only           bool
	from           FromExp
	alias          string
	withOrdinality bool
	columnAliases  []string
}

var ErrFromItemLateralAndOnly = errors.New("cannot specify both LATERAL and ONLY")

func (i fromItem) WriteSQL(sb *SQLBuilder) {
	if i.lateral && i.only {
		sb.AddError(ErrFromItemLateralAndOnly)
		return
	}
	if i.only {
		sb.WriteString("ONLY ")
	}
	if i.lateral {
		sb.WriteString("LATERAL ")
	}
	i.from.WriteSQL(sb)
	if i.alias != "" {
		sb.WriteString(" AS ")
		sb.WriteString(i.alias)
	}
	if len(i.columnAliases) > 0 {
		if i.alias == "" {
			sb.WriteString(" AS")
		}
		sb.WriteString(" (")
		for i, name := range i.columnAliases {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(name)
		}
		sb.WriteString(")")
	}
}

type FromSelectBuilder struct {
	SelectBuilder
}

// As sets the alias for the last added from item.
func (b FromSelectBuilder) As(alias string) FromSelectBuilder {
	newBuilder := b
	newBuilder.parts.from = make([]fromItem, len(b.parts.from))
	copy(newBuilder.parts.from, b.parts.from)

	lastIdx := len(newBuilder.parts.from) - 1
	newBuilder.parts.from[lastIdx].alias = alias

	return newBuilder
}

// ColumnAliases sets the column aliases for the last added from item.
func (b FromSelectBuilder) ColumnAliases(aliases ...string) FromSelectBuilder {
	newBuilder := b
	newBuilder.parts.from = make([]fromItem, len(b.parts.from))
	copy(newBuilder.parts.from, b.parts.from)

	lastIdx := len(newBuilder.parts.from) - 1
	newBuilder.parts.from[lastIdx].columnAliases = aliases

	return newBuilder
}

func NewRowsFromBuilder(fns ...FuncBuilder) RowsFromBuilder {
	return RowsFromBuilder{
		fns: fns,
	}
}

type RowsFromBuilder struct {
	fns            []FuncBuilder
	withOrdinality bool
}

func (r RowsFromBuilder) isFromExp()        {}
func (r RowsFromBuilder) isFromLateralExp() {}

func (r RowsFromBuilder) WithOrdinality() RowsFromBuilder {
	newBuilder := r

	newBuilder.withOrdinality = true

	return newBuilder
}

func (r RowsFromBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("ROWS FROM (")
	for i, fn := range r.fns {
		if i > 0 {
			sb.WriteString(",")
		}
		fn.WriteSQL(sb)
	}
	sb.WriteString(")")
	if r.withOrdinality {
		sb.WriteString(" WITH ORDINALITY")
	}
}

type joinType string

const (
	joinTypeInner joinType = "JOIN"
	joinTypeLeft  joinType = "LEFT JOIN"
	joinTypeRight joinType = "RIGHT JOIN"
	joinTypeFull  joinType = "FULL JOIN"
)

type join struct {
	joinType joinType
	lateral  bool
	from     FromExp
	alias    string
	on       Exp
	using    []string
}

func (l join) WriteSQL(sb *SQLBuilder) {
	sb.WriteString(string(l.joinType))
	if l.lateral {
		sb.WriteString(" LATERAL")
	}
	sb.WriteRune(' ')
	l.from.WriteSQL(sb)
	if l.alias != "" {
		sb.WriteString(" AS ")
		sb.WriteString(l.alias)
	}
	if l.on != nil {
		sb.WriteString(" ON ")
		l.on.WriteSQL(sb)
	} else if len(l.using) > 0 {
		sb.WriteString(" USING (")
		for i, col := range l.using {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(col)
		}
		sb.WriteString(")")
	}
}

func (l join) isFromExp() {}
func (l join) IsExp()     {}

func (b SelectBuilder) Join(from FromExp) JoinSelectBuilder {
	return b.addJoin(joinTypeInner, from, false)
}

func (b SelectBuilder) JoinLateral(from FromExp) JoinSelectBuilder {
	return b.addJoin(joinTypeInner, from, true)
}

func (b SelectBuilder) LeftJoin(from FromExp) JoinSelectBuilder {
	return b.addJoin(joinTypeLeft, from, false)
}

func (b SelectBuilder) LeftJoinLateral(from FromExp) JoinSelectBuilder {
	return b.addJoin(joinTypeLeft, from, true)
}

func (b SelectBuilder) RightJoin(from FromExp) JoinSelectBuilder {
	return b.addJoin(joinTypeRight, from, false)
}

func (b SelectBuilder) FullJoin(from FromExp) JoinSelectBuilder {
	return b.addJoin(joinTypeFull, from, false)
}

// TODO NATURAL
// TODO CROSS JOIN

func (b SelectBuilder) addJoin(joinType joinType, from FromExp, lateral bool) JoinSelectBuilder {
	newBuilder := b

	newBuilder.parts.from = make([]fromItem, len(b.parts.from), len(b.parts.from)+1)
	copy(newBuilder.parts.from, b.parts.from)

	newBuilder.parts.from = append(newBuilder.parts.from, fromItem{
		from: join{
			joinType: joinType,
			lateral:  lateral,
			from:     from,
		},
	})
	return JoinSelectBuilder{
		SelectBuilder: newBuilder,
	}
}

type JoinSelectBuilder struct {
	SelectBuilder
}

func (b JoinSelectBuilder) As(alias string) JoinSelectBuilder {
	newBuilder := b
	newBuilder.parts.from = make([]fromItem, len(b.parts.from))
	copy(newBuilder.parts.from, b.parts.from)

	lastIdx := len(newBuilder.parts.from) - 1
	lastFromItem := newBuilder.parts.from[lastIdx]
	join := lastFromItem.from.(join)

	newJoin := join
	newJoin.alias = alias
	newBuilder.parts.from[lastIdx].from = newJoin

	return newBuilder
}

func (b JoinSelectBuilder) On(cond Exp, rest ...Exp) SelectBuilder {
	newBuilder := b.SelectBuilder
	newBuilder.parts.from = make([]fromItem, len(b.parts.from))
	copy(newBuilder.parts.from, b.parts.from)

	lastIdx := len(newBuilder.parts.from) - 1
	lastFromItem := newBuilder.parts.from[lastIdx]
	join := lastFromItem.from.(join)

	newJoin := join

	var on Exp
	if len(rest) == 0 {
		on = cond
	} else {
		on = And(append([]Exp{cond}, rest...)...)
	}

	newJoin.on = on
	newBuilder.parts.from[lastIdx].from = newJoin

	return newBuilder
}

func (b JoinSelectBuilder) Using(columns ...string) SelectBuilder {
	newBuilder := b.SelectBuilder
	newBuilder.parts.from = make([]fromItem, len(b.parts.from))
	copy(newBuilder.parts.from, b.parts.from)

	lastIdx := len(newBuilder.parts.from) - 1
	lastFromItem := newBuilder.parts.from[lastIdx]
	join := lastFromItem.from.(join)

	newJoin := join
	newJoin.using = columns
	newBuilder.parts.from[lastIdx].from = newJoin

	return newBuilder
}

// Where adds a WHERE condition to the query.
// Multiple calls to Where are joined with AND.
func (b SelectBuilder) Where(cond Exp) SelectBuilder {
	newBuilder := b

	newBuilder.parts.whereConjunction = make([]Exp, len(b.parts.whereConjunction), len(b.parts.whereConjunction)+1)
	copy(newBuilder.parts.whereConjunction, b.parts.whereConjunction)

	newBuilder.parts.whereConjunction = append(newBuilder.parts.whereConjunction, cond)
	return newBuilder
}

// [ GROUP BY [ ALL | DISTINCT ] grouping_element [, ...] ]
// and grouping_element can be one of:
//    ( )
//    expression
//    ( expression [, ...] )
//    ROLLUP ( { expression | ( expression [, ...] ) } [, ...] )
//    CUBE ( { expression | ( expression [, ...] ) } [, ...] )
//    GROUPING SETS ( grouping_element [, ...] )

// GroupBy adds a grouping element for the given expressions to the GROUP BY clause.
// If no expressions are given, special grouping elements can be added via GroupyBySelectBuilder.
// Use GroupyBySelectBuilder.Empty to add an empty grouping element.
func (b SelectBuilder) GroupBy(exps ...Exp) GroupyBySelectBuilder {
	if len(exps) == 0 {
		return GroupyBySelectBuilder{b}
	}

	return b.groupByAdd(groupingElement{
		sets: [][]Exp{exps},
	})
}

// Distinct adds the DISTINCT keyword to the GROUP BY clause.
func (b GroupyBySelectBuilder) Distinct() GroupyBySelectBuilder {
	newBuilder := b
	newBuilder.parts.groupByDistinct = true
	return newBuilder
}

// Empty adds an empty grouping element to the GROUP BY clause.
func (b GroupyBySelectBuilder) Empty() GroupyBySelectBuilder {
	return b.groupByAdd(groupingElement{
		sets: [][]Exp{nil},
	})
}

// Rollup adds a ROLLUP grouping element for the given expression sets to the GROUP BY clause.
func (b GroupyBySelectBuilder) Rollup(sets ...[]Exp) GroupyBySelectBuilder {
	return b.groupByAdd(groupingElement{
		groupingType: groupingTypeRollup,
		sets:         sets,
	})
}

// Cube adds a CUBE grouping element for the given expression sets to the GROUP BY clause.
func (b GroupyBySelectBuilder) Cube(sets ...[]Exp) GroupyBySelectBuilder {
	return b.groupByAdd(groupingElement{
		groupingType: groupingTypeCube,
		sets:         sets,
	})
}

// GroupingSets adds a GROUPING SETS grouping element for the given expression sets to the GROUP BY clause.
func (b GroupyBySelectBuilder) GroupingSets(sets ...[]Exp) GroupyBySelectBuilder {
	return b.groupByAdd(groupingElement{
		groupingType: groupingTypeGroupingSets,
		sets:         sets,
	})
}

type GroupyBySelectBuilder struct {
	SelectBuilder
}

func (b SelectBuilder) groupByAdd(el groupingElement) GroupyBySelectBuilder {
	newBuilder := b

	newBuilder.parts.groupBys = make([]groupingElement, len(b.parts.groupBys), len(b.parts.groupBys)+1)
	copy(newBuilder.parts.groupBys, b.parts.groupBys)

	newBuilder.parts.groupBys = append(newBuilder.parts.groupBys, el)

	return GroupyBySelectBuilder{newBuilder}
}

type groupingElement struct {
	groupingType groupingType
	sets         [][]Exp
}

type groupingType string

const (
	groupingTypeRollup       groupingType = "ROLLUP"
	groupingTypeCube         groupingType = "CUBE"
	groupingTypeGroupingSets groupingType = "GROUPING SETS"
)

func (e groupingElement) WriteSQL(sb *SQLBuilder) {
	if e.groupingType == "" {
		e.writeSet(sb, e.sets[0])
		return
	}

	sb.WriteString(string(e.groupingType))
	if len(e.sets) > 1 {
		sb.WriteString(" (")
	} else {
		sb.WriteRune(' ')
	}
	for i, set := range e.sets {
		if i > 0 {
			sb.WriteString(",")
		}
		e.writeSet(sb, set)
	}
	if len(e.sets) > 1 {
		sb.WriteRune(')')
	}
}

func (e groupingElement) writeSet(sb *SQLBuilder, exps []Exp) {
	if len(exps) == 1 {
		exps[0].WriteSQL(sb)
		return
	}

	sb.WriteString("(")
	for i, exp := range exps {
		if i > 0 {
			sb.WriteString(",")
		}
		exp.WriteSQL(sb)
	}
	sb.WriteString(")")
}

// HAVING condition

// Having adds a HAVING condition to the query.
// Multiple calls to Having are joined with AND.
func (b SelectBuilder) Having(cond Exp) SelectBuilder {
	newBuilder := b

	newBuilder.parts.havingConjunction = make([]Exp, len(b.parts.havingConjunction), len(b.parts.havingConjunction)+1)
	copy(newBuilder.parts.havingConjunction, b.parts.havingConjunction)

	newBuilder.parts.havingConjunction = append(newBuilder.parts.havingConjunction, cond)
	return newBuilder
}

// TODO: [ WINDOW window_name AS ( window_definition ) [, ...] ]

// select_statement UNION [ ALL | DISTINCT ] select_statement
// select_statement INTERSECT [ ALL | DISTINCT ] select_statement
// select_statement EXCEPT [ ALL | DISTINCT ] select_statement

func (b SelectBuilder) Union() CombinationBuilder {
	return b.addCombination(combinationTypeUnion)
}

func (b SelectBuilder) Intersect() CombinationBuilder {
	return b.addCombination(combinationTypeIntersect)
}

func (b SelectBuilder) Except() CombinationBuilder {
	return b.addCombination(combinationTypeExcept)
}

func (b SelectBuilder) addCombination(typ combinationType) CombinationBuilder {
	newBuilder := b

	newBuilder.combinations = make([]selectCombination, len(b.combinations), len(b.combinations)+1)
	copy(newBuilder.combinations, b.combinations)

	newBuilder.combinations = append(newBuilder.combinations, selectCombination{
		parts:           b.parts,
		combinationType: typ,
	})

	newBuilder.parts = selectQueryParts{}

	return CombinationBuilder{
		SelectBuilder: newBuilder,
	}
}

type CombinationBuilder struct {
	SelectBuilder
}

func (b CombinationBuilder) All() CombinationBuilder {
	newBuilder := b

	newBuilder.combinations = make([]selectCombination, len(b.combinations))
	copy(newBuilder.combinations, b.combinations)

	lastIdx := len(newBuilder.combinations) - 1
	newBuilder.combinations[lastIdx].all = true

	return newBuilder
}

// [ ORDER BY expression [ ASC | DESC | USING operator ] [ NULLS { FIRST | LAST } ] [, ...] ]

func (b SelectBuilder) OrderBy(exp Exp) OrderBySelectBuilder {
	newBuilder := b

	newBuilder.parts.orderBys = make([]orderByClause, len(b.parts.orderBys), len(b.parts.orderBys)+1)
	copy(newBuilder.parts.orderBys, b.parts.orderBys)

	newBuilder.parts.orderBys = append(newBuilder.parts.orderBys, orderByClause{
		exp: exp,
	})

	return OrderBySelectBuilder{
		SelectBuilder: newBuilder,
	}
}

// LIMIT { count | ALL }
// OFFSET start

func (b SelectBuilder) Limit(exp Exp) SelectBuilder {
	newBuilder := b
	newBuilder.parts.limit = exp
	return newBuilder
}

func (b SelectBuilder) Offset(exp Exp) SelectBuilder {
	newBuilder := b
	newBuilder.parts.offset = exp
	return newBuilder
}

type OrderBySelectBuilder struct {
	SelectBuilder
}

func (b OrderBySelectBuilder) Asc() OrderBySelectBuilder {
	return b.setOrder(sortOrderAsc)
}

func (b OrderBySelectBuilder) Desc() OrderBySelectBuilder {
	return b.setOrder(sortOrderDesc)
}

// TODO: Support ORDER BY expression USING operator

func (b OrderBySelectBuilder) setOrder(order sortOrder) OrderBySelectBuilder {
	newBuilder := b
	newBuilder.parts.orderBys = make([]orderByClause, len(b.parts.orderBys))
	copy(newBuilder.parts.orderBys, b.parts.orderBys)

	lastIdx := len(newBuilder.parts.orderBys) - 1
	newBuilder.parts.orderBys[lastIdx].order = order

	return newBuilder
}

func (b OrderBySelectBuilder) NullsFirst() OrderBySelectBuilder {
	return b.setNullsOrder(sortNullsFirst)
}

func (b OrderBySelectBuilder) NullsLast() OrderBySelectBuilder {
	return b.setNullsOrder(sortNullsLast)
}

func (b OrderBySelectBuilder) setNullsOrder(nulls sortNulls) OrderBySelectBuilder {
	newBuilder := b
	newBuilder.parts.orderBys = make([]orderByClause, len(b.parts.orderBys))
	copy(newBuilder.parts.orderBys, b.parts.orderBys)

	lastIdx := len(newBuilder.parts.orderBys) - 1
	newBuilder.parts.orderBys[lastIdx].nulls = nulls

	return newBuilder
}

// TODO: [ FOR { UPDATE | NO KEY UPDATE | SHARE | KEY SHARE } [ OF table_name [, ...] ] [ NOWAIT | SKIP LOCKED ] [...] ]

// WriteSQL writes the select as an expression.
func (b SelectBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteRune('(')
	b.innerWriteSQL(sb)
	sb.WriteRune(')')
}

// innerWriteSQL writes the select without the surrounding parentheses.
func (b SelectBuilder) innerWriteSQL(sb *SQLBuilder) {
	if len(b.withQueries) > 0 {
		sb.WriteString("WITH ")
		if b.hasRecursiveWith() {
			// from the docs: When there are multiple queries in the WITH clause, RECURSIVE should be written only once, immediately after WITH. It applies to all queries in the WITH clause, though it has no effect on queries that do not use recursion or forward references.
			sb.WriteString("RECURSIVE ")
		}
		for i, w := range b.withQueries {
			if i > 0 {
				sb.WriteString(",")
			}
			w.writeSQL(sb)
		}
		sb.WriteRune(' ')
	}

	// Write any previous select with combination via UNION, INTERSECT or EXCEPT
	for _, c := range b.combinations {
		writeSelectParts(sb, c.parts)
		sb.WriteRune(' ')
		sb.WriteString(string(c.combinationType))
		if c.all {
			sb.WriteString(" ALL")
		}
		sb.WriteRune(' ')
	}

	// Write the current select
	writeSelectParts(sb, b.parts)

	if len(b.parts.orderBys) > 0 {
		sb.WriteString(" ORDER BY ")
		for i, clause := range b.parts.orderBys {
			if i > 0 {
				sb.WriteRune(',')
			}
			clause.WriteSQL(sb)
		}
	}

	if b.parts.limit != nil {
		sb.WriteString(" LIMIT ")
		b.parts.limit.WriteSQL(sb)
	}

	if b.parts.offset != nil {
		sb.WriteString(" OFFSET ")
		b.parts.offset.WriteSQL(sb)
	}
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
	w.builder.WriteSQL(sb)
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

func writeSelectParts(sb *SQLBuilder, parts selectQueryParts) {
	sb.WriteString("SELECT ")
	if parts.distinct {
		sb.WriteString("DISTINCT ")
		if len(parts.distinctOn) > 0 {
			sb.WriteString("ON (")
			for i, exp := range parts.distinctOn {
				if i > 0 {
					sb.WriteString(",")
				}
				exp.WriteSQL(sb)
			}
			sb.WriteString(") ")
		}
	}
	if parts.selectJson != nil {
		parts.selectJson.WriteSQL(sb)
		if len(parts.selectList) > 0 {
			sb.WriteString(",")
		}
	}
	for i, exp := range parts.selectList {
		if i > 0 {
			sb.WriteString(",")
		}
		exp.exp.WriteSQL(sb)
		if exp.alias != "" {
			sb.WriteString(" AS ")
			sb.WriteString(exp.alias)
		}
	}

	if len(parts.from) > 0 {
		sb.WriteString(" FROM ")
		for i, f := range parts.from {
			if i > 0 {
				if _, isJoin := f.from.(join); !isJoin {
					sb.WriteString(",")
				} else {
					sb.WriteRune(' ')
				}
			}
			f.WriteSQL(sb)
		}
	}

	if len(parts.whereConjunction) > 0 {
		sb.WriteString(" WHERE ")
		And(parts.whereConjunction...).WriteSQL(sb)
	}

	if len(parts.groupBys) > 0 {
		sb.WriteString(" GROUP BY ")
		if parts.groupByDistinct {
			sb.WriteString("DISTINCT ")
		}
		for i, groupBy := range parts.groupBys {
			if i > 0 {
				sb.WriteString(",")
			}
			groupBy.WriteSQL(sb)
		}
	}

	if len(parts.havingConjunction) > 0 {
		sb.WriteString(" HAVING ")
		And(parts.havingConjunction...).WriteSQL(sb)
	}
}

// ApplyIf applies the given function to the builder if the condition is true.
// It returns the builder itself if the condition is false, otherwise it returns the result of the function.
// It' especially helpful for building a query conditionally.
func (b SelectBuilder) ApplyIf(cond bool, f func(q SelectBuilder) SelectBuilder) SelectBuilder {
	if cond && f != nil {
		return f(b)
	}
	return b
}

func (b SelectBuilder) hasRecursiveWith() bool {
	for _, w := range b.withQueries {
		if w.recursive {
			return true
		}
	}
	return false
}
