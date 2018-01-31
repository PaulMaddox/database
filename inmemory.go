package database

import (
	"reflect"
	"strings"
)

// InMemoryDBAdapter provides helper methods for in memory database operations
type InMemoryDBAdapter struct {
	db map[string][]interface{}
}

// InMemoryDBTable provides helper methods for in memory database operations
type InMemoryDBTable struct {
	table string
	db    map[string][]interface{}
}

// NewInMemoryDatabaseAdapter creates a new database adapter for in memory database interactions
func NewInMemoryDatabaseAdapter() (*InMemoryDBAdapter, error) {
	InMemoryDB := map[string][]interface{}{}
	return &InMemoryDBAdapter{
		db: InMemoryDB,
	}, nil
}

// Table returns a in memory table that has methods to interact with rows
func (m *InMemoryDBAdapter) Table(name string) Table {
	return &InMemoryDBTable{
		table: name,
		db:    m.db,
	}
}

// Get an item from the database
func (m *InMemoryDBTable) Get(id string, v interface{}) error {

	for i, item := range m.db[m.table] {

		itemID, err := GetID(item)
		if err != nil {
			continue
		}

		if itemID == id {
			val := reflect.ValueOf(m.db[m.table][i])
			destVal := reflect.ValueOf(v)
			if destVal.Kind() == reflect.Ptr && destVal.Kind() == val.Kind() {
				if destElem := destVal.Elem(); destElem.CanSet() {
					destElem.Set(val.Elem())
					return nil
				}
			}
		}
	}

	return ErrRecordNotFound

}

// Put item in the database
func (m *InMemoryDBTable) Put(v interface{}) error {

	// Initalize the table if it doesn't exist
	if m.db[m.table] == nil {
		m.db[m.table] = []interface{}{}
	}

	itemID, err := GetID(v)
	if err != nil {
		return ErrInvalidRecordID
	}

	// If the item with the same ID exists, overwrite it
	for i, existing := range m.db[m.table] {
		existingID, err := GetID(existing)
		if err != nil {
			continue
		}
		if itemID == existingID {
			m.db[m.table][i] = v
			return nil
		}
	}

	m.db[m.table] = append(m.db[m.table], v)
	return nil

}

// Update (overwrite) an item in the database
func (m *InMemoryDBTable) Update(v interface{}) error {
	return m.Put(v)
}

// Delete an item from the database
func (m *InMemoryDBTable) Delete(id string) error {

	for i, item := range m.db[m.table] {
		itemID, err := GetID(item)
		if err != nil {
			continue
		}
		if itemID == id {
			m.db[m.table] = append(m.db[m.table][:i], m.db[m.table][i+1:]...)
			return nil
		}
	}

	return ErrRecordNotFound

}

// GetID uses reflection to get the 'id' field of an object
// or returns ErrInvalidRecordID if no 'id' field is found
func GetID(item interface{}) (string, error) {

	ref := reflect.ValueOf(item).Elem()
	for i := 0; i < ref.NumField(); i++ {
		if strings.ToLower(ref.Type().Field(i).Name) == "id" {
			return ref.Field(i).String(), nil
		}
	}

	return "", ErrInvalidRecordID

}
