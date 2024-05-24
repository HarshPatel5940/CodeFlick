package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Count struct {
	Count int
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	count := &Count{Count: 0}
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {

		return c.Render(200, "index", map[string]interface{}{
			"title": "Hello World",
			"Count": count.Count,
		})
	})

	e.POST("/api/increment", func(c echo.Context) error {
		count.Count++
		return c.Render(200, "counter", map[string]interface{}{
			"Count": count.Count,
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
