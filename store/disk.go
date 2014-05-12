package store

import (
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Miniand/txtbox.io/str"
)

const (
	LatestSymlink = "latest"
	MetaDir       = "meta"
	UsersDir      = "users"
	PublicFile    = "public"
)

var filenameRegexp = regexp.MustCompile(`^(\d+);(.*);(.*)$`)

type DiskStore struct {
	Location string
}

func NewDiskStore(location string) *DiskStore {
	return &DiskStore{
		Location: location,
	}
}

func ParseFilename(filename string) (rev int, title, by string, err error) {
	parts := strings.Split(filename, ";")
	rev, err = strconv.Atoi(parts[0])
	if len(parts) > 1 {
		title = parts[1]
	}
	if len(parts) > 2 {
		by = parts[2]
	}
	return
}

func ParseUser(filename string) (user string, perm int, name string, err error) {
	parts := strings.Split(filename, ";")
	user = parts[0]
	if len(parts) > 1 {
		perm, err = strconv.Atoi(parts[1])
	}
	if len(parts) > 2 {
		name = parts[2]
	}
	return
}

func (ds *DiskStore) Create() (File, error) {
	var (
		id, fileLoc string
	)
	for {
		id = str.RandomBase62Str(8)
		fileLoc = path.Join(ds.Location, id)
		if _, err := os.Stat(fileLoc); os.IsNotExist(err) {
			break
		} else if err != nil {
			return nil, err
		}
	}
	if err := os.MkdirAll(path.Join(fileLoc, MetaDir), 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(path.Join(fileLoc, UsersDir), 0755); err != nil {
		return nil, err
	}
	return &DiskFile{
		id:       id,
		location: fileLoc,
	}, nil
}

func (ds *DiskStore) Find(id string) (File, error) {
	fileLoc := path.Join(ds.Location, id)
	if _, err := os.Stat(fileLoc); os.IsNotExist(err) {
		return nil, ErrNotFound
	}
	return &DiskFile{
		id:       id,
		location: fileLoc,
	}, nil
}

func (ds *DiskStore) Delete(id string) error {
	return os.RemoveAll(path.Join(ds.Location, id))
}

type DiskFile struct {
	id       string
	location string
}

func (df *DiskFile) Id() (string, error) {
	return df.id, nil
}

func (df *DiskFile) PublicFilename() string {
	return path.Join(df.location, MetaDir, PublicFile)
}

func (df *DiskFile) RevisionFilename(rev int) string {
	return path.Join(df.location, strconv.Itoa(rev))
}

func (df *DiskFile) Visibility() (int, error) {
	if _, err := os.Stat(df.PublicFilename()); err != nil {
		if os.IsNotExist(err) {
			return VisPrivate, nil
		}
		return 0, err
	}
	return VisPublic, nil
}

func (df *DiskFile) SetVisibility(vis int) (err error) {
	if vis == 1 {
		f, err := os.Create(df.PublicFilename())
		if err != nil {
			f.Close()
			if os.IsExist(err) {
				err = nil
			}
		}
	} else {
		err := os.Remove(df.PublicFilename())
		if os.IsNotExist(err) {
			err = nil
		}
	}
	return
}

func (df *DiskFile) Revision(num int) (Revision, error) {
	return &DiskRevision{
		filename: df.RevisionFilename(num),
	}, nil
}

func (df *DiskFile) latestRevisionNum() (int, error) {
	f, err := os.Readlink(df.latestPath())
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	num, _, _, err := ParseFilename(path.Base(f))
	return num, err
}

func (df *DiskFile) Latest() (Revision, error) {
	filename, err := os.Readlink(df.latestPath())
	return &DiskRevision{
		filename: filename,
	}, err
}

func (df *DiskFile) latestPath() string {
	return path.Join(df.location, LatestSymlink)
}

func (df *DiskFile) Save(body io.Reader, title, by string) (Revision, error) {
	num, err := df.latestRevisionNum()
	if err != nil {
		return nil, fmt.Errorf("could not get latest revision num, %v", err)
	}
	num += 1
	filename := path.Join(df.location, fmt.Sprintf("%d;%s;%s", num, title, by))
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err := io.Copy(f, body); err != nil {
		return nil, err
	}
	if err := os.Symlink(filename, df.RevisionFilename(num)); err != nil {
		return nil, err
	}
	if err := os.RemoveAll(df.latestPath()); err != nil {
		return nil, err
	}
	if err := os.Symlink(filename, df.latestPath()); err != nil {
		return nil, err
	}
	return &DiskRevision{
		filename: filename,
	}, nil
}

func (df *DiskFile) Users() ([]User, error) {
	return nil, nil
}

func (df *DiskFile) User(id string) (User, error) {
	return nil, nil
}

func (df *DiskFile) CreateUser(perm int) (User, error) {
	return nil, nil
}

func (df *DiskFile) UpdateUserPerm(id string, perm int) error {
	return nil
}

func (df *DiskFile) UpdateUserName(id, name string) error {
	return nil
}

func (df *DiskFile) DeleteUser(id string) error {
	return nil
}

type DiskRevision struct {
	filename string
}

func (dr *DiskRevision) Title() (string, error) {
	_, title, _, err := ParseFilename(path.Base(dr.filename))
	return title, err
}

func (dr *DiskRevision) Body() (io.Reader, error) {
	return os.Open(dr.filename)
}

func (dr *DiskRevision) By() (string, error) {
	_, _, by, err := ParseFilename(path.Base(dr.filename))
	return by, err
}

func (dr *DiskRevision) Num() (int, error) {
	num, _, _, err := ParseFilename(path.Base(dr.filename))
	return num, err
}

func (dr *DiskRevision) At() (time.Time, error) {
	fi, err := os.Stat(dr.filename)
	if err != nil {
		return time.Now(), err
	}
	return fi.ModTime(), nil
}

type DiskUser struct{}

func (du *DiskUser) Id() (string, error) {
	return "", nil
}

func (du *DiskUser) Name() (string, error) {
	return "", nil
}

func (du *DiskUser) Perm() (int, error) {
	return 0, nil
}
