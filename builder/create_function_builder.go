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
	returns      string
	language     string
	body         string
	dollarTag    string
}

// OrReplace adds OR REPLACE to the CREATE FUNCTION statement.
func (b CreateFunctionBuilder) OrReplace() CreateFunctionBuilder {
	newBuilder := b
	newBuilder.orReplace = true
	return newBuilder
}

// Returns sets the return type of the function.
func (b CreateFunctionBuilder) Returns(typeName string) CreateFunctionBuilder {
	newBuilder := b
	newBuilder.returns = typeName
	return newBuilder
}

// Language sets the language of the function.
func (b CreateFunctionBuilder) Language(lang string) CreateFunctionBuilder {
	newBuilder := b
	newBuilder.language = lang
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
	sb.WriteString("()")
	if b.returns != "" {
		sb.WriteString(" RETURNS ")
		sb.WriteString(b.returns)
	}
	if b.language != "" {
		sb.WriteString(" LANGUAGE ")
		sb.WriteString(b.language)
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
