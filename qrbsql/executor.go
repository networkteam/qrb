package qrbsql

import (
	"context"
	"database/sql"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/builder"
)

type QueryBuilder struct {
	*builder.QueryBuilder
}

type ExecutiveQueryBuilder struct {
	*builder.QueryBuilder
	executor Executor
}

type ExecutorBuilder struct {
	executor Executor
}

// Executor is the interface that wraps the basic Query, QueryRow and Exec methods.
// It allows to use *sql.DB and *sql.TX as an executor.
type Executor interface {
	QueryContext(ctx context.Context, sql string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, sql string, args ...any) *sql.Row
	ExecContext(ctx context.Context, sql string, args ...any) (sql.Result, error)
}

// Build a new QueryBuilder that can execute queries via pgx conn / tx after calling QueryBuilder.WithExecutor.
func Build(builder builder.SQLWriter) *QueryBuilder {
	return &QueryBuilder{
		QueryBuilder: qrb.Build(builder),
	}
}

// WithExecutor sets the executor to use for the query.
// It can be a pgx conn or tx.
func (b *QueryBuilder) WithExecutor(executor Executor) *ExecutiveQueryBuilder {
	return &ExecutiveQueryBuilder{
		QueryBuilder: b.QueryBuilder,
		executor:     executor,
	}
}

// NewExecutorBuilder starts with an Executor and allows to build an executable query via ExecutorBuilder.Build.
func NewExecutorBuilder(executor Executor) *ExecutorBuilder {
	return &ExecutorBuilder{
		executor: executor,
	}
}

func (b *ExecutorBuilder) Build(builder builder.SQLWriter) *ExecutiveQueryBuilder {
	return &ExecutiveQueryBuilder{
		QueryBuilder: qrb.Build(builder),
		executor:     b.executor,
	}
}

func (b *ExecutiveQueryBuilder) WithNamedArgs(args map[string]any) *ExecutiveQueryBuilder {
	b.QueryBuilder.WithNamedArgs(args)
	return b
}

func (b *ExecutiveQueryBuilder) WithoutValidation() *ExecutiveQueryBuilder {
	b.QueryBuilder.WithoutValidation()
	return b
}

func (b *ExecutiveQueryBuilder) Query(ctx context.Context) (rows *sql.Rows, err error) {
	sql, args, err := b.ToSQL()
	if err != nil {
		return rows, err
	}
	return b.executor.QueryContext(ctx, sql, args...)
}

func (b *ExecutiveQueryBuilder) QueryRow(ctx context.Context) (row *sql.Row, err error) {
	sql, args, err := b.ToSQL()
	if err != nil {
		return row, err
	}
	return b.executor.QueryRowContext(ctx, sql, args...), nil
}

func (b *ExecutiveQueryBuilder) Exec(ctx context.Context) (result sql.Result, err error) {
	sql, args, err := b.ToSQL()
	if err != nil {
		return result, err
	}
	return b.executor.ExecContext(ctx, sql, args...)
}
