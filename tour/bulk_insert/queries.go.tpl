{{ define "schema" }}
    CREATE TABLE IF NOT EXISTS books (
        id INTEGER PRIMARY KEY,
        title TEXT NOT NULL,
        author TEXT NOT NULL
    )
{{ end }}

{{/* sqlt uses github.com/jba/templatecheck to check input parameters. */}}
{{ define "insert_book" }}
    INSERT INTO books (title, author) VALUES ({{ .Title }}, {{ .Author }}) RETURNING id;
{{ end }}

{{/* You can use range/if to create dynamic queries. */}}
{{ define "insert_books" }}
    INSERT INTO books (title, author) VALUES
    {{ range $i, $c := . }}
        {{ if $i }}, {{ end }}
        ({{ $c.Title }}, {{ $c.Author }})
    {{ end }}
    RETURNING id;
{{ end }}

{{ define "get_book" }}
    SELECT
        id              {{ ScanInt "ID" }}
        , title         {{ ScanString "Title" }}
        , author        {{ ScanString "Author" }}
    FROM books
    WHERE id = {{ . }};
{{ end }}
