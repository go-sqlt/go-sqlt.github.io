package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/go-sqlt/sqlt"
	_ "modernc.org/sqlite"
)

type Genre int

const (
	Adventure Genre = iota + 1
	Tragedy
	Allegorical
)

type Book struct {
	ID      int64
	Title   string
	Genre   Genre
	Author  string
	AddedAt time.Time
}

type Params struct {
	Title  string
	Genre  Genre
	Author string
}

type Insert struct {
	Title    string
	Genre    Genre
	AuthorID int64
}

var (
	config = sqlt.Config{
		Templates: []sqlt.Template{
			sqlt.Funcs(sprig.TxtFuncMap()),
			sqlt.Funcs(template.FuncMap{
				"ValueGenre": func(g Genre) (string, error) {
					switch g {
					case Adventure:
						return "Adventure", nil
					case Tragedy:
						return "Tragedy", nil
					case Allegorical:
						return "Allegorical", nil
					default:
						return "", fmt.Errorf("unknown genre: %d", g)
					}
				},
				"ScanGenre": func() sqlt.Scanner[Book] {
					return func() (any, func(dest *Book) error) {
						var txt string

						return &txt, func(dest *Book) error {
							switch txt {
							case "Adventure":
								dest.Genre = Adventure
							case "Tragedy":
								dest.Genre = Tragedy
							case "Allegorical":
								dest.Genre = Allegorical
							default:
								return fmt.Errorf("unknown genre: %s", txt)
							}

							return nil
						}
					}
				},
			}),
			sqlt.ParseFiles("custom_functions/queries.go.tpl"),
		},
	}

	schema = sqlt.Exec[any](config, sqlt.Lookup("schema"))

	upsertAuthor = sqlt.One[string, int64](config, sqlt.Lookup("upsert_author"))
	insertBook   = sqlt.One[Insert, int64](config, sqlt.Lookup("insert_book"))

	upsertAuthors = sqlt.All[[]Params, int64](config, sqlt.Lookup("upsert_authors"))
	insertBooks   = sqlt.All[[]Insert, int64](config, sqlt.Lookup("insert_books"))

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

func (r Repository) Create(ctx context.Context, params Params) (id int64, err error) {
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

	author, err := upsertAuthor.Exec(ctx, tx, params.Author)
	if err != nil {
		return 0, err
	}

	return insertBook.Exec(ctx, tx, Insert{
		Title:    params.Title,
		Genre:    params.Genre,
		AuthorID: author,
	})
}

func (r Repository) CreateMany(ctx context.Context, params []Params) ([]int64, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = tx.Commit()
		}
	}()

	authors, err := upsertAuthors.Exec(ctx, tx, params)
	if err != nil {
		return nil, err
	}

	insert := make([]Insert, len(authors))

	for i, p := range authors {
		insert[i] = Insert{
			Title:    params[i].Title,
			Genre:    params[i].Genre,
			AuthorID: p,
		}
	}

	return insertBooks.Exec(ctx, tx, insert)
}

func (r Repository) Get(ctx context.Context, id int64) (Book, error) {
	return getBook.Exec(ctx, r.DB, id)
}
