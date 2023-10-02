package cfghive

import (
	"errors"
	"fmt"
	"strings"
)

// A hive that is memory resident.
type MemHive struct {
	// The root hive.
	data       map[string]interface{}
	hasChanges bool
	inMemory   bool
}

// Creates a new file hive.
func NewMemHive() (*MemHive, error) {
	h := &MemHive{hasChanges: false, inMemory: true}
	h.data = make(map[string]interface{})
	return h, nil
}

// Gets the characteristics of the hive.
func (h *MemHive) Characteristics() HiveCharacteristics {
	return HiveCharacteristics{false, false, false}
}

// Loads the hive from the file.
// the file is a JSON file.
func (h *MemHive) Load() error {
	return errors.New("not implemented")
}

func (h *MemHive) Get(key string) (interface{}, error) {
	path := strings.Split(key, "/")
	if len(path) < 1 {
		return nil, errors.New("a key must have at least one path element")

	}
	search := h.data
	for i, pf := range path {
		if i == len(path)-1 {
			return search[pf], nil
		}
		next, ok := search[pf]
		if !ok {
			return nil, fmt.Errorf("key %s does not exist", key)
		}
		search, ok = next.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%s is not at the path leaf, and is not a subhive", pf)
		}
	}
	return nil, fmt.Errorf("key %s does not exist", key)
}

func (h *MemHive) Set(key string, value interface{}) error {
	path := strings.Split(key, "/")
	if len(path) < 1 {
		return errors.New("a key must have at least one path element")
	}
	search := h.data
	for i, pf := range path {
		if i == len(path)-1 {
			search[pf] = value
			return nil
		}
		next, ok := search[pf]
		if !ok {
			return fmt.Errorf("key %s does not exist", key)
		}
		search, ok = next.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%s is not at the path leaf, and is not a subhive", pf)
		}
	}
	h.hasChanges = true
	return nil
}

func (h *MemHive) Delete(key string) {
	panic("not implemented")
}

func (h *MemHive) NewSub(key string) {
	h.Set(key, make(map[string]interface{}))
}

func (h *MemHive) Rollback() (bool, error) {
	return false, nil
}

func (h *MemHive) Commit() (bool, error) {
	return false, nil
}
