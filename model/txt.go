package model

import (
	"bytes"
	"math/rand"
	"time"
)

const (
	TXT_ACCESS_NONE = iota
	TXT_ACCESS_READ
	TXT_ACCESS_WRITE
	TXT_ACCESS_ADMIN

	TXT_USER_TOKEN_LEN = 12
)

var txtUserNameChars = []byte(
	"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var txtUserNameCharsLen = len(txtUserNameChars)
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type Txt struct {
	Id            string        `gorethink:"id,omitempty"`
	Current       TxtRevision   `gorethink:"current"`
	Revisions     []TxtRevision `gorethink:"revisions"`
	DefaultAccess int           `gorethink:"defaultAccess"`
	Users         []TxtUser     `gorethink:"users"`
}

func NewTxt() Txt {
	return Txt{
		Revisions: []TxtRevision{},
		Users:     []TxtUser{},
	}
}

type TxtRevision struct {
	Title   string    `gorethink:"title"`
	Content string    `gorethink:"content"`
	At      time.Time `gorethink:"at"`
	By      string    `gorethink:"by"`
}

type TxtUser struct {
	Token  string `gorethink:"token"`
	Name   string `gorethink:"name"`
	Access int    `gorethink:"access"`
}

func NewTxtUser() TxtUser {
	return TxtUser{
		Token: NewTxtUserToken(TXT_USER_TOKEN_LEN),
	}
}

func NewTxtUserToken(len int) string {
	buf := &bytes.Buffer{}
	for i := 0; i < len; i++ {
		buf.WriteByte(txtUserNameChars[r.Int()%txtUserNameCharsLen])
	}
	return buf.String()
}
