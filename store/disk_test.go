package store

import (
	"bytes"
	"os"
	"testing"
)

func TestParseFilename(t *testing.T) {
	rev, title, by, err := ParseFilename("123;This is the file name.json;bob")
	if err != nil {
		t.Fatal(err)
	}
	if rev != 123 {
		t.Fatalf("Got %d", rev)
	}
	if title != "This is the file name.json" {
		t.Fatalf("Got %s", title)
	}
	if by != "bob" {
		t.Fatalf("Got %s", title)
	}
}

func TestDiskStore_Create(t *testing.T) {
	var s Store
	s = NewDiskStore(os.TempDir())
	f, err := s.Create()
	if err != nil {
		t.Fatal(err)
	}
	id, err := f.Id()
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("Expected the file to have an ID.")
	}
}

func TestDiskFile_Save(t *testing.T) {
	var s Store
	s = NewDiskStore(os.TempDir())
	f, err := s.Create()
	if err != nil {
		t.Fatal(err)
	}
	id, err := f.Id()
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("Expected the file to have an ID.")
	}
	rev, err := f.Save(bytes.NewBufferString("blah"), "blah file", "egg")
	if err != nil {
		t.Fatal(err)
	}
	title, err := rev.Title()
	if err != nil {
		t.Fatal(err)
	}
	if title != "blah file" {
		t.Fatalf("Expected blah file, got %s", title)
	}
}

func TestDiskFile_Save_Subsequent(t *testing.T) {
	var s Store
	s = NewDiskStore(os.TempDir())
	f, err := s.Create()
	if err != nil {
		t.Fatal(err)
	}
	id, err := f.Id()
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("Expected the file to have an ID.")
	}
	_, err = f.Save(bytes.NewBufferString("blah"), "blah file", "egg")
	if err != nil {
		t.Fatal(err)
	}
	rev2, err := f.Save(bytes.NewBufferString("blah2"), "blah file2", "egg2")
	if err != nil {
		t.Fatal(err)
	}
	num, err := rev2.Num()
	if err != nil {
		t.Fatal(err)
	}
	if num != 2 {
		t.Fatal(num)
	}
}

func TestDiskFile_Find(t *testing.T) {
	var s Store
	s = NewDiskStore(os.TempDir())
	f, err := s.Create()
	if err != nil {
		t.Fatal(err)
	}
	id, err := f.Id()
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("Expected the file to have an ID.")
	}
	rev, err := f.Save(bytes.NewBufferString("blah"), "blah file", "egg")
	if err != nil {
		t.Fatal(err)
	}
	f2, err := s.Find(id)
	if err != nil {
		t.Fatal(err)
	}
	rev, err = f2.Latest()
	if err != nil {
		t.Fatal(err)
	}
	title, err := rev.Title()
	if err != nil {
		t.Fatal(err)
	}
	if title != "blah file" {
		t.Fatalf("Expected blah file, got %s", title)
	}
}

func TestDiskStore_Delete(t *testing.T) {
	var s Store
	s = NewDiskStore(os.TempDir())
	f, err := s.Create()
	if err != nil {
		t.Fatal(err)
	}
	id, err := f.Id()
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("Expected the file to have an ID.")
	}
	if err := s.Delete(id); err != nil {
		t.Fatal(err)
	}
}
