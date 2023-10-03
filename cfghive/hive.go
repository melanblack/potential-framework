package cfghive

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

	// Load the hive.
	Load() error

	// Get Gets a value from the hive.
	// Returns nil, false if the value does not exist.
	Get(key string) (interface{}, error)

	// Set Sets a value in the hive.
	// If the value already exists, it's replaced.
	Set(key string, value interface{}) error

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
}
