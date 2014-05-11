package store

import (
	"errors"
	"io"
	"time"
)

const (
	VisPrivate = iota
	VisPublic
)

const (
	PermRead = iota
	PermWrite
	PermAdmin
)

var (
	ErrNotFound = errors.New("could not find file")
)

type Store interface {
	Create() (File, error)
	Find(id string) (File, error)
	Delete(id string) error
}

type File interface {
	Id() (string, error)
	Visibility() (int, error)
	SetVisibility(v int) error
	Revision(num int) (Revision, error)
	Latest() (Revision, error)
	Save(body io.Reader, title, by string) (Revision, error)
	Users() ([]User, error)
	User(id string) (User, error)
	CreateUser(perm int) (User, error)
	UpdateUserPerm(id string, perm int) error
	UpdateUserName(id, name string) error
	DeleteUser(id string) error
}

type Revision interface {
	Title() (string, error)
	Body() (io.Reader, error)
	By() (string, error)
	Num() (int, error)
	At() (time.Time, error)
}

type User interface {
	Id() (string, error)
	Name() (string, error)
	Perm() (int, error)
}
