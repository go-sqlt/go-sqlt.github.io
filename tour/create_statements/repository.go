package repository

import (
	"context"
	"database/sql"

	"github.com/go-sqlt/sqlt"
	_ "modernc.org/sqlite"
)

type Book struct {
	ID    int64
	Title string
}

var (
	schema = `
		CREATE TABLE IF NOT EXISTS books (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL
		);
	`

	// Question is the default placeholder.
	config = sqlt.Question()

	// Single column queries do not need mapping definition.
	// Params are always parameterized preventing SQL injection.
	insertBook = sqlt.First[string, int64](config, sqlt.Parse(`
		INSERT INTO books (title) VALUES ({{ . }}) RETURNING id;
	`))

	// Define the mapping with Scan functions.
	// One ensures that only one row is returned by the query (else: sqlt.ErrTooManyRows).
	getBook = sqlt.One[int64, Book](config, sqlt.Parse(`
		SELECT
			id              {{ ScanInt "ID" }}
			, title         {{ ScanString "Title" }}
		FROM books
		WHERE id = {{ . }};
	`))
)

func NewRepository() (Repository, error) {
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(1)")
	if err != nil {
		return Repository{}, err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return Repository{}, err
	}

	return Repository{
		DB: db,
	}, nil
}

type Repository struct {
	DB *sql.DB
}

func (r Repository) Create(ctx context.Context, title string) (int64, error) {
	return insertBook.Exec(ctx, r.DB, title)
}

func (r Repository) Get(ctx context.Context, id int64) (Book, error) {
	return getBook.Exec(ctx, r.DB, id)
}
