package jrm

import (
	"strings"
)

type QueryBuilder struct {
	builder *SelectBuilder
}

func (b *QueryBuilder) ToSQL() (sql string, args []any) {
	sb := new(strings.Builder)
	args = b.builder.WriteSQL(sb)
	return sb.String(), args
}

func Query(builder *SelectBuilder) *QueryBuilder {
	return &QueryBuilder{
		builder: builder,
	}
}

// ---

func With(queryName string, builder *SelectBuilder) *SelectBuilder {
	selectBuilder := newSelectBuilder()
	return selectBuilder.With(queryName, builder)
}

func Select(exps ...Exp) *SelectBuilder {
	selectBuilder := newSelectBuilder()
	return selectBuilder.Select(exps...)
}

type SelectBuilder struct {
	parts selectQueryParts
}

type withQuery struct {
	queryName string
	builder   *SelectBuilder
}

type outputExp struct {
	exp Exp
	as  string
}

type selectQueryParts struct {
	with       []withQuery
	selectList []outputExp
	from       []FromItem
}

func newSelectBuilder() *SelectBuilder {
	return &SelectBuilder{}
}

func (b *SelectBuilder) With(queryName string, builder *SelectBuilder) *SelectBuilder {
	b.parts.with = append(b.parts.with, withQuery{
		queryName: queryName,
		builder:   builder,
	})
	return b
}

func (b *SelectBuilder) Select(exps ...Exp) *SelectBuilder {
	for _, exp := range exps {
		b.parts.selectList = append(b.parts.selectList, outputExp{
			exp: exp,
		})
	}
	return b
}

func (b *SelectBuilder) SelectAs(exp Exp, as string) *SelectBuilder {
	b.parts.selectList = append(b.parts.selectList, outputExp{
		exp: exp,
		as:  as,
	})
	return b
}

type Ident string

func (i Ident) isFromItem() {}
func (i Ident) isExp()      {}

func (i Ident) WriteSQL(sb *strings.Builder) (args []any) {
	sb.WriteString(string(i))
	return nil
}

type FromItem interface {
	Exp
	isFromItem()
}

func (b *SelectBuilder) From(fromItem ...FromItem) *SelectBuilder {
	b.parts.from = append(b.parts.from, fromItem...)
	return b
}

type join interface {
	isJoin()
}

type leftJoin struct {
	from          FromItem
	joinCondition Exp
}

func (l leftJoin) WriteSQL(sb *strings.Builder) (args []any) {
	sb.WriteString("LEFT JOIN ")
	args = l.from.WriteSQL(sb)
	if l.joinCondition != nil {
		sb.WriteString(" ON ")
		args = append(args, l.joinCondition.WriteSQL(sb)...)
	}
	return args
}

func (l leftJoin) isFromItem() {}
func (l leftJoin) isJoin()     {}
func (l leftJoin) isExp()      {}

func (b *SelectBuilder) LeftJoin(from FromItem, on ...Exp) *SelectBuilder {
	var joinCondition Exp
	if len(on) == 1 {
		joinCondition = on[0]
	} else if len(on) > 1 {
		joinCondition = And(on...)
	}

	b.parts.from = append(b.parts.from, leftJoin{
		from:          from,
		joinCondition: joinCondition,
	})
	return b
}

func (b *SelectBuilder) WriteSQL(sb *strings.Builder) (args []any) {
	if len(b.parts.with) > 0 {
		sb.WriteString("WITH ")
		for i, w := range b.parts.with {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(w.queryName)
			sb.WriteString(" AS (")
			args = append(args, w.builder.WriteSQL(sb)...)
			sb.WriteString(")")
		}
		sb.WriteString(" ")
	}

	sb.WriteString("SELECT ")
	for i, exp := range b.parts.selectList {
		if i > 0 {
			sb.WriteString(",")
		}
		args = append(args, exp.exp.WriteSQL(sb)...)
		if exp.as != "" {
			sb.WriteString(" AS ")
			sb.WriteString(exp.as)
		}
	}

	if len(b.parts.from) > 0 {
		sb.WriteString(" FROM ")
		for i, fromItem := range b.parts.from {
			if i > 0 {
				if _, isJoin := fromItem.(join); !isJoin {
					sb.WriteString(",")
				} else {
					sb.WriteRune(' ')
				}
			}
			args = append(args, fromItem.WriteSQL(sb)...)
		}
	}
	return args
}
