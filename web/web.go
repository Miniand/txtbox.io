package web

import (
	"os"

	"github.com/Miniand/txtbox.io/controller"
	"github.com/Miniand/txtbox.io/store"
	"github.com/go-martini/martini"
)

func New() (*martini.ClassicMartini, error) {
	storage := os.Getenv("STORAGE")
	if storage == "" {
		storage = os.TempDir()
	}
	var st store.Store
	st = store.NewDiskStore(storage)
	m := martini.Classic()
	m.Map(st)
	m.Get("/:id/:user\\.:extension", controller.TxtShow)
	m.Get("/:id/:user", controller.TxtShow)
	m.Get("/:id\\.:extension", controller.TxtShow)
	m.Post("/:id", controller.TxtUpdate)
	m.Get("/:id", controller.TxtShow)
	m.Post("/", controller.TxtCreate)
	m.Get("/", controller.Home)
	return m, nil
}
