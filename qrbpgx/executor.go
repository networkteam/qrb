package qrbpgx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

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
// It allows to use pgx conn, pool and tx as an executor.
type Executor interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
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

func (b *ExecutiveQueryBuilder) Query(ctx context.Context) (rows pgx.Rows, err error) {
	sql, args, err := b.ToSQL()
	if err != nil {
		return rows, err
	}
	return b.executor.Query(ctx, sql, args...)
}

func (b *ExecutiveQueryBuilder) QueryRow(ctx context.Context) (row pgx.Row, err error) {
	sql, args, err := b.ToSQL()
	if err != nil {
		return row, err
	}
	return b.executor.QueryRow(ctx, sql, args...), nil
}

func (b *ExecutiveQueryBuilder) Exec(ctx context.Context) (commandTag pgconn.CommandTag, err error) {
	sql, args, err := b.ToSQL()
	if err != nil {
		return commandTag, err
	}
	return b.executor.Exec(ctx, sql, args...)
}
