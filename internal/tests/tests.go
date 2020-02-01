package tests

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vikramcse/the-service/internal/platform/database"
	"github.com/vikramcse/the-service/internal/platform/database/databasetest"
	"github.com/vikramcse/the-service/internal/schema"
)

// NewUnit creates a test database in docker container. It created the required
// database structure. The databases are empty at this stage.

// It does not return error if anything fails, as this is intended for testing
// only. Instead it will call Fatal on provided testing.T

// It returns the database object to use as well as a function to call at the
// end of the test
func NewUnit(t *testing.T) (*sqlx.DB, func()) {
	t.Helper()

	c := databasetest.StartContainer(t)

	db, err := database.Open(database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       c.Host,
		Name:       "postgres",
		DisableTLS: true,
	})
	if err != nil {
		t.Fatalf("opening database connection: %v", err)
	}
	t.Log("waiting for database to be ready")

	// wait for the database to be ready. Wait 100ms longer between each
	// attempt. Do not try more than 20 times
	var pingError error
	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		pingError := db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}

	if pingError != nil {
		databasetest.DumpContainerLogs(t, c)
		databasetest.StopContainer(t, c)
		t.Fatalf("waiting for database to be ready: %v", pingError)
	}

	if err := schema.Migrate(db); err != nil {
		databasetest.StopContainer(t, c)
		t.Fatalf("migrating: %s", err)
	}

	// teradown is the funtion that should be invoked when the caller is done
	// with the database
	teradown := func() {
		t.Helper()
		db.Close()
		databasetest.StopContainer(t, c)
	}

	return db, teradown
}
