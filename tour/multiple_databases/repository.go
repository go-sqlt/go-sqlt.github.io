package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/go-sqlt/sqlt"
	_ "github.com/jackc/pgx/v5/stdlib"
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

func NewRepository(db *sql.DB, opts ...sqlt.Option) (Repository, error) {
	config := sqlt.Configure(append(opts,
		sqlt.Funcs(sprig.TxtFuncMap()),
		sqlt.ParseFiles("multiple_databases/queries.go.tpl"),
	)...)

	schema := sqlt.Exec[any](config, sqlt.Lookup("schema"))

	_, err := schema.Exec(context.Background(), db, nil)
	if err != nil {
		return Repository{}, err
	}

	return Repository{
		DB:              db,
		Schema:          sqlt.Exec[any](config, sqlt.Lookup("schema")),
		UpsertAuthors:   sqlt.All[[]string, int64](config, sqlt.Lookup("upsert_authors")),
		InsertBook:      sqlt.One[string, int64](config, sqlt.Lookup("insert_book")),
		LinkBookAuthors: sqlt.Exec[Link](config, sqlt.Lookup("link_book_authors")),
		QueryBooks:      sqlt.All[Query, Book](config, sqlt.NoExpirationCache(512), sqlt.Lookup("query_books")),
	}, nil
}

func NewSqlite() (Repository, error) {
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(1)")
	if err != nil {
		return Repository{}, err
	}

	return NewRepository(db, sqlt.Sqlite())
}

func NewPostgres(conn string) (Repository, error) {
	db, err := sql.Open("pgx", conn)
	if err != nil {
		return Repository{}, err
	}

	return NewRepository(db, sqlt.Postgres())
}

type Repository struct {
	DB              *sql.DB
	Schema          sqlt.Statement[any, sql.Result]
	UpsertAuthors   sqlt.Statement[[]string, []int64]
	InsertBook      sqlt.Statement[string, int64]
	LinkBookAuthors sqlt.Statement[Link, sql.Result]
	QueryBooks      sqlt.Statement[Query, []Book]
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

	authorIDs, err := r.UpsertAuthors.Exec(ctx, tx, params.Authors)
	if err != nil {
		return 0, err
	}

	id, err = r.InsertBook.Exec(ctx, tx, params.Title)
	if err != nil {
		return 0, err
	}

	_, err = r.LinkBookAuthors.Exec(ctx, tx, Link{
		BookID:    id,
		AuthorIDs: authorIDs,
	})

	return id, err
}

func (r Repository) Query(ctx context.Context, params Query) ([]Book, error) {
	return r.QueryBooks.Exec(ctx, r.DB, params)
}
