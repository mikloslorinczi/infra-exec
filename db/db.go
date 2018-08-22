package db

import (
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/mikloslorinczi/infra-exec/common"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

// I is the interface of the JSON DB.
type I interface {
	Load() error
	save() error
	Add(t common.Task) (string, error)
	Remove(id string) (bool, error)
	Query(id string) (common.Task, bool)
	QueryAll() ([]common.Task, bool)
	Update(id string, t common.Task) (bool, error)
}

type jsonDB struct {
	rwMutex sync.RWMutex
	path    string
	data    common.Tasks
}

// ConnectJSONDB returns a pointer to a jsonDB, working with the given path.
func ConnectJSONDB(path string) I {
	return &jsonDB{path: path}
}

// Load opens the JSON file and loads it to data
func (db *jsonDB) Load() error {
	db.rwMutex.Lock()
	defer db.rwMutex.Unlock()
	file, err := os.OpenFile(db.path, os.O_RDWR|os.O_CREATE, 0600)
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
		err = common.FromJSON(&db.data, bytes)
		if err != nil {
			return errors.Wrapf(err, "Cannot decode JSON file %v", db.path)
		}
	}
	return nil
}

// save saves the data to the JSON file.
func (db *jsonDB) save() error {
	var jsonData []byte
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
	err = common.ToJSON(&db.data, &jsonData)
	if err != nil {
		return errors.Wrap(err, "Cannot enncode JSON")
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return errors.Wrapf(err, "Cannot write to DB file %v", db.path)
	}

	return nil
}

// Add appends the data with a Task and saves it to the JSON file.
func (db *jsonDB) Add(t common.Task) (string, error) {
	t.ID = xid.New().String()
	t.Node = "None"
	t.Status = "Created"
	db.data = append(db.data, t)
	err := db.save()
	if err != nil {
		return "", err
	}
	return t.ID, nil
}

// Remove deletes a task by ID.
// Returns a bool indicating if the remove was success, and an optional io error.
func (db *jsonDB) Remove(id string) (bool, error) {
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

// Query returns a Task by ID. And a bool indicating if,
// the Task was found or not.
func (db *jsonDB) Query(id string) (common.Task, bool) {
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

// QueryAll returns the whole db and a bool indicating if it contains any element.
func (db *jsonDB) QueryAll() ([]common.Task, bool) {
	return db.data, len(db.data) > 0
}

// Update overwrites a Task (found by ID) with the argument Task's fileds (except ID).
func (db *jsonDB) Update(id string, t common.Task) (bool, error) {
	updated := false
	for i, task := range db.data {
		if task.ID == id {
			db.data[i].Node = t.Node
			db.data[i].Tags = t.Tags
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
