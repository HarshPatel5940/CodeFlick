package main

import (
	"html/template"
	"io"
	"log"
	"strconv"

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

type Person struct {
	Name  string
	Email string
}

type H = map[string]interface{}

type Contact struct {
	// ? can't name it .id cause it will then have problem with templating
	Id    int
	Name  string
	Email string
}

type contacts = []Contact

type Data struct {
	Contacts contacts
}

func newContact(id int, name string, email string) Contact {
	return Contact{
		Id:    id,
		Name:  name,
		Email: email,
	}
}

func newData() Data {
	return Data{
		Contacts: []Contact{
			newContact(0, "harsh", "test@gmail.com"),
			newContact(1, "harsh-2", "test-2@gmail.com"),
		},
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = newTemplate()

	count := &Count{Count: 0}
	data := newData()

	e.GET("/", func(c echo.Context) error {

		return c.Render(200, "index", H{
			"title": "Hello World",
			"Count": count.Count,
		})
	})

	e.POST("/api/increment", func(c echo.Context) error {
		count.Count++
		return c.Render(200, "counter", H{
			"Count": count.Count,
		})
	})

	e.GET("/contacts", func(c echo.Context) error {
		return c.Render(200, "contact-index", data)
	})

	e.POST("/api/contacts", func(c echo.Context) error {
		name := c.FormValue("Name")
		email := c.FormValue("Email")
		newContact := newContact(count.Count+2, name, email)
		count.Count++
		data.Contacts = append(data.Contacts, newContact)

		return c.Render(200, "contacts", data)
	})

	e.DELETE("/api/contacts/:id", func(c echo.Context) error {
		fId := c.Param("id")
		log.Println("-->", fId)
		Id, err := strconv.Atoi(fId)
		if err != nil {
			return c.JSON(400, H{
				"error": "Invalid Id",
			})
		}

		// filter out the contact with the given id
		contacts := data.Contacts[:0]
		for _, contact := range data.Contacts {
			if contact.Id != Id {
				contacts = append(contacts, contact)
			}
		}
		count.Count--
		data.Contacts = contacts

		return c.Render(200, "contacts", data)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
