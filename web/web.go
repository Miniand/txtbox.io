package web

import (
	"os"

	"github.com/Miniand/txtbox.io/controller"

	r "github.com/dancannon/gorethink"
	"github.com/go-martini/martini"
)

func New() (*martini.ClassicMartini, error) {
	dbAddr := os.Getenv("DB_ADDR")
	if dbAddr == "" {
		dbAddr = "localhost:28015"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "txtbox"
	}
	session, err := r.Connect(r.ConnectOpts{
		Address:  dbAddr,
		Database: dbName,
	})
	if err != nil {
		return nil, err
	}
	m := martini.Classic()
	m.Map(session)
	m.Get("/:id/:user\\.:extension", controller.TxtShow)
	m.Get("/:id/:user", controller.TxtShow)
	m.Get("/:id\\.:extension", controller.TxtShow)
	m.Get("/:id", controller.TxtShow)
	m.Post("/", controller.TxtCreate)
	m.Get("/", controller.Home)
	return m, nil
}
