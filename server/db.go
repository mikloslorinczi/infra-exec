package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

// DBI is the interface of the JSON DB.
type DBI interface {
	Load() error
	save() error
	add(t common.Task) (string, error)
	remove(id string) (bool, error)
	query(id string) (common.Task, bool)
	queryAll() ([]common.Task, bool)
	update(id string, t common.Task) (bool, error)
}

type jsonDB struct {
	rwMutex sync.RWMutex
	path    string
	data    []common.Task
}

// ConnectJSONDB returns a pointer to a jsonDB, working with the given path.
func connectJSONDB(path string) DBI {
	return &jsonDB{path: path}
}

func (db *jsonDB) Load() error {
	db.rwMutex.Lock()
	defer db.rwMutex.Unlock()
	file, err := os.OpenFile(db.path, os.O_RDONLY|os.O_CREATE, 0600)
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
	if len(bytes) > 0 {
		err = json.Unmarshal(bytes, &db.data)
		if err != nil {
			return errors.Wrapf(err, "Cannot decode JSON file %v", db.path)
		}
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

func (db *jsonDB) add(t common.Task) (string, error) {
	t.ID = xid.New().String()
	t.Node = "None"
	t.Status = "added"
	db.data = append(db.data, t)
	err := db.save()
	if err != nil {
		return "", err
	}
	return t.ID, nil
}

func (db *jsonDB) remove(id string) (bool, error) {
	removed := false
	for i, task := range db.data {
		if task.ID == id {
			db.data = append(db.data[:i], db.data[i+1:]...)
			err := db.save()
			if err != nil {
				return removed, err
			}
			removed = true
		}
	}
	return removed, nil
}

func (db *jsonDB) query(id string) (common.Task, bool) {
	t := common.Task{}
	found := false
	for _, task := range db.data {
		if task.ID == id {
			t = task
			found = true
		}
	}
	return t, found
}

func (db *jsonDB) queryAll() ([]common.Task, bool) {
	return db.data, len(db.data) > 0
}

func (db *jsonDB) update(id string, t common.Task) (bool, error) {
	updated := false
	for i, task := range db.data {
		if task.ID == id {
			db.data[i].Node = t.Node
			db.data[i].Status = t.Status
			db.data[i].Status = t.Status
			db.data[i].Command = t.Command
			err := db.save()
			if err != nil {
				return updated, err
			}
			updated = true
		}
	}
	return updated, nil
}
