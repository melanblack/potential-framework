package cfghive_test

import (
	"testing"

	"github.com/melanblack/potential-framework/cfghive"
)

func TestFitsInterface(t *testing.T) {
	h, _ := cfghive.NewMemHive()
	var _ cfghive.Hive = h
}

func TestNewMemHive(t *testing.T) {
	h, err := cfghive.NewMemHive()
	if err != nil {
		t.Fatal(err)
	}
	if h == nil {
		t.Fatal("h is nil")
	}
}

func TestSet(t *testing.T) {
	h, err := cfghive.NewMemHive()
	if err != nil {
		t.Fatal(err)
	}
	if h == nil {
		t.Fatal("h is nil")
	}
	err = h.Set("foo", "bar")
	if err != nil {
		t.Fatal(err)
	}

	err = h.Set("foo/baz", "bar")
	if err == nil {
		t.Fatal("No error when setting a value in a non-existent sub-hive")
	}
	h.NewSub("fez")
	err = h.Set("fez/baz", "bar")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	h, err := cfghive.NewMemHive()
	if err != nil {
		t.Fatal(err)
	}
	if h == nil {
		t.Fatal("h is nil")
	}
	err = h.Set("foo", "bar")
	if err != nil {
		t.Fatal(err)
	}
	v, err := h.Get("foo")
	if err != nil {
		t.Fatal(err)
	}
	if v == nil {
		t.Fatal("v is nil")
	}
	if v.(string) != "bar" {
		t.Fatal("v is not bar")
	}
}
