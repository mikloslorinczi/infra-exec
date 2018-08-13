package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/pkg/errors"
)

type task struct {
	ID      string `json:"id"`
	Node    string `json:"node"`
	Tags    string `json:"tags"`
	Status  string `json:"status"`
	Command string `json:"command"`
}

// DBI is the interface of the JSON DB.
type DBI interface {
	open() error
	save() error
	add(body []byte) (string, error)
	remove(id string) error
	query(id string) ([]byte, error)
	update(id string, body []byte) error
}

type jsonDB struct {
	rwMutex sync.RWMutex
	path    string
	data    []task
}

// ConnectJSONDB returns a pointer to a jsonDB, working with the given path.
func ConnectJSONDB(path string) DBI {
	return &jsonDB{path: path}
}

func (db *jsonDB) open() error {
	db.rwMutex.Lock()
	defer db.rwMutex.Unlock()
	file, err := os.OpenFile(db.path, os.O_RDONLY|os.O_CREATE, 0440)
	if err != nil {
		return errors.Wrapf(err, "Cannot open DB File %v for Read", db.path)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatalf("Error closing DB file %v\n%v", db.path, err)
		}
	}()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.Wrapf(err, "Cannot read from DB file %v", db.path)
	}
	err = json.Unmarshal(bytes, &db.data)
	if err != nil {
		return errors.Wrapf(err, "Cannot decode JSON file %v", db.path)
	}

	return nil
}

func (db *jsonDB) save() error {
	db.rwMutex.Lock()
	defer db.rwMutex.Unlock()
	file, err := os.OpenFile(db.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0220)
	if err != nil {
		return errors.Wrapf(err, "Cannot open DB File %v for Write", db.path)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatalf("Error closing DB file %v\n%v", db.path, err)
		}
	}()

	jsonData, err := json.Marshal(db.data)
	if err != nil {
		return errors.Wrap(err, "Cannot enncode JSON")
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return errors.Wrapf(err, "Cannot write to DB file %v", db.path)
	}

	return nil
}

func (db *jsonDB) add(body []byte) (string, error) {
	return "", nil
}

func (db *jsonDB) remove(id string) error {
	return nil
}

func (db *jsonDB) query(id string) ([]byte, error) {
	return []byte{}, nil
}

func (db *jsonDB) update(id string, body []byte) error {
	return nil
}
