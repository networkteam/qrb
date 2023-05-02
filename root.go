package qrb

import (
	"github.com/networkteam/qrb/builder"
)

// This file exports the root level functions for building queries.
// All functions are just wrappers around the builder package, so all the builder types and interfaces don't clutter the root package exports.

// Build starts a new query builder based on the given SQLWriter.
// For executing the query, use qrbpgx.Build or qrbsql.Build which can set an executor specific to a driver.
func Build(w builder.SQLWriter) *builder.QueryBuilder {
	return builder.Build(w)
}

// ---

// With starts a new builder with the given WITH query.
// Call WithBuilder.As to define the query.
func With(queryName string) builder.WithWithBuilder {
	return builder.With(queryName)
}

// WithRecursive starts a new builder with the given WITH RECURSIVE query.
func WithRecursive(queryName string) builder.WithWithBuilder {
	return builder.WithRecursive(queryName)
}

// Select the given output expressions for the select list and start a new SelectBuilder.
func Select(exps ...builder.Exp) builder.SelectSelectBuilder {
	var selectBuilder builder.SelectBuilder
	return selectBuilder.Select(exps...)
}

// SelectJson sets the JSON selection for this select builder.
//
// It will always be the first element in the select list.
// It can be modified later by SelectBuilder.ApplySelectJson.
func SelectJson(obj builder.JsonBuildObjectBuilder) builder.SelectJsonSelectBuilder {
	var selectBuilder builder.SelectBuilder
	return selectBuilder.ApplySelectJson(func(builder builder.JsonBuildObjectBuilder) builder.JsonBuildObjectBuilder { return obj })
}

// Agg builds an aggregate function expression.
func Agg(name string, exps []builder.Exp) builder.AggExpBuilder {
	return builder.Agg(name, exps)
}

// Func is a function call expression.
func Func(name string, args ...builder.Exp) builder.FuncBuilder {
	return builder.Func(name, args...)
}

func RowsFrom(fn builder.FuncBuilder, fns ...builder.FuncBuilder) builder.RowsFromBuilder {
	return builder.NewRowsFromBuilder(
		append([]builder.FuncBuilder{fn}, fns...)...,
	)
}

func And(exps ...builder.Exp) builder.Exp {
	return builder.And(exps...)
}

func Or(exps ...builder.Exp) builder.Exp {
	return builder.Or(exps...)
}

func Case(exp ...builder.Exp) builder.CaseBuilder {
	return builder.Case(exp...)
}

func Coalesce(exp builder.Exp, rest ...builder.Exp) builder.FuncExp {
	return builder.Coalesce(exp, rest...)
}

func NullIf(value1, value2 builder.Exp) builder.FuncExp {
	return builder.NullIf(value1, value2)
}

func Greatest(exp builder.Exp, rest ...builder.Exp) builder.FuncExp {
	return builder.Greatest(exp, rest...)
}

func Least(exp builder.Exp, rest ...builder.Exp) builder.FuncExp {
	return builder.Least(exp, rest...)
}

// Arg creates an expression that represents an argument that will be bound to a placeholder with the given value.
// Each call to Arg will create a new placeholder and emit the argument when writing the query.
func Arg(argument any) builder.ExpBase {
	return builder.Arg(argument)
}

// Args creates argument expressions for the given arguments (of the same type).
func Args[T any](arguments ...T) builder.Expressions {
	return builder.Args(arguments...)
}

// Bind creates an expression that represents an argument that will be bound to a placeholder with the given value.
// If Bind is called again with the same name, the same placeholder will be used.
func Bind(argName string) builder.Exp {
	return builder.Bind(argName)
}

// N writes the given name / identifier.
//
// It will validate the identifier when writing the query,
// but it will not detect all invalid identifiers that are invalid in PostgreSQL (especially considering reserved keywords).
func N(s string) builder.IdentExp {
	return builder.N(s)
}

func String(s string) builder.Exp {
	return builder.String(s)
}

func Float(f float64) builder.Exp {
	return builder.Float(f)
}

func Int(s int) builder.Exp {
	return builder.Int(s)
}

func Bool(b bool) builder.Exp {
	return builder.Bool(b)
}

func Array(exps ...builder.Exp) builder.Exp {
	return builder.Array(exps...)
}

func Null() builder.Exp {
	return builder.Null()
}

func Default() builder.Exp {
	return builder.Default()
}

func Interval(s string) builder.Exp {
	return builder.Interval(s)
}

// Exps returns a slice of expressions, just for syntactic sugar.
// TODO We could use this as a way to express a scalar list of expressions e.g. for IN by using a custom slice type
func Exps(exps ...builder.Exp) builder.Expressions {
	return exps
}

// --- Subquery Expressions

func Exists(subquery builder.SelectExp) builder.Exp {
	return builder.Exists(subquery)
}

// --- Commands

func InsertInto(tableName builder.Identer) builder.InsertBuilder {
	return builder.InsertInto(tableName)
}

func Update(tableName builder.Identer) builder.UpdateBuilder {
	return builder.Update(tableName)
}

func DeleteFrom(tableName builder.Identer) builder.DeleteBuilder {
	return builder.DeleteFrom(tableName)
}
