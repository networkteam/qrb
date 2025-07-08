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
- **Expressions**: Use `N()` for identifiers, `Arg()` for parameters, and `String()`, `Int()`, etc. for literals
- **Fluent API**: Chain method calls naturally following SQL structure
- **Type Safety**: Builders guide you through valid SQL construction with appropriate method availability

### Recommended Import Pattern

For the best development experience, we recommend using a dot import for the main qrb package:

```go
import (
    . "github.com/networkteam/qrb"
    "github.com/networkteam/qrb/fn"
)
```

This allows you to write clean, readable queries without the `qrb.` prefix:

```go
// With dot import (recommended)
q := Select(N("name")).From(N("users")).Where(N("active").Eq(Bool(true)))

// Without dot import (more verbose)
q := qrb.Select(qrb.N("name")).From(qrb.N("users")).Where(qrb.N("active").Eq(qrb.Bool(true)))
```

All examples in this README use the dot import pattern for improved readability.

## Examples

### Basic Queries

#### Simple SELECT

```go
q := Select(N("*")).From(N("users"))
```

```sql
SELECT * FROM users
```

#### SELECT with WHERE

```go
q := Select(N("name"), N("email")).
    From(N("users")).
    Where(N("active").Eq(Bool(true)))
```

```sql
SELECT name, email FROM users WHERE active = true
```

#### SELECT with multiple conditions

```go
q := Select(N("*")).
    From(N("employees")).
    Where(And(
        Or(
            N("firstname").ILike(Arg("John%")),
            N("lastname").ILike(Arg("John%")),
        ),
        N("active").Eq(Bool(true)),
    ))
```

```sql
SELECT * FROM employees 
WHERE ((firstname ILIKE $1) OR (lastname ILIKE $1)) AND (active = $2)
```

#### SELECT DISTINCT

```go
q := Select().Distinct().
    Select(N("department")).
    From(N("employees"))
```

```sql
SELECT DISTINCT department FROM employees
```

#### SELECT with ORDER BY and LIMIT

```go
q := Select(N("name"), N("salary")).
    From(N("employees")).
    OrderBy(N("salary")).Desc().NullsLast().
    Limit(Int(10)).
    Offset(Int(20))
```

```sql
SELECT name, salary FROM employees 
ORDER BY salary DESC NULLS LAST 
LIMIT 10 OFFSET 20
```

### CRUD Operations

#### INSERT with VALUES

```go
q := InsertInto(N("users")).
    ColumnNames("name", "email", "active").
    Values(String("John Doe"), String("john@example.com"), Bool(true))
```

```sql
INSERT INTO users (name, email, active) 
VALUES ('John Doe', 'john@example.com', true)
```

#### INSERT multiple rows

```go
q := InsertInto(N("products")).
    ColumnNames("name", "price", "category").
    Values(String("Laptop"), Float(999.99), String("Electronics")).
    Values(String("Book"), Float(19.99), String("Literature"))
```

```sql
INSERT INTO products (name, price, category) VALUES
    ('Laptop', 999.99, 'Electronics'),
    ('Book', 19.99, 'Literature')
```

#### INSERT with SELECT

```go
q := InsertInto(N("archived_users")).
    Query(Select(N("*")).From(N("users")).Where(N("active").Eq(Bool(false))))
```

```sql
INSERT INTO archived_users SELECT * FROM users WHERE active = false
```

#### INSERT with RETURNING

```go
q := InsertInto(N("users")).
    ColumnNames("name", "email").
    Values(String("Jane Doe"), String("jane@example.com")).
    Returning(N("id"), N("created_at"))
```

```sql
INSERT INTO users (name, email) VALUES ('Jane Doe', 'jane@example.com')
RETURNING id, created_at
```

#### UPSERT (INSERT with ON CONFLICT)

```go
q := InsertInto(N("users")).
    ColumnNames("email", "name").
    Values(String("john@example.com"), String("John Updated")).
    OnConflict(N("email")).DoUpdate().
    Set("name", N("EXCLUDED.name")).
    Set("updated_at", N("NOW()"))
```

```sql
INSERT INTO users (email, name) VALUES ('john@example.com', 'John Updated')
ON CONFLICT (email) DO UPDATE 
SET name = EXCLUDED.name, updated_at = NOW()
```

#### UPDATE

```go
q := Update(N("users")).
    Set("name", String("Updated Name")).
    Set("updated_at", N("NOW()")).
    Where(N("id").Eq(Arg(123)))
```

```sql
UPDATE users SET name = 'Updated Name', updated_at = NOW() WHERE id = $1
```

#### UPDATE with FROM

```go
q := Update(N("employees")).
    Set("department_name", N("d.name")).
    From(N("departments")).As("d").
    Where(N("employees.department_id").Eq(N("d.id")))
```

```sql
UPDATE employees SET department_name = d.name 
FROM departments AS d 
WHERE employees.department_id = d.id
```

#### DELETE

```go
q := DeleteFrom(N("users")).
    Where(N("active").Eq(Bool(false)))
```

```sql
DELETE FROM users WHERE active = false
```

#### DELETE with USING

```go
q := DeleteFrom(N("orders")).
    Using(N("customers")).
    Where(And(
        N("orders.customer_id").Eq(N("customers.id")),
        N("customers.status").Eq(String("inactive")),
    ))
```

```sql
DELETE FROM orders USING customers 
WHERE orders.customer_id = customers.id AND customers.status = 'inactive'
```

### Joins

#### INNER JOIN

```go
q := Select(N("u.name"), N("p.title")).
    From(N("users")).As("u").
    Join(N("posts")).As("p").On(N("u.id").Eq(N("p.user_id")))
```

```sql
SELECT u.name, p.title FROM users AS u 
JOIN posts AS p ON u.id = p.user_id
```

#### LEFT JOIN

```go
q := Select(N("u.name"), N("p.title")).
    From(N("users")).As("u").
    LeftJoin(N("posts")).As("p").On(N("u.id").Eq(N("p.user_id")))
```

```sql
SELECT u.name, p.title FROM users AS u 
LEFT JOIN posts AS p ON u.id = p.user_id
```

#### JOIN with USING

```go
q := Select(N("u.name"), N("p.title")).
    From(N("users")).As("u").
    Join(N("posts")).As("p").Using("user_id")
```

```sql
SELECT u.name, p.title FROM users AS u 
JOIN posts AS p USING (user_id)
```

#### Multiple JOINs

```go
q := Select(N("u.name"), N("p.title"), N("c.name")).
    From(N("users")).As("u").
    Join(N("posts")).As("p").On(N("u.id").Eq(N("p.user_id"))).
    Join(N("categories")).As("c").On(N("p.category_id").Eq(N("c.id")))
```

```sql
SELECT u.name, p.title, c.name FROM users AS u 
JOIN posts AS p ON u.id = p.user_id 
JOIN categories AS c ON p.category_id = c.id
```

### Aggregation & Grouping

#### GROUP BY with aggregate functions

```go
q := Select(N("department")).
    Select(fn.Count(N("*"))).As("employee_count").
    From(N("employees")).
    GroupBy(N("department"))
```

```sql
SELECT department, count(*) AS employee_count 
FROM employees 
GROUP BY department
```

#### GROUP BY with HAVING

```go
q := Select(N("department")).
    Select(fn.Avg(N("salary"))).As("avg_salary").
    From(N("employees")).
    GroupBy(N("department")).
    Having(fn.Avg(N("salary")).Gt(Int(50000)))
```

```sql
SELECT department, avg(salary) AS avg_salary 
FROM employees 
GROUP BY department 
HAVING avg(salary) > 50000
```

#### GROUP BY with ROLLUP

```go
q := Select(N("department"), N("job_title"), fn.Sum(N("salary"))).
    From(N("employees")).
    GroupBy().
    Rollup(
        Exps(N("department")),
        Exps(N("job_title")),
    )
```

```sql
SELECT department, job_title, sum(salary) 
FROM employees 
GROUP BY ROLLUP (department, job_title)
```

#### GROUP BY with GROUPING SETS

```go
q := Select(N("department"), N("job_title"), fn.Sum(N("salary"))).
    From(N("employees")).
    GroupBy().
    GroupingSets(
        Exps(N("department")),
        Exps(N("job_title")),
        Exps(),
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
q := Select(
    N("name"),
    N("salary"),
    fn.RowNumber().Over().PartitionBy(N("department")).OrderBy(N("salary")).Desc(),
).From(N("employees"))
```

```sql
SELECT name, salary, row_number() OVER (PARTITION BY department ORDER BY salary DESC) 
FROM employees
```

#### RANK and DENSE_RANK

```go
q := Select(
    N("name"),
    N("salary"),
    fn.Rank().Over().PartitionBy(N("department")).OrderBy(N("salary")).Desc(),
    fn.DenseRank().Over().PartitionBy(N("department")).OrderBy(N("salary")).Desc(),
).From(N("employees"))
```

```sql
SELECT name, salary, 
       rank() OVER (PARTITION BY department ORDER BY salary DESC),
       dense_rank() OVER (PARTITION BY department ORDER BY salary DESC)
FROM employees
```

#### Named Windows

```go
q := Select(
    fn.Sum(N("salary")).Over("w"),
    fn.Avg(N("salary")).Over("w"),
).From(N("employees")).
Window("w").As().PartitionBy(N("department")).OrderBy(N("salary")).Desc().
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
q := Select(
    fn.JsonBuildObject().
        Prop("id", N("id")).
        Prop("name", N("name")).
        Prop("email", N("email")),
).From(N("users"))
```

```sql
SELECT json_build_object('id', id, 'name', name, 'email', email) 
FROM users
```

#### JSON with aggregation

```go
q := Select(
    N("department"),
    fn.JsonAgg(
        fn.JsonBuildObject().
            Prop("name", N("name")).
            Prop("salary", N("salary")),
    ).OrderBy(N("name")),
).From(N("employees")).
GroupBy(N("department"))
```

```sql
SELECT department, 
       json_agg(json_build_object('name', name, 'salary', salary) ORDER BY name)
FROM employees 
GROUP BY department
```

#### Complex nested JSON with CTEs

```go
q := With("author_json").As(
    Select(N("authors.author_id")).
        Select(
            fn.JsonBuildObject().
                Prop("id", N("authors.author_id")).
                Prop("name", N("authors.name")),
        ).As("json").
        From(N("authors")),
).
Select(
    N("posts.post_id"),
    fn.JsonBuildObject().
        Prop("title", N("posts.title")).
        Prop("author", N("author_json.json")),
).
From(N("posts")).
LeftJoin(N("author_json")).On(N("posts.author_id").Eq(N("author_json.author_id")))
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
q := Select(Array(String("a"), String("b"), String("c")))
```

```sql
SELECT ARRAY['a', 'b', 'c']
```

#### Array functions

```go
q := Select(
    fn.ArrayAppend(Array(Int(1), Int(2)), Int(3)),
    fn.ArrayLength(Array(Int(1), Int(2), Int(3)), Int(1)),
)
```

```sql
SELECT array_append(ARRAY[1, 2], 3), array_length(ARRAY[1, 2, 3], 1)
```

#### UNNEST

```go
q := Select(N("*")).
    From(fn.Unnest(Array(String("a"), String("b"), String("c")))).
    As("t").ColumnAliases("value")
```

```sql
SELECT * FROM unnest(ARRAY['a', 'b', 'c']) AS t (value)
```

#### Array aggregation

```go
q := Select(
    N("department"),
    fn.ArrayAgg(N("name")).OrderBy(N("name")),
).From(N("employees")).
GroupBy(N("department"))
```

```sql
SELECT department, array_agg(name ORDER BY name) 
FROM employees 
GROUP BY department
```

### Subqueries

#### EXISTS

```go
q := Select(N("name")).
    From(N("users")).
    Where(Exists(
        Select(Int(1)).
            From(N("posts")).
            Where(N("posts.user_id").Eq(N("users.id"))),
    ))
```

```sql
SELECT name FROM users 
WHERE EXISTS (SELECT 1 FROM posts WHERE posts.user_id = users.id)
```

#### IN with subquery

```go
q := Select(N("name")).
    From(N("users")).
    Where(N("id").In(
        Select(N("user_id")).
            From(N("posts")).
            Where(N("published").Eq(Bool(true))),
    ))
```

```sql
SELECT name FROM users 
WHERE id IN (SELECT user_id FROM posts WHERE published = true)
```

#### Correlated subquery

```go
q := Select(N("name"), N("salary")).
    From(N("employees")).As("e1").
    Where(N("salary").Gt(
        Select(fn.Avg(N("salary"))).
            From(N("employees")).As("e2").
            Where(N("e1.department").Eq(N("e2.department"))),
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
q := Select(N("avg_salary")).
    From(
        Select(fn.Avg(N("salary"))).As("avg_salary").
            From(N("employees")).
            GroupBy(N("department")),
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
q := With("recent_orders").As(
    Select(N("*")).
        From(N("orders")).
        Where(N("created_at").Gt(String("2023-01-01"))),
).
Select(N("customer_name"), fn.Count(N("*"))).
From(N("recent_orders")).
GroupBy(N("customer_name"))
```

```sql
WITH recent_orders AS (
    SELECT * FROM orders WHERE created_at > '2023-01-01'
)
SELECT customer_name, count(*) FROM recent_orders GROUP BY customer_name
```

#### Recursive CTE

```go
q := WithRecursive("employee_hierarchy").
    ColumnNames("employee_id", "name", "manager_id", "level").As(
    Select(N("employee_id"), N("name"), N("manager_id"), Int(1)).
        From(N("employees")).
        Where(N("manager_id").IsNull()).
        Union().All().
        Select(N("e.employee_id"), N("e.name"), N("e.manager_id"), N("eh.level").Plus(Int(1))).
        From(N("employees")).As("e").
        Join(N("employee_hierarchy")).As("eh").On(N("e.manager_id").Eq(N("eh.employee_id"))),
).
Select(N("*")).From(N("employee_hierarchy"))
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
q := Select(N("*")).
    From(RowsFrom(
        fn.JsonToRecordset(String(`[{"name":"John","age":30},{"name":"Jane","age":25}]`)).
            ColumnDefinition("name", "TEXT").
            ColumnDefinition("age", "INTEGER"),
        fn.GenerateSeries(Int(1), Int(2)),
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
q := Select(
    fn.Upper(N("name")),
    fn.Lower(N("email")),
    fn.Initcap(N("title")),
).From(N("users"))
```

```sql
SELECT upper(name), lower(email), initcap(title) 
FROM users
```

#### Date/time functions

```go
q := Select(
    fn.Extract("year", N("created_at")),
    fn.Extract("month", N("created_at")),
    N("created_at").Plus(Interval("1 day")),
).From(N("orders"))
```

```sql
SELECT extract(year from created_at), extract(month from created_at), created_at + INTERVAL '1 day'
FROM orders
```

#### Mathematical operators

```go
q := Select(
    N("price").Op("*", N("quantity")).As("total"),
    N("price").Op("*", Float(1.08)).As("price_with_tax"),
).From(N("order_items"))
```

```sql
SELECT price * quantity AS total, price * 1.08 AS price_with_tax 
FROM order_items
```

#### CASE expressions

```go
q := Select(
    N("name"),
    Case().
        When(N("salary").Lt(Int(30000)), String("Low")).
        When(N("salary").Lt(Int(70000)), String("Medium")).
        Else(String("High")).
        As("salary_grade"),
).From(N("employees"))
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
q := Select(N("*")).
    From(N("users")).
    Where(And(
        N("name").Like(Bind("search_term")),
        N("active").Eq(Bind("is_active")),
    ))

sql, args, err := Build(q).
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
q := Select(N("*")).
    From(N("users")).
    Where(And(
        N("name").Like(Arg("John%")),
        N("active").Eq(Arg(true)),
    ))

sql, args, err := Build(q).ToSQL()
```

```sql
SELECT * FROM users WHERE name LIKE $1 AND active = $2
-- args: ["John%", true]
```

#### Mixing named and positional parameters

```go
q := Select(N("*")).
    From(N("users")).
    Where(And(
        N("name").Like(Bind("search_term")),
        N("active").Eq(Arg(true)),
    ))

sql, args, err := Build(q).
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
	. "github.com/networkteam/qrb"
	"github.com/networkteam/qrb/qrbpgx"
)

func main() {
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	q := Select(N("name"), N("email")).
		From(N("users")).
		Where(N("active").Eq(Bool(true)))

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
	. "github.com/networkteam/qrb"
	"github.com/networkteam/qrb/qrbsql"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	q := Select(N("name"), N("email")).
		From(N("users")).
		Where(N("active").Eq(Bool(true)))

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
    UserID    = N("users.id")
    UserName  = N("users.name")
    UserEmail = N("users.email")
)

q := Select(UserName, UserEmail).From(N("users"))
```

### 2. Leverage Immutability for Reusable Queries

```go
// Base query that can be reused
baseQuery := Select(N("*")).From(N("users"))

// Create specific variations
activeUsers := baseQuery.Where(N("active").Eq(Bool(true)))
recentUsers := baseQuery.Where(N("created_at").Gt(String("2023-01-01")))
```

### 3. Use Named Parameters for Dynamic Queries

```go
q := Select(N("*")).
    From(N("users")).
    Where(N("name").Like(Bind("search")))

// Easy to reuse with different parameters
sql, args, _ := Build(q).WithNamedArgs(map[string]any{
    "search": "John%",
}).ToSQL()
```

### 4. Organize Complex Queries with CTEs

```go
// Break complex queries into readable parts
userData := Select(N("id"), N("name")).From(N("users"))
postData := Select(N("user_id"), fn.Count(N("*"))).From(N("posts")).GroupBy(N("user_id"))

q := With("user_data").As(userData).
    With("post_counts").As(postData).
    Select(N("ud.name"), N("pc.count")).
    From(N("user_data")).As("ud").
    LeftJoin(N("post_counts")).As("pc").Using("user_id")
```

### 5. Use JSON for Complex Data Structures

```go
// Build hierarchical data efficiently
q := Select(
    fn.JsonBuildObject().
        Prop("user", fn.JsonBuildObject().
            Prop("id", N("u.id")).
            Prop("name", N("u.name"))).
        Prop("posts", fn.JsonAgg(
            fn.JsonBuildObject().
                Prop("id", N("p.id")).
                Prop("title", N("p.title")))),
).From(N("users")).As("u").
LeftJoin(N("posts")).As("p").On(N("u.id").Eq(N("p.user_id"))).
GroupBy(N("u.id"), N("u.name"))
```

## License

[MIT](./LICENSE)