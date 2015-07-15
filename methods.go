package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const colName = "objects"

// Initializes a session with mongodb.
func initDB() *mgo.Session {
	session, err := mgo.Dial(mongoURI)

	if err != nil {
		log.Fatal(err)
	}

	if err = session.DB("").C(colName).EnsureIndexKey("key"); err != nil {
		log.Fatal(err)
	}

	return session
}

func Get(s *mgo.Session, k string) (*Object, error) {
	c := s.DB("").C(colName)

	// Query.
	q := bson.M{
		"key": k,
	}

	// Projection.
	p := bson.M{
		"_id":     0,
		"history": 0,
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

func Log(s *mgo.Session, k string) ([]*Diff, error) {
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

	d := diff(v, nil)
	d.Version = 1
	d.Time = time.Now().UTC().Unix()

	o := Object{
		ID:      bson.NewObjectId(),
		Key:     k,
		Value:   v,
		Version: d.Version,
		Time:    d.Time,
		History: []*Diff{d},
	}

	err := c.Insert(&o)

	return &o, true, err
}

// Updates an existing objects.
func update(s *mgo.Session, o *Object, v map[string]interface{}) (*Diff, bool, error) {
	d := diff(v, o.Value)

	if d == nil {
		return nil, false, nil
	}

	c := s.DB("").C(colName)

	// Increment the version.
	d.Version = o.Version + 1
	d.Time = time.Now().UTC().Unix()

	// Keys to update.
	chg := mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": bson.M{
				"version": d.Version,
				"time":    d.Time,
				"value":   v,
			},
			"$push": bson.M{
				"history": d,
			},
		},
	}

	// Apply the change.
	if _, err := c.Find(bson.M{"_id": o.ID}).Apply(chg, o); err != nil {
		return d, true, err
	}

	return d, true, nil
}

func Put(s *mgo.Session, k string, v map[string]interface{}) (*Diff, error) {
	c := s.DB("").C(colName)

	// Query.
	q := bson.M{
		"key": k,
	}

	var (
		d       *Diff
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

		return o.History[0], nil
	}

	if err != nil {
		return nil, err
	}

	d, changed, err = update(s, o, v)

	if err != nil {
		return nil, err
	}

	if changed {
		return d, nil
	}

	return nil, nil
}
