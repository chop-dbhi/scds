package main

import (
	"reflect"
	"testing"
)

func TestEmpty(t *testing.T) {
	if Diff(nil, nil) != nil {
		t.Error("two nils should equal nil")
	}
}

func TestCreate(t *testing.T) {
	doc := map[string]interface{}{
		"foo": 1,
	}

	if !reflect.DeepEqual(Diff(nil, doc).Additions, doc) {
		t.Error("create should be an addition")
	}
}

func TestDelete(t *testing.T) {
	doc := map[string]interface{}{
		"foo": 1,
	}

	if !reflect.DeepEqual(Diff(doc, nil).Removals, doc) {
		t.Error("delete should be a removal")
	}
}

func TestAdd(t *testing.T) {
	b := map[string]interface{}{
		"foo": 1,
	}

	a := map[string]interface{}{
		"foo": 1,
		"bar": 2,
	}

	d := map[string]interface{}{
		"bar": 2,
	}

	if !reflect.DeepEqual(Diff(b, a).Additions, d) {
		t.Error("failed to add")
	}
}

func TestRemove(t *testing.T) {
	b := map[string]interface{}{
		"foo": 1,
		"bar": 2,
	}

	a := map[string]interface{}{
		"foo": 1,
	}

	d := map[string]interface{}{
		"bar": 2,
	}

	if !reflect.DeepEqual(Diff(b, a).Removals, d) {
		t.Error("failed to remove")
	}
}

func TestChange(t *testing.T) {
	b := map[string]interface{}{
		"foo": 1,
	}

	a := map[string]interface{}{
		"foo": 2,
	}

	d := map[string]Change{
		"foo": Change{1, 2},
	}

	r := Diff(b, a)

	if !reflect.DeepEqual(r.Changes, d) {
		t.Errorf("failed to change: %v", r.Changes)
	}
}
