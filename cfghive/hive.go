package cfghive

import (
	"strings"
)

type HiveCharacteristics struct {
	// Is the hive implementation transactional?
	IsTxn bool
	// Is the hive implementation persistent?
	IsPersistent bool
	// Is the hive implementation thread-safe?
	IsThreadSafe bool
}

type Hive interface {
	// Characteristics Gets the characteristics of the hive.
	Characteristics() HiveCharacteristics

	// Load the hive from some persistent storage.
	Load() error

	// Get Gets a value from the hive.
	// Returns nil, false if the value does not exist.
	Get(key string) (*HiveValue, error)
	GetBool(key string) (bool, error)
	GetInt(key string) (int, error)
	GetFloat(key string) (float64, error)
	GetString(key string) (*string, error)

	// Set Sets a value in the hive.
	// If the value already exists, it's replaced.
	Set(key string, value interface{}) error
	SetBool(key string, value bool) error
	SetInt(key string, value int) error
	SetFloat(key string, value float64) error
	SetString(key string, value string) error

	// Delete Deletes a value from the hive.
	// Returns the old value, or nil if the value did not exist.
	Delete(key string)

	// NewSub Create a sub-hive with the given key.
	NewSub(key string)

	// Rollback Rolls back the hive to the last commit.
	// Returns true if the hive was rolled back.
	// For some hives, this is a no-op.
	Rollback() (bool, error)

	// Commit Commits the hive.
	// Returns true if the hive was committed.
	// For some hives, this is a no-op.
	Commit() (bool, error)

	// Save Saves the hive to some persistent storage.
	// Unlike Commit this always saves the hive.
	Save() error

	// GetData Get the data of the hive.
	GetData() *map[string]HiveValue
}

func pathToKeys(path string) []string {
	path = strings.TrimSuffix(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}
