package controller

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/Miniand/txtbox.io/model"
	"github.com/go-martini/martini"

	r "github.com/dancannon/gorethink"
)

var (
	hexChars    = "0123456789abcdef"
	base62Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

func BaseAtof(input, chars string) float64 {
	base := float64(len(chars))
	value := 0.0
	valueMap := map[rune]float64{}
	for i, a := range input {
		valueMap[a] = float64(i)
	}
	for _, a := range input {
		value = base*value + valueMap[a]
	}
	return value
}

func BaseFtoa(input float64, chars string) string {
	base := float64(len(chars))
	buf := &bytes.Buffer{}
	for input > 0 {
		buf.WriteByte(chars[int(math.Mod(input, base))])
		input = math.Floor(input / base)
	}
	return buf.String()
}

func ConvBaseStr(in, inBase, outBase string) string {
	return BaseFtoa(BaseAtof(in, inBase), outBase)
}

func EncodeId(input string) string {
	return ConvBaseStr(strings.Replace(input, "-", "", -1), hexChars,
		base62Chars)
}

func DecodeId(input string) string {
	raw := ConvBaseStr(input, base62Chars, hexChars)
	if len(raw) <= 20 {
		return raw
	}
	return raw[0:8] + "-" + raw[8:12] + "-" + raw[12:16] + "-" + raw[16:20] +
		"-" + raw[20:]
}

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
