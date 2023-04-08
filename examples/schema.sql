DROP TABLE IF EXISTS authors, genres, books, book_genre;

CREATE TABLE authors
(
    author_id SERIAL PRIMARY KEY,
    name      TEXT NOT NULL
);

CREATE TABLE genres
(
    genre_id SERIAL PRIMARY KEY,
    name     TEXT NOT NULL
);

CREATE TABLE books
(
    book_id          SERIAL PRIMARY KEY,
    title            TEXT        NOT NULL,
    author_id        INTEGER     NOT NULL REFERENCES authors (author_id),
    publication_year INTEGER,
    created_at       TIMESTAMPTZ NOT NULL,
    updated_at       TIMESTAMPTZ NOT NULL
);

CREATE TABLE book_genre
(
    book_id  INTEGER REFERENCES books (book_id),
    genre_id INTEGER REFERENCES genres (genre_id),
    PRIMARY KEY (book_id, genre_id)
);

INSERT INTO authors (author_id, name)
VALUES (1, 'J.D. Salinger'),
       (2, 'Harper Lee'),
       (3, 'Jane Austen'),
       (4, 'F. Scott Fitzgerald');

INSERT INTO genres (genre_id, name)
VALUES (1, 'Classic'),
       (2, 'Fiction'),
       (3, 'Romance');

INSERT INTO books (book_id, title, author_id, publication_year, created_at, updated_at)
VALUES (1, 'The Catcher in the Rye', 1, 1951, '2023-01-01 00:00:00', '2023-01-01 00:00:00'),
       (2, 'To Kill a Mockingbird', 2, 1960, '2023-01-02 00:00:00', '2023-01-02 00:00:00'),
       (3, 'Pride and Prejudice', 3, 1813, '2023-01-03 00:00:00', '2023-01-03 00:00:00'),
       (4, 'Go Set a Watchman', 2, 2015, '2023-01-04 00:00:00', '2023-01-04 00:00:00'),
       (5, 'The Great Gatsby', 4, 1925, '2023-01-05 00:00:00', '2023-01-05 00:00:00');

INSERT INTO book_genre (book_id, genre_id)
VALUES (1, 1),
       (1, 2),
       (2, 1),
       (2, 2),
       (3, 1),
       (3, 3),
       (4, 1),
       (4, 2),
       (5, 1),
       (5, 2);
