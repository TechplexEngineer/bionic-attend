package data

import (
	"context"
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestQueries_GetUser_no_rows(t *testing.T) {
	is := is.New(t)

	dataFile := "./testdb.db"
	is.NoErr(os.RemoveAll(dataFile))

	driverName := "sqlite" //https://gitlab.com/cznic/sqlite/blob/v1.13.1/examples/example1/main.go#L30
	db, err := sql.Open(driverName, dataFile)
	is.NoErr(err)
	file, err := ioutil.ReadFile("../db/schema.sql")
	is.NoErr(err)
	_, err = db.Exec(string(file))
	is.NoErr(err)

	q := New(db)

	user, err := q.GetUser(context.Background(), "123")

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		log.Printf("no rows") // desired outcome
	} else {
		is.NoErr(err)
		log.Printf("%#v", user)
	}

	is.NoErr(db.Close())

	is.NoErr(os.Remove(dataFile))
}
