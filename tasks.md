# Project Task List (Gocket)

A concise, practical sequence to build a local Pocket-like app in Go with a SQLite DB and a simple web UI. Aim to complete 1–2 steps per session.

1. Initialize Go module and tooling
   - go mod init, add basic .gitignore, choose minimal dependencies (HTTP router, SQLite driver, HTML parsing)

2. Create project structure
   - cmd/gocket/main.go, internal/server, internal/storage, internal/snapshot, web/templates, web/static

3. Add configuration and env handling
   - Port, DB path (use `gocket.db` by default), log level

4. Set up SQLite access layer
   - Use modernc.org/sqlite or mattn/go-sqlite3; create connection helper

5. Define DB schema and migrations
   - Tables: articles(id, url, title, saved_at), snapshots(article_id, html, content_text)
   - Provide an auto-run migration on startup

6. Implement storage repository methods
   - Create, get, list articles; upsert/get snapshots; simple search by title/text

7. Build minimal HTTP server and routing
   - Health route, index route, static files, template renderer

8. Implement “Add URL” endpoint and form
   - POST /articles with a URL; validate and persist article record

9. Implement snapshot/fetch pipeline
   - Fetch URL, parse title, extract text, store raw HTML + text snapshot

10. List and browse saved articles
   - GET /articles shows list; GET /articles/{id} shows detail with snapshot

11. Basic HTML templates and styling
   - Simple, readable templates and CSS for list/detail/add form

12. Background snapshotting (optional first pass)
   - If fetch fails inline, enqueue retry or provide a manual “refresh snapshot” action

13. Add basic tests
   - Unit tests for storage and snapshot parsing; a few handler tests

14. Containerization
   - Dockerfile with multi-stage build; mount DB file; expose port

15. Polish and docs
   - Update README with run/build instructions, env vars, and Docker usage

Nice-to-haves (later):
- Full-text search (FTS5), tags, read/unread, pagination, favicon fetching
- Kubernetes manifest, graceful shutdown, metrics


