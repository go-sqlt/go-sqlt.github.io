package repository

import (
	"context"
	"database/sql"

	"github.com/go-sqlt/sqlt"
	_ "modernc.org/sqlite"
)

type Book struct {
	ID     int64
	Title  string
	Author string
}

type Params struct {
	Title  string
	Author string
}

var (
	config = sqlt.ParseFiles("bulk_insert/queries.go.tpl")

	schema = sqlt.Exec[any](config, sqlt.Lookup("schema"))

	insertBook = sqlt.First[Params, int64](config, sqlt.Lookup("insert_book"))

	insertBooks = sqlt.All[[]Params, int64](config, sqlt.Lookup("insert_books"))

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

func (r Repository) Create(ctx context.Context, params Params) (int64, error) {
	return insertBook.Exec(ctx, r.DB, params)
}

func (r Repository) CreateMany(ctx context.Context, params []Params) ([]int64, error) {
	return insertBooks.Exec(ctx, r.DB, params)
}

func (r Repository) Get(ctx context.Context, id int64) (Book, error) {
	return getBook.Exec(ctx, r.DB, id)
}
