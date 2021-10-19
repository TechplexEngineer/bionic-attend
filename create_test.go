package server

import (
	"database/sql"
	"github.com/matryer/is"
	"github.com/techplexengineer/go-frc-attend/data"
	"io/ioutil"
	"log"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func SetupTestDB(t *testing.T) *data.Queries {
	is := is.New(t)

	dataFile := "./testdb.queries"
	err := os.RemoveAll(dataFile)
	is.NoErr(err)

	driverName := "sqlite" //https://gitlab.com/cznic/sqlite/blob/v1.13.1/examples/example1/main.go#L30
	db, err := sql.Open(driverName, dataFile)
	is.NoErr(err)
	file, err := ioutil.ReadFile("./queries/schema.sql")
	is.NoErr(err)
	_, err = db.Exec(string(file))
	is.NoErr(err)

	return data.New(db)
}

func TestCreateNewUser(t *testing.T) {
	is := is.New(t)
	type args struct {
		db            *data.Queries
		handleUserMsg func(msg string, path Route) error
		firstName     string
		lastName      string
		userID        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "happy", args: args{
			db: SetupTestDB(t),
			handleUserMsg: func(msg string, path Route) error {
				log.Printf("Msg: %s - route: %s", msg, path)
				is.Equal(msg, "User created")
				is.Equal(string(path), "/")
				return nil
			},
			firstName: "Davy",
			lastName:  "Jones",
			userID:    "123456",
		}, wantErr: false},
		{name: "missing firstname", args: args{
			db: SetupTestDB(t),
			handleUserMsg: func(msg string, path Route) error {
				log.Printf("Msg: %s - route: %s", msg, path)
				is.Equal(msg, "First Name must be longer than 3 characters.")
				is.Equal(string(path), "")
				return nil
			},
			firstName: "",
			lastName:  "Jones",
			userID:    "123456",
		}, wantErr: true},
		{name: "missing lastname", args: args{
			db: SetupTestDB(t),
			handleUserMsg: func(msg string, path Route) error {
				log.Printf("Msg: %s - route: %s", msg, path)
				is.Equal(msg, "Last Name must be longer than 3 characters.")
				is.Equal(string(path), "")
				return nil
			},
			firstName: "Davy",
			lastName:  "",
			userID:    "123456",
		}, wantErr: true},
		{name: "missing userid", args: args{
			db: SetupTestDB(t),
			handleUserMsg: func(msg string, path Route) error {
				log.Printf("Msg: %s - route: %s", msg, path)
				is.Equal(msg, "UserID must be longer than 3 characters.")
				is.Equal(string(path), "")
				return nil
			},
			firstName: "Davy",
			lastName:  "Jones",
			userID:    "123456",
		}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateNewUser(tt.args.db, tt.args.handleUserMsg, tt.args.firstName, tt.args.lastName, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("CreateNewUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
