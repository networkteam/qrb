package fn

import (
	"errors"

	"github.com/networkteam/qrb/builder"
)

// --- 9.47.JSON Creation Functions

/*
TODO Implement these JSON creation functions:

to_json ( anyelement ) → json

to_jsonb ( anyelement ) → jsonb

Converts any SQL value to json or jsonb. Arrays and composites are converted recursively to arrays and objects (multidimensional arrays become arrays of arrays in JSON). Otherwise, if there is a cast from the SQL data type to json, the cast function will be used to perform the conversion;[a] otherwise, a scalar JSON value is produced. For any scalar other than a number, a Boolean, or a null value, the text representation will be used, with escaping as necessary to make it a valid JSON string value.

to_json('Fred said "Hi."'::text) → "Fred said \"Hi.\""

to_jsonb(row(42, 'Fred said "Hi."'::text)) → {"f1": 42, "f2": "Fred said \"Hi.\""}

array_to_json ( anyarray [, boolean ] ) → json

Converts an SQL array to a JSON array. The behavior is the same as to_json except that line feeds will be added between top-level array elements if the optional boolean parameter is true.

array_to_json('{{1,5},{99,100}}'::int[]) → [[1,5],[99,100]]

row_to_json ( record [, boolean ] ) → json

Converts an SQL composite value to a JSON object. The behavior is the same as to_json except that line feeds will be added between top-level elements if the optional boolean parameter is true.

row_to_json(row(1,'foo')) → {"f1":1,"f2":"foo"}

json_build_array ( VARIADIC "any" ) → json

jsonb_build_array ( VARIADIC "any" ) → jsonb

Builds a possibly-heterogeneously-typed JSON array out of a variadic argument list. Each argument is converted as per to_json or to_jsonb.

json_build_array(1, 2, 'foo', 4, 5) → [1, 2, "foo", 4, 5]

*/

// JsonBuildObject builds the json_build_object function.
// It is based on a builder pattern to specify properties (see builder.JsonBuildObjectBuilder).
//
//	( VARIADIC "any" ) → json
//
// Builds a JSON object out of a variadic argument list. By convention, the argument list consists of alternating keys and values. Key arguments are coerced to text; value arguments are converted as per to_json.
func JsonBuildObject() builder.JsonBuildObjectBuilder {
	return builder.JsonBuildObject(false)
}

// JsonbBuildObject builds the json_build_object function.
// It is based on a builder pattern to specify properties (see builder.JsonBuildObjectBuilder).
//
//	( VARIADIC "any" ) → jsonb
//
// Builds a JSON object out of a variadic argument list. By convention, the argument list consists of alternating keys and values. Key arguments are coerced to text; value arguments are converted as per to_jsonb.
func JsonbBuildObject() builder.JsonBuildObjectBuilder {
	return builder.JsonBuildObject(true)
}

/*
TODO Implement these JSON creation functions:

json_object ( text[] ) → json

jsonb_object ( text[] ) → jsonb

Builds a JSON object out of a text array. The array must have either exactly one dimension with an even number of members, in which case they are taken as alternating key/value pairs, or two dimensions such that each inner array has exactly two elements, which are taken as a key/value pair. All values are converted to JSON strings.

json_object('{a, 1, b, "def", c, 3.5}') → {"a" : "1", "b" : "def", "c" : "3.5"}

json_object('{{a, 1}, {b, "def"}, {c, 3.5}}') → {"a" : "1", "b" : "def", "c" : "3.5"}

json_object ( keys text[], values text[] ) → json

jsonb_object ( keys text[], values text[] ) → jsonb

This form of json_object takes keys and values pairwise from separate text arrays. Otherwise it is identical to the one-argument form.

json_object('{a,b}', '{1,2}') → {"a": "1", "b": "2"}
*/

// --- 9.48. JSON Processing Functions

// JsonArrayElements builds the json_array_elements function.
//
//	( json ) → setof json
//
// Expands the top-level JSON array into a set of JSON values.
func JsonArrayElements(exp builder.Exp) builder.ExpBase {
	return builder.FuncExp("json_array_elements", []builder.Exp{exp})
}

// JsonbArrayElements builds the jsonb_array_elements function.
//
//	( jsonb ) → setof jsonb
//
// Expands the top-level JSON array into a set of JSON values.
func JsonbArrayElements(exp builder.Exp) builder.ExpBase {
	return builder.FuncExp(" jsonb_array_elements", []builder.Exp{exp})
}

// JsonArrayElementsText builds the json_array_elements_text function.
//
//	( json ) → setof text
//
// Expands the top-level JSON array into a set of text values.
func JsonArrayElementsText(exp builder.Exp) builder.ExpBase {
	return builder.FuncExp(" json_array_elements_text", []builder.Exp{exp})
}

// JsonbArrayElementsText builds the jsonb_array_elements_text function.
//
//	( jsonb ) → setof text
//
// Expands the top-level JSON array into a set of text values.
func JsonbArrayElementsText(exp builder.Exp) builder.ExpBase {
	return builder.FuncExp("  jsonb_array_elements_text", []builder.Exp{exp})
}

// JsonArrayLength builds the json_array_length function.
//
//	( json ) → integer
//
// Returns the number of elements in the top-level JSON array.
func JsonArrayLength(exp builder.Exp) builder.ExpBase {
	return builder.FuncExp("  json_array_length", []builder.Exp{exp})
}

// JsonbArrayLength builds the jsonb_array_length function.
//
//	( jsonb ) → integer
//
// Returns the number of elements in the top-level JSON array.
func JsonbArrayLength(exp builder.Exp) builder.ExpBase {
	return builder.FuncExp("   jsonb_array_length", []builder.Exp{exp})
}

// JsonEach builds the json_each function.
//
//	( json ) → setof record ( key text, value json )
//
// Expands the top-level JSON object into a set of key/value pairs.
func JsonEach(exp builder.Exp) builder.ExpBase {
	return builder.FuncExp("json_each", []builder.Exp{exp})
}

// JsonbEach builds the jsonb_each function.
//
//	( jsonb ) → setof record ( key text, value jsonb )
//
// Expands the top-level JSON object into a set of key/value pairs.
func JsonbEach(exp builder.Exp) builder.ExpBase {
	return builder.FuncExp("jsonb_each", []builder.Exp{exp})
}

// JsonEachText builds the json_each_text function.
//
//	( json ) → setof record ( key text, value text )
//
// Expands the top-level JSON object into a set of key/value pairs. The returned values will be of type text.
func JsonEachText(exp builder.Exp) builder.ExpBase {
	return builder.FuncExp("json_each_text", []builder.Exp{exp})
}

// JsonbEachText builds the jsonb_each_text function.
//
//	( jsonb ) → setof record ( key text, value text )
//
// Expands the top-level JSON object into a set of key/value pairs. The returned values will be of type text.
func JsonbEachText(exp builder.Exp) builder.ExpBase {
	return builder.FuncExp("jsonb_each_text", []builder.Exp{exp})
}

// JsonExtractPath builds the json_extract_path function.
//
//	( from_json json, VARIADIC path_elems text[] ) → json
//
// Extracts JSON sub-object at the specified path. (This is functionally equivalent to the #> operator, but writing the path out as a variadic list can be more convenient in some cases.)
func JsonExtractPath(fromJson builder.Exp, pathElems ...builder.Exp) builder.ExpBase {
	return builder.FuncExp("json_extract_path", append([]builder.Exp{fromJson}, pathElems...))
}

// JsonbExtractPath builds the json_extract_path function.
//
//	( from_json jsonb, VARIADIC path_elems text[] ) → jsonb
//
// Extracts JSON sub-object at the specified path. (This is functionally equivalent to the #> operator, but writing the path out as a variadic list can be more convenient in some cases.)
func JsonbExtractPath(fromJson builder.Exp, pathElems ...builder.Exp) builder.ExpBase {
	return builder.FuncExp("jsonb_extract_path", append([]builder.Exp{fromJson}, pathElems...))
}

// JsonExtractPathText builds the json_extract_path_text function.
//
//	( from_json json, VARIADIC path_elems text[] ) → text
//
// Extracts JSON sub-object at the specified path as text. (This is functionally equivalent to the #>> operator.)
func JsonExtractPathText(fromJson builder.Exp, pathElems ...builder.Exp) builder.ExpBase {
	return builder.FuncExp("json_extract_path_text", append([]builder.Exp{fromJson}, pathElems...))
}

// JsonbExtractPathText builds the jsonb_extract_path_text function.
//
//	( from_json jsonb, VARIADIC path_elems text[] ) → text
//
// Extracts JSON sub-object at the specified path as text. (This is functionally equivalent to the #>> operator.)
func JsonbExtractPathText(fromJson builder.Exp, pathElems ...builder.Exp) builder.ExpBase {
	return builder.FuncExp("jsonb_extract_path_text", append([]builder.Exp{fromJson}, pathElems...))
}

// JsonObjectKeys builds the json_object_keys function.
//
//	( json ) → setof text
//
// Returns the set of keys in the top-level JSON object.
func JsonObjectKeys(exp builder.Exp) builder.ExpBase {
	return builder.FuncExp("json_object_keys", []builder.Exp{exp})
}

// JsonbObjectKeys builds the jsonb_object_keys function.
//
//	( jsonb ) → setof text
//
// Returns the set of keys in the top-level JSON object.
func JsonbObjectKeys(exp builder.Exp) builder.ExpBase {
	return builder.FuncExp("jsonb_object_keys", []builder.Exp{exp})
}

// JsonPopulateRecord builds the json_populate_record function.
//
//	( base anyelement, from_json json ) → anyelement
//
// Expands the top-level JSON object to a row having the composite type of the base argument. The JSON object is scanned for fields whose names match column names of the output row type, and their values are inserted into those columns of the output. (Fields that do not correspond to any output column name are ignored.) In typical use, the value of base is just NULL, which means that any output columns that do not match any object field will be filled with nulls. However, if base isn't NULL then the values it contains will be used for unmatched columns.
func JsonPopulateRecord(base builder.Exp, fromJson builder.Exp) builder.ExpBase {
	return builder.FuncExp("json_populate_record", []builder.Exp{base, fromJson})
}

// JsonbPopulateRecord builds the jsonb_populate_record function.
//
//	( base anyelement, from_json jsonb ) → anyelement
//
// Expands the top-level JSON object to a row having the composite type of the base argument. The JSON object is scanned for fields whose names match column names of the output row type, and their values are inserted into those columns of the output. (Fields that do not correspond to any output column name are ignored.) In typical use, the value of base is just NULL, which means that any output columns that do not match any object field will be filled with nulls. However, if base isn't NULL then the values it contains will be used for unmatched columns.
func JsonbPopulateRecord(base builder.Exp, fromJson builder.Exp) builder.ExpBase {
	return builder.FuncExp("jsonb_populate_record", []builder.Exp{base, fromJson})
}

// JsonPopulateRecordset builds the json_populate_recordset function.
//
//	( base anyelement, from_json json ) → setof anyelement
//
// Expands the top-level JSON array of objects to a set of rows having the composite type of the base argument. Each element of the JSON array is processed as described above for json[b]_populate_record.
func JsonPopulateRecordset(base builder.Exp, fromJson builder.Exp) builder.ExpBase {
	return builder.FuncExp("json_populate_recordset", []builder.Exp{base, fromJson})
}

// JsonbPopulateRecordset builds the jsonb_populate_recordset function.
//
//	( base anyelement, from_json jsonb ) → setof anyelement
//
// Expands the top-level JSON array of objects to a set of rows having the composite type of the base argument. Each element of the JSON array is processed as described above for json[b]_populate_record.
func JsonbPopulateRecordset(base builder.Exp, fromJson builder.Exp) builder.ExpBase {
	return builder.FuncExp("jsonb_populate_recordset", []builder.Exp{base, fromJson})
}

// JsonToRecord builds the json_to_record function.
//
//	( json ) → record
//
// Expands the top-level JSON object to a row having the composite type defined by an AS clause. (As with all functions returning record, the calling query must explicitly define the structure of the record with an AS clause.) The output record is filled from fields of the JSON object, in the same way as described above for json[b]_populate_record. Since there is no input record value, unmatched columns are always filled with nulls.
func JsonToRecord(exp builder.Exp) builder.FuncBuilder {
	return builder.Func("json_to_record", exp)
}

// JsonbToRecord builds the jsonb_to_record function.
//
//	( jsonb ) → record
//
// Expands the top-level JSON object to a row having the composite type defined by an AS clause. (As with all functions returning record, the calling query must explicitly define the structure of the record with an AS clause.) The output record is filled from fields of the JSON object, in the same way as described above for json[b]_populate_record. Since there is no input record value, unmatched columns are always filled with nulls.
func JsonbToRecord(exp builder.Exp) builder.FuncBuilder {
	return builder.Func("jsonb_to_record", exp)
}

// JsonToRecordset builds the json_to_recordset function.
//
//	( json ) → setof record
//
// Expands the top-level JSON array of objects to a set of rows having the composite type defined by an AS clause. (As with all functions returning record, the calling query must explicitly define the structure of the record with an AS clause.) Each element of the JSON array is processed as described above for json[b]_populate_record.
func JsonToRecordset(exp builder.Exp) builder.FuncBuilder {
	return builder.Func("json_to_recordset", exp)
}

// JsonbToRecordset builds the jsonb_to_recordset function.
//
//	( jsonb ) → setof record
//
// Expands the top-level JSON array of objects to a set of rows having the composite type defined by an AS clause. (As with all functions returning record, the calling query must explicitly define the structure of the record with an AS clause.) Each element of the JSON array is processed as described above for json[b]_populate_record.
func JsonbToRecordset(exp builder.Exp) builder.FuncBuilder {
	return builder.Func("jsonb_to_recordset", exp)
}

// JsonbSet builds the jsonb_set function.
//
//	( target jsonb, path text[], new_value jsonb [, create_if_missing boolean ] ) → jsonb
//
// Returns target with the item designated by path replaced by new_value, or with new_value added if create_if_missing is true (which is the default) and the item designated by path does not exist. All earlier steps in the path must exist, or the target is returned unchanged. As with the path oriented operators, negative integers that appear in the path count from the end of JSON arrays. If the last path step is an array index that is out of range, and create_if_missing is true, the new value is added at the beginning of the array if the index is negative, or at the end of the array if it is positive.
func JsonbSet(target builder.Exp, path builder.Exp, newValue builder.Exp, createIfMissing ...builder.Exp) builder.ExpBase {
	args := []builder.Exp{target, path, newValue}
	if len(createIfMissing) > 1 {
		panic(errors.New("too many arguments"))
	}
	if len(createIfMissing) > 0 {
		args = append(args, createIfMissing[0])
	}
	return builder.FuncExp("jsonb_set", args)
}

// JsonbSetLax builds the jsonb_set_lax function.
//
//	( target jsonb, path text[], new_value jsonb [, create_if_missing boolean [, null_value_treatment text ]] ) → jsonb
//
// If new_value is not NULL, behaves identically to jsonb_set. Otherwise behaves according to the value of null_value_treatment which must be one of 'raise_exception', 'use_json_null', 'delete_key', or 'return_target'. The default is 'use_json_null'.
func JsonbSetLax(target builder.Exp, path builder.Exp, newValue builder.Exp, options ...builder.Exp) builder.ExpBase {
	args := []builder.Exp{target, path, newValue}
	if len(options) > 2 {
		panic(errors.New("too many arguments"))
	}
	if len(options) > 0 {
		args = append(args, options...)
	}
	return builder.FuncExp("jsonb_set_lax", args)
}

// JsonbInsert builds the jsonb_insert function.
//
//	( target jsonb, path text[], new_value jsonb [, insert_after boolean ] ) → jsonb
//
// Returns target with new_value inserted. If the item designated by the path is an array element, new_value will be inserted before that item if insert_after is false (which is the default), or after it if insert_after is true. If the item designated by the path is an object field, new_value will be inserted only if the object does not already contain that key. All earlier steps in the path must exist, or the target is returned unchanged. As with the path oriented operators, negative integers that appear in the path count from the end of JSON arrays. If the last path step is an array index that is out of range, the new value is added at the beginning of the array if the index is negative, or at the end of the array if it is positive.
func JsonbInsert(target builder.Exp, path builder.Exp, newValue builder.Exp, insertAfter ...builder.Exp) builder.ExpBase {
	args := []builder.Exp{target, path, newValue}
	if len(insertAfter) > 1 {
		panic(errors.New("too many arguments"))
	}
	if len(insertAfter) > 0 {
		args = append(args, insertAfter[0])
	}
	return builder.FuncExp("jsonb_insert", args)
}

// JsonStripNulls builds the json_strip_nulls function.
//
//	( json ) → json
//
// Deletes all object fields that have null values from the given JSON value, recursively. Null values that are not object fields are untouched.
func JsonStripNulls(exp builder.Exp) builder.ExpBase {
	return builder.FuncExp("json_strip_nulls", []builder.Exp{exp})
}

// JsonbStripNulls builds the jsonb_strip_nulls function.
//
//	( jsonb ) → jsonb
//
// Deletes all object fields that have null values from the given JSON value, recursively. Null values that are not object fields are untouched.
func JsonbStripNulls(exp builder.Exp) builder.ExpBase {
	return builder.FuncExp("jsonb_strip_nulls", []builder.Exp{exp})
}

// JsonbPathExists builds the jsonb_path_exists function.
//
//	( target jsonb, path jsonpath [, vars jsonb [, silent boolean ]] ) → boolean
//
// Checks whether the JSON path returns any item for the specified JSON value. If the vars argument is specified, it must be a JSON object, and its fields provide named values to be substituted into the jsonpath expression. If the silent argument is specified and is true, the function suppresses the same errors as the @? and @@ operators do.
func JsonbPathExists(target builder.Exp, path builder.Exp, options ...builder.Exp) builder.ExpBase {
	args := []builder.Exp{target, path}
	if len(options) > 2 {
		panic(errors.New("too many arguments"))
	}
	if len(options) > 0 {
		args = append(args, options...)
	}
	return builder.FuncExp("jsonb_path_exists", args)
}

// JsonbPathMatch builds the jsonb_path_match function.
//
//	( target jsonb, path jsonpath [, vars jsonb [, silent boolean ]] ) → boolean
//
// Returns the result of a JSON path predicate check for the specified JSON value. Only the first item of the result is taken into account. If the result is not Boolean, then NULL is returned. The optional vars and silent arguments act the same as for jsonb_path_exists.
func JsonbPathMatch(target builder.Exp, path builder.Exp, options ...builder.Exp) builder.ExpBase {
	args := []builder.Exp{target, path}
	if len(options) > 2 {
		panic(errors.New("too many arguments"))
	}
	if len(options) > 0 {
		args = append(args, options...)
	}
	return builder.FuncExp("jsonb_path_match", args)
}

// JsonbPathQuery builds the jsonb_path_query function.
//
//	( target jsonb, path jsonpath [, vars jsonb [, silent boolean ]] ) → setof jsonb
//
// Returns all JSON items returned by the JSON path for the specified JSON value. The optional vars and silent arguments act the same as for jsonb_path_exists.
func JsonbPathQuery(target builder.Exp, path builder.Exp, options ...builder.Exp) builder.ExpBase {
	args := []builder.Exp{target, path}
	if len(options) > 2 {
		panic(errors.New("too many arguments"))
	}
	if len(options) > 0 {
		args = append(args, options...)
	}
	return builder.FuncExp("jsonb_path_query", args)
}

// JsonbPathQueryArray builds the jsonb_path_query_array function.
//
//	( target jsonb, path jsonpath [, vars jsonb [, silent boolean ]] ) → jsonb
//
// Returns all JSON items returned by the JSON path for the specified JSON value, as a JSON array. The optional vars and silent arguments act the same as for jsonb_path_exists.
func JsonbPathQueryArray(target builder.Exp, path builder.Exp, options ...builder.Exp) builder.ExpBase {
	args := []builder.Exp{target, path}
	if len(options) > 2 {
		panic(errors.New("too many arguments"))
	}
	if len(options) > 0 {
		args = append(args, options...)
	}
	return builder.FuncExp("jsonb_path_query_array", args)
}

// JsonbPathQueryFirst builds the jsonb_path_query_first function.
//
//	( target jsonb, path jsonpath [, vars jsonb [, silent boolean ]] ) → jsonb
//
// Returns the first JSON item returned by the JSON path for the specified JSON value. Returns NULL if there are no results. The optional vars and silent arguments act the same as for jsonb_path_exists.
func JsonbPathQueryFirst(target builder.Exp, path builder.Exp, options ...builder.Exp) builder.ExpBase {
	args := []builder.Exp{target, path}
	if len(options) > 2 {
		panic(errors.New("too many arguments"))
	}
	if len(options) > 0 {
		args = append(args, options...)
	}
	return builder.FuncExp("jsonb_path_query_first", args)
}

// JsonbPathExistsTZ builds the jsonb_path_exists function.
//
//	( target jsonb, path jsonpath [, vars jsonb [, silent boolean ]] ) → boolean
//
// Checks whether the JSON path returns any item for the specified JSON value. If the vars argument is specified, it must be a JSON object, and its fields provide named values to be substituted into the jsonpath expression. If the silent argument is specified and is true, the function suppresses the same errors as the @? and @@ operators do.
// This function acts like its counterpart without the _tz suffix, except that this functions supports comparisons of date/time values that require timezone-aware conversions. Due to this dependency, this function is marked as stable, which means this function cannot be used in indexes. The counterpart is immutable, and so can be used in indexes; but it will throw errors if asked to make such comparisons.
func JsonbPathExistsTZ(target builder.Exp, path builder.Exp, options ...builder.Exp) builder.ExpBase {
	args := []builder.Exp{target, path}
	if len(options) > 2 {
		panic(errors.New("too many arguments"))
	}
	if len(options) > 0 {
		args = append(args, options...)
	}
	return builder.FuncExp("jsonb_path_exists_tz", args)
}

// JsonbPathMatchTZ builds the jsonb_path_match function.
//
//	( target jsonb, path jsonpath [, vars jsonb [, silent boolean ]] ) → boolean
//
// Returns the result of a JSON path predicate check for the specified JSON value. Only the first item of the result is taken into account. If the result is not Boolean, then NULL is returned. The optional vars and silent arguments act the same as for jsonb_path_exists.
// This function acts like its counterpart without the _tz suffix, except that this functions supports comparisons of date/time values that require timezone-aware conversions. Due to this dependency, this function is marked as stable, which means this function cannot be used in indexes. The counterpart is immutable, and so can be used in indexes; but it will throw errors if asked to make such comparisons.
func JsonbPathMatchTZ(target builder.Exp, path builder.Exp, options ...builder.Exp) builder.ExpBase {
	args := []builder.Exp{target, path}
	if len(options) > 2 {
		panic(errors.New("too many arguments"))
	}
	if len(options) > 0 {
		args = append(args, options...)
	}
	return builder.FuncExp("jsonb_path_match_tz", args)
}

// JsonbPathQueryTZ builds the jsonb_path_query function.
//
//	( target jsonb, path jsonpath [, vars jsonb [, silent boolean ]] ) → setof jsonb
//
// Returns all JSON items returned by the JSON path for the specified JSON value. The optional vars and silent arguments act the same as for jsonb_path_exists.
// This function acts like its counterpart without the _tz suffix, except that this functions supports comparisons of date/time values that require timezone-aware conversions. Due to this dependency, this function is marked as stable, which means this function cannot be used in indexes. The counterpart is immutable, and so can be used in indexes; but it will throw errors if asked to make such comparisons.
func JsonbPathQueryTZ(target builder.Exp, path builder.Exp, options ...builder.Exp) builder.ExpBase {
	args := []builder.Exp{target, path}
	if len(options) > 2 {
		panic(errors.New("too many arguments"))
	}
	if len(options) > 0 {
		args = append(args, options...)
	}
	return builder.FuncExp("jsonb_path_query_tz", args)
}

// JsonbPathQueryArrayTZ builds the jsonb_path_query_array function.
//
//	( target jsonb, path jsonpath [, vars jsonb [, silent boolean ]] ) → jsonb
//
// Returns all JSON items returned by the JSON path for the specified JSON value, as a JSON array. The optional vars and silent arguments act the same as for jsonb_path_exists.
// This function acts like its counterpart without the _tz suffix, except that this functions supports comparisons of date/time values that require timezone-aware conversions. Due to this dependency, this function is marked as stable, which means this function cannot be used in indexes. The counterpart is immutable, and so can be used in indexes; but it will throw errors if asked to make such comparisons.
func JsonbPathQueryArrayTZ(target builder.Exp, path builder.Exp, options ...builder.Exp) builder.ExpBase {
	args := []builder.Exp{target, path}
	if len(options) > 2 {
		panic(errors.New("too many arguments"))
	}
	if len(options) > 0 {
		args = append(args, options...)
	}
	return builder.FuncExp("jsonb_path_query_array_tz", args)
}

// JsonbPathQueryFirstTZ builds the jsonb_path_query_first function.
//
//	( target jsonb, path jsonpath [, vars jsonb [, silent boolean ]] ) → jsonb
//
// Returns the first JSON item returned by the JSON path for the specified JSON value. Returns NULL if there are no results. The optional vars and silent arguments act the same as for jsonb_path_exists.
// This function acts like its counterpart without the _tz suffix, except that this functions supports comparisons of date/time values that require timezone-aware conversions. Due to this dependency, this function is marked as stable, which means this function cannot be used in indexes. The counterpart is immutable, and so can be used in indexes; but it will throw errors if asked to make such comparisons.
func JsonbPathQueryFirstTZ(target builder.Exp, path builder.Exp, options ...builder.Exp) builder.ExpBase {
	args := []builder.Exp{target, path}
	if len(options) > 2 {
		panic(errors.New("too many arguments"))
	}
	if len(options) > 0 {
		args = append(args, options...)
	}
	return builder.FuncExp("jsonb_path_query_first_tz", args)
}

// JsonbPretty builds the jsonb_pretty function.
//
//	( jsonb ) → text
//
// Converts the given JSON value to pretty-printed, indented text.
func JsonbPretty(exp builder.Exp) builder.ExpBase {
	return builder.FuncExp("jsonb_pretty", []builder.Exp{exp})
}

// JsonTypeof builds the json_typeof function.
//
//	( json ) → text
//
// Returns the type of the top-level JSON value as a text string. Possible types are object, array, string, number, boolean, and null. (The null result should not be confused with an SQL NULL; see the examples.)
func JsonTypeof(exp builder.Exp) builder.ExpBase {
	return builder.FuncExp("json_typeof", []builder.Exp{exp})
}

// JsonbTypeof builds the jsonb_typeof function.
//
//	( jsonb ) → text
//
// Returns the type of the top-level JSON value as a text string. Possible types are object, array, string, number, boolean, and null. (The null result should not be confused with an SQL NULL; see the examples.)
func JsonbTypeof(exp builder.Exp) builder.ExpBase {
	return builder.FuncExp("jsonb_typeof", []builder.Exp{exp})
}
