package fn

import "github.com/networkteam/qrb/builder"

// Unnest builds an unnest functional call.
//
//	unnest ( anyarray, anyarray [, ... ] ) → setof anyelement, anyelement [, ... ]
//
// With single argument: Expands an array into a set of rows. The array's elements are read out in storage order.
// With multiple arguments: Expands multiple arrays (possibly of different data types) into a set of rows. If the arrays are not all the same length then the shorter ones are padded with NULLs. This form is only allowed in a query's FROM clause.
func Unnest(anyarray builder.Exp, anyarrays ...builder.Exp) builder.FuncBuilder {
	return builder.Func("unnest", append([]builder.Exp{anyarray}, anyarrays...)...)
}
