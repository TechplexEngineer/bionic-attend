package server

import (
	"fmt"
	"github.com/gorilla/sessions"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"
)

const (
	FlashCookieName = "flash"
)

func SetFlash(w http.ResponseWriter, value string) {
	log.Printf("FLASH: %s", value)
	cookie := sessions.NewCookie(FlashCookieName, value, &sessions.Options{
		Path:   "/",
		Domain: "",
		MaxAge: 60, // seconds
	})
	//c := &http.Cookie{Name: FlashCookieName, Value: value}
	http.SetCookie(w, cookie)
}

func GetFlash(w http.ResponseWriter, r *http.Request) (string, error) {
	c, err := r.Cookie(FlashCookieName)
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return "", nil
		default:
			return "", err
		}
	}
	cookie := sessions.NewCookie(FlashCookieName, c.Value, &sessions.Options{
		Path:   "/",
		Domain: "",
		MaxAge: -1, // delete now
	})
	http.SetCookie(w, cookie) // delete cookie
	log.Printf("GetFlash!")
	return c.Value, nil
}

func LoadBaseTemplates(fs fs.FS, funcs *template.FuncMap) (*template.Template, error) {
	//@note: if performance becomes a problem, we could load these once, instead of every request

	tmpl := template.New("layout.html")

	tmpl = tmpl.Funcs(template.FuncMap{
		"now": time.Now,
		"toHTML": func(s string) (template.HTML, error) {
			return template.HTML(s), nil
		},
		"toJS": func(s string) template.JS {
			return template.JS(s)
		},
	})

	if funcs != nil {
		tmpl = tmpl.Funcs(*funcs)
	}

	baseLayoutTemplate := "partials/layout.html"
	tmpl, err := tmpl.ParseFS(fs, baseLayoutTemplate)
	if err != nil {
		return nil, fmt.Errorf("unable to parse baseLayoutTemplate (%s): %w", baseLayoutTemplate, err)
	}

	partialGlob := "partials/*.html"
	tmpl, err = tmpl.ParseFS(fs, partialGlob)
	if err != nil {
		return nil, fmt.Errorf("unable to parse glob (%s): %w", partialGlob, err)
	}

	return tmpl, nil
}
