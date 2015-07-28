package main

import (
	"errors"
	"reflect"

	"gopkg.in/mgo.v2/bson"
)

var (
	ErrUnknownRevision = errors.New("unknown revision")
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

func applyRevision(o *Object, r *Revision) {
	var (
		key string
		val interface{}
		chg Change
	)

	o.Version = r.Version
	o.Time = r.Time

	if r.Additions != nil {
		for key, val = range r.Additions {
			o.Value[key] = val
		}
	}

	if r.Removals != nil {
		for key, val = range r.Removals {
			delete(o.Value, key)
		}
	}

	if r.Changes != nil {
		for key, chg = range r.Changes {
			o.Value[key] = chg.After
		}
	}

}

// AtVersion reverts the objects to the specified version.
func (o *Object) AtVersion(v int) *Object {
	if v == 0 {
		return nil
	}

	n := Object{
		ID:    o.ID,
		Key:   o.Key,
		Value: make(map[string]interface{}),
	}

	for i, rev := range o.History {
		if rev.Version > v {
			n.History = o.History[:i]
			break
		}

		applyRevision(&n, rev)
	}

	return &n
}

// AtTime reverts the object to the state as of the specified time.
func (o *Object) AtTime(t int64) *Object {
	n := Object{
		ID:    o.ID,
		Key:   o.Key,
		Value: make(map[string]interface{}),
	}

	for i, rev := range o.History {
		if rev.Time > t {
			n.History = o.History[:i]
			break
		}

		applyRevision(&n, rev)
	}

	// The time is earlier than the first revision of this object.
	if n.Version == 0 {
		return nil
	}

	return &n
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
