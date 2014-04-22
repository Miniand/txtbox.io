package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Miniand/txtbox.io/model"

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
	http.Redirect(res, req, fmt.Sprintf(
		"/%s/%s", wr.GeneratedKeys[0], txtUser.Token), http.StatusFound)
}
