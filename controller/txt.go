package controller

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Miniand/txtbox.io/store"
	"github.com/go-martini/martini"
)

func TxtCreate(st store.Store, res http.ResponseWriter,
	req *http.Request) {
	fi, err := st.Create()
	if err != nil {
		panic(err)
	}
	_, err = fi.Save(req.Body, "", "")
	if err != nil {
		panic(err)
	}
	id, err := fi.Id()
	if err != nil {
		panic(err)
	}
	http.Redirect(res, req, fmt.Sprintf("/%s", id), http.StatusFound)
}

func TxtShow(st store.Store, res http.ResponseWriter,
	params martini.Params) {
	fi, err := st.Find(params["id"])
	if err != nil {
		panic(err)
	}
	rev, err := fi.Latest()
	if err != nil {
		panic(err)
	}
	body, err := rev.Body()
	if err != nil {
		panic(err)
	}
	res.WriteHeader(http.StatusOK)
	io.Copy(res, body)
}
