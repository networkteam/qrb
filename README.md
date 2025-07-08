# QRB - PostgreSQL Query Builder

[![GoDoc](https://godoc.org/github.com/networkteam/qrb?status.svg)](https://godoc.org/github.com/networkteam/qrb)
[![Build Status](https://github.com/networkteam/qrb/actions/workflows/test.yml/badge.svg)](https://github.com/networkteam/qrb/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/networkteam/qrb)](https://goreportcard.com/report/github.com/networkteam/qrb)
[![codecov](https://codecov.io/gh/networkteam/qrb/branch/main/graph/badge.svg?token=S8X8TMLQ9O)](https://codecov.io/gh/networkteam/qrb)

A comprehensive PostgreSQL query builder in Go with extensive support for PostgreSQL-specific features.

## Why QRB?

- **Pure PostgreSQL focus** - No compromises for multi-database support
- **JSON-first design** - First-class support for `json_build_object` and `json_agg` to build hierarchical data
- **Complete feature set** - Supports advanced PostgreSQL features like CTEs, window functions, arrays, and more
- **Type-safe** - Uses explicit types instead of `any` where possible, no reflection
- **Immutable builders** - All data structures are immutable by design
- **Multiple drivers** - Works with both `database/sql` and `pgx` via `qrbsql` and `qrbpgx` packages

## Installation

```bash
go get github.com/networkteam/qrb
```

## Quick Start

```go
package main

import (
    "fmt"
    . "github.com/networkteam/qrb"
	"github.com/networkteam/qrb/qrbsql"
)

func main() {
	filter := true
    q := Select(N("name"), N("email")).
        From(N("users")).
        Where(N("active").Eq(Arg(filter))).
        OrderBy(N("name"))

    sql, args, _ := qrbsql.Build(q).ToSQL()
    fmt.Println(sql) // SELECT name, email FROM users WHERE active = $1 ORDER BY name
	fmt.Println(args) // [true]
}
```

## Core Concepts

- **Immutable Builders**: All builders return new instances, making them safe for reuse
- **Expressions**: Use `qrb.N()` for identifiers, `qrb.Arg()` for parameters, and `qrb.String()`, `qrb.Int()`, etc. for literals
- **Fluent API**: Chain method calls naturally following SQL structure
- **Type Safety**: Builders guide you through valid SQL construction with appropriate method availability

## Examples

### Basic Queries

#### Simple SELECT

```go
q := qrb.Select(qrb.N("*")).From(qrb.N("users"))
```

```sql
SELECT * FROM users
```

#### SELECT with WHERE

```go
q := qrb.Select(qrb.N("name"), qrb.N("email")).
    From(qrb.N("users")).
    Where(qrb.N("active").Eq(qrb.Bool(true)))
```

```sql
SELECT name, email FROM users WHERE active = true
```

#### SELECT with multiple conditions

```go
q := qrb.Select(qrb.N("*")).
    From(qrb.N("employees")).
    Where(qrb.And(
        qrb.Or(
            qrb.N("firstname").ILike(qrb.Arg("John%")),
            qrb.N("lastname").ILike(qrb.Arg("John%")),
        ),
        qrb.N("active").Eq(qrb.Bool(true)),
    ))
```

```sql
SELECT * FROM employees 
WHERE ((firstname ILIKE $1) OR (lastname ILIKE $1)) AND (active = $2)
```

#### SELECT DISTINCT

```go
q := qrb.Select().Distinct().
    Select(qrb.N("department")).
    From(qrb.N("employees"))
```

```sql
SELECT DISTINCT department FROM employees
```

#### SELECT with ORDER BY and LIMIT

```go
q := qrb.Select(qrb.N("name"), qrb.N("salary")).
    From(qrb.N("employees")).
    OrderBy(qrb.N("salary")).Desc().NullsLast().
    Limit(qrb.Int(10)).
    Offset(qrb.Int(20))
```

```sql
SELECT name, salary FROM employees 
ORDER BY salary DESC NULLS LAST 
LIMIT 10 OFFSET 20
```

### CRUD Operations

#### INSERT with VALUES

```go
q := qrb.InsertInto(qrb.N("users")).
    ColumnNames("name", "email", "active").
    Values(qrb.String("John Doe"), qrb.String("john@example.com"), qrb.Bool(true))
```

```sql
INSERT INTO users (name, email, active) 
VALUES ('John Doe', 'john@example.com', true)
```

#### INSERT multiple rows

```go
q := qrb.InsertInto(qrb.N("products")).
    ColumnNames("name", "price", "category").
    Values(qrb.String("Laptop"), qrb.Float(999.99), qrb.String("Electronics")).
    Values(qrb.String("Book"), qrb.Float(19.99), qrb.String("Literature"))
```

```sql
INSERT INTO products (name, price, category) VALUES
    ('Laptop', 999.99, 'Electronics'),
    ('Book', 19.99, 'Literature')
```

#### INSERT with SELECT

```go
q := qrb.InsertInto(qrb.N("archived_users")).
    Query(qrb.Select(qrb.N("*")).From(qrb.N("users")).Where(qrb.N("active").Eq(qrb.Bool(false))))
```

```sql
INSERT INTO archived_users SELECT * FROM users WHERE active = false
```

#### INSERT with RETURNING

```go
q := qrb.InsertInto(qrb.N("users")).
    ColumnNames("name", "email").
    Values(qrb.String("Jane Doe"), qrb.String("jane@example.com")).
    Returning(qrb.N("id"), qrb.N("created_at"))
```

```sql
INSERT INTO users (name, email) VALUES ('Jane Doe', 'jane@example.com')
RETURNING id, created_at
```

#### UPSERT (INSERT with ON CONFLICT)

```go
q := qrb.InsertInto(qrb.N("users")).
    ColumnNames("email", "name").
    Values(qrb.String("john@example.com"), qrb.String("John Updated")).
    OnConflict(qrb.N("email")).DoUpdate().
    Set("name", qrb.N("EXCLUDED.name")).
    Set("updated_at", qrb.N("NOW()"))
```

```sql
INSERT INTO users (email, name) VALUES ('john@example.com', 'John Updated')
ON CONFLICT (email) DO UPDATE 
SET name = EXCLUDED.name, updated_at = NOW()
```

#### UPDATE

```go
q := qrb.Update(qrb.N("users")).
    Set("name", qrb.String("Updated Name")).
    Set("updated_at", qrb.N("NOW()")).
    Where(qrb.N("id").Eq(qrb.Arg(123)))
```

```sql
UPDATE users SET name = 'Updated Name', updated_at = NOW() WHERE id = $1
```

#### UPDATE with FROM

```go
q := qrb.Update(qrb.N("employees")).
    Set("department_name", qrb.N("d.name")).
    From(qrb.N("departments")).As("d").
    Where(qrb.N("employees.department_id").Eq(qrb.N("d.id")))
```

```sql
UPDATE employees SET department_name = d.name 
FROM departments AS d 
WHERE employees.department_id = d.id
```

#### DELETE

```go
q := qrb.DeleteFrom(qrb.N("users")).
    Where(qrb.N("active").Eq(qrb.Bool(false)))
```

```sql
DELETE FROM users WHERE active = false
```

#### DELETE with USING

```go
q := qrb.DeleteFrom(qrb.N("orders")).
    Using(qrb.N("customers")).
    Where(qrb.And(
        qrb.N("orders.customer_id").Eq(qrb.N("customers.id")),
        qrb.N("customers.status").Eq(qrb.String("inactive")),
    ))
```

```sql
DELETE FROM orders USING customers 
WHERE orders.customer_id = customers.id AND customers.status = 'inactive'
```

### Joins

#### INNER JOIN

```go
q := qrb.Select(qrb.N("u.name"), qrb.N("p.title")).
    From(qrb.N("users")).As("u").
    Join(qrb.N("posts")).As("p").On(qrb.N("u.id").Eq(qrb.N("p.user_id")))
```

```sql
SELECT u.name, p.title FROM users AS u 
JOIN posts AS p ON u.id = p.user_id
```

#### LEFT JOIN

```go
q := qrb.Select(qrb.N("u.name"), qrb.N("p.title")).
    From(qrb.N("users")).As("u").
    LeftJoin(qrb.N("posts")).As("p").On(qrb.N("u.id").Eq(qrb.N("p.user_id")))
```

```sql
SELECT u.name, p.title FROM users AS u 
LEFT JOIN posts AS p ON u.id = p.user_id
```

#### JOIN with USING

```go
q := qrb.Select(qrb.N("u.name"), qrb.N("p.title")).
    From(qrb.N("users")).As("u").
    Join(qrb.N("posts")).As("p").Using("user_id")
```

```sql
SELECT u.name, p.title FROM users AS u 
JOIN posts AS p USING (user_id)
```

#### Multiple JOINs

```go
q := qrb.Select(qrb.N("u.name"), qrb.N("p.title"), qrb.N("c.name")).
    From(qrb.N("users")).As("u").
    Join(qrb.N("posts")).As("p").On(qrb.N("u.id").Eq(qrb.N("p.user_id"))).
    Join(qrb.N("categories")).As("c").On(qrb.N("p.category_id").Eq(qrb.N("c.id")))
```

```sql
SELECT u.name, p.title, c.name FROM users AS u 
JOIN posts AS p ON u.id = p.user_id 
JOIN categories AS c ON p.category_id = c.id
```

### Aggregation & Grouping

#### GROUP BY with aggregate functions

```go
q := qrb.Select(qrb.N("department")).
    Select(fn.Count(qrb.N("*"))).As("employee_count").
    From(qrb.N("employees")).
    GroupBy(qrb.N("department"))
```

```sql
SELECT department, count(*) AS employee_count 
FROM employees 
GROUP BY department
```

#### GROUP BY with HAVING

```go
q := qrb.Select(qrb.N("department")).
    Select(fn.Avg(qrb.N("salary"))).As("avg_salary").
    From(qrb.N("employees")).
    GroupBy(qrb.N("department")).
    Having(fn.Avg(qrb.N("salary")).Gt(qrb.Int(50000)))
```

```sql
SELECT department, avg(salary) AS avg_salary 
FROM employees 
GROUP BY department 
HAVING avg(salary) > 50000
```

#### GROUP BY with ROLLUP

```go
q := qrb.Select(qrb.N("department"), qrb.N("job_title"), fn.Sum(qrb.N("salary"))).
    From(qrb.N("employees")).
    GroupBy().
    Rollup(
        qrb.Exps(qrb.N("department")),
        qrb.Exps(qrb.N("job_title")),
    )
```

```sql
SELECT department, job_title, sum(salary) 
FROM employees 
GROUP BY ROLLUP (department, job_title)
```

#### GROUP BY with GROUPING SETS

```go
q := qrb.Select(qrb.N("department"), qrb.N("job_title"), fn.Sum(qrb.N("salary"))).
    From(qrb.N("employees")).
    GroupBy().
    GroupingSets(
        qrb.Exps(qrb.N("department")),
        qrb.Exps(qrb.N("job_title")),
        qrb.Exps(),
    )
```

```sql
SELECT department, job_title, sum(salary) 
FROM employees 
GROUP BY GROUPING SETS (department, job_title, ())
```

### Window Functions

#### ROW_NUMBER

```go
q := qrb.Select(
    qrb.N("name"),
    qrb.N("salary"),
    fn.RowNumber().Over().PartitionBy(qrb.N("department")).OrderBy(qrb.N("salary")).Desc(),
).From(qrb.N("employees"))
```

```sql
SELECT name, salary, row_number() OVER (PARTITION BY department ORDER BY salary DESC) 
FROM employees
```

#### RANK and DENSE_RANK

```go
q := qrb.Select(
    qrb.N("name"),
    qrb.N("salary"),
    fn.Rank().Over().PartitionBy(qrb.N("department")).OrderBy(qrb.N("salary")).Desc(),
    fn.DenseRank().Over().PartitionBy(qrb.N("department")).OrderBy(qrb.N("salary")).Desc(),
).From(qrb.N("employees"))
```

```sql
SELECT name, salary, 
       rank() OVER (PARTITION BY department ORDER BY salary DESC),
       dense_rank() OVER (PARTITION BY department ORDER BY salary DESC)
FROM employees
```

#### Named Windows

```go
q := qrb.Select(
    fn.Sum(qrb.N("salary")).Over("w"),
    fn.Avg(qrb.N("salary")).Over("w"),
).From(qrb.N("employees")).
Window("w").As().PartitionBy(qrb.N("department")).OrderBy(qrb.N("salary")).Desc().
SelectBuilder
```

```sql
SELECT sum(salary) OVER w, avg(salary) OVER w
FROM employees
WINDOW w AS (PARTITION BY department ORDER BY salary DESC)
```

### JSON Operations

#### Simple JSON object

```go
q := qrb.Select(
    fn.JsonBuildObject().
        Prop("id", qrb.N("id")).
        Prop("name", qrb.N("name")).
        Prop("email", qrb.N("email")),
).From(qrb.N("users"))
```

```sql
SELECT json_build_object('id', id, 'name', name, 'email', email) 
FROM users
```

#### JSON with aggregation

```go
q := qrb.Select(
    qrb.N("department"),
    fn.JsonAgg(
        fn.JsonBuildObject().
            Prop("name", qrb.N("name")).
            Prop("salary", qrb.N("salary")),
    ).OrderBy(qrb.N("name")),
).From(qrb.N("employees")).
GroupBy(qrb.N("department"))
```

```sql
SELECT department, 
       json_agg(json_build_object('name', name, 'salary', salary) ORDER BY name)
FROM employees 
GROUP BY department
```

#### Complex nested JSON with CTEs

```go
q := qrb.With("author_json").As(
    qrb.Select(qrb.N("authors.author_id")).
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
LeftJoin(qrb.N("author_json")).On(qrb.N("posts.author_id").Eq(qrb.N("author_json.author_id")))
```

```sql
WITH author_json AS (
    SELECT authors.author_id, 
           json_build_object('id', authors.author_id, 'name', authors.name) AS json
    FROM authors
)
SELECT posts.post_id, 
       json_build_object('title', posts.title, 'author', author_json.json)
FROM posts
LEFT JOIN author_json ON posts.author_id = author_json.author_id
```

### Array Operations

#### Array construction

```go
q := qrb.Select(qrb.Array(qrb.String("a"), qrb.String("b"), qrb.String("c")))
```

```sql
SELECT ARRAY['a', 'b', 'c']
```

#### Array functions

```go
q := qrb.Select(
    fn.ArrayAppend(qrb.Array(qrb.Int(1), qrb.Int(2)), qrb.Int(3)),
    fn.ArrayLength(qrb.Array(qrb.Int(1), qrb.Int(2), qrb.Int(3)), qrb.Int(1)),
)
```

```sql
SELECT array_append(ARRAY[1, 2], 3), array_length(ARRAY[1, 2, 3], 1)
```

#### UNNEST

```go
q := qrb.Select(qrb.N("*")).
    From(fn.Unnest(qrb.Array(qrb.String("a"), qrb.String("b"), qrb.String("c")))).
    As("t").ColumnAliases("value")
```

```sql
SELECT * FROM unnest(ARRAY['a', 'b', 'c']) AS t (value)
```

#### Array aggregation

```go
q := qrb.Select(
    qrb.N("department"),
    fn.ArrayAgg(qrb.N("name")).OrderBy(qrb.N("name")),
).From(qrb.N("employees")).
GroupBy(qrb.N("department"))
```

```sql
SELECT department, array_agg(name ORDER BY name) 
FROM employees 
GROUP BY department
```

### Subqueries

#### EXISTS

```go
q := qrb.Select(qrb.N("name")).
    From(qrb.N("users")).
    Where(qrb.Exists(
        qrb.Select(qrb.Int(1)).
            From(qrb.N("posts")).
            Where(qrb.N("posts.user_id").Eq(qrb.N("users.id"))),
    ))
```

```sql
SELECT name FROM users 
WHERE EXISTS (SELECT 1 FROM posts WHERE posts.user_id = users.id)
```

#### IN with subquery

```go
q := qrb.Select(qrb.N("name")).
    From(qrb.N("users")).
    Where(qrb.N("id").In(
        qrb.Select(qrb.N("user_id")).
            From(qrb.N("posts")).
            Where(qrb.N("published").Eq(qrb.Bool(true))),
    ))
```

```sql
SELECT name FROM users 
WHERE id IN (SELECT user_id FROM posts WHERE published = true)
```

#### Correlated subquery

```go
q := qrb.Select(qrb.N("name"), qrb.N("salary")).
    From(qrb.N("employees")).As("e1").
    Where(qrb.N("salary").Gt(
        qrb.Select(fn.Avg(qrb.N("salary"))).
            From(qrb.N("employees")).As("e2").
            Where(qrb.N("e1.department").Eq(qrb.N("e2.department"))),
    ))
```

```sql
SELECT name, salary FROM employees AS e1 
WHERE salary > (
    SELECT avg(salary) FROM employees AS e2 
    WHERE e1.department = e2.department
)
```

#### Subquery in FROM

```go
q := qrb.Select(qrb.N("avg_salary")).
    From(
        qrb.Select(fn.Avg(qrb.N("salary"))).As("avg_salary").
            From(qrb.N("employees")).
            GroupBy(qrb.N("department")),
    ).As("dept_averages")
```

```sql
SELECT avg_salary FROM (
    SELECT avg(salary) AS avg_salary FROM employees GROUP BY department
) AS dept_averages
```

### Advanced Features

#### Common Table Expressions (WITH)

```go
q := qrb.With("recent_orders").As(
    qrb.Select(qrb.N("*")).
        From(qrb.N("orders")).
        Where(qrb.N("created_at").Gt(qrb.String("2023-01-01"))),
).
Select(qrb.N("customer_name"), fn.Count(qrb.N("*"))).
From(qrb.N("recent_orders")).
GroupBy(qrb.N("customer_name"))
```

```sql
WITH recent_orders AS (
    SELECT * FROM orders WHERE created_at > '2023-01-01'
)
SELECT customer_name, count(*) FROM recent_orders GROUP BY customer_name
```

#### Recursive CTE

```go
q := qrb.WithRecursive("employee_hierarchy").
    ColumnNames("employee_id", "name", "manager_id", "level").As(
    qrb.Select(qrb.N("employee_id"), qrb.N("name"), qrb.N("manager_id"), qrb.Int(1)).
        From(qrb.N("employees")).
        Where(qrb.N("manager_id").IsNull()).
        Union().All().
        Select(qrb.N("e.employee_id"), qrb.N("e.name"), qrb.N("e.manager_id"), qrb.N("eh.level").Plus(qrb.Int(1))).
        From(qrb.N("employees")).As("e").
        Join(qrb.N("employee_hierarchy")).As("eh").On(qrb.N("e.manager_id").Eq(qrb.N("eh.employee_id"))),
).
Select(qrb.N("*")).From(qrb.N("employee_hierarchy"))
```

```sql
WITH RECURSIVE employee_hierarchy (employee_id, name, manager_id, level) AS (
    SELECT employee_id, name, manager_id, 1 FROM employees WHERE manager_id IS NULL
    UNION ALL
    SELECT e.employee_id, e.name, e.manager_id, eh.level + 1
    FROM employees AS e
    JOIN employee_hierarchy AS eh ON e.manager_id = eh.employee_id
)
SELECT * FROM employee_hierarchy
```

#### ROWS FROM

```go
q := qrb.Select(qrb.N("*")).
    From(qrb.RowsFrom(
        fn.JsonToRecordset(qrb.String(`[{"name":"John","age":30},{"name":"Jane","age":25}]`)).
            ColumnDefinition("name", "TEXT").
            ColumnDefinition("age", "INTEGER"),
        fn.GenerateSeries(qrb.Int(1), qrb.Int(2)),
    ).WithOrdinality()).
    As("t").ColumnAliases("name", "age", "series_value", "ordinality")
```

```sql
SELECT * FROM ROWS FROM (
    json_to_recordset('[{"name":"John","age":30},{"name":"Jane","age":25}]') AS (name TEXT, age INTEGER),
    generate_series(1, 2)
) WITH ORDINALITY AS t (name, age, series_value, ordinality)
```

### Functions & Operators

#### String functions

```go
q := qrb.Select(
    fn.Upper(qrb.N("name")),
    fn.Lower(qrb.N("email")),
    fn.Initcap(qrb.N("title")),
).From(qrb.N("users"))
```

```sql
SELECT upper(name), lower(email), initcap(title) 
FROM users
```

#### Date/time functions

```go
q := qrb.Select(
    fn.Extract("year", qrb.N("created_at")),
    fn.Extract("month", qrb.N("created_at")),
    qrb.N("created_at").Plus(qrb.Interval("1 day")),
).From(qrb.N("orders"))
```

```sql
SELECT extract(year from created_at), extract(month from created_at), created_at + INTERVAL '1 day'
FROM orders
```

#### Mathematical operators

```go
q := qrb.Select(
    qrb.N("price").Op("*", qrb.N("quantity")).As("total"),
    qrb.N("price").Op("*", qrb.Float(1.08)).As("price_with_tax"),
).From(qrb.N("order_items"))
```

```sql
SELECT price * quantity AS total, price * 1.08 AS price_with_tax 
FROM order_items
```

#### CASE expressions

```go
q := qrb.Select(
    qrb.N("name"),
    qrb.Case().
        When(qrb.N("salary").Lt(qrb.Int(30000)), qrb.String("Low")).
        When(qrb.N("salary").Lt(qrb.Int(70000)), qrb.String("Medium")).
        Else(qrb.String("High")).
        As("salary_grade"),
).From(qrb.N("employees"))
```

```sql
SELECT name, 
       CASE 
           WHEN salary < 30000 THEN 'Low'
           WHEN salary < 70000 THEN 'Medium'
           ELSE 'High'
       END AS salary_grade
FROM employees
```

### Placeholders & Parameters

#### Named parameters

```go
q := qrb.Select(qrb.N("*")).
    From(qrb.N("users")).
    Where(qrb.And(
        qrb.N("name").Like(qrb.Bind("search_term")),
        qrb.N("active").Eq(qrb.Bind("is_active")),
    ))

sql, args, err := qrb.Build(q).
    WithNamedArgs(map[string]any{
        "search_term": "John%",
        "is_active":   true,
    }).
    ToSQL()
```

```sql
SELECT * FROM users WHERE name LIKE $1 AND active = $2
-- args: ["John%", true]
```

#### Positional parameters

```go
q := qrb.Select(qrb.N("*")).
    From(qrb.N("users")).
    Where(qrb.And(
        qrb.N("name").Like(qrb.Arg("John%")),
        qrb.N("active").Eq(qrb.Arg(true)),
    ))

sql, args, err := qrb.Build(q).ToSQL()
```

```sql
SELECT * FROM users WHERE name LIKE $1 AND active = $2
-- args: ["John%", true]
```

#### Mixing named and positional parameters

```go
q := qrb.Select(qrb.N("*")).
    From(qrb.N("users")).
    Where(qrb.And(
        qrb.N("name").Like(qrb.Bind("search_term")),
        qrb.N("active").Eq(qrb.Arg(true)),
    ))

sql, args, err := qrb.Build(q).
    WithNamedArgs(map[string]any{
        "search_term": "John%",
    }).
    ToSQL()
```

```sql
SELECT * FROM users WHERE name LIKE $1 AND active = $2
-- args: ["John%", true]
```

## Execution

### With pgx

```go
package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/networkteam/qrb/qrbpgx"
)

func main() {
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	q := qrb.Select(qrb.N("name"), qrb.N("email")).
		From(qrb.N("users")).
		Where(qrb.N("active").Eq(qrb.Bool(true)))

	rows, err := qrbpgx.Build(q).WithExecutor(pool).Query(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name, email string
		err := rows.Scan(&name, &email)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Name: %s, Email: %s\n", name, email)
	}
}
```

### With database/sql

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/qrbsql"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	q := qrb.Select(qrb.N("name"), qrb.N("email")).
		From(qrb.N("users")).
		Where(qrb.N("active").Eq(qrb.Bool(true)))

	rows, err := qrbsql.Build(q).WithExecutor(db).Query(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name, email string
		err := rows.Scan(&name, &email)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Name: %s, Email: %s\n", name, email)
	}
}

```

## Best Practices

### 1. Use Type-Safe Column References

```go
// Good - define reusable column references
var (
    UserID    = qrb.N("users.id")
    UserName  = qrb.N("users.name")
    UserEmail = qrb.N("users.email")
)

q := qrb.Select(UserName, UserEmail).From(qrb.N("users"))
```

### 2. Leverage Immutability for Reusable Queries

```go
// Base query that can be reused
baseQuery := qrb.Select(qrb.N("*")).From(qrb.N("users"))

// Create specific variations
activeUsers := baseQuery.Where(qrb.N("active").Eq(qrb.Bool(true)))
recentUsers := baseQuery.Where(qrb.N("created_at").Gt(qrb.String("2023-01-01")))
```

### 3. Use Named Parameters for Dynamic Queries

```go
q := qrb.Select(qrb.N("*")).
    From(qrb.N("users")).
    Where(qrb.N("name").Like(qrb.Bind("search")))

// Easy to reuse with different parameters
sql, args, _ := qrb.Build(q).WithNamedArgs(map[string]any{
    "search": "John%",
}).ToSQL()
```

### 4. Organize Complex Queries with CTEs

```go
// Break complex queries into readable parts
userData := qrb.Select(qrb.N("id"), qrb.N("name")).From(qrb.N("users"))
postData := qrb.Select(qrb.N("user_id"), fn.Count(qrb.N("*"))).From(qrb.N("posts")).GroupBy(qrb.N("user_id"))

q := qrb.With("user_data").As(userData).
    With("post_counts").As(postData).
    Select(qrb.N("ud.name"), qrb.N("pc.count")).
    From(qrb.N("user_data")).As("ud").
    LeftJoin(qrb.N("post_counts")).As("pc").Using("user_id")
```

### 5. Use JSON for Complex Data Structures

```go
// Build hierarchical data efficiently
q := qrb.Select(
    fn.JsonBuildObject().
        Prop("user", fn.JsonBuildObject().
            Prop("id", qrb.N("u.id")).
            Prop("name", qrb.N("u.name"))).
        Prop("posts", fn.JsonAgg(
            fn.JsonBuildObject().
                Prop("id", qrb.N("p.id")).
                Prop("title", qrb.N("p.title")))),
).From(qrb.N("users")).As("u").
LeftJoin(qrb.N("posts")).As("p").On(qrb.N("u.id").Eq(qrb.N("p.user_id"))).
GroupBy(qrb.N("u.id"), qrb.N("u.name"))
```

## License

[MIT](./LICENSE)