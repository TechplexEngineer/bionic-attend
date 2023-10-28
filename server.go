package server

import (
	"database/sql"
	"embed"
	_ "embed" // for schema file
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"io/fs"
	_ "modernc.org/sqlite" // for database driver
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/gorilla/mux"
	"github.com/techplexengineer/bionic-attend/data"
)

type Server struct {
	// shared state between routes
	router     *mux.Router
	templateFs fs.FS
	queries    *data.Queries
	db         *sql.DB
}

// ServeHTTP implements the http.Handler interface which allows the server to be passed to http.ListenAndServe
func (s Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func (s *Server) SetupTemplateFs(fs fs.FS) {
	s.templateFs = fs
}

//go:embed db/schema.sql
var schema string

//go:embed db/migrations/*.sql
var migrationsFs embed.FS

func (s *Server) SetupDB(dbFile string) error {
	driverName := "sqlite" //https://gitlab.com/cznic/sqlite/blob/v1.13.1/examples/example1/main.go#L30
	var err error
	s.db, err = sql.Open(driverName, dbFile)
	if err != nil {
		return err
	}

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		// dbFile needs to be created
		_, err = s.db.Exec(schema)
		if err != nil {
			return fmt.Errorf("unable to create tables [[%s]] - %w", schema, err)
		}
	}

	migrations, err := iofs.New(migrationsFs, "db/migrations")
	if err != nil {
		return fmt.Errorf("unable to migrate: %w", err)
	}

	migrationDb, err := sqlite.WithInstance(s.db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("unable to wrap db for migration: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", migrations, "sqlite", migrationDb)
	if err != nil {
		return fmt.Errorf("unable to prepare db for migration: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("unable to migrate up: %w", err)
	}

	s.queries = data.New(s.db)
	return nil
}

//func (s *Server) SetupSession() error {
//	authKey := securecookie.GenerateRandomKey(32)
//	encKey := securecookie.GenerateRandomKey(32)
//	sess := sessions.NewCookieStore(authKey, encKey)
//
//	sess.
//}
