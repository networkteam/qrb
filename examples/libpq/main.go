package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/urfave/cli/v2"
)

func main() {
	var db *sql.DB

	app := cli.NewApp()
	app.Name = "example-libpq"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "database-url",
			Usage:   "Database URL",
			EnvVars: []string{"DATABASE_URL"},
			Value:   "dbname=qrb-examples sslmode=disable",
		},
	}
	app.Before = func(c *cli.Context) error {
		var err error
		db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		return err
	}
	app.Commands = []*cli.Command{
		{
			Name: "books",
			Subcommands: []*cli.Command{
				{
					Name: "list",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "genre",
							Usage: "Filter by genre",
						},
						&cli.StringFlag{
							Name:  "author",
							Usage: "Filter by author",
						},
					},
					Action: func(c *cli.Context) error {
						books, err := findAllBooks(c.Context, db, booksFilter{
							GenreName:  c.String("genre"),
							AuthorName: c.String("author"),
						})
						if err != nil {
							return err
						}

						for _, book := range books {
							fmt.Printf(
								"%s by %s (%d) [%s]\n",
								book.Title,
								book.AuthorName,
								book.PublicationYear,
								strings.Join(book.Genres, ", "),
							)
						}

						return nil
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
