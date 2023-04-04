package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/networkteam/qrb"
	"github.com/networkteam/qrb/qrbpgx"
)

// E.g. run via `DATABASE_URL="dbname=chinook" go run ./examples/pgx`
func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
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

	{
		row, err := qrbpgx.
			Build(q).
			WithExecutor(conn).
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

	{
		row, err := qrbpgx.
			Build(q).
			WithExecutor(pool).
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

		log.Printf("(pool) Album title: %s, Artist name: %s, Longest track: %s\n", albumTitle, artistName, time.Duration(length)*time.Millisecond)
	}

	(func() {
		tx, err := pool.BeginTx(ctx, pgx.TxOptions{
			IsoLevel: pgx.ReadCommitted,
		})
		if err != nil {
			log.Fatalf("Error starting transaction: %v", err)
		}
		defer func(tx pgx.Tx, ctx context.Context) {
			err := tx.Commit(ctx)
			if err != nil {
				log.Fatalf("Error committing transaction: %v", err)
			}
		}(tx, ctx)

		row, err := qrbpgx.
			Build(q).
			WithExecutor(tx).
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
