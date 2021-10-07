package server

import (
	"net/http"
)

func (s Server) handleIndex() http.HandlerFunc {
	// one time handler setup work can go here

	return func(w http.ResponseWriter, req *http.Request) {
		// generate data

		// handle the request
		s.handleTemplate("index.html", nil)(w, req)
	}
}

//func (s *server) handleIndex(templateFs fs.FS) http.HandlerFunc {
//	var (
//		init   sync.Once
//		tpl    *template.Template
//		tplerr error
//	)
//	return func(w http.ResponseWriter, r *http.Request) {
//		init.Do(func() {
//			tpl, tplerr = template.ParseFS(templateFs)
//		})
//		if tplerr != nil {
//			http.Error(w, tplerr.Error(), http.StatusInternalServerError)
//			return
//		}
//		w.WriteHeader(http.StatusOK)
//		err := tpl.Execute(w, nil)
//		if err != nil {
//			// w.WriteHeader(http.StatusInternalServerError)
//			log.Panicf("unable to execute template: %s", err)
//			return
//		}
//	}
//}
