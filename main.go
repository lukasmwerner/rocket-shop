package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var db *pgxpool.Pool

func main() {
	poolconfig, err := pgxpool.ParseConfig("postgres://postgres:imAbadPassword@127.0.0.1:5432")
	if err != nil {
		os.Exit(1)
	}

	db, err = pgxpool.ConnectConfig(context.Background(), poolconfig)

	t := &Template{templates: template.Must(template.ParseGlob("public/views/*.tmpl"))}

	e := echo.New()
	e.Renderer = t
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.tmpl", "")
	})

	e.GET("/rockets", getRockets)
	e.GET("/rockets/:id", getRocket)
	e.Debug = true
	e.Logger.Fatal(e.Start(":1323"))
}

func getRocket(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return err
	}
	Name := ""
	Price := ""
	Description := ""
	Company := ""
	err = db.QueryRow(context.Background(), "SELECT * from stock WHERE id=$1", id).Scan(nil, &Name, &Price, &Description, &Company)
	if err != nil {
		return err
	}
	rocket := Rocket{
		Id:          id,
		Name:        Name,
		Price:       Price,
		Description: Description,
		Company:     Company,
	}
	return c.Render(http.StatusOK, "rocketListing.tmpl", rocket)
}

func getRockets(c echo.Context) error {
	var rockets []Rocket = make([]Rocket, 0)
	rows, err := db.Query(context.Background(), "SELECT * FROM stock;")
	if err != nil {
		return err
	}
	if err == pgx.ErrNoRows {
		return c.NoContent(http.StatusNotFound)
	}
	for rows.Next() {
		id := 0
		Name := ""
		Price := ""
		Description := ""
		Company := ""
		rows.Scan(&id, &Name, &Price, &Description, &Company)
		rockets = append(rockets, Rocket{
			Id:          id,
			Name:        Name,
			Price:       Price,
			Description: Description,
			Company:     Company,
		})
	}
	fmt.Println(rockets)

	return c.Render(http.StatusOK, "rocketsList.tmpl", rockets)
}
