package builder

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type SQLWriter interface {
	WriteSQL(sb *SQLBuilder)
}

type innerSQLWriter interface {
	innerWriteSQL(sb *SQLBuilder)
}

// TODO It would be better to separate building the SQL and getting arguments to allow for re-use of queries.

// Write the actual SQL and generate arguments.
// This is internal, it is exposed via qrb.Build.
func writeToSQLString(w SQLWriter, namedArgs map[string]any, opts sqlBuilderOpts) (sql string, args []any, err error) {
	sb := newSqlBuilder(opts)

	if iw, ok := w.(innerSQLWriter); ok {
		iw.innerWriteSQL(sb)
	} else {
		w.WriteSQL(sb)
	}

	sql = sb.sb.String()
	args = sb.args

	for argName, argIdx := range sb.namedArgs {
		argValue, exists := namedArgs[argName]
		if !exists {
			return "", nil, fmt.Errorf("missing named argument %q", argName)
		}
		args[argIdx-1] = argValue
	}

	return sql, args, errors.Join(sb.errs...)
}

type SQLBuilder struct {
	opts sqlBuilderOpts
	sb   *strings.Builder
	// Current positional placeholder index
	argIdx int
	// List of arguments created by CreatePlaceholder or BindPlaceholder (it only adds a nil value for later binding).
	args []any
	// Map of named arguments to positional placeholder index.
	namedArgs map[string]int
	// List of errors that occurred while building the actual SQL.
	errs []error
}

type sqlBuilderOpts struct {
	validating  bool
	prettyPrint bool
}

func newSqlBuilder(opts sqlBuilderOpts) *SQLBuilder {
	return &SQLBuilder{
		sb:   new(strings.Builder),
		opts: opts,
	}
}

func (b *SQLBuilder) WriteRune(r rune) {
	b.sb.WriteRune(r)
}

func (b *SQLBuilder) WriteString(s string) {
	b.sb.WriteString(s)
}

func (b *SQLBuilder) CreatePlaceholder(argument any) string {
	b.args = append(b.args, argument)
	b.argIdx++
	return "$" + strconv.Itoa(b.argIdx)
}

func (b *SQLBuilder) BindPlaceholder(name string) string {
	if b.namedArgs == nil {
		b.namedArgs = make(map[string]int)
	}
	argIdx, exists := b.namedArgs[name]
	if !exists {
		// Add an empty argument, it will be replaced later by the named argument.
		b.args = append(b.args, nil)
		b.argIdx++
		argIdx = b.argIdx
		b.namedArgs[name] = argIdx
	}
	return "$" + strconv.Itoa(argIdx)
}

func (b *SQLBuilder) Validating() bool {
	return b.opts.validating
}

func (b *SQLBuilder) AddError(err error) {
	b.errs = append(b.errs, err)
}
