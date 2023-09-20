package builder

func JsonBuildObject(isJsonB bool) JsonBuildObjectBuilder {
	return JsonBuildObjectBuilder{
		isJsonB: isJsonB,
		props:   newImmutableSliceMap[string, Exp](),
	}
}

type JsonBuildObjectBuilder struct {
	isJsonB bool
	props   immutableSliceMap[string, Exp]
}

var _ Exp = JsonBuildObjectBuilder{}

func (b JsonBuildObjectBuilder) IsExp() {}

func (b JsonBuildObjectBuilder) WriteSQL(sb *SQLBuilder) {
	if b.isJsonB {
		sb.WriteString("jsonb_build_object(")
	} else {
		sb.WriteString("json_build_object(")
	}

	i := 0
	for _, entry := range b.props {
		if i > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(pqQuoteLiteral(entry.k))
		sb.WriteRune(',')
		entry.v.WriteSQL(sb)

		i++
	}

	sb.WriteRune(')')
}

func (b JsonBuildObjectBuilder) Prop(key string, value Exp) JsonBuildObjectBuilder {
	newProps := b.props.Set(key, value)
	return JsonBuildObjectBuilder{
		isJsonB: b.isJsonB,
		props:   newProps,
	}
}

func (b JsonBuildObjectBuilder) PropIf(condition bool, key string, value Exp) JsonBuildObjectBuilder {
	if condition {
		return b.Prop(key, value)
	}
	return b
}

func (b JsonBuildObjectBuilder) Unset(key string) JsonBuildObjectBuilder {
	newProps := b.props.Delete(key)
	return JsonBuildObjectBuilder{
		isJsonB: b.isJsonB,
		props:   newProps,
	}
}

type JsonBuildObjectBuilderBuilder struct {
	builder immutableSliceMap[string, Exp]
}

func (b JsonBuildObjectBuilder) Start() *JsonBuildObjectBuilderBuilder {
	return &JsonBuildObjectBuilderBuilder{
		builder: b.props.clone(),
	}
}

func (bb *JsonBuildObjectBuilderBuilder) Prop(key string, value Exp) *JsonBuildObjectBuilderBuilder {
	bb.builder.mutatingSet(key, value)
	return bb
}

func (bb *JsonBuildObjectBuilderBuilder) End() JsonBuildObjectBuilder {
	return JsonBuildObjectBuilder{
		props: bb.builder,
	}
}

func (bb *JsonBuildObjectBuilderBuilder) PropIf(condition bool, key string, value Exp) *JsonBuildObjectBuilderBuilder {
	if condition {
		bb.Prop(key, value)
	}
	return bb
}
