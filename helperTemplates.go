package server

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cast"
)

const (
	FlashCookieName = "flash"
)

func SetFlash(w http.ResponseWriter, value string) {
	log.Printf("FLASH: %s", value)
	//c := &http.Cookie{Name: FlashCookieName, Value: value}
	//http.SetCookie(w, c)
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
	dc := &http.Cookie{Name: FlashCookieName, MaxAge: -1}
	http.SetCookie(w, dc)
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
		"seq": Seq,
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

// Seq creates a sequence of integers.  It's named and used as GNU's seq.
//
// Examples:
//     3 => 1, 2, 3
//     1 2 4 => 1, 3
//     -3 => -1, -2, -3
//     1 4 => 1, 2, 3, 4
//     1 -2 => 1, 0, -1, -2
func Seq(args ...interface{}) ([]int, error) {
	if len(args) < 1 || len(args) > 3 {
		return nil, errors.New("invalid number of arguments to Seq")
	}

	intArgs := cast.ToIntSlice(args)
	if len(intArgs) < 1 || len(intArgs) > 3 {
		return nil, errors.New("invalid arguments to Seq")
	}

	inc := 1
	var last int
	first := intArgs[0]

	if len(intArgs) == 1 {
		last = first
		if last == 0 {
			return []int{}, nil
		} else if last > 0 {
			first = 1
		} else {
			first = -1
			inc = -1
		}
	} else if len(intArgs) == 2 {
		last = intArgs[1]
		if last < first {
			inc = -1
		}
	} else {
		inc = intArgs[1]
		last = intArgs[2]
		if inc == 0 {
			return nil, errors.New("'increment' must not be 0")
		}
		if first < last && inc < 0 {
			return nil, errors.New("'increment' must be > 0")
		}
		if first > last && inc > 0 {
			return nil, errors.New("'increment' must be < 0")
		}
	}

	// sanity check
	if last < -100000 {
		return nil, errors.New("size of result exceeds limit")
	}
	size := ((last - first) / inc) + 1

	// sanity check
	if size <= 0 || size > 2000 {
		return nil, errors.New("size of result exceeds limit")
	}

	seq := make([]int, size)
	val := first
	for i := 0; ; i++ {
		seq[i] = val
		val += inc
		if (inc < 0 && val < last) || (inc > 0 && val > last) {
			break
		}
	}

	return seq, nil
}
