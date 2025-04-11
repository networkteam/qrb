package fn

import "github.com/networkteam/qrb/builder"

// RowNumber builds the row_number window function.
//
//	row_number () → bigint
//
// Returns the number of the current row within its partition, counting from 1.
func RowNumber() builder.WindowFuncBuilder {
	return builder.WindowFuncBuilder{
		FuncCall: builder.FuncExp("row_number", nil),
	}
}

// Ntile builds the ntile window function.
//
//	ntile ( num_buckets integer ) → integer
//
// Returns an integer ranging from 1 to the argument value, dividing the partition as equally as possible.
func Ntile(numBuckets builder.Exp) builder.WindowFuncBuilder {
	return builder.WindowFuncBuilder{
		FuncCall: builder.FuncExp("ntile", []builder.Exp{numBuckets}),
	}
}

// Lag builds the lag window function.
//
//	lag ( value anycompatible [, offset integer [, default anycompatible ]] ) → anycompatible
//
// Returns value evaluated at the row that is offset rows before the current row within the partition; if there is no such row, instead returns default (which must be of a type compatible with value). Both offset and default are evaluated with respect to the current row. If omitted, offset defaults to 1 and default to NULL.
func Lag(value builder.Exp, offsetAndDefault ...builder.Exp) builder.WindowFuncBuilder {
	return builder.WindowFuncBuilder{
		FuncCall: builder.FuncExp("lag", append([]builder.Exp{value}, offsetAndDefault...)),
	}
}

// Lead builds the lead window function.
//
//	lead ( value anycompatible [, offset integer [, default anycompatible ]] ) → anycompatible
//
// Returns value evaluated at the row that is offset rows after the current row within the partition; if there is no such row, instead returns default (which must be of a type compatible with value). Both offset and default are evaluated with respect to the current row. If omitted, offset defaults to 1 and default to NULL.
func Lead(value builder.Exp, offsetAndDefault ...builder.Exp) builder.WindowFuncBuilder {
	return builder.WindowFuncBuilder{
		FuncCall: builder.FuncExp("lead", append([]builder.Exp{value}, offsetAndDefault...)),
	}
}

// FirstValue builds the first_value window function.
//
//	first_value ( value anyelement ) → anyelement
//
// Returns value evaluated at the row that is the first row of the window frame.
func FirstValue(value builder.Exp) builder.WindowFuncBuilder {
	return builder.WindowFuncBuilder{
		FuncCall: builder.FuncExp("first_value", []builder.Exp{value}),
	}
}

// LastValue builds the last_value window function.
//
//	last_value ( value anyelement ) → anyelement
//
// Returns value evaluated at the row that is the last row of the window frame.
func LastValue(value builder.Exp, offsetAndDefault ...builder.Exp) builder.WindowFuncBuilder {
	return builder.WindowFuncBuilder{
		FuncCall: builder.FuncExp("last_value", []builder.Exp{value}),
	}
}

// NthValue builds the  nth_value window function.
//
//	nth_value ( value anyelement, n integer ) → anyelement
//
// Returns value evaluated at the row that is the n'th row of the window frame (counting from 1); returns NULL if there is no such row.
func NthValue(value builder.Exp, n builder.Exp) builder.WindowFuncBuilder {
	return builder.WindowFuncBuilder{
		FuncCall: builder.FuncExp("nth_value", []builder.Exp{value, n}),
	}
}
