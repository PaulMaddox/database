package database_test

import (
	"testing"

	"github.com/paulmaddox/database"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

type TestRecord struct {
	ID        string
	Name      string
	Url       string
	Picture   string
	StartDate int64
	EndDate   int64
}

// TestCreate tests that creating new events works as expected
func TestPutReadUpdateDelete(t *testing.T) {

	db, err := database.NewInMemoryDatabaseAdapter()
	assert.NoError(t, err)
	table := db.Table("test-table")

	event := &TestRecord{
		ID:        uuid.New(),
		Name:      "Test Event",
		Url:       "https://test.event.example.com",
		Picture:   "https://test.event.example.com/event.jpg",
		StartDate: 1516785100,
		EndDate:   1516785200,
	}

	// Write to database
	assert.Nil(t, table.Put(event))

	// Read from database
	read := &TestRecord{}
	assert.Nil(t, table.Get(event.ID, read))
	assert.EqualValues(t, event, read)

	// Update database
	event.Name = "Updated Event"
	assert.Nil(t, table.Update(event))

	// Read from database
	updated := &TestRecord{}
	assert.Nil(t, table.Get(event.ID, updated))
	assert.EqualValues(t, event, updated)

	// Delete from database
	assert.Nil(t, table.Delete(event.ID))

	// Check it's actually deleted
	deleted := &TestRecord{}
	assert.Equal(t, database.ErrRecordNotFound, table.Get(event.ID, deleted))

}
