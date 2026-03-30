package tinydb

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// DB represents the local JSON database engine.
type DB struct {
	filepath string
	mu       sync.RWMutex
	data     map[string]map[string]any
}

// New initializes the database and loads existing data if the file exists.
func New(filepath string) (*DB, error) {
	db := &DB{
		filepath: filepath,
		data:     make(map[string]map[string]any),
	}

	if err := db.load(); err != nil {
		return nil, err
	}

	return db, nil
}

// load reads the JSON file from disk into memory.
func (db *DB) load() error {
	bytes, err := os.ReadFile(db.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, &db.data)
}

// save writes the in-memory map back to the JSON file on disk.
func (db *DB) save() error {
	bytes, err := json.MarshalIndent(db.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(db.filepath, bytes, 0644)
}

// generateID creates a random 16-character string.
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Insert adds a new record to a collection and returns the generated ID.
func (db *DB) Insert(collection string, record any) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	id := generateID()

	if db.data[collection] == nil {
		db.data[collection] = make(map[string]any)
	}

	db.data[collection][id] = record

	if err := db.save(); err != nil {
		return "", err
	}

	return id, nil
}

// Read fetches a specific record by its collection and ID.
func (db *DB) Read(collection string, id string) (any, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.data[collection] == nil {
		return nil, fmt.Errorf("collection '%s' not found", collection)
	}

	record, exists := db.data[collection][id]
	if !exists {
		return nil, fmt.Errorf("record '%s' not found", id)
	}

	return record, nil
}

// Update modifies an existing record by its ID and saves the changes.
func (db *DB) Update(collection string, id string, record any) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.data[collection] == nil {
		return fmt.Errorf("collection '%s' not found", collection)
	}

	if _, exists := db.data[collection][id]; !exists {
		return fmt.Errorf("record '%s' not found", id)
	}

	db.data[collection][id] = record

	return db.save()
}

// Delete removes a specific record from a collection and saves the file.
func (db *DB) Delete(collection string, id string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.data[collection] == nil {
		return fmt.Errorf("collection '%s' not found", collection)
	}

	delete(db.data[collection], id)

	return db.save()
}
