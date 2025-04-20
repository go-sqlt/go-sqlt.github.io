{{ define "schema" }}
    CREATE TABLE IF NOT EXISTS books (
        id INTEGER PRIMARY KEY,
        title TEXT NOT NULL
    );
{{ end }}

{{ define "insert_book" }}
    INSERT INTO books (title) VALUES ({{ . }}) RETURNING id;
{{ end }}

{{ define "get_book" }}
    SELECT
        id          {{ ScanInt "ID" }}
        , title     {{ ScanString "Title" }}
    FROM books
    WHERE id = {{ . }};
{{ end }}
