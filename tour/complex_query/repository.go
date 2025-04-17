package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/go-sqlt/sqlt"
	_ "modernc.org/sqlite"
)

type Book struct {
	ID      int64
	Title   string
	Authors []string
	AddedAt time.Time
}

type Query struct {
	Title       string
	Author      string
	AddedBefore time.Time
}

type Insert struct {
	Title   string
	Authors []string
}

type Link struct {
	BookID    int64
	AuthorIDs []int64
}

var (
	config = sqlt.Config{
		Templates: []sqlt.Template{
			sqlt.Funcs(sprig.TxtFuncMap()),
			sqlt.ParseFiles("complex_query/queries.go.tpl"),
		},
		Log: func(ctx context.Context, info sqlt.Info) {
			if info.Cached {
				fmt.Println(info.SQL, info.Args)
			}
		},
	}

	schema = sqlt.Exec[any](config, sqlt.Lookup("schema"))

	upsertAuthors = sqlt.All[[]string, int64](config, sqlt.Lookup("upsert_authors"))

	insertBook = sqlt.One[string, int64](config, sqlt.Lookup("insert_book"))

	linkBookAuthors = sqlt.Exec[Link](config, sqlt.Lookup("link_book_authors"))

	queryBooks = sqlt.All[Query, Book](config, sqlt.NoExpirationCache(512), sqlt.Lookup("query_books"))
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

func (r Repository) Create(ctx context.Context, params Insert) (id int64, err error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = tx.Commit()
		}
	}()

	authorIDs, err := upsertAuthors.Exec(ctx, tx, params.Authors)
	if err != nil {
		return 0, err
	}

	id, err = insertBook.Exec(ctx, tx, params.Title)
	if err != nil {
		return 0, err
	}

	_, err = linkBookAuthors.Exec(ctx, tx, Link{
		BookID:    id,
		AuthorIDs: authorIDs,
	})

	return id, err
}

func (r Repository) Query(ctx context.Context, params Query) ([]Book, error) {
	return queryBooks.Exec(ctx, r.DB, params)
}
