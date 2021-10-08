package server

import (
	"database/sql"
	_ "embed" // for schema file
	"fmt"
	"io/fs"
	_ "modernc.org/sqlite" // for databae driver
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/techplexengineer/go-frc-attend/data"
)

type Server struct {
	// shared state between routes
	router     *mux.Router
	templateFs fs.FS
	db         *data.Queries
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

func (s *Server) SetupDB(dbFile string) error {
	driverName := "sqlite" //https://gitlab.com/cznic/sqlite/blob/v1.13.1/examples/example1/main.go#L30
	db, err := sql.Open(driverName, dbFile)
	if err != nil {
		return err
	}

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		// dbFile needs to be created
		_, err = db.Exec(schema)
		if err != nil {
			return fmt.Errorf("unable to create tables [[%s]] - %w", schema, err)
		}
	}

	s.db = data.New(db)
	return nil
}
