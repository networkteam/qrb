package qrb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/builder"
	"github.com/networkteam/qrb/fn"
	"github.com/networkteam/qrb/internal/testhelper"
)

func TestSelectBuilder(t *testing.T) {
	t.Run("with / json", func(t *testing.T) {
		myCategory := "SQL Hacks"

		q := qrb.With("author_json").As(
			qrb.
				Select(
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
			LeftJoin(
				qrb.N("author_json"),
			).On(qrb.N("posts.author_id").Eq(qrb.N("author_json.author_id"))).
			Where(qrb.N("posts.category").Eq(qrb.Arg(myCategory))).
			OrderBy(qrb.N("posts.created_at")).Desc().NullsLast()

		testhelper.AssertSQLWriterEquals(
			t,
			// language=PostgreSQL
			`
			WITH author_json AS (
				SELECT
					authors.author_id,
					json_build_object('id', authors.author_id, 'name', authors.name) AS json
				FROM
					authors
			)
			SELECT
				posts.post_id,
				json_build_object('title', posts.title, 'author', author_json.json)
			FROM
				posts
				LEFT JOIN author_json ON posts.author_id = author_json.author_id
			WHERE
				posts.category = $1
			ORDER BY
				posts.created_at DESC NULLS LAST
			`,
			[]any{myCategory},
			q,
		)
	})

	t.Run("complex nested JSON with CTEs", func(t *testing.T) {
		bookJSON := fn.JsonBuildObject().
			Prop("Title", qrb.N("books.title")).
			Prop("AuthorID", qrb.N("books.author_id")).
			Prop("PublicationYear", qrb.N("books.publication_year")).
			Prop("CreatedAt", qrb.N("books.created_at")).
			Prop("UpdatedAt", qrb.N("books.updated_at")).
			Prop("ID", qrb.N("books.book_id"))

		authorJSON := fn.JsonBuildObject().
			Prop("AuthorID", qrb.N("authors.author_id")).
			Prop("Name", qrb.N("authors.name"))

		genreJSON := fn.JsonBuildObject().
			Prop("GenreID", qrb.N("genres.genre_id")).
			Prop("Name", qrb.N("genres.name"))

		// Book and genre are joined via book_genre table

		type authorQueryOpts struct {
			IncludeBooks bool
		}

		type bookQueryOpts struct {
			IncludeGenres bool
			IncludeAuthor bool
			AuthorOpts    authorQueryOpts
		}

		opts := bookQueryOpts{
			IncludeGenres: true,
			IncludeAuthor: true,
			AuthorOpts: authorQueryOpts{
				IncludeBooks: true,
			},
		}

		q := qrb.
			SelectJson(bookJSON).
			From(qrb.N("books")).
			LeftJoin(qrb.N("authors")).Using("author_id").
			Where(qrb.N("books.book_id").Eq(qrb.Arg(2)))

		if opts.IncludeAuthor {
			// See below for a nicer way to do this via subselects
			q = q.ApplySelectJson(func(obj builder.JsonBuildObjectBuilder) builder.JsonBuildObjectBuilder {
				return obj.Prop("Author", authorJSON.
					PropIf(opts.AuthorOpts.IncludeBooks, "Books", qrb.N("author_books.books")),
				)
			}).SelectBuilder

			if opts.AuthorOpts.IncludeBooks {
				q = q.AppendWith(qrb.With("author_books").As(
					qrb.Select(qrb.N("author_id")).
						Select(
							qrb.Coalesce(
								fn.JsonAgg(bookJSON).OrderBy(qrb.N("publication_year")),
								qrb.String("[]"),
							),
						).As("books").
						From(qrb.N("books")).
						GroupBy(qrb.N("author_id")),
				)).
					LeftJoin(qrb.N("author_books")).Using("author_id")
			}
		}

		if opts.IncludeGenres {
			q = q.AppendWith(qrb.With("book_genres").As(
				qrb.Select(qrb.N("book_id")).
					Select(
						qrb.Coalesce(
							fn.JsonAgg(genreJSON).OrderBy(qrb.N("name")),
							qrb.String("[]"),
						),
					).As("genres").
					From(qrb.N("book_genre")).
					Join(qrb.N("genres")).Using("genre_id").
					GroupBy(qrb.N("book_id")),
			)).
				ApplySelectJson(func(obj builder.JsonBuildObjectBuilder) builder.JsonBuildObjectBuilder {
					return obj.Prop("Genres", qrb.N("book_genres.genres"))
				}).
				LeftJoin(qrb.N("book_genres")).Using("book_id")
		}

		testhelper.AssertSQLWriterEquals(
			t,
			// language=PostgreSQL
			`
			WITH author_books AS (SELECT author_id,
										 COALESCE(json_agg(json_build_object('Title', books.title, 'AuthorID', books.author_id,
																			 'PublicationYear', books.publication_year, 'CreatedAt',
																			 books.created_at, 'UpdatedAt', books.updated_at, 'ID',
																			 books.book_id) ORDER BY publication_year),
												  '[]') AS books
								  FROM books
								  GROUP BY author_id),
				 book_genres AS (SELECT book_id,
										COALESCE(json_agg(json_build_object('GenreID', genres.genre_id, 'Name', genres.name)
														  ORDER BY name), '[]') AS genres
								 FROM book_genre
										  JOIN genres USING (genre_id)
								 GROUP BY book_id)
			SELECT json_build_object('Title', books.title, 'AuthorID', books.author_id, 'PublicationYear', books.publication_year,
									 'CreatedAt', books.created_at, 'UpdatedAt', books.updated_at, 'ID', books.book_id, 'Author',
									 json_build_object('AuthorID', authors.author_id, 'Name', authors.name, 'Books',
													   author_books.books), 'Genres', book_genres.genres)
			FROM books
					 LEFT JOIN authors USING (author_id)
					 LEFT JOIN author_books USING (author_id)
					 LEFT JOIN book_genres USING (book_id)
			WHERE books.book_id = $1
    		`,
			[]any{2},
			q,
		)
	})

	t.Run("complex nested JSON with subselects", func(t *testing.T) {
		type authorQueryOpts struct {
			IncludeBooks bool
		}

		type bookQueryOpts struct {
			IncludeGenres bool
			IncludeAuthor bool
			AuthorOpts    authorQueryOpts
		}

		genreJSON := fn.JsonBuildObject().
			Prop("GenreID", qrb.N("genres.genre_id")).
			Prop("Name", qrb.N("genres.name"))

		baseBookJSON := fn.JsonBuildObject().
			Prop("Title", qrb.N("books.title")).
			Prop("AuthorID", qrb.N("books.author_id")).
			Prop("PublicationYear", qrb.N("books.publication_year")).
			Prop("CreatedAt", qrb.N("books.created_at")).
			Prop("UpdatedAt", qrb.N("books.updated_at")).
			Prop("ID", qrb.N("books.book_id"))

		selectAuthorBooks := qrb.
			Select(qrb.Coalesce(fn.JsonAgg(baseBookJSON).OrderBy(qrb.N("books.publication_year")), qrb.String("[]"))).
			From(qrb.N("books")).
			Where(qrb.N("books.author_id").Eq(qrb.N("authors.author_id")))

		buildAuthorJSON := func(opts authorQueryOpts) builder.JsonBuildObjectBuilder {
			return fn.JsonBuildObject().
				Prop("AuthorID", qrb.N("authors.author_id")).
				Prop("Name", qrb.N("authors.name")).
				PropIf(opts.IncludeBooks, "Books", selectAuthorBooks)
		}

		selectAuthors := func(opts authorQueryOpts) builder.SelectBuilder {
			return qrb.
				SelectJson(buildAuthorJSON(opts)).
				From(qrb.N("authors")).
				SelectBuilder
		}

		selectBookAuthor := func(opts bookQueryOpts) builder.SelectBuilder {
			return selectAuthors(opts.AuthorOpts).Where(qrb.N("authors.author_id").Eq(qrb.N("books.author_id")))
		}

		buildBookJSON := func(opts bookQueryOpts) builder.JsonBuildObjectBuilder {
			return baseBookJSON.
				PropIf(opts.IncludeAuthor, "Author", selectBookAuthor(opts)).
				ApplyIf(opts.IncludeGenres, func(b builder.JsonBuildObjectBuilder) builder.JsonBuildObjectBuilder {
					return b.Prop("Genres", qrb.Select(qrb.Coalesce(fn.JsonAgg(genreJSON).OrderBy(qrb.N("genres.name")), qrb.String("[]"))).
						From(qrb.N("book_genre")).
						LeftJoin(qrb.N("genres")).Using("genre_id").
						Where(qrb.N("book_genre.book_id").Eq(qrb.N("books.book_id"))))
				})
		}

		selectBook := func(opts bookQueryOpts) builder.SelectBuilder {
			return qrb.
				SelectJson(buildBookJSON(opts)).
				From(qrb.N("books")).
				SelectBuilder
		}

		opts := bookQueryOpts{
			IncludeGenres: true,
			IncludeAuthor: true,
			AuthorOpts: authorQueryOpts{
				IncludeBooks: true,
			},
		}

		t.Run("with all options", func(t *testing.T) {
			q := selectBook(opts).
				Where(qrb.N("books.book_id").Eq(qrb.Arg(2)))

			testhelper.AssertSQLWriterEquals(
				t,
				// language=PostgreSQL
				`
				SELECT json_build_object(
							   'Title', books.title,
							   'AuthorID', books.author_id,
							   'PublicationYear', books.publication_year,
							   'CreatedAt', books.created_at,
							   'UpdatedAt', books.updated_at,
							   'ID', books.book_id,
							   'Author', (SELECT json_build_object(
														 'AuthorID', authors.author_id,
														 'Name', authors.name,
														 'Books', (SELECT COALESCE(
																				  json_agg(
																						  json_build_object(
																								  'Title', books.title,
																								  'AuthorID', books.author_id,
																								  'PublicationYear',
																								  books.publication_year,
																								  'CreatedAt', books.created_at,
																								  'UpdatedAt', books.updated_at,
																								  'ID', books.book_id
																							  )
																						  ORDER BY books.publication_year),
																				  '[]'
																			  )
																   FROM books
																   WHERE books.author_id = authors.author_id)
													 )
										  FROM authors
										  WHERE authors.author_id = books.author_id),
							   'Genres', (SELECT COALESCE(
														 json_agg(
																 json_build_object(
																		 'GenreID', genres.genre_id,
																		 'Name', genres.name
																	 )
																 ORDER BY genres.name),
														 '[]'
													 )
										  FROM book_genre
												   LEFT JOIN genres USING (genre_id)
										  WHERE book_genre.book_id = books.book_id)
						   )
				FROM books
				WHERE books.book_id = $1
				`,
				[]any{2},
				q,
			)
		})

		t.Run("without options", func(t *testing.T) {
			q := selectBook(bookQueryOpts{}).
				Where(qrb.N("books.book_id").Eq(qrb.Arg(2)))

			testhelper.AssertSQLWriterEquals(
				t,
				// language=PostgreSQL
				`
				SELECT json_build_object('Title', books.title, 'AuthorID', books.author_id, 'PublicationYear', books.publication_year,
                         'CreatedAt', books.created_at, 'UpdatedAt', books.updated_at, 'ID', books.book_id)
				FROM books
				WHERE books.book_id = $1
    			`,
				[]any{2},
				q,
			)
		})

		t.Run("with modified JSON selection", func(t *testing.T) {
			q := selectBook(bookQueryOpts{}).
				ApplySelectJson(func(obj builder.JsonBuildObjectBuilder) builder.JsonBuildObjectBuilder {
					return obj.
						Unset("CreatedAt").
						Unset("UpdatedAt")
				}).
				OrderBy(qrb.N("books.publication_year")).
				Limit(qrb.Int(10)).
				Offset(qrb.Arg(5))

			testhelper.AssertSQLWriterEquals(
				t,
				// language=PostgreSQL
				`
				SELECT json_build_object('Title',books.title,'AuthorID',books.author_id,'PublicationYear',books.publication_year,'ID',books.book_id) FROM books ORDER BY books.publication_year LIMIT 10 OFFSET $1
    			`,
				[]any{5},
				q,
			)
		})
	})

	// These examples ar taken from https://www.postgresql.org/docs/14/sql-select.html#id-1.9.3.171.9
	t.Run("examples", func(t *testing.T) {
		t.Run("example 1", func(t *testing.T) {
			q := qrb.Select(qrb.N("f.title"), qrb.N("f.did"), qrb.N("d.name"), qrb.N("f.date_prod"), qrb.N("f.kind")).
				From(qrb.N("distributors")).As("d").Join(qrb.N("films")).As("f").Using("did")

			testhelper.AssertSQLWriterEquals(
				t,
				// language=PostgreSQL
				`
				SELECT f.title, f.did, d.name, f.date_prod, f.kind
					FROM distributors AS d JOIN films AS f USING (did)
				`,
				nil,
				q,
			)
		})

		t.Run("example 2", func(t *testing.T) {
			q := qrb.Select(qrb.N("kind")).Select(fn.Sum(qrb.N("len"))).As("total").From(qrb.N("films")).GroupBy(qrb.N("kind"))

			testhelper.AssertSQLWriterEquals(
				t,
				// language=PostgreSQL
				`
				SELECT kind, sum(len) AS total FROM films GROUP BY kind
				`,
				nil,
				q,
			)
		})

		t.Run("example 3", func(t *testing.T) {
			q := qrb.Select(qrb.N("kind")).Select(fn.Sum(qrb.N("len"))).As("total").
				From(qrb.N("films")).
				GroupBy(qrb.N("kind")).
				Having(fn.Sum(qrb.N("len")).Lt(qrb.Interval("5 hours")))

			testhelper.AssertSQLWriterEquals(
				t,
				// language=PostgreSQL
				`
				SELECT kind, sum(len) AS total
				FROM films
				GROUP BY kind
				HAVING sum(len) < INTERVAL '5 hours'
				`,
				nil,
				q,
			)
		})

		t.Run("example 4", func(t *testing.T) {
			t.Run("query 1", func(t *testing.T) {
				q := qrb.Select(qrb.N("*")).
					From(qrb.N("distributors")).
					OrderBy(qrb.N("name"))

				testhelper.AssertSQLWriterEquals(
					t,
					// language=PostgreSQL
					`
					SELECT * FROM distributors ORDER BY name
					`,
					nil,
					q,
				)
			})

			t.Run("query 2", func(t *testing.T) {
				q := qrb.Select(qrb.N("*")).
					From(qrb.N("distributors")).
					OrderBy(qrb.Int(2))

				testhelper.AssertSQLWriterEquals(
					t,
					// language=PostgreSQL
					`
					SELECT * FROM distributors ORDER BY 2
					`,
					nil,
					q,
				)
			})
		})

		t.Run("example 5", func(t *testing.T) {
			q := qrb.Select(qrb.N("distributors.name")).
				From(qrb.N("distributors")).
				Where(qrb.N("distributors.name").Like(qrb.String("W%"))).
				Union().
				Select(qrb.N("actors.name")).
				From(qrb.N("actors")).
				Where(qrb.N("actors.name").Like(qrb.String("W%")))

			testhelper.AssertSQLWriterEquals(
				t,
				// language=PostgreSQL
				`
				SELECT distributors.name
					FROM distributors
					WHERE distributors.name LIKE 'W%'
				UNION
				SELECT actors.name
					FROM actors
					WHERE actors.name LIKE 'W%'
				`,
				nil,
				q,
			)
		})

		t.Run("example 6", func(t *testing.T) {
			t.Run("query 1", func(t *testing.T) {
				q := qrb.Select(qrb.N("*")).
					From(qrb.Func("distributors", qrb.Int(111)))

				testhelper.AssertSQLWriterEquals(
					t,
					// language=PostgreSQL
					`
					SELECT * FROM distributors(111)
					`,
					nil,
					q,
				)
			})

			t.Run("query 2", func(t *testing.T) {
				q := qrb.Select(qrb.N("*")).
					From(
						qrb.Func("distributors_2", qrb.Int(111)).
							As("d").
							ColumnDefinition("f1", "int").
							ColumnDefinition("f2", "text"),
					)

				testhelper.AssertSQLWriterEquals(
					t,
					// language=PostgreSQL
					`
					SELECT * FROM distributors_2(111) AS d (f1 int, f2 text)
					`,
					nil,
					q,
				)
			})
		})

		t.Run("example 7", func(t *testing.T) {
			q := qrb.Select(qrb.N("*")).
				From(qrb.Func("unnest", qrb.Array(
					qrb.String("a"),
					qrb.String("b"),
					qrb.String("c"),
					qrb.String("d"),
					qrb.String("e"),
					qrb.String("f"),
				)).WithOrdinality())

			testhelper.AssertSQLWriterEquals(
				t,
				// language=PostgreSQL
				`
				SELECT * FROM unnest(ARRAY['a','b','c','d','e','f']) WITH ORDINALITY
				`,
				nil,
				q,
			)
		})

		t.Run("example 7b", func(t *testing.T) {
			q := qrb.Select(qrb.N("*")).
				From(qrb.Func("unnest", qrb.Array(
					qrb.String("a"),
					qrb.String("b"),
					qrb.String("c"),
					qrb.String("d"),
					qrb.String("e"),
					qrb.String("f"),
				)).WithOrdinality().ColumnDefinition("x", "text"))

			// From the docs: To use ORDINALITY together with a column definition list, you must use the ROWS FROM( ... ) syntax and put the column definition list inside ROWS FROM( ... ).

			_, _, err := qrb.Build(q).ToSQL()
			require.Error(t, err)
		})

		t.Run("example 8", func(t *testing.T) {
			q := qrb.With("t").As(
				qrb.Select(qrb.Func("random")).As("x").
					From(qrb.Func("generate_series", qrb.Int(1), qrb.Int(3))),
			).
				Select(qrb.N("*")).From(qrb.N("t")).
				Union().All().
				Select(qrb.N("*")).From(qrb.N("t"))

			testhelper.AssertSQLWriterEquals(
				t,
				// language=PostgreSQL
				`
				WITH t AS (
					SELECT random() AS x FROM generate_series(1, 3)
				)
				SELECT * FROM t
				UNION ALL
				SELECT * FROM t
				`,
				nil,
				q,
			)
		})

		t.Run("example 9", func(t *testing.T) {
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

			testhelper.AssertSQLWriterEquals(
				t,
				// language=PostgreSQL
				`
				WITH RECURSIVE employee_recursive(distance, employee_name, manager_name) AS (
					SELECT 1, employee_name, manager_name
					FROM employee
					WHERE manager_name = 'Mary'
				  UNION ALL
					SELECT er.distance + 1, e.employee_name, e.manager_name
					FROM employee_recursive AS er, employee AS e
					WHERE er.employee_name = e.manager_name
				  )
				SELECT distance, employee_name FROM employee_recursive
				`,
				nil,
				q,
			)
		})

		t.Run("example 10", func(t *testing.T) {
			q := qrb.Select(qrb.N("m.name")).As("mname").Select(qrb.N("pname")).
				From(qrb.N("manufacturers")).As("m").
				FromLateral(qrb.Func("get_product_names", qrb.N("m.id"))).As("pname")

			testhelper.AssertSQLWriterEquals(
				t,
				// language=PostgreSQL
				`
				SELECT m.name AS mname, pname
				FROM manufacturers AS m, LATERAL get_product_names(m.id) AS pname
				`,
				nil,
				q,
			)
		})

		t.Run("example 11", func(t *testing.T) {
			q := qrb.Select(qrb.N("m.name")).As("mname").Select(qrb.N("pname")).
				From(qrb.N("manufacturers")).As("m").
				LeftJoinLateral(qrb.Func("get_product_names", qrb.N("m.id"))).As("pname").On(qrb.Bool(true))

			testhelper.AssertSQLWriterEquals(
				t,
				// language=PostgreSQL
				`
				SELECT m.name AS mname, pname
				FROM manufacturers AS m LEFT JOIN LATERAL get_product_names(m.id) AS pname ON true
				`,
				nil,
				q,
			)
		})
	})
}

func TestSelectBuilder_From(t *testing.T) {
	t.Run("table", func(t *testing.T) {
		q1 := qrb.Select(qrb.Int(1)).From(qrb.N("foo"))
		q2 := q1.From(qrb.N("bar"))

		testhelper.AssertSQLWriterEquals(t, "SELECT 1 FROM foo", nil, q1)

		testhelper.AssertSQLWriterEquals(t, "SELECT 1 FROM foo,bar", nil, q2)
	})

	t.Run("only table", func(t *testing.T) {
		q1 := qrb.Select(qrb.Int(1)).FromOnly(qrb.N("foo"))
		q2 := q1.From(qrb.N("bar"))

		testhelper.AssertSQLWriterEquals(t, "SELECT 1 FROM ONLY foo", nil, q1)

		testhelper.AssertSQLWriterEquals(t, "SELECT 1 FROM ONLY foo,bar", nil, q2)
	})

	t.Run("select", func(t *testing.T) {
		q := qrb.Select(qrb.N("avg_quantity")).From(
			qrb.Select(fn.Avg(qrb.N("quantity"))).As("avg_quantity").From(qrb.N("sales")).GroupBy(qrb.N("brand")),
		).As("t")

		testhelper.AssertSQLWriterEquals(
			t,
			`SELECT avg_quantity FROM (SELECT avg(quantity) AS avg_quantity FROM sales GROUP BY brand) AS t`,
			nil,
			q,
		)
	})

	t.Run("lateral select", func(t *testing.T) {
		q := qrb.Select(qrb.N("avg_quantity")).FromLateral(
			qrb.Select(fn.Avg(qrb.N("quantity"))).As("avg_quantity").From(qrb.N("sales")).GroupBy(qrb.N("brand")),
		).As("t")

		testhelper.AssertSQLWriterEquals(
			t,
			`SELECT avg_quantity FROM LATERAL (SELECT avg(quantity) AS avg_quantity FROM sales GROUP BY brand) AS t`,
			nil,
			q,
		)
	})

	t.Run("rows from", func(t *testing.T) {
		// Example from https://www.postgresql.org/docs/current/queries-table-expressions.html#QUERIES-TABLEFUNCTIONS

		// Note: Added WithOrdinality to test that it can be applied to RowsFrom.

		q := qrb.Select(qrb.N("*")).
			From(qrb.RowsFrom(
				qrb.Func("json_to_recordset", qrb.String(`[{"a":40,"b":"foo"},{"a":"100","b":"bar"}]`)).ColumnDefinition("a", "INTEGER").ColumnDefinition("b", "TEXT"),
				qrb.Func("generate_series", qrb.Int(1), qrb.Int(3)),
			).WithOrdinality()).As("x").ColumnAliases("p", "q", "s").
			OrderBy(qrb.N("p"))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT *
			FROM ROWS FROM
				(
					json_to_recordset('[{"a":40,"b":"foo"},{"a":"100","b":"bar"}]')
						AS (a INTEGER, b TEXT),
					generate_series(1, 3)
				) WITH ORDINALITY AS x (p, q, s)
			ORDER BY p
			`,
			nil,
			q,
		)
	})

	t.Run("use embedded IdentExp", func(t *testing.T) {
		var films = struct {
			builder.Identer
		}{
			Identer: qrb.N("films"),
		}

		q := qrb.
			Select(qrb.N("*")).
			From(films)

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT * FROM films
			`,
			nil,
			q,
		)
	})
}

func TestSelectBuilder_LeftJoin(t *testing.T) {
	q1 := qrb.Select(qrb.Int(1)).From(qrb.N("foo")).LeftJoin(qrb.N("bar")).On(qrb.N("foo.id").Eq(qrb.N("bar.id")))
	q2 := q1.LeftJoin(qrb.N("baz")).Using("id")

	testhelper.AssertSQLWriterEquals(t, "SELECT 1 FROM foo LEFT JOIN bar ON foo.id = bar.id", nil, q1)
	testhelper.AssertSQLWriterEquals(t, "SELECT 1 FROM foo LEFT JOIN bar ON foo.id = bar.id LEFT JOIN baz USING (id)", nil, q2)
}

func TestSelectBuilder_CrossJoin(t *testing.T) {
	q1 := qrb.Select(qrb.Int(1)).From(qrb.N("foo")).CrossJoin(qrb.N("bar")).On(qrb.N("foo.id").Eq(qrb.N("bar.id")))
	q2 := q1.CrossJoinLateral(qrb.N("baz")).Using("id")

	testhelper.AssertSQLWriterEquals(t, "SELECT 1 FROM foo CROSS JOIN bar ON foo.id = bar.id", nil, q1)
	testhelper.AssertSQLWriterEquals(t, "SELECT 1 FROM foo CROSS JOIN bar ON foo.id = bar.id CROSS JOIN LATERAL baz USING (id)", nil, q2)
}

func TestSelectBuilder_Select(t *testing.T) {
	t.Run("immutability", func(t *testing.T) {
		q1 := qrb.Select(qrb.Int(1))
		q2 := q1.Select(qrb.Int(2))

		testhelper.AssertSQLWriterEquals(t, "SELECT 1", nil, q1)
		testhelper.AssertSQLWriterEquals(t, "SELECT 1,2", nil, q2)
	})

	t.Run("select distinct", func(t *testing.T) {
		q := qrb.Select().Distinct().
			Select(qrb.N("foo")).
			From(qrb.N("bar"))

		testhelper.AssertSQLWriterEquals(t, "SELECT DISTINCT foo FROM bar", nil, q)
	})

	t.Run("select distinct on", func(t *testing.T) {
		q := qrb.Select().Distinct().On(qrb.N("name"), qrb.Func("lower", qrb.N("email"))).
			Select(qrb.N("foo")).
			From(qrb.N("bar"))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT DISTINCT ON (name, lower(email)) foo
			FROM bar
			`,
			nil,
			q,
		)
	})
}

func TestSelectBuilder_SelectAs(t *testing.T) {
	q1 := qrb.Select().Select(qrb.Int(1)).As("foo")
	q2 := q1.Select(qrb.Int(2)).As("bar")

	testhelper.AssertSQLWriterEquals(t, "SELECT 1 AS foo", nil, q1)
	testhelper.AssertSQLWriterEquals(t, "SELECT 1 AS foo,2 AS bar", nil, q2)
}

func TestSelectBuilder_Where(t *testing.T) {
	t.Run("immutability", func(t *testing.T) {
		q1 := qrb.Select(qrb.N("foo")).Select().Where(qrb.N("is_active").Eq(qrb.Bool(true)))
		q2 := q1.Where(qrb.N("username").Eq(qrb.Arg("admin")))

		testhelper.AssertSQLWriterEquals(t, "SELECT foo WHERE is_active = true", nil, q1)
		testhelper.AssertSQLWriterEquals(t, "SELECT foo WHERE is_active = true AND username = $1", []any{"admin"}, q2)
	})

	t.Run("where exists", func(t *testing.T) {
		q := qrb.Select(qrb.N("col1")).
			From(qrb.N("tab1")).
			Where(qrb.Exists(
				qrb.Select(qrb.Int(1)).From(qrb.N("tab2")).Where(qrb.N("col2").Eq(qrb.N("tab1.col2"))),
			))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT col1
			FROM tab1
			WHERE EXISTS (SELECT 1 FROM tab2 WHERE col2 = tab1.col2)
			`,
			nil,
			q,
		)
	})

	t.Run("where not exists", func(t *testing.T) {
		q := qrb.Select(qrb.N("col1")).
			From(qrb.N("tab1")).
			Where(qrb.Not(qrb.Exists(
				qrb.Select(qrb.Int(1)).From(qrb.N("tab2")).Where(qrb.N("col2").Eq(qrb.N("tab1.col2"))),
			)))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT col1
			FROM tab1
			WHERE NOT EXISTS (SELECT 1 FROM tab2 WHERE col2 = tab1.col2)
			`,
			nil,
			q,
		)
	})

	t.Run("where in args", func(t *testing.T) {
		ids := []int{1, 2, 3}

		q := qrb.Select(qrb.N("username")).
			From(qrb.N("accounts")).
			Where(qrb.N("id").In(qrb.Args(ids...)))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT username
			FROM accounts
			WHERE id IN ($1, $2, $3)
			`,
			[]any{1, 2, 3},
			q,
		)
	})

	t.Run("where in exps", func(t *testing.T) {
		q := qrb.Select(qrb.N("username")).
			From(qrb.N("accounts")).
			Where(qrb.N("id").In(qrb.Exps(qrb.Int(42), qrb.String("abc"))))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT username
			FROM accounts
			WHERE id IN (42, 'abc')
			`,
			nil,
			q,
		)
	})

	t.Run("where with negated junction", func(t *testing.T) {
		q := qrb.Select(qrb.N("*")).
			From(qrb.N("accounts")).
			Where(qrb.Not(qrb.And(
				qrb.N("is_active").Eq(qrb.Bool(true)),
				qrb.N("username").Eq(qrb.Arg("admin")),
			)))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT *
			FROM accounts
			WHERE NOT (is_active = true AND username = $1)
			`,
			[]any{"admin"},
			q,
		)
	})

	t.Run("where with negated comparison", func(t *testing.T) {
		q := qrb.Select(qrb.N("*")).
			From(qrb.N("accounts")).
			Where(qrb.Not(
				qrb.N("is_active").Eq(qrb.Bool(true)),
			))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT *
			FROM accounts
			WHERE NOT is_active = true
			`,
			nil,
			q,
		)
	})

	t.Run("where with equal is not null", func(t *testing.T) {
		isActive := true

		q := qrb.Select(qrb.N("*")).
			From(qrb.N("accounts")).
			Where(qrb.Arg(isActive).Eq(
				qrb.N("deactivated_at").IsNull(),
			))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT *
			FROM accounts
			WHERE $1 = (deactivated_at IS NULL)
			`,
			[]any{true},
			q,
		)
	})

	t.Run("where all with subselect", func(t *testing.T) {
		q := qrb.Select(qrb.N("*")).
			From(qrb.N("employees")).
			Where(qrb.N("salary").Gt(qrb.All(qrb.Select(qrb.N("salary")).From(qrb.N("managers")))))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT *
			FROM employees 
			WHERE salary > ALL (SELECT salary FROM managers)
			`,
			nil,
			q,
		)
	})

	t.Run("where any with array", func(t *testing.T) {
		q := qrb.Select(qrb.N("*")).
			From(qrb.N("table")).
			Where(qrb.N("column").Eq(qrb.Any(qrb.Array(qrb.Int(1), qrb.Int(2), qrb.Int(3)))))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			SELECT *
			FROM table
			WHERE column = ANY (ARRAY[1, 2, 3])
			`,
			nil,
			q,
		)
	})
}

func TestSelectBuilder_GroupBy(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		q1 := qrb.
			Select(fn.Sum(qrb.N("y"))).
			From(qrb.N("test1")).
			GroupBy().
			Empty()

		testhelper.AssertSQLWriterEquals(t, "SELECT sum(y) FROM test1 GROUP BY ()", nil, q1)
	})

	t.Run("single", func(t *testing.T) {
		q1 := qrb.
			Select(qrb.N("x"), fn.Sum(qrb.N("y"))).
			From(qrb.N("test1")).
			GroupBy(qrb.N("x"))

		testhelper.AssertSQLWriterEquals(t, "SELECT x,sum(y) FROM test1 GROUP BY x", nil, q1)
	})

	t.Run("multiple", func(t *testing.T) {
		q1 := qrb.
			Select(
				qrb.N("product_id"),
				qrb.N("p.name"),
			).
			Select(fn.Sum(qrb.N("s.units")).Op("*", qrb.N("p.price"))).As("sales").
			From(qrb.N("products")).As("p").
			LeftJoin(qrb.N("sales")).As("s").Using("product_id").
			GroupBy(qrb.N("product_id"), qrb.N("p.name"), qrb.N("p.price"))

		testhelper.AssertSQLWriterEquals(
			t,
			"SELECT product_id,p.name,sum(s.units) * p.price AS sales FROM products AS p LEFT JOIN sales AS s USING (product_id) GROUP BY (product_id,p.name,p.price)",
			nil,
			q1,
		)
	})

	t.Run("rollup", func(t *testing.T) {
		q1 := qrb.
			Select(qrb.N("a"), qrb.N("b"), qrb.N("c"), qrb.N("d")).
			From(qrb.N("test1")).
			GroupBy().
			Rollup(
				qrb.Exps(qrb.N("a")),
				qrb.Exps(qrb.N("b"), qrb.N("c")),
				qrb.Exps(qrb.N("d")),
			)

		testhelper.AssertSQLWriterEquals(
			t,
			"SELECT a,b,c,d FROM test1 GROUP BY ROLLUP (a,(b,c),d)",
			nil,
			q1,
		)
	})

	t.Run("distinct rollup", func(t *testing.T) {
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

		testhelper.AssertSQLWriterEquals(
			t,
			"SELECT a,b,c FROM test1 GROUP BY DISTINCT ROLLUP (a, b), ROLLUP (a, c)",
			nil,
			q,
		)
	})

	t.Run("cube", func(t *testing.T) {
		q1 := qrb.
			Select(qrb.N("a"), qrb.N("b"), qrb.N("c"), qrb.N("d")).
			From(qrb.N("test1")).
			GroupBy().
			Cube(
				qrb.Exps(qrb.N("a"), qrb.N("b")),
				qrb.Exps(qrb.N("c"), qrb.N("d")),
			)

		testhelper.AssertSQLWriterEquals(
			t,
			"SELECT a,b,c,d FROM test1 GROUP BY CUBE ((a,b),(c,d))",
			nil,
			q1,
		)
	})

	t.Run("grouping sets", func(t *testing.T) {
		q1 := qrb.
			Select(qrb.N("brand"), qrb.N("size"), fn.Sum(qrb.N("sales"))).
			From(qrb.N("items_sold")).
			GroupBy().
			GroupingSets(
				qrb.Exps(qrb.N("brand")),
				qrb.Exps(qrb.N("size")),
				qrb.Exps(),
			)

		testhelper.AssertSQLWriterEquals(
			t,
			"SELECT brand,size,sum(sales) FROM items_sold GROUP BY GROUPING SETS (brand,size,())",
			nil,
			q1,
		)
	})
}

func TestSelectBuilder_OrderBy(t *testing.T) {
	q1 := qrb.
		Select(qrb.N("foo")).
		OrderBy(qrb.N("foo")).Desc()
	q2 := q1.
		Select(qrb.N("bar")).
		OrderBy(qrb.N("bar")).Asc().NullsLast()

	testhelper.AssertSQLWriterEquals(t, "SELECT foo ORDER BY foo DESC", nil, q1)

	testhelper.AssertSQLWriterEquals(t, "SELECT foo,bar ORDER BY foo DESC,bar ASC NULLS LAST", nil, q2)
}

func TestSelectBuilder_With(t *testing.T) {
	t.Run("immutability", func(t *testing.T) {
		q1 := qrb.With("foo").As(qrb.Select(qrb.Int(1))).Select(qrb.N("foo"))
		q2 := q1.AppendWith(qrb.With("bar").As(qrb.Select(qrb.Int(2))))

		testhelper.AssertSQLWriterEquals(t, "WITH foo AS (SELECT 1) SELECT foo", nil, q1)

		testhelper.AssertSQLWriterEquals(t, "WITH foo AS (SELECT 1),bar AS (SELECT 2) SELECT foo", nil, q2)
	})

	t.Run("materialized", func(t *testing.T) {
		q := qrb.With("foo").AsMaterialized(qrb.Select(qrb.Int(1))).Select(qrb.N("foo"))

		testhelper.AssertSQLWriterEquals(t, "WITH foo AS MATERIALIZED (SELECT 1) SELECT foo", nil, q)
	})

	t.Run("not materialized", func(t *testing.T) {
		q := qrb.With("foo").AsNotMaterialized(qrb.Select(qrb.Int(1))).Select(qrb.N("foo"))

		testhelper.AssertSQLWriterEquals(t, "WITH foo AS NOT MATERIALIZED (SELECT 1) SELECT foo", nil, q)
	})

	t.Run("multiple recursive", func(t *testing.T) {
		q := qrb.With("foo").As(qrb.Select(qrb.Int(1))).
			WithRecursive("bar").As(qrb.Select(qrb.Int(2))).
			Select(qrb.N("foo"))

		testhelper.AssertSQLWriterEquals(t, "WITH RECURSIVE foo AS (SELECT 1),bar AS (SELECT 2) SELECT foo", nil, q)
	})

	t.Run("recursive with search depth", func(t *testing.T) {
		q := qrb.
			WithRecursive("search_tree").ColumnNames("id", "link", "data").As(
			qrb.Select(qrb.N("t.id"), qrb.N("t.link"), qrb.N("t.data")).
				From(qrb.N("tree")).As("t").
				Union().All().
				Select(qrb.N("t.id"), qrb.N("t.link"), qrb.N("t.data")).
				From(qrb.N("tree")).As("t").
				From(qrb.N("search_tree")).As("st").
				Where(qrb.N("t.id").Eq(qrb.N("st.link"))),
		).SearchDepthFirst().By(qrb.N("id")).Set("ordercol").
			Select(qrb.N("*")).From(qrb.N("search_tree")).OrderBy(qrb.N("ordercol"))

		testhelper.AssertSQLWriterEquals(
			t,
			`
			WITH RECURSIVE search_tree(id, link, data) AS (
				SELECT t.id, t.link, t.data
				FROM tree AS t
			  UNION ALL
				SELECT t.id, t.link, t.data
				FROM tree AS t, search_tree AS st
				WHERE t.id = st.link
			) SEARCH DEPTH FIRST BY id SET ordercol
			SELECT * FROM search_tree ORDER BY ordercol
			`,
			nil,
			q,
		)
	})
}

func TestSelectBuilder_For(t *testing.T) {
	t.Run("for update", func(t *testing.T) {
		var q builder.SelectBuilder = qrb.Select(qrb.N("foo")).From(qrb.N("bar")).ForUpdate().
			SelectBuilder

		testhelper.AssertSQLWriterEquals(t, "SELECT foo FROM bar FOR UPDATE", nil, q)
	})

	t.Run("for key share of table1, table2 skip locked", func(t *testing.T) {
		var q builder.SelectBuilder = qrb.Select(qrb.N("foo")).From(qrb.N("bar")).ForKeyShare().Of("table1", "table2").SkipLocked().
			SelectBuilder

		testhelper.AssertSQLWriterEquals(t, "SELECT foo FROM bar FOR KEY SHARE OF table1,table2 SKIP LOCKED", nil, q)
	})
}

func TestSelectBuilder_IsEmpty(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		query := builder.SelectBuilder{}
		assert.Equal(t, true, query.IsEmpty())
	})

	t.Run("not empty", func(t *testing.T) {
		query := qrb.Select(qrb.N("foo")).From(qrb.N("bar")).SelectBuilder
		assert.Equal(t, false, query.IsEmpty())
	})
}
