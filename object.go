package main

import (
	"reflect"

	"gopkg.in/mgo.v2/bson"
)

type Diff struct {
	Version   int
	Time      int64
	Additions map[string]interface{}    `bson:",omitempty"`
	Removals  map[string]interface{}    `bson:",omitempty"`
	Changes   map[string][2]interface{} `bson:",omitempty"`
}

type Object struct {
	ID      bson.ObjectId `bson:"_id" json:"_id,omitempty"`
	Key     string
	Value   map[string]interface{}
	Version int
	Time    int64
	History []*Diff `json:",omitempty"`
}

// Diff returns the set of changes representing the different between two
// documents. The `a` document is being compared against the `b` document,
// therefore the output will be relative to `a`. Currently this only diffs
// the top-level keys and does not recurse into sub-documents.
func diff(a, b map[string]interface{}) *Diff {
	// No existing document to compare, a is an addition.
	if b == nil {
		return &Diff{
			Additions: a,
		}
	}

	var (
		ok     bool
		ak, bk string
		av, bv interface{}

		adds    = make(map[string]interface{})
		removes = make(map[string]interface{})
		changes = make(map[string][2]interface{})
	)

	// Additions and changes.
	for ak, av = range a {
		// Key does not exist in b, mark as addition.
		if bv, ok = b[ak]; !ok {
			adds[ak] = av

			// Compare a and b values.
		} else if !reflect.DeepEqual(av, bv) {
			changes[ak] = [2]interface{}{bv, av}
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
	d := Diff{}

	if len(adds) > 0 {
		d.Additions = adds
	}

	if len(removes) > 0 {
		d.Removals = removes
	}

	if len(changes) > 0 {
		d.Changes = changes
	}

	return &d
}
