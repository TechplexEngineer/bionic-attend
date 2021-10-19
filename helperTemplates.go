package server

import (
	"fmt"
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
	maxAgeSec := 60 // seconds

	// The Expires field calculated based on the MaxAge value, for Internet
	// Explorer compatibility.
	d := time.Duration(maxAgeSec) * time.Second
	expires := time.Now().Add(d)

	cookie := &http.Cookie{
		Name:    FlashCookieName,
		Value:   value,
		Path:    "/",
		Domain:  "",
		MaxAge:  maxAgeSec, // seconds
		Expires: expires,
	}
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

	cookie := &http.Cookie{
		Name:    FlashCookieName,
		Path:    "/",
		Domain:  "",
		MaxAge:  -1, // remove now
		Expires: time.Unix(1, 0),
	}

	http.SetCookie(w, cookie) // delete cookie
	log.Printf("GetFlash!")
	return c.Value, nil
}

func LoadBaseTemplates(fs fs.FS, funcs *template.FuncMap) (*template.Template, error) {
	// @note: if performance becomes a problem, we could load these once, instead of every request

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
