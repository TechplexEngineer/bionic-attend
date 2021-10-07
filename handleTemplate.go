package server

import (
	"fmt"
	"log"
	"net/http"
)

func (s Server) handleTemplate(templatePath string, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		tmpl, err := LoadBaseTemplates(s.templateFs)
		if err != nil {
			log.Printf("error LoadBaseTemplates - %s", err)
			s.handleInternalError(err)(w, req)
			return
		}

		tmpl, err = tmpl.ParseFS(s.templateFs, templatePath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panicf("2 unable to parse template (%s): %s", templatePath, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = tmpl.Execute(w, data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Panicf("unable to execute template: %s", err)
			return
		}
	}
}

func (s Server) handleInternalError(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("ERROR: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = fmt.Fprintf(w, "Internal Server Error")
		if err != nil {
			log.Printf("Unable to write to http response")
		}

	}
}
