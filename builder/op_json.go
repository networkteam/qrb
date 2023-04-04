package builder

// JSON operators

const (
	opExtract         Operator = "->"
	opExtractText     Operator = "->>"
	opExtractPath     Operator = "#>"
	opExtractPathText Operator = ">>"
	opContains        Operator = "@>"
	opContainedBy     Operator = "<@"
)

// JsonExtract builds the -> operator for json / jsonb.
//
//	json -> text → json
//	jsonb -> text → jsonb
//
// Extracts JSON object field with the given key.
//
//	json -> integer → json
//	jsonb -> integer → jsonb
//
// Extracts nth element of JSON array (array elements are indexed from zero, but negative integers count from the end).
func (b ExpBase) JsonExtract(rgt Exp) Exp {
	return b.Op(opExtract, rgt)
}

// JsonExtractText builds the ->> operator for json / jsonb.
//
//	json ->> text → text
//	jsonb ->> text → text
//
// Extracts JSON object field with the given key, as text.
//
//	json ->> integer → text
//	jsonb ->> integer → text
//
// Extracts nth element of JSON array, as text.
func (b ExpBase) JsonExtractText(rgt Exp) Exp {
	return b.Op(opExtractText, rgt)
}

// JsonExtractPath builds the #> operator for json / jsonb.
//
//	json #> text[] → json
//	jsonb #> text[] → jsonb
//
// Extracts JSON sub-object at the specified path, where path elements can be either field keys or array indexes.
//
// Example:
//
//	fn.JsonExtractPath(qrb.String(`{"a": {"b": ["foo","bar"]}}`).Cast("jsonb"), qrb.Array(qrb.String("a"), qrb.String("b")))
func (b ExpBase) JsonExtractPath(rgt Exp) Exp {
	return b.Op(opExtractPath, rgt)
}

// JsonExtractPathText builds the #>> operator for json / jsonb.
//
//	json #>> text[] → text
//	jsonb #>> text[] → text
//
// Extracts JSON sub-object at the specified path as text.
//
// Example:
//
//	fn.JsonExtractPathText(qrb.String(`{"a": {"b": ["foo","bar"]}}`).Cast("jsonb"), qrb.Array(qrb.String("a"), qrb.String("b"), qrb.String("1)))
func (b ExpBase) JsonExtractPathText(rgt Exp) Exp {
	return b.Op(opExtractPathText, rgt)
}

func (b ExpBase) JsonContains(rgt Exp) Exp {
	return b.Op(opContains, rgt)
}

func (b ExpBase) JsonContainedBy(rgt Exp) Exp {
	return b.Op(opContainedBy, rgt)
}
