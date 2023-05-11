package builder

import "errors"

// CASE WHEN condition THEN result
// [WHEN ...]
// [ELSE result]
// END

func Case(exp ...Exp) CaseBuilder {
	b := CaseBuilder{}
	if len(exp) > 0 {
		b.expression = exp[0]
	}
	return b
}

func (b CaseBuilder) When(condition Exp) CaseWhenBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.conditions, b.conditions, 1)

	newBuilder.conditions = append(newBuilder.conditions, caseCondition{
		condition: condition,
	})

	return CaseWhenBuilder{newBuilder}
}

func (b CaseBuilder) Else(result Exp) CaseBuilder {
	newBuilder := b

	newBuilder.elseResult = result

	return newBuilder
}

type CaseExp struct {
	ExpBase
	expression Exp
	conditions []caseCondition
	elseResult Exp
}

func (b CaseBuilder) End() CaseExp {
	exp := CaseExp{
		expression: b.expression,
		conditions: b.conditions,
		elseResult: b.elseResult,
	}
	exp.Exp = exp // self-reference for base methods
	return exp
}

type CaseWhenBuilder struct {
	builder CaseBuilder
}

func (b CaseWhenBuilder) Then(result Exp) CaseBuilder {
	newBuilder := b.builder
	cloneSlice(&newBuilder.conditions, b.builder.conditions, 0)

	newBuilder.conditions[len(newBuilder.conditions)-1].result = result

	return newBuilder
}

type CaseBuilder struct {
	expression Exp
	conditions []caseCondition
	elseResult Exp
}

type caseCondition struct {
	condition Exp
	result    Exp
}

var ErrNoConditionsGiven = errors.New("case: no conditions given")

func (c CaseExp) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("CASE")
	if c.expression != nil {
		sb.WriteString(" ")
		c.expression.WriteSQL(sb)
	}
	if sb.Validating() && len(c.conditions) == 0 {
		sb.AddError(ErrNoConditionsGiven)
	}
	for _, condition := range c.conditions {
		sb.WriteString(" WHEN ")
		condition.condition.WriteSQL(sb)
		sb.WriteString(" THEN ")
		condition.result.WriteSQL(sb)
	}
	if c.elseResult != nil {
		sb.WriteString(" ELSE ")
		c.elseResult.WriteSQL(sb)
	}
	sb.WriteString(" END")
}

// COALESCE(value [, ...])

func Coalesce(exp Exp, rest ...Exp) ExpBase {
	return FuncExp("COALESCE", append([]Exp{exp}, rest...))
}

// NULLIF(value1, value2)

func NullIf(value1, value2 Exp) ExpBase {
	return FuncExp("NULLIF", []Exp{value1, value2})
}

// GREATEST(value [, ...])

func Greatest(exp Exp, rest ...Exp) ExpBase {
	return FuncExp("GREATEST", append([]Exp{exp}, rest...))
}

// LEAST(value [, ...])

func Least(exp Exp, rest ...Exp) ExpBase {
	return FuncExp("LEAST", append([]Exp{exp}, rest...))
}
