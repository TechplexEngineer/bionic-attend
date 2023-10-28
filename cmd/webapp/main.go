package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	server "github.com/techplexengineer/bionic-attend"
)

func run() error {
	svr := server.Server{}
	svr.SetupTemplateFs(server.Resources)

	dbPath, dbPathFound := os.LookupEnv("DBPATH")
	if !dbPathFound {
		dbPath = "./attendance.db"
	}

	if err := svr.SetupDB(dbPath); err != nil {
		return fmt.Errorf("unable to setupdb: %w", err)
	}
	svr.SetupRoutes()

	port, portFound := os.LookupEnv("PORT")
	if !portFound {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Server listening on %s", addr)
	// this blocks until the server exits
	err := http.ListenAndServe(addr, svr)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Printf("Error: %s", err)
		getwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Unable to get CWD: %s", err)
		} else {
			fmt.Printf("CWD IS: " + getwd)
		}
		os.Exit(1)
	}
}
