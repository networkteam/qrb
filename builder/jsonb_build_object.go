package builder

func JsonbBuildObject() JsonbBuildObjectBuilder {
	return JsonbBuildObjectBuilder{
		props: newImmutableSliceMap[string, Exp](),
	}
}

type JsonbBuildObjectBuilder struct {
	props immutableSliceMap[string, Exp]
}

var _ Exp = JsonbBuildObjectBuilder{}

func (b JsonbBuildObjectBuilder) IsExp()       {}
func (b JsonbBuildObjectBuilder) NoParensExp() {}

func (b JsonbBuildObjectBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("jsonb_build_object(")

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

func (b JsonbBuildObjectBuilder) Prop(key string, value Exp) JsonbBuildObjectBuilder {
	newProps := b.props.Set(key, value)
	return JsonbBuildObjectBuilder{
		props: newProps,
	}
}

func (b JsonbBuildObjectBuilder) PropIf(condition bool, key string, value Exp) JsonbBuildObjectBuilder {
	if condition {
		return b.Prop(key, value)
	}
	return b
}

func (b JsonbBuildObjectBuilder) Unset(key string) JsonbBuildObjectBuilder {
	newProps := b.props.Delete(key)
	return JsonbBuildObjectBuilder{
		props: newProps,
	}
}

type JsonbBuildObjectBuilderBuilder struct {
	builder immutableSliceMap[string, Exp]
}

func (b JsonbBuildObjectBuilder) Start() *JsonbBuildObjectBuilderBuilder {
	return &JsonbBuildObjectBuilderBuilder{
		builder: b.props.clone(),
	}
}

func (bb *JsonbBuildObjectBuilderBuilder) Prop(key string, value Exp) *JsonbBuildObjectBuilderBuilder {
	bb.builder.mutatingSet(key, value)
	return bb
}

func (bb *JsonbBuildObjectBuilderBuilder) End() JsonbBuildObjectBuilder {
	return JsonbBuildObjectBuilder{
		props: bb.builder,
	}
}

func (bb *JsonbBuildObjectBuilderBuilder) PropIf(condition bool, key string, value Exp) *JsonbBuildObjectBuilderBuilder {
	if condition {
		bb.Prop(key, value)
	}
	return bb
}
