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

	newBuilder.conditions = make([]caseCondition, len(newBuilder.conditions), len(newBuilder.conditions)+1)
	copy(newBuilder.conditions, b.conditions)

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

func (c CaseWhenBuilder) Then(result Exp) CaseBuilder {
	newBuilder := c.builder

	newBuilder.conditions = make([]caseCondition, len(newBuilder.conditions))
	copy(newBuilder.conditions, c.builder.conditions)

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

func Coalesce(exp Exp, rest ...Exp) FuncExp {
	return funcExp("COALESCE", append([]Exp{exp}, rest...))
}

// NULLIF(value1, value2)

func NullIf(value1, value2 Exp) FuncExp {
	return funcExp("NULLIF", []Exp{value1, value2})
}

// GREATEST(value [, ...])

func Greatest(exp Exp, rest ...Exp) FuncExp {
	return funcExp("GREATEST", append([]Exp{exp}, rest...))
}

// LEAST(value [, ...])

func Least(exp Exp, rest ...Exp) FuncExp {
	return funcExp("LEAST", append([]Exp{exp}, rest...))
}
