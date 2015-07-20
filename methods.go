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

func Log(s *mgo.Session, k string) ([]*Revision, error) {
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

func Put(s *mgo.Session, k string, v map[string]interface{}) (*Revision, error) {
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

		return o.History[0], nil
	}

	if err != nil {
		return nil, err
	}

	r, changed, err = update(s, o, v)

	if err != nil {
		return nil, err
	}

	if changed {
		return r, nil
	}

	return nil, nil
}
