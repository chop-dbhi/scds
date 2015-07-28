package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const colName = "objects"

func Get(cfg *Config, k string) (*Object, error) {
	return get(cfg, k, false)
}

func get(cfg *Config, k string, history bool) (*Object, error) {
	s := cfg.Mongo.Session()

	c := s.DB("").C(colName)

	// Query.
	q := bson.M{
		"key": k,
	}

	// Projection.
	p := bson.M{
		"_id": 0,
	}

	if !history {
		p["history"] = 0
	}

	var o Object

	err := c.Find(q).Select(p).One(&o)

	if err == mgo.ErrNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &o, nil
}

func Log(cfg *Config, k string) ([]*Revision, error) {
	s := cfg.Mongo.Session()

	c := s.DB("").C(colName)

	// Query.
	q := bson.M{
		"key": k,
	}

	// Projection.
	p := bson.M{
		"_id":     0,
		"history": 1,
	}

	var o Object

	err := c.Find(q).Select(p).One(&o)

	if err == mgo.ErrNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return o.History, nil

}

// Inserts an object into the store.
func insert(s *mgo.Session, k string, v map[string]interface{}) (*Object, bool, error) {
	c := s.DB("").C(colName)

	r := Diff(nil, v)
	r.Version = 1
	r.Time = time.Now().UTC().Unix()

	o := Object{
		ID:      bson.NewObjectId(),
		Key:     k,
		Value:   v,
		Version: r.Version,
		Time:    r.Time,
		History: []*Revision{r},
	}

	err := c.Insert(&o)

	return &o, true, err
}

// Updates an existing objects.
func update(s *mgo.Session, o *Object, v map[string]interface{}) (*Revision, bool, error) {
	r := Diff(o.Value, v)

	if r == nil {
		return nil, false, nil
	}

	c := s.DB("").C(colName)

	// Increment the version.
	r.Version = o.Version + 1
	r.Time = time.Now().UTC().Unix()

	// Keys to update.
	chg := mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": bson.M{
				"version": r.Version,
				"time":    r.Time,
				"value":   v,
			},
			"$push": bson.M{
				"history": r,
			},
		},
	}

	// Apply the change.
	if _, err := c.Find(bson.M{"_id": o.ID}).Apply(chg, o); err != nil {
		return r, true, err
	}

	return r, true, nil
}

func Put(cfg *Config, k string, v map[string]interface{}) (*Revision, error) {
	s := cfg.Mongo.Session()

	c := s.DB("").C(colName)

	// Query.
	q := bson.M{
		"key": k,
	}

	var (
		r       *Revision
		err     error
		changed bool
	)

	o := &Object{
		Value: make(map[string]interface{}),
	}

	err = c.Find(q).One(&o)

	// Does not exist. Insert it.
	if err == mgo.ErrNotFound {
		o, changed, err = insert(s, k, v)

		if err != nil {
			return nil, err
		}

		r = o.History[0]

		if err = SendNotifications(cfg, o, r); err != nil {
			fmt.Fprintln(os.Stderr, "[smtp] error sending email:", err)
		}

		return o.History[0], nil
	}

	if err != nil {
		return nil, err
	}

	r, changed, err = update(s, o, v)

	if err != nil {
		return nil, err
	}

	// Object changed.
	if changed {
		if err = SendNotifications(cfg, o, r); err != nil {
			fmt.Fprintln(os.Stderr, "[smtp] error sending email:", err)
		}

		return r, nil
	}

	return nil, nil
}
