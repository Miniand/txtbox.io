package store

import (
	"errors"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"

	"github.com/Miniand/txtbox.io/str"
)

const (
	LatestSymlink = "latest"
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
	matches := filenameRegexp.FindStringSubmatch(filename)
	if matches == nil {
		err = errors.New("could not parse filename")
	} else {
		rev, err = strconv.Atoi(matches[1])
		title = matches[2]
		by = matches[3]
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
	if err := os.MkdirAll(fileLoc, 0755); err != nil {
		return nil, err
	}
	return &DiskFile{
		id:       id,
		location: fileLoc,
	}, nil
}

func (ds *DiskStore) Find(id string) (File, error) {
	return nil, nil
}

func (ds *DiskStore) Delete(id string) error {
	return nil
}

type DiskFile struct {
	id       string
	location string
}

func (df *DiskFile) Id() (string, error) {
	return df.id, nil
}

func (df *DiskFile) Visibility() (int, error) {
	return 0, nil
}

func (df *DiskFile) SetVisibility(vis int) error {
	return nil
}

func (df *DiskFile) Revision(num int) (Revision, error) {
	return nil, nil
}

func (df *DiskFile) latestRevisionNum() (int, error) {
	f, err := os.Readlink(df.latestPath())
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	num, _, _, err := ParseFilename(f)
	return num, err
}

func (df *DiskFile) Latest() (Revision, error) {
	return nil, nil
}

func (df *DiskFile) latestPath() string {
	return path.Join(df.location, LatestSymlink)
}

func (df *DiskFile) Save(body io.Reader, title, by string) (Revision, error) {
	_, err := df.latestRevisionNum()
	if err != nil {
		return nil, err
	}
	return nil, nil
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

type DiskRevision struct{}

func (dr *DiskRevision) Title() (string, error) {
	return "", nil
}

func (dr *DiskRevision) Body() (io.Reader, error) {
	return nil, nil
}

func (dr *DiskRevision) By() (string, error) {
	return "", nil
}

func (dr *DiskRevision) Num() (int, error) {
	return 0, nil
}

func (dr *DiskRevision) At() (time.Time, error) {
	return time.Now(), nil
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
