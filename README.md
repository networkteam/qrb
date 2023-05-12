# QRB

[![GoDoc](https://godoc.org/github.com/networkteam/qrb?status.svg)](https://godoc.org/github.com/networkteam/qrb)
[![Build Status](https://github.com/networkteam/qrb/workflows/Go/badge.svg)](https://github.com/networkteam/qrb/actions?workflow=run%20tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/networkteam/qrb)](https://goreportcard.com/report/github.com/networkteam/qrb)
[![codecov](https://codecov.io/gh/networkteam/qrb/branch/main/graph/badge.svg?token=S8X8TMLQ9O)](https://codecov.io/gh/networkteam/qrb)

A PostgreSQL query builder in Go.

## Why?

* Pure focus on PostgreSQL dialect
* Include support for selecting JSON from queries (e.g. `json_build_object`) to fetch nested data
* With other query builders, there is often a need to write raw SQL to achieve the desired result
* Supports `database/sql` and [`pgx`](https://github.com/jackc/pgx) drivers via `qrbsql`and `qrbpgx` packages

## Design goals

* All builder data structures are immutable by design
* Implement the full PostgreSQL feature set, including lesser used features
* Use explicit types instead of `any` where possible and do not use reflection
* First-hand support for JSON selection (i.e. use `json_build_object` and `json_agg` to select hierarchical data via JSON)
* Write SQL as Go code following the natural order of the query parts
* Guide the developer by providing builder types with methods appropriate for the current context 

## Install

```bash
go get github.com/networkteam/qrb
```

## Examples

### Select JSON with common table expression

```go
myCategory := "SQL Hacks"

q := qrb.With("author_json").As(
        qrb.Select(
                qrb.N("authors.author_id"),
        ).
        Select(
            fn.JsonBuildObject().
                Prop("id", qrb.N("authors.author_id")).
                Prop("name", qrb.N("authors.name")),
        ).As("json").
        From(qrb.N("authors")),
    ).
    Select(
        qrb.N("posts.post_id"),
        fn.JsonBuildObject().
            Prop("title", qrb.N("posts.title")).
            Prop("author", qrb.N("author_json.json")),
    ).
    From(qrb.N("posts")).
    LeftJoin(qrb.N("author_json")).On(qrb.N("posts.author_id").Eq(qrb.N("author_json.author_id"))).
    Where(qrb.N("posts.category").Eq(qrb.Arg(myCategory))).
    OrderBy(qrb.N("posts.created_at")).Desc().NullsLast()
}
```

<details>
<summary>Generated SQL</summary>

```sql
WITH author_json AS (
    SELECT authors.author_id, json_build_object('id', authors.author_id, 'name', authors.name) AS json
    FROM authors
)
SELECT posts.post_id, json_build_object('title', posts.title, 'author', author_json.json)
FROM posts
    LEFT JOIN author_json ON posts.author_id = author_json.author_id
WHERE posts.category = $1
ORDER BY posts.created_at DESC NULLS LAST
```
</details>

### Placeholders

`qrb.Bind` supports named arguments that can be supplied after the query value has been built.

`qrb.Arg` will generate positional arguments with the specified value.

Both can be combined.

```go
q := qrb.
    Select(qrb.N("*")).
    From(qrb.N("employees")).
    Where(qrb.And(
        qrb.Or(
            qrb.N("firstname").ILike(qrb.Bind("search")),
            qrb.N("lastname").ILike(qrb.Bind("search")),
        ),
        qrb.N("active").Eq(qrb.Arg(true)),
    ))

sql, args, err := qrb.
    Build(q).
    WithNamedArgs(map[string]any{"search": "Jo%"}).
    ToSQL()

// args: []any{"Jo%", true}
```

<details>
<summary>Generated SQL</summary>

```sql
SELECT *
FROM employees
WHERE ((firstname ILIKE $1) OR (lastname ILIKE $1))
  AND (active = $2)
```
</details>

### Select `WITH RECURSIVE`

```go
q := qrb.WithRecursive("employee_recursive").ColumnNames("distance", "employee_name", "manager_name").As(
        qrb.Select(qrb.Int(1), qrb.N("employee_name"), qrb.N("manager_name")).
            From(qrb.N("employee")).
            Where(qrb.N("manager_name").Eq(qrb.String("Mary"))).
            Union().All().
            Select(qrb.N("er.distance").Op("+", qrb.Int(1)), qrb.N("e.employee_name"), qrb.N("e.manager_name")).
            From(qrb.N("employee_recursive")).As("er").
            From(qrb.N("employee")).As("e").
            Where(qrb.N("er.employee_name").Eq(qrb.N("e.manager_name"))),
    ).
    Select(qrb.N("distance"), qrb.N("employee_name")).From(qrb.N("employee_recursive"))
```

<details>
<summary>Generated SQL</summary>

```sql
WITH RECURSIVE employee_recursive (distance, employee_name, manager_name) AS (
    SELECT 1, employee_name, manager_name
    FROM employee
    WHERE manager_name = 'Mary'
    UNION ALL
    SELECT er.distance + 1, e.employee_name, e.manager_name
    FROM employee_recursive AS er, employee AS e
    WHERE er.employee_name = e.manager_name
)
SELECT distance, employee_name
FROM employee_recursive
```
</details>

### Functions including `ROWS FROM`

```go
q := qrb.Select(qrb.N("*")).
    From(qrb.RowsFrom(
        qrb.Func("json_to_recordset", qrb.String(`[{"a":40,"b":"foo"},{"a":"100","b":"bar"}]`)).ColumnDefinition("a", "INTEGER").ColumnDefinition("b", "TEXT"),
        qrb.Func("generate_series", qrb.Int(1), qrb.Int(3)),
    ).WithOrdinality()).As("x").ColumnAliases("p", "q", "s").
    OrderBy(qrb.N("p"))
```

<details>
<summary>Generated SQL</summary>

```sql
SELECT *
FROM ROWS FROM (
         json_to_recordset('[{"a":40,"b":"foo"},{"a":"100","b":"bar"}]') AS (a INTEGER, b TEXT),
         generate_series(1, 3)
         ) WITH ORDINALITY AS x (p, q, s)
ORDER BY
    p
```
</details>

### Group by with grouping sets, rollup and cube

```go
q := qrb.
    Select(qrb.N("a"), qrb.N("b"), qrb.N("c")).
    From(qrb.N("test1")).
    GroupBy().Distinct().
    Rollup(
        qrb.Exps(qrb.N("a"), qrb.N("b")),
    ).
    Rollup(
        qrb.Exps(qrb.N("a"), qrb.N("c")),
    )
```

<details>
<summary>Generated SQL</summary>

```sql
SELECT a, b, c
FROM test1
GROUP BY DISTINCT ROLLUP (a, b), ROLLUP (a, c)
```
</details>

### Executing queries with pgx

```go
conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
if err != nil {
    log.Fatalf("Unable to connect to database: %v", err)
}

q := qrb.
	With("longest_track").As(
        qrb.Select(qrb.N(`"AlbumId"`), qrb.N(`"Milliseconds"`)).
            From(qrb.N(`"Track"`)).
            OrderBy(qrb.N(`"Milliseconds"`)).Desc().
            Limit(qrb.Int(1)),
    ).
    Select(qrb.N(`"Title"`)).As(`"AlbumTitle"`).
    Select(qrb.N(`"Name"`)).As(`"ArtistName"`).
    Select(qrb.N(`"Milliseconds"`)).As(`"Length"`).
    From(qrb.N(`"Album"`)).
    Join(qrb.N(`"longest_track"`)).Using(`"AlbumId"`).
    Join(qrb.N(`"Artist"`)).Using(`"ArtistId"`)

row, err := qrbpgx.
    Build(q).
    WithExecutor(conn).
    QueryRow(ctx)
```

The executor can either be a `*pgx.Conn`, `*pgxpool.Pool` or `pgx.Tx` (or any other type implementing `qrbpgx.Executor`).
It can be specified after building the query with `WithExecutor` or in advance via `qrbpgx.NewExecutorBuilder`. 

### Executing queries with database/sql

(e.g. github.com/lib/pq or pgx with adapter)

```go
db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
if err != nil {
    panic(err)
}

q := qrb.
	With("longest_track").As(
        qrb.Select(qrb.N(`"AlbumId"`), qrb.N(`"Milliseconds"`)).
            From(qrb.N(`"Track"`)).
            OrderBy(qrb.N(`"Milliseconds"`)).Desc().
            Limit(qrb.Int(1)),
    ).
    Select(qrb.N(`"Title"`)).As(`"AlbumTitle"`).
    Select(qrb.N(`"Name"`)).As(`"ArtistName"`).
    Select(qrb.N(`"Milliseconds"`)).As(`"Length"`).
    From(qrb.N(`"Album"`)).
    Join(qrb.N(`"longest_track"`)).Using(`"AlbumId"`).
    Join(qrb.N(`"Artist"`)).Using(`"ArtistId"`)

row, err := qrbsql.
    Build(q).
    WithExecutor(db).
    QueryRow(ctx)
```

The executor can either be a `*sql.DB` or `*sql.Tx` (or any other type implementing `qrbsql.Executor`).
It can be specified after building the query with `WithExecutor` or in advance via `qrbpgx.NewExecutorBuilder`.

## License

[MIT](./LICENSE)
