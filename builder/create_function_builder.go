package builder

// CreateFunction starts building a CREATE FUNCTION statement.
func CreateFunction(functionName Identer) CreateFunctionBuilder {
	return CreateFunctionBuilder{
		functionName: functionName,
	}
}

// CreateFunctionBuilder builds a CREATE FUNCTION statement.
type CreateFunctionBuilder struct {
	functionName Identer
	orReplace    bool
	params       []functionParam
	returns      string
	returnsTable []functionReturnColumn
	language     string
	volatility   string
	nullHandling string
	security     string
	parallel     string
	body         string
	dollarTag    string
}

type functionParam struct {
	mode       string // "", "IN", "OUT", "INOUT", "VARIADIC"
	name       string
	typeName   string
	defaultExp Exp
}

type functionReturnColumn struct {
	name     string
	typeName string
}

// OrReplace adds OR REPLACE to the CREATE FUNCTION statement.
func (b CreateFunctionBuilder) OrReplace() CreateFunctionBuilder {
	newBuilder := b
	newBuilder.orReplace = true
	return newBuilder
}

// Param adds an input parameter to the function.
func (b CreateFunctionBuilder) Param(name string, typeName string) ParamCreateFunctionBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.params, b.params, 1)
	newBuilder.params = append(newBuilder.params, functionParam{
		name:     name,
		typeName: typeName,
	})
	return ParamCreateFunctionBuilder{CreateFunctionBuilder: newBuilder}
}

// InParam adds an explicit IN parameter to the function.
func (b CreateFunctionBuilder) InParam(name string, typeName string) ParamCreateFunctionBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.params, b.params, 1)
	newBuilder.params = append(newBuilder.params, functionParam{
		mode:     "IN",
		name:     name,
		typeName: typeName,
	})
	return ParamCreateFunctionBuilder{CreateFunctionBuilder: newBuilder}
}

// OutParam adds an OUT parameter to the function.
func (b CreateFunctionBuilder) OutParam(name string, typeName string) ParamCreateFunctionBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.params, b.params, 1)
	newBuilder.params = append(newBuilder.params, functionParam{
		mode:     "OUT",
		name:     name,
		typeName: typeName,
	})
	return ParamCreateFunctionBuilder{CreateFunctionBuilder: newBuilder}
}

// InOutParam adds an INOUT parameter to the function.
func (b CreateFunctionBuilder) InOutParam(name string, typeName string) ParamCreateFunctionBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.params, b.params, 1)
	newBuilder.params = append(newBuilder.params, functionParam{
		mode:     "INOUT",
		name:     name,
		typeName: typeName,
	})
	return ParamCreateFunctionBuilder{CreateFunctionBuilder: newBuilder}
}

// VariadicParam adds a VARIADIC parameter to the function.
func (b CreateFunctionBuilder) VariadicParam(name string, typeName string) ParamCreateFunctionBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.params, b.params, 1)
	newBuilder.params = append(newBuilder.params, functionParam{
		mode:     "VARIADIC",
		name:     name,
		typeName: typeName,
	})
	return ParamCreateFunctionBuilder{CreateFunctionBuilder: newBuilder}
}

// Returns sets the return type of the function.
func (b CreateFunctionBuilder) Returns(typeName string) CreateFunctionBuilder {
	newBuilder := b
	newBuilder.returns = typeName
	return newBuilder
}

// ReturnsTable starts a RETURNS TABLE clause.
func (b CreateFunctionBuilder) ReturnsTable() ReturnsTableCreateFunctionBuilder {
	return ReturnsTableCreateFunctionBuilder{CreateFunctionBuilder: b}
}

// Language sets the language of the function.
func (b CreateFunctionBuilder) Language(lang string) CreateFunctionBuilder {
	newBuilder := b
	newBuilder.language = lang
	return newBuilder
}

// Immutable sets the function volatility to IMMUTABLE.
func (b CreateFunctionBuilder) Immutable() CreateFunctionBuilder {
	newBuilder := b
	newBuilder.volatility = "IMMUTABLE"
	return newBuilder
}

// Stable sets the function volatility to STABLE.
func (b CreateFunctionBuilder) Stable() CreateFunctionBuilder {
	newBuilder := b
	newBuilder.volatility = "STABLE"
	return newBuilder
}

// Volatile sets the function volatility to VOLATILE.
func (b CreateFunctionBuilder) Volatile() CreateFunctionBuilder {
	newBuilder := b
	newBuilder.volatility = "VOLATILE"
	return newBuilder
}

// Strict sets the function to STRICT (RETURNS NULL ON NULL INPUT).
func (b CreateFunctionBuilder) Strict() CreateFunctionBuilder {
	newBuilder := b
	newBuilder.nullHandling = "STRICT"
	return newBuilder
}

// CalledOnNullInput sets the function to CALLED ON NULL INPUT.
func (b CreateFunctionBuilder) CalledOnNullInput() CreateFunctionBuilder {
	newBuilder := b
	newBuilder.nullHandling = "CALLED ON NULL INPUT"
	return newBuilder
}

// ReturnsNullOnNullInput sets the function to RETURNS NULL ON NULL INPUT.
func (b CreateFunctionBuilder) ReturnsNullOnNullInput() CreateFunctionBuilder {
	newBuilder := b
	newBuilder.nullHandling = "RETURNS NULL ON NULL INPUT"
	return newBuilder
}

// SecurityDefiner sets the function to SECURITY DEFINER.
func (b CreateFunctionBuilder) SecurityDefiner() CreateFunctionBuilder {
	newBuilder := b
	newBuilder.security = "SECURITY DEFINER"
	return newBuilder
}

// SecurityInvoker sets the function to SECURITY INVOKER.
func (b CreateFunctionBuilder) SecurityInvoker() CreateFunctionBuilder {
	newBuilder := b
	newBuilder.security = "SECURITY INVOKER"
	return newBuilder
}

// ParallelSafe sets the function to PARALLEL SAFE.
func (b CreateFunctionBuilder) ParallelSafe() CreateFunctionBuilder {
	newBuilder := b
	newBuilder.parallel = "PARALLEL SAFE"
	return newBuilder
}

// ParallelRestricted sets the function to PARALLEL RESTRICTED.
func (b CreateFunctionBuilder) ParallelRestricted() CreateFunctionBuilder {
	newBuilder := b
	newBuilder.parallel = "PARALLEL RESTRICTED"
	return newBuilder
}

// ParallelUnsafe sets the function to PARALLEL UNSAFE.
func (b CreateFunctionBuilder) ParallelUnsafe() CreateFunctionBuilder {
	newBuilder := b
	newBuilder.parallel = "PARALLEL UNSAFE"
	return newBuilder
}

// Body sets the function body.
func (b CreateFunctionBuilder) Body(body string) CreateFunctionBuilder {
	newBuilder := b
	newBuilder.body = body
	return newBuilder
}

// DollarTag sets the dollar-quoting tag for the function body.
// An empty string (the default) produces $$...$$, "fn" produces $fn$...$fn$.
func (b CreateFunctionBuilder) DollarTag(tag string) CreateFunctionBuilder {
	newBuilder := b
	newBuilder.dollarTag = tag
	return newBuilder
}

// WriteSQL writes the CREATE FUNCTION statement.
func (b CreateFunctionBuilder) WriteSQL(sb *SQLBuilder) {
	sb.WriteString("CREATE ")
	if b.orReplace {
		sb.WriteString("OR REPLACE ")
	}
	sb.WriteString("FUNCTION ")
	b.functionName.WriteSQL(sb)
	sb.WriteRune('(')
	for i, p := range b.params {
		if i > 0 {
			sb.WriteString(",")
		}
		if p.mode != "" {
			sb.WriteString(p.mode)
			sb.WriteRune(' ')
		}
		sb.WriteString(p.name)
		sb.WriteRune(' ')
		sb.WriteString(p.typeName)
		if p.defaultExp != nil {
			sb.WriteString(" DEFAULT ")
			p.defaultExp.WriteSQL(sb)
		}
	}
	sb.WriteRune(')')
	if len(b.returnsTable) > 0 {
		sb.WriteString(" RETURNS TABLE (")
		for i, col := range b.returnsTable {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(col.name)
			sb.WriteRune(' ')
			sb.WriteString(col.typeName)
		}
		sb.WriteRune(')')
	} else if b.returns != "" {
		sb.WriteString(" RETURNS ")
		sb.WriteString(b.returns)
	}
	if b.language != "" {
		sb.WriteString(" LANGUAGE ")
		sb.WriteString(b.language)
	}
	if b.volatility != "" {
		sb.WriteRune(' ')
		sb.WriteString(b.volatility)
	}
	if b.nullHandling != "" {
		sb.WriteRune(' ')
		sb.WriteString(b.nullHandling)
	}
	if b.security != "" {
		sb.WriteRune(' ')
		sb.WriteString(b.security)
	}
	if b.parallel != "" {
		sb.WriteRune(' ')
		sb.WriteString(b.parallel)
	}
	if b.body != "" {
		sb.WriteString(" AS ")
		dollarQuote := "$" + b.dollarTag + "$"
		sb.WriteString(dollarQuote)
		sb.WriteString("\n")
		sb.WriteString(b.body)
		sb.WriteString("\n")
		sb.WriteString(dollarQuote)
	}
}

// --- ParamCreateFunctionBuilder ---

// ParamCreateFunctionBuilder is returned after adding a parameter, allowing Default to be chained.
type ParamCreateFunctionBuilder struct {
	CreateFunctionBuilder
}

// Default adds a DEFAULT expression to the last parameter.
func (b ParamCreateFunctionBuilder) Default(exp Exp) ParamCreateFunctionBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.params, b.params, 0)
	lastIdx := len(newBuilder.params) - 1
	newBuilder.params[lastIdx].defaultExp = exp
	return newBuilder
}

// --- ReturnsTableCreateFunctionBuilder ---

// ReturnsTableCreateFunctionBuilder is returned after ReturnsTable(), providing Column for defining return columns.
type ReturnsTableCreateFunctionBuilder struct {
	CreateFunctionBuilder
}

// Column adds a column to the RETURNS TABLE clause.
func (b ReturnsTableCreateFunctionBuilder) Column(name string, typeName string) ReturnsTableCreateFunctionBuilder {
	newBuilder := b
	cloneSlice(&newBuilder.returnsTable, b.returnsTable, 1)
	newBuilder.returnsTable = append(newBuilder.returnsTable, functionReturnColumn{
		name:     name,
		typeName: typeName,
	})
	return newBuilder
}
