package app

import (
	"bytes"
	"fmt"
	"forum/models"
	"forum/ui"
	"html/template"
	"io/fs"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Local().Format("02 Jan 2006 at 15:04")
}

func sequence(start, end int) []int {
	var seq []int
	for i := start; i <= end; i++ {
		seq = append(seq, i)
	}
	return seq
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"add": func(a, b int) int {
		return a + b
	},
	"sub": func(a, b int) int {
		return a - b
	},
	"sequence": sequence,
	"toLower":  strings.ToLower,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/*.layout.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}

var Quotes = []string{"Strength is not in the grandmothers. After all, grandmothers are already old.", "Out of the 64 battles I fought, I had 64 victories. All battles were with shadows.", "Took a knife - cut, took a doshik - eat", "I live as the cards fall. You live as your mom says.", "Never give up, go towards your goal! And if it's difficult - give up.", "If you get lost in the forest, go home.", "Remember: just one mistake - and you're wrong.", "Do it the right way. If it's not the right way, don't do it.", "As my grandfather used to say, \"I'm your grandfather.\"", "Work is not a wolf. Nobody is a wolf. Only a wolf is a wolf."}

func (app *Application) Render(w http.ResponseWriter, status int, page string, data *models.TemplateData) {
	i := rand.Intn(10)
	data.Quote = Quotes[i]
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.ServerError(w, err)
		return
	}
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}
