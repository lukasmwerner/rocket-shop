package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

type Rocket struct {
	Id          int
	Name        string
	Price       string
	Description string
	Company     string
}

var Falcon9 = Rocket{
	Id:          0,
	Name:        "Falcon 9",
	Price:       "$36 Million/Seat",
	Description: "Testing testing 123",
	Company:     "SpaceX",
}
var DeltaIV = Rocket{
	Id:          0,
	Name:        "Delta IV",
	Price:       "$56 Million/Seat",
	Description: "~Retired~",
	Company:     "ULA",
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	t := &Template{templates: template.Must(template.ParseGlob("public/views/*.tmpl"))}

	e := echo.New()
	e.Renderer = t
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.tmpl", "")
	})

	e.GET("/rockets", getRockets)
	e.GET("/rockets/:id", getRocket)

	e.Logger.Fatal(e.Start(":1323"))
}

func getRocket(c echo.Context) error {
	//idParam := c.Param("id")
	// id, err := strconv.Aoti(idParam)

	return c.Render(http.StatusOK, "rocketListing.tmpl", Falcon9)
}

func getRockets(c echo.Context) error {
	var rockets []Rocket = []Rocket{Falcon9, DeltaIV}
	return c.Render(http.StatusOK, "rocketsList.tmpl", rockets)
}
