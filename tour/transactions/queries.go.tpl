{{ define "schema" }}
    CREATE TABLE IF NOT EXISTS books (
        id INTEGER PRIMARY KEY,
        title TEXT NOT NULL,
        author_id INTEGER REFERENCES authors(id),
        added_at DATE NOT NULL
    );

    CREATE TABLE IF NOT EXISTS authors (
        id INTEGER PRIMARY KEY,
        name TEXT UNIQUE
    );
{{ end }}

{{ define "upsert_author" }}
    INSERT INTO authors (name) VALUES ({{ . }})
    ON CONFLICT (name) DO UPDATE SET
        name = EXCLUDED.name
    RETURNING id;
{{ end }}

{{/* "now" is a template function imported from Masterminds/sprig. */}}
{{ define "insert_book" }}
    INSERT INTO books (title, author_id, added_at) VALUES ({{ .Title }}, {{ .AuthorID }}, {{ now }}) RETURNING id;
{{ end }}

{{ define "upsert_authors" }}
    INSERT INTO authors (name) VALUES
    {{ range $i, $p := . }}
        {{ if $i }}, {{ end }}
        ({{ $p.Author }})
    {{ end }}
    ON CONFLICT (name) DO UPDATE SET
        name = EXCLUDED.name
    RETURNING id;
{{ end }}

{{ define "insert_books" }}
    {{ $now := now }}
    INSERT INTO books (title, author_id, added_at) VALUES
    {{ range $i, $p := . }}
        {{ if $i }}, {{ end }}
        ({{ $p.Title }}, {{ $p.AuthorID }}, {{ $now }})
    {{ end }}
    RETURNING id;
{{ end }}

{{ define "get_book" }}
    SELECT
        books.id                {{ Scan "ID" }}
        , books.title           {{ Scan "Title" }}
        , authors.name          {{ Scan "Author" }}
        , books.added_at        {{ Scan "AddedAt" }}
    FROM books
    LEFT JOIN authors ON authors.id = books.author_id
    WHERE books.id = {{ . }};
{{ end }}
