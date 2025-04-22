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
	config = sqlt.ParseFiles("load_from_file/queries.go.tpl")

	// sqlt.Config and Lookup implement the sqlt.Option interface.
	schema = sqlt.Exec[any](config, sqlt.Lookup("schema"))

	// Statements panic if a template error occurs or a type-safety check fails.
	// Therefore, they should be created at application startup.
	insertBook = sqlt.First[string, int64](config, sqlt.Lookup("insert_book"))

	getBook = sqlt.First[int64, Book](config, sqlt.Lookup("get_book"))
)

func NewRepository() (Repository, error) {
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(1)")
	if err != nil {
		return Repository{}, err
	}

	_, err = schema.Exec(context.Background(), db, nil)
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
