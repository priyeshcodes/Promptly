// internal/task/store.go
package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.etcd.io/bbolt"
	bolt "go.etcd.io/bbolt"
	"time"
)

const taskBucket = "tasks"

type TaskStore struct {
	db *bbolt.DB
}

func NewTaskStore(dbPath string) (*TaskStore, error) {
	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	// Ensure the bucket exists
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(taskBucket))
		return err
	})

	if err != nil {
		return nil, err
	}

	return &TaskStore{db: db}, nil
}

func (s *TaskStore) SaveTask(t *Task) error {
	t.ID = uuid.New().String()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))
		return b.Put([]byte(t.ID), data)
	})
}

func (s *TaskStore) GetAllTasks() ([]Task, error) {
	var tasks []Task

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))
		return b.ForEach(func(k, v []byte) error {
			var t Task
			if err := json.Unmarshal(v, &t); err != nil {
				return err
			}
			tasks = append(tasks, t)
			return nil
		})
	})

	return tasks, err
}

func (s *TaskStore) Close() error {
	return s.db.Close()
}

func (s *TaskStore) MarkTaskComplete(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))
		v := b.Get([]byte(id))
		if v == nil {
			return errors.New("task not found")
		}

		var t Task
		if err := json.Unmarshal(v, &t); err != nil {
			return err
		}

		t.Completed = true
		t.UpdatedAt = time.Now()

		updatedData, err := json.Marshal(&t)
		if err != nil {
			return err
		}

		return b.Put([]byte(id), updatedData)
	})
}

// GetTaskByID fetches a task by its ID
func (s *TaskStore) GetTaskByID(id string) (*Task, error) {
	var t Task

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))
		if b == nil {
			return fmt.Errorf("task bucket not found")
		}

		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("task with ID %s not found", id)
		}

		return json.Unmarshal(data, &t)
	})

	if err != nil {
		return nil, err
	}
	return &t, nil
}

// DeleteTask deletes a task by its ID
func (s *TaskStore) DeleteTask(id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))
		if b == nil {
			return fmt.Errorf("task bucket not found")
		}
		return b.Delete([]byte(id))
	})
}

func (s *TaskStore) UpdateTask(t *Task) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		data, err := json.Marshal(t)
		if err != nil {
			return err
		}
		return b.Put([]byte(t.ID), data)
	})
}
