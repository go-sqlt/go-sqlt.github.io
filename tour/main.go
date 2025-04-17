package main

import (
	"context"
	"fmt"
	"time"

	bulk_insert "docs/tour/bulk_insert"
	complex_query "docs/tour/complex_query"
	create_statements "docs/tour/create_statements"
	custom_functions "docs/tour/custom_functions"
	load_from_file "docs/tour/load_from_file"
	multiple_databases "docs/tour/multiple_databases"
	transactions "docs/tour/transactions"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func main() {
	create_statements_example()
	load_from_file_example()
	bulk_insert_example()
	transactions_example()
	custom_functions_example()
	complex_query_example()
	multiple_db_sqlite_example()
	multiple_db_postgres_example()
}

func multiple_db_sqlite_example() {
	r, err := multiple_databases.NewSqlite()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	_, err = r.Create(ctx, multiple_databases.Insert{
		Title:   "Moby-Dick",
		Authors: []string{"Herman Melville"},
	})
	if err != nil {
		panic(err)
	}

	_, err = r.Create(ctx, multiple_databases.Insert{
		Title:   "Good Omens",
		Authors: []string{"Neil Gaiman", "Terry Pratchett"},
	})
	if err != nil {
		panic(err)
	}

	_, err = r.Create(ctx, multiple_databases.Insert{
		Title:   "Discworld",
		Authors: []string{"Terry Pratchett"},
	})
	if err != nil {
		panic(err)
	}

	for range 2 {
		books, err := r.Query(ctx, multiple_databases.Query{
			Author: "Terry Pratchett",
		})
		if err != nil {
			panic(err)
		}

		fmt.Println(books)
	}

	fmt.Println("multiple_db_sqlite_example ✅")
}

func multiple_db_postgres_example() {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:15.3-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			panic(err)
		}
	}()

	conn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	r, err := multiple_databases.NewPostgres(conn)
	if err != nil {
		panic(err)
	}

	_, err = r.Create(ctx, multiple_databases.Insert{
		Title:   "Moby-Dick",
		Authors: []string{"Herman Melville"},
	})
	if err != nil {
		panic(err)
	}

	_, err = r.Create(ctx, multiple_databases.Insert{
		Title:   "Good Omens",
		Authors: []string{"Neil Gaiman", "Terry Pratchett"},
	})
	if err != nil {
		panic(err)
	}

	_, err = r.Create(ctx, multiple_databases.Insert{
		Title:   "Discworld",
		Authors: []string{"Terry Pratchett"},
	})
	if err != nil {
		panic(err)
	}

	for range 2 {
		books, err := r.Query(ctx, multiple_databases.Query{
			Author: "Terry Pratchett",
		})
		if err != nil {
			panic(err)
		}

		fmt.Println(books)
	}

	fmt.Println("multiple_db_postgres_example ✅")
}

func complex_query_example() {
	r, err := complex_query.NewRepository()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	_, err = r.Create(ctx, complex_query.Insert{
		Title:   "Moby-Dick",
		Authors: []string{"Herman Melville"},
	})
	if err != nil {
		panic(err)
	}

	_, err = r.Create(ctx, complex_query.Insert{
		Title:   "Good Omens",
		Authors: []string{"Neil Gaiman", "Terry Pratchett"},
	})
	if err != nil {
		panic(err)
	}

	_, err = r.Create(ctx, complex_query.Insert{
		Title:   "Discworld",
		Authors: []string{"Terry Pratchett"},
	})
	if err != nil {
		panic(err)
	}

	for range 2 {
		books, err := r.Query(ctx, complex_query.Query{
			Author: "Terry Pratchett",
		})
		if err != nil {
			panic(err)
		}

		fmt.Println(books)
	}

	fmt.Println("complex_query_example ✅")
}

func transactions_example() {
	r, err := transactions.NewRepository()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	id, err := r.Create(ctx, transactions.Params{
		Title:  "Moby-Dick",
		Author: "Herman Melville",
	})
	if err != nil {
		panic(err)
	}

	book, err := r.Get(ctx, id)
	if err != nil {
		panic(err)
	}

	fmt.Println(book)

	ids, err := r.CreateMany(ctx, []transactions.Params{
		{
			Title:  "The Great Gatsby",
			Author: "F. Scott Fitzgerald",
		},
		{
			Title:  "Lord of the Flies",
			Author: "William Golding",
		},
	})
	if err != nil {
		panic(err)
	}

	for _, id := range ids {
		book, err := r.Get(ctx, id)
		if err != nil {
			panic(err)
		}

		fmt.Println(book)
	}

	fmt.Println("transactions_example ✅")
}

func custom_functions_example() {
	r, err := custom_functions.NewRepository()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	id, err := r.Create(ctx, custom_functions.Params{
		Title:  "Moby-Dick",
		Author: "Herman Melville",
		Genre:  custom_functions.Adventure,
	})
	if err != nil {
		panic(err)
	}

	book, err := r.Get(ctx, id)
	if err != nil {
		panic(err)
	}

	fmt.Println(book)

	ids, err := r.CreateMany(ctx, []custom_functions.Params{
		{
			Title:  "The Great Gatsby",
			Author: "F. Scott Fitzgerald",
			Genre:  custom_functions.Tragedy,
		},
		{
			Title:  "Lord of the Flies",
			Author: "William Golding",
			Genre:  custom_functions.Allegorical,
		},
	})
	if err != nil {
		panic(err)
	}

	for _, id := range ids {
		book, err := r.Get(ctx, id)
		if err != nil {
			panic(err)
		}

		fmt.Println(book)
	}

	fmt.Println("custom_functions_example ✅")
}

func bulk_insert_example() {
	r, err := bulk_insert.NewRepository()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	id, err := r.Create(ctx, bulk_insert.Params{
		Title:  "Moby-Dick",
		Author: "Herman Melville",
	})
	if err != nil {
		panic(err)
	}

	book, err := r.Get(ctx, id)
	if err != nil {
		panic(err)
	}

	fmt.Println(book)

	ids, err := r.CreateMany(ctx, []bulk_insert.Params{
		{
			Title:  "The Great Gatsby",
			Author: "F. Scott Fitzgerald",
		},
		{
			Title:  "Lord of the Flies",
			Author: "William Golding",
		},
	})
	if err != nil {
		panic(err)
	}

	for _, id := range ids {
		book, err := r.Get(ctx, id)
		if err != nil {
			panic(err)
		}

		fmt.Println(book)
	}

	fmt.Println("bulk_insert_example ✅")
}

func load_from_file_example() {
	r, err := load_from_file.NewRepository()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	id, err := r.Create(ctx, "Moby-Dick")
	if err != nil {
		panic(err)
	}

	book, err := r.Get(ctx, id)
	if err != nil {
		panic(err)
	}

	fmt.Println(book)
	fmt.Println("load_from_file_example ✅")
}

func create_statements_example() {
	r, err := create_statements.NewRepository()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	id, err := r.Create(ctx, "Moby-Dick")
	if err != nil {
		panic(err)
	}

	book, err := r.Get(ctx, id)
	if err != nil {
		panic(err)
	}

	fmt.Println(book)
	fmt.Println("create_statements_example ✅")
}
