package server

//func (s *server) handleHome(templateFs fs.FS) http.HandlerFunc {
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
