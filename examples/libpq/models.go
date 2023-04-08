package main

import (
	"time"

	"github.com/lib/pq"
)

type Book struct {
	ID              int
	Title           string
	AuthorID        int
	AuthorName      string
	PublicationYear int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Genres          pq.StringArray // Note: this is specific to the lib/pq driver, we could as well return JSON from the query
}
