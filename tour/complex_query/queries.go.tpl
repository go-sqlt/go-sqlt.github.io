{{ define "schema" }}
    CREATE TABLE IF NOT EXISTS books (
        id INTEGER PRIMARY KEY,
        title TEXT NOT NULL,
        added_at DATE NOT NULL
    );

    CREATE TABLE IF NOT EXISTS authors (
        id INTEGER PRIMARY KEY,
        name TEXT UNIQUE
    );

    CREATE TABLE IF NOT EXISTS book_authors (
        book_id INTEGER REFERENCES books(id),
        author_id INTEGER REFERENCES authors(id),
        PRIMARY KEY (book_id, author_id)
    );
{{ end }}

{{ define "upsert_authors" }}
    INSERT INTO authors (name) VALUES
    {{ range $i, $a := . }}
            {{ if $i }}, {{ end }}
            ({{ $a }})
    {{ end }}
    ON CONFLICT (name) DO UPDATE SET
            id = authors.id,
            name = EXCLUDED.name
    RETURNING id;
{{ end }}

{{ define "insert_book" }}
    INSERT INTO books (title, added_at) VALUES ({{ . }}, {{ now }}) RETURNING id;
{{ end }}

{{ define "link_book_authors" }}
    INSERT INTO book_authors (book_id, author_id) VALUES
    {{ range $i, $a := .AuthorIDs }}
            {{ if $i }}, {{ end }}
            ({{ $.BookID }}, {{ $a }})
    {{ end }}
    ON CONFLICT DO NOTHING;
{{ end }}

{{ define "query_books" }}
    SELECT
        books.id                        {{ ScanInt "ID" }}
        , books.title                   {{ ScanString "Title" }}
        {{/* ScanStringSlice scans the column as a string and splits it into a slice of strings */}}
        , GROUP_CONCAT(authors.name)    {{ ScanStringSlice "Authors" "," }}
        , books.added_at                {{ ScanTime "AddedAt" }}
    FROM books
    LEFT JOIN book_authors ON books.id = book_authors.book_id
    LEFT JOIN authors ON authors.id = book_authors.author_id
    WHERE 1=1
    {{ with .Title }}
        AND books.title = {{ . }}
    {{ end }}
    {{ with .Author }}
        AND books.id IN (
            SELECT ba.book_id
            FROM book_authors ba
            JOIN authors a ON a.id = ba.author_id
            WHERE a.name = {{ . }}
        )
    {{ end }}
    {{ if not .AddedBefore.IsZero }}
        AND books.added_at < {{ .AddedBefore }}
    {{ end }}
    GROUP BY books.id, books.title, books.added_at;
{{ end }}
