package cfghive

import (
	"errors"
	"fmt"
)

// MemHive A hive that is memory resident.
type MemHive struct {
	// The root hive.
	data       map[string]HiveValue
	hasChanges bool
	inMemory   bool
}

func (h *MemHive) GetString(key string) (*string, error) {
	//TODO implement me
	panic("implement me")
}

// NewMemHive Creates a new file hive.
func NewMemHive() (*MemHive, error) {
	h := &MemHive{hasChanges: false, inMemory: true}
	h.data = make(map[string]HiveValue)
	return h, nil
}

// Characteristics Gets the characteristics of the hive.
func (h *MemHive) Characteristics() HiveCharacteristics {
	return HiveCharacteristics{false, false, false}
}

// Load Loads the hive from the file.
// the file is a JSON file.
func (h *MemHive) Load() error {
	return errors.New("not implemented")
}

func (h *MemHive) Get(key string) (*HiveValue, error) {
	path := pathToKeys(key)
	if len(path) == 0 {
		return nil, errors.New("a key must have at least one path element")
	}
	search := h.data
	for i, pf := range path {
		if i == len(path)-1 {
			val, ok := search[pf]
			if !ok {
				return nil, fmt.Errorf("key %s does not exist", key)
			}
			return &val, nil
		}
		next, ok := search[pf]
		if !ok {
			return nil, fmt.Errorf("key %s does not exist", key)
		}
		if next.IsStoredType(HiveTypeSub) {
			search, _ = next.Sub()
		} else {
			return nil, fmt.Errorf("%s is not at the path leaf, and is not a subhive", pf)
		}
	}
	return nil, fmt.Errorf("key %s does not exist", key)
}

func (h *MemHive) GetBool(key string) (bool, error) {
	v, err := h.Get(key)
	if err != nil {
		return false, err
	}
	b, err := v.Bool()
	if err != nil {
		return false, err
	}
	return b, nil
}

func (h *MemHive) GetInt(key string) (int, error) {
	v, err := h.Get(key)
	if err != nil {
		return 0, err
	}
	i, err := v.Int()
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (h *MemHive) GetFloat(key string) (float64, error) {
	v, err := h.Get(key)
	if err != nil {
		return 0, err
	}
	f, err := v.Float64()
	if err != nil {
		return 0, err
	}
	return f, nil
}

func (h *MemHive) Set(key string, value interface{}) error {
	path := pathToKeys(key)
	if len(path) == 0 {
		return errors.New("a key must have at least one path element")
	}
	search := h.data
	for i, pf := range path {
		if i == len(path)-1 {
			val, err := NewHiveValue(value)
			if err != nil {
				return err
			}
			search[pf] = val
			return nil
		}
		next, ok := search[pf]
		if !ok {
			return fmt.Errorf("key %s does not exist", key)
		}
		sub, err := next.Sub()
		if err != nil {
			return fmt.Errorf("%s is not at the path leaf, and is not a subhive", pf)
		}
		search = sub
	}
	h.hasChanges = true
	return nil
}

func (h *MemHive) SetBool(key string, value bool) error {
	err := h.Set(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (h *MemHive) SetInt(key string, value int) error {
	err := h.Set(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (h *MemHive) SetFloat(key string, value float64) error {
	err := h.Set(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (h *MemHive) SetString(key string, value string) error {
	err := h.Set(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (h *MemHive) Delete(key string) {
	path := pathToKeys(key)
	if len(path) == 0 {
		return
	}
	search := h.data
	for i, pf := range path {
		if i == len(path)-1 {
			delete(search, pf)
			return
		}
		next, ok := search[pf]
		if !ok {
			return
		}
		sub, err := next.Sub()
		if err != nil {
			return
		}
		search = sub
	}
	h.hasChanges = true
}

func (h *MemHive) NewSub(key string) {
	v, _ := NewHiveValue(make(map[string]HiveValue))
	h.Set(key, v)
}

func (h *MemHive) Rollback() (bool, error) {
	return false, nil
}

func (h *MemHive) Commit() (bool, error) {
	return false, nil
}

func (h *MemHive) Save() error {
	return nil
}

func (h *MemHive) GetData() *map[string]HiveValue {
	return &h.data
}
