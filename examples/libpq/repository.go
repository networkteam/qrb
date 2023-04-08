package main

import (
	"context"

	. "github.com/networkteam/qrb"
	"github.com/networkteam/qrb/builder"
	"github.com/networkteam/qrb/fn"
	"github.com/networkteam/qrb/qrbsql"
)

type booksFilter struct {
	GenreName  string
	AuthorName string
}

// findAllBooks returns all books that match the given filter.
// If the filter is empty, all books are returned.
//
// This is a simple implementation of a repository function using a dynamically built query.
func findAllBooks(ctx context.Context, executor qrbsql.Executor, filter booksFilter) ([]Book, error) {
	query := Select(
		N("book_id"),
		N("title"),
		N("author_id"),
		N("publication_year"),
		N("created_at"),
		N("updated_at"),
	).
		Select(N("authors.name")).As("author_name").
		Select(selectBookGenresArray()).As("genres").
		From(N("books")).
		LeftJoin(N("authors")).Using("author_id")

	if filter.GenreName != "" {
		query = query.Where(
			Exists(
				Select().
					From(N("book_genre")).
					LeftJoin(N("genres")).Using("genre_id").
					Where(
						And(
							N("book_genre.book_id").Eq(N("books.book_id")),
							N("genres.name").ILike(Arg(filter.GenreName).Concat(String("%"))),
						),
					),
			),
		)
	}
	if filter.AuthorName != "" {
		query = query.Where(N("authors.name").ILike(Arg(filter.AuthorName).Concat(String("%"))))
	}

	query = query.OrderBy(N("title")).SelectBuilder // Note: return the embedded field to return a builder.SelectBuilder instead of a builder.OrderBySelectBuilder

	return queryAndScanBooks(ctx, executor, query)
}

// selectBookGenresArray builds a select expression to select a genres of a book as an array.
//
// It serves as an example on how parts of a query can be extracted into functions.
func selectBookGenresArray() builder.Exp {
	return Select(fn.ArrayAgg(N("genres.name")).OrderBy(N("genres.name"))).
		From(N("book_genre")).
		LeftJoin(N("genres")).Using("genre_id").
		Where(N("book_genre.book_id").Eq(N("books.book_id")))
}

func queryAndScanBooks(ctx context.Context, executor qrbsql.Executor, query builder.SelectBuilder) ([]Book, error) {
	rows, err := qrbsql.Build(query).WithExecutor(executor).Query(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err = rows.Scan(
			&book.ID,
			&book.Title,
			&book.AuthorID,
			&book.PublicationYear,
			&book.CreatedAt,
			&book.UpdatedAt,
			&book.AuthorName,
			&book.Genres,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, rows.Err()
}
