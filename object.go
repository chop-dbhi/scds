package main

import (
	"reflect"

	"gopkg.in/mgo.v2/bson"
)

type Change struct {
	Before interface{}
	After  interface{}
}

type Revision struct {
	Version   int
	Time      int64
	Additions map[string]interface{} `bson:",omitempty"`
	Removals  map[string]interface{} `bson:",omitempty"`
	Changes   map[string]Change      `bson:",omitempty"`
}

type Object struct {
	ID      bson.ObjectId `bson:"_id" json:"_id,omitempty"`
	Key     string
	Value   map[string]interface{}
	Version int
	Time    int64
	History []*Revision `json:",omitempty" yaml:",omitempty"`
}

// Diff returns the set of changes representing the different between two
// documents. Compares the before (`b`) and after (`a`) state of the document.
// Currently this only diffs the top-level keys and does not recurse into
// sub-documents.
func Diff(b, a map[string]interface{}) *Revision {
	if (a == nil || len(a) == 0) && (b == nil || len(b) == 0) {
		return nil
	}

	// No existing document to compare, a is an addition.
	if b == nil || len(b) == 0 {
		return &Revision{
			Additions: a,
		}
	}

	// Next state is nil, b is a removal.
	if a == nil || len(a) == 0 {
		return &Revision{
			Removals: b,
		}
	}

	var (
		ok     bool
		ak, bk string
		av, bv interface{}

		adds    = make(map[string]interface{})
		removes = make(map[string]interface{})
		changes = make(map[string]Change)
	)

	// Additions and changes.
	for ak, av = range a {
		// Key does not exist in b, mark as addition.
		if bv, ok = b[ak]; !ok {
			adds[ak] = av

			// Compare a and b values.
		} else if !reflect.DeepEqual(av, bv) {
			changes[ak] = Change{
				Before: bv,
				After:  av,
			}
		}
	}

	// Removals.
	for bk, bv = range b {
		// Keys in b that no longer exist in a.
		if _, ok = a[bk]; !ok {
			removes[bk] = bv
		}
	}

	// No difference.
	if len(adds) == 0 && len(removes) == 0 && len(changes) == 0 {
		return nil
	}

	// Build the diff.
	r := Revision{}

	if len(adds) > 0 {
		r.Additions = adds
	}

	if len(removes) > 0 {
		r.Removals = removes
	}

	if len(changes) > 0 {
		r.Changes = changes
	}

	return &r
}
