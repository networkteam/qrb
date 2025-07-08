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

// ArrayAppend appends an element to the end of an array.
//
//	array_append ( anyarray, anyelement ) → anyarray
func ArrayAppend(arr builder.Exp, elem builder.Exp) builder.ExpBase {
	return builder.FuncExp("array_append", []builder.Exp{arr, elem})
}

// ArrayPrepend prepends an element to the beginning of an array.
//
//	array_prepend ( anyelement, anyarray ) → anyarray
func ArrayPrepend(elem builder.Exp, arr builder.Exp) builder.ExpBase {
	return builder.FuncExp("array_prepend", []builder.Exp{elem, arr})
}

// ArrayCat concatenates two arrays.
//
//	array_cat ( anyarray, anyarray ) → anyarray
func ArrayCat(arr1 builder.Exp, arr2 builder.Exp) builder.ExpBase {
	return builder.FuncExp("array_cat", []builder.Exp{arr1, arr2})
}

// ArrayDims returns a text representation of an array's dimensions.
//
//	array_dims ( anyarray ) → text
func ArrayDims(arr builder.Exp) builder.ExpBase {
	return builder.FuncExp("array_dims", []builder.Exp{arr})
}

// ArrayNdims returns the number of dimensions of an array.
//
//	array_ndims ( anyarray ) → integer
func ArrayNdims(arr builder.Exp) builder.ExpBase {
	return builder.FuncExp("array_ndims", []builder.Exp{arr})
}

// ArrayLength returns the length of the requested array dimension.
//
//	array_length ( anyarray, integer ) → integer
func ArrayLength(arr builder.Exp, dim builder.Exp) builder.ExpBase {
	return builder.FuncExp("array_length", []builder.Exp{arr, dim})
}

// ArrayLower returns the lower bound of the requested array dimension.
//
//	array_lower ( anyarray, integer ) → integer
func ArrayLower(arr builder.Exp, dim builder.Exp) builder.ExpBase {
	return builder.FuncExp("array_lower", []builder.Exp{arr, dim})
}

// ArrayUpper returns the upper bound of the requested array dimension.
//
//	array_upper ( anyarray, integer ) → integer
func ArrayUpper(arr builder.Exp, dim builder.Exp) builder.ExpBase {
	return builder.FuncExp("array_upper", []builder.Exp{arr, dim})
}

// ArrayRemove removes all occurrences of the given value from the array.
//
//	array_remove ( anyarray, anyelement ) → anyarray
func ArrayRemove(arr builder.Exp, elem builder.Exp) builder.ExpBase {
	return builder.FuncExp("array_remove", []builder.Exp{arr, elem})
}

// ArrayReplace replaces each array element equal to the second argument with the third argument.
//
//	array_replace ( anyarray, anyelement, anyelement ) → anyarray
func ArrayReplace(arr builder.Exp, from builder.Exp, to builder.Exp) builder.ExpBase {
	return builder.FuncExp("array_replace", []builder.Exp{arr, from, to})
}

// ArrayPosition returns the subscript of the first occurrence of the second argument in the array, or NULL if it's not present.
//
//	array_position ( anyarray, anyelement [, integer ] ) → integer
func ArrayPosition(arr builder.Exp, elem builder.Exp, start ...builder.Exp) builder.ExpBase {
	args := []builder.Exp{arr, elem}
	if len(start) > 0 {
		args = append(args, start[0])
	}
	return builder.FuncExp("array_position", args)
}

// ArrayPositions returns an array of the subscripts of all occurrences of the second argument in the array given as first argument.
//
//	array_positions ( anyarray, anyelement ) → integer[]
func ArrayPositions(arr builder.Exp, elem builder.Exp) builder.ExpBase {
	return builder.FuncExp("array_positions", []builder.Exp{arr, elem})
}

// ArrayToString concatenates array elements using the supplied delimiter and optional null string.
//
//	array_to_string ( anyarray, text [, text ] ) → text
func ArrayToString(arr builder.Exp, delim builder.Exp, nullString ...builder.Exp) builder.ExpBase {
	args := []builder.Exp{arr, delim}
	if len(nullString) > 0 {
		args = append(args, nullString[0])
	}
	return builder.FuncExp("array_to_string", args)
}

// StringToArray splits string into array elements using supplied delimiter and optional null string.
//
//	string_to_array ( text, text [, text ] ) → text[]
func StringToArray(text builder.Exp, delim builder.Exp, nullString ...builder.Exp) builder.ExpBase {
	args := []builder.Exp{text, delim}
	if len(nullString) > 0 {
		args = append(args, nullString[0])
	}
	return builder.FuncExp("string_to_array", args)
}

// ArrayFill returns an array filled with copies of the given value, having dimensions of the lengths specified by the second argument.
//
//	array_fill ( anyelement, integer[] [, integer[] ] ) → anyarray
func ArrayFill(value builder.Exp, dims builder.Exp, lowerBounds ...builder.Exp) builder.ExpBase {
	args := []builder.Exp{value, dims}
	if len(lowerBounds) > 0 {
		args = append(args, lowerBounds[0])
	}
	return builder.FuncExp("array_fill", args)
}
