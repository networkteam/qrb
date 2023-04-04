package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/qrbsql"
)

var (
	Track              = qrb.N(`"Track"`)
	Track_AlbumId      = qrb.N(`"Track"."AlbumId"`)
	Track_Milliseconds = qrb.N(`"Track"."Milliseconds"`)

	Album       = qrb.N(`"Album"`)
	Album_Title = qrb.N(`"Album"."Title"`)
	Artist      = qrb.N(`"Artist"`)
	Artist_Name = qrb.N(`"Artist"."Name"`)
)

// E.g. run via `DATABASE_URL="dbname=chinook sslmode=disable" go run ./examples/libpq`
func main() {
	ctx := context.Background()

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	// TODO The naming is horrible with uppercased letters in the schema!

	q := qrb.With("longest_track").As(
		qrb.Select(Track_AlbumId, Track_Milliseconds).
			From(Track).
			OrderBy(Track_Milliseconds).Desc().
			Limit(qrb.Int(1)),
	).
		Select(Album_Title).As("album_title").
		Select(Artist_Name).As("artist_name").
		Select(qrb.N(`"Milliseconds"`)).As("length").
		From(Album).
		Join(qrb.N(`"longest_track"`)).Using(`"AlbumId"`).
		Join(Artist).Using(`"ArtistId"`)

	{
		row, err := qrbsql.
			Build(q).
			WithExecutor(db).
			QueryRow(ctx)

		var (
			albumTitle string
			artistName string
			length     int64
		)
		err = row.Scan(&albumTitle, &artistName, &length)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		log.Printf("(conn) Album title: %s, Artist name: %s, Longest track: %s\n", albumTitle, artistName, time.Duration(length)*time.Millisecond)
	}

	(func() {
		tx, err := db.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
		})
		if err != nil {
			log.Fatalf("Error starting transaction: %v", err)
		}
		defer func(tx *sql.Tx) {
			err := tx.Commit()
			if err != nil {
				log.Fatalf("Error committing transaction: %v", err)
			}
		}(tx)

		ex := qrbsql.NewExecutorBuilder(tx)

		row, err := ex.
			Build(q).
			QueryRow(ctx)

		var (
			albumTitle string
			artistName string
			length     int64
		)
		err = row.Scan(&albumTitle, &artistName, &length)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		log.Printf("(tx) Album title: %s, Artist name: %s, Longest track: %s\n", albumTitle, artistName, time.Duration(length)*time.Millisecond)
	})()
}
