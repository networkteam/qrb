package fn

import (
	"github.com/networkteam/qrb/builder"
)

// --- General-Purpose Aggregate Functions

// ArrayAgg builds the array_agg aggregate function.
//
//	array_agg ( anynonarray ) → anyarray
//
// Collects all the input values, including nulls, into an array.
//
//	array_agg ( anyarray ) → anyarray
//
// Concatenates all the input arrays into an array of one higher dimension.
func ArrayAgg(exp builder.Exp) builder.AggBuilder {
	return builder.Agg("array_agg", []builder.Exp{exp})
}

// Avg builds the avg aggregate function.
//
//	avg ( T ) → T
//
// Computes the average (arithmetic mean) of all the non-null input values.
func Avg(exp builder.Exp) builder.AggBuilder {
	return builder.Agg("avg", []builder.Exp{exp})
}

// BitAnd builds the bit_and aggregate function.
//
//	bit_and ( T ) → T
//
// Computes the bitwise AND of all non-null input values.
func BitAnd(exp builder.Exp) builder.AggBuilder {
	return builder.Agg("bit_and", []builder.Exp{exp})
}

// BitOr builds the bit_or aggregate function.
//
//	bit_or ( T ) → T
//
// Computes the bitwise OR of all non-null input values.
func BitOr(exp builder.Exp) builder.AggBuilder {
	return builder.Agg("bit_or", []builder.Exp{exp})
}

// BitXor builds the bit_xor aggregate function.
//
//	bit_xor ( T ) → T
//
// Computes the bitwise exclusive OR of all non-null input values. Can be useful as a checksum for an unordered set of values.
func BitXor(exp builder.Exp) builder.AggBuilder {
	return builder.Agg("bit_xor", []builder.Exp{exp})
}

// BoolAnd builds the bool_and aggregate function.
//
//	bool_and ( boolean ) → boolean
//
// Returns true if all non-null input values are true, otherwise false.
func BoolAnd(exp builder.Exp) builder.AggBuilder {
	return builder.Agg("bool_and", []builder.Exp{exp})
}

// BoolOr builds the bool_or aggregate function.
//
//	bool_or ( boolean ) → boolean
//
// Returns true if any non-null input value is true, otherwise false.
func BoolOr(exp builder.Exp) builder.AggBuilder {
	return builder.Agg("bool_or", []builder.Exp{exp})
}

// Count builds the count aggregate function.
//
//	count ( * ) → bigint
//
// Computes the number of input rows.
//
//	count ( expression ) → bigint
//
// Computes the number of input rows in which the input value is not null.
//
// Example:
//
//	builder.Select(fn.Count(builder.N("*"))).From("table")
func Count(exp builder.Exp) builder.AggBuilder {
	return builder.Agg("count", []builder.Exp{exp})
}

// JsonAgg builds the json_agg aggregate function.
//
//	json_agg ( anyelement ) → json
//
// Collects all the input values, including nulls, into a JSON array. Values are converted to JSON as per to_json.
func JsonAgg(exp builder.Exp) builder.AggBuilder {
	return builder.Agg("json_agg", []builder.Exp{exp})
}

// JsonbAgg builds the json_agg aggregate function.
//
//	jsonb_agg ( anyelement ) → jsonb
//
// Collects all the input values, including nulls, into a JSON array. Values are converted to JSON as per to_jsonb.
func JsonbAgg(exp builder.Exp) builder.AggBuilder {
	return builder.Agg("jsonb_agg", []builder.Exp{exp})
}

// JsonObjectAgg builds the json_object_agg aggregate function.
//
//	json_object_agg ( key "any", value "any" ) → json
//
// Collects all the key/value pairs into a JSON object. Key arguments are coerced to text; value arguments are converted as per to_json. Values can be null, but not keys.
func JsonObjectAgg(key, value builder.Exp) builder.AggBuilder {
	return builder.Agg("json_object_agg", []builder.Exp{key, value})
}

// JsonbObjectAgg builds the jsonb_object_agg aggregate function.
//
//	jsonb_object_agg ( key "any", value "any" ) → jsonb
//
// Collects all the key/value pairs into a JSON object. Key arguments are coerced to text; value arguments are converted as per to_jsonb. Values can be null, but not keys.
func JsonbObjectAgg(key, value builder.Exp) builder.AggBuilder {
	return builder.Agg("jsonb_object_agg", []builder.Exp{key, value})
}

// StringAgg builds the string_agg aggregate function.
//
//	string_agg ( value text, delimiter text ) → text
//	string_agg ( value bytea, delimiter bytea ) → bytea
//
// Concatenates the non-null input values into a string. Each value after the first is preceded by the corresponding delimiter (if it's not null).
func StringAgg(value, delimiter builder.Exp) builder.AggBuilder {
	return builder.Agg("string_agg", []builder.Exp{value, delimiter})
}

// Max builds the max aggregate function.
//
//	max ( see text ) → same as input type
//
// Computes the maximum of the non-null input values. Available for any numeric, string, date/time, or enum type, as well as inet, interval, money, oid, pg_lsn, tid, xid8, and arrays of any of these types.
func Max(exp builder.Exp) builder.AggBuilder {
	return builder.Agg("max", []builder.Exp{exp})
}

// Min builds the min aggregate function.
//
//	min ( see text ) → same as input type
//
// Computes the minimum of the non-null input values. Available for any numeric, string, date/time, or enum type, as well as inet, interval, money, oid, pg_lsn, tid, xid8, and arrays of any of these types.
func Min(exp builder.Exp) builder.AggBuilder {
	return builder.Agg("min", []builder.Exp{exp})
}

// RangeAgg builds the range_agg aggregate function.
//
//	range_agg ( value anyrange ) → anymultirange
//	range_agg ( value anymultirange ) → anymultirange
//
// Computes the union of the non-null input values.
func RangeAgg(value builder.Exp) builder.AggBuilder {
	return builder.Agg("range_agg", []builder.Exp{value})
}

// RangeIntersectAgg builds the range_intersect_agg aggregate function.
//
//	range_intersect_agg ( value anyrange ) → anyrange
//	range_intersect_agg ( value anymultirange ) → anymultirange
//
// Computes the intersection of the non-null input values.
func RangeIntersectAgg(value builder.Exp) builder.AggBuilder {
	return builder.Agg("range_intersect_agg", []builder.Exp{value})
}

// Sum builds the sum aggregate function.
//
//	sum ( T ) → T
//
// Computes the sum of the non-null input values.
func Sum(exp builder.Exp) builder.AggBuilder {
	return builder.Agg("sum", []builder.Exp{exp})
}

// Xmlagg builds the xmlagg aggregate function.
//
//	xmlagg ( xml ) → xml
//
// Concatenates the non-null XML input values.
func Xmlagg(exp builder.Exp) builder.AggBuilder {
	return builder.Agg("xmlagg", []builder.Exp{exp})
}

// --- TODO Aggregate Functions for Statistics

// --- Ordered-Set Aggregate Functions

// Mode builds the mode aggregate function.
//
//	mode () WITHIN GROUP ( ORDER BY anyelement ) → anyelement
//
// Computes the mode, the most frequent value of the aggregated argument (arbitrarily choosing the first one if there are multiple equally-frequent values). The aggregated argument must be of a sortable type.
func Mode() builder.AggBuilder {
	return builder.Agg("mode", nil)
}

// PercentileCont builds the percentile_cont aggregate function.
//
//	percentile_cont ( fraction double precision ) WITHIN GROUP ( ORDER BY double precision ) → double precision
//	percentile_cont ( fraction double precision ) WITHIN GROUP ( ORDER BY interval ) → interval
//	percentile_cont ( fractions double precision[] ) WITHIN GROUP ( ORDER BY double precision ) → double precision[]
//	percentile_cont ( fractions double precision[] ) WITHIN GROUP ( ORDER BY interval ) → interval[]
//
// Computes the continuous percentile, a value corresponding to the specified fraction within the ordered set of aggregated argument values. This will interpolate between adjacent input items if needed.
func PercentileCont(fraction builder.Exp) builder.AggBuilder {
	return builder.Agg("percentile_cont", []builder.Exp{fraction})
}

// PercentileDisc builds the percentile_disc aggregate function.
//
//	percentile_disc ( fraction double precision ) WITHIN GROUP ( ORDER BY anyelement ) → anyelement
//
// Computes the discrete percentile, the first value within the ordered set of aggregated argument values whose position in the ordering equals or exceeds the specified fraction. The aggregated argument must be of a sortable type.
//
//	percentile_disc ( fractions double precision[] ) WITHIN GROUP ( ORDER BY anyelement ) → anyarray
//
// Computes multiple discrete percentiles. The result is an array of the same dimensions as the fractions parameter, with each non-null element replaced by the input value corresponding to that percentile. The aggregated argument must be of a sortable type.
func PercentileDisc(fraction builder.Exp) builder.AggBuilder {
	return builder.Agg("percentile_disc", []builder.Exp{fraction})
}

// --- Hypothetical-Set Aggregate Functions

// Rank builds the rank aggregate function.
//
//	rank ( args ) WITHIN GROUP ( ORDER BY sorted_args ) → bigint
//
// Computes the rank of the hypothetical row, with gaps; that is, the row number of the first row in its peer group.
func Rank(args ...builder.Exp) builder.AggBuilder {
	return builder.Agg("rank", args)
}

// DenseRank builds the dense_rank aggregate function.
//
//	dense_rank ( args ) WITHIN GROUP ( ORDER BY sorted_args ) → bigint
//
// Computes the rank of the hypothetical row, without gaps; this function effectively counts peer groups.
func DenseRank(args ...builder.Exp) builder.AggBuilder {
	return builder.Agg("dense_rank", args)
}

// PercentRank builds the percent_rank aggregate function.
//
//	percent_rank ( args ) WITHIN GROUP ( ORDER BY sorted_args ) → double precision
//
// Computes the relative rank of the hypothetical row, that is (rank - 1) / (total rows - 1). The value thus ranges from 0 to 1 inclusive.
func PercentRank(args ...builder.Exp) builder.AggBuilder {
	return builder.Agg("percent_rank", args)
}

// CumeDist builds the cume_dist aggregate function.
//
//	cume_dist ( args ) WITHIN GROUP ( ORDER BY sorted_args ) → double precision
//
// Computes the cumulative distribution, that is (number of rows preceding or peers with hypothetical row) / (total rows). The value thus ranges from 1/N to 1.
func CumeDist(args ...builder.Exp) builder.AggBuilder {
	return builder.Agg("cume_dist", args)
}

// --- Grouping Operations

// Grouping builds the grouping operation.
//
//	GROUPING ( group_by_expression(s) ) → integer
//
// Returns a bit mask indicating which GROUP BY expressions are not included in the current grouping set. Bits are assigned with the rightmost argument corresponding to the least-significant bit; each bit is 0 if the corresponding expression is included in the grouping criteria of the grouping set generating the current result row, and 1 if it is not included.
func Grouping(exps ...builder.Exp) builder.AggBuilder {
	return builder.Agg("GROUPING", exps)
}
