package controller

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Miniand/txtbox.io/model"
	"github.com/go-martini/martini"

	r "github.com/dancannon/gorethink"
)

func TxtCreate(session *r.Session, res http.ResponseWriter,
	req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	txt := model.NewTxt()
	txtUser := model.NewTxtUser()
	txtUser.Access = model.TXT_ACCESS_ADMIN
	txt.Users = append(txt.Users, txtUser)
	txtRev := model.TxtRevision{
		Title:   "Blah",
		Content: string(body),
		At:      time.Now(),
		By:      txtUser.Token,
	}
	txt.Current = txtRev
	txt.Revisions = append(txt.Revisions, txtRev)
	wr, err := r.Table("txt").Insert(txt).RunWrite(session)
	if err != nil {
		panic(err)
	}
	http.Redirect(res, req, fmt.Sprintf("/%s/%s",
		wr.GeneratedKeys[0], txtUser.Token), http.StatusFound)
}

func TxtShow(session *r.Session, res http.ResponseWriter,
	params martini.Params) {
	id := params["id"]
	txt := model.Txt{}
	row, err := r.Table("txt").Get(id).RunRow(session)
	if err != nil {
		panic(err)
	}
	if err := row.Scan(&txt); err != nil {
		panic(err)
	}
	res.WriteHeader(http.StatusOK)
	bw := bufio.NewWriter(res)
	bw.WriteString(txt.Current.Content)
	bw.Flush()
}
