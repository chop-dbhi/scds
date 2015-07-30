package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func ErrInvalidKey(k string) error {
	return fmt.Errorf("Key contains invalid chars: %s", k)
}

var keyRegexp *regexp.Regexp

func init() {
	keyRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)*$`)
}

func checkKey(key string) bool {
	return keyRegexp.MatchString(key)
}

func Keys(cfg *Config) ([]string, error) {
	c := cfg.Mongo.Objects()

	// Projection.
	p := bson.M{
		"_id": 0,
		"key": 1,
	}

	var objs []*Object

	if err := c.Find(nil).Select(p).All(&objs); err != nil {
		return nil, err
	}

	keys := make([]string, len(objs))

	for i, obj := range objs {
		keys[i] = obj.Key
	}

	return keys, nil
}

func Get(cfg *Config, k string) (*Object, error) {
	return get(cfg, k, false)
}

func get(cfg *Config, k string, history bool) (*Object, error) {
	if !checkKey(k) {
		return nil, ErrInvalidKey(k)
	}

	c := cfg.Mongo.Objects()

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
	if !checkKey(k) {
		return nil, ErrInvalidKey(k)
	}

	c := cfg.Mongo.Objects()

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
func insert(c *mgo.Collection, k string, v map[string]interface{}) (*Object, bool, error) {
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
func update(c *mgo.Collection, o *Object, v map[string]interface{}) (*Revision, bool, error) {
	r := Diff(o.Value, v)

	if r == nil {
		return nil, false, nil
	}

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
	if !checkKey(k) {
		return nil, ErrInvalidKey(k)
	}

	c := cfg.Mongo.Objects()

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
		o, changed, err = insert(c, k, v)

		if err != nil {
			return nil, err
		}

		r = o.History[0]

		if err = NotifyEmail(cfg, o, r); err != nil {
			fmt.Fprintln(os.Stderr, "[smtp] error sending email:", err)
		}

		return o.History[0], nil
	}

	if err != nil {
		return nil, err
	}

	r, changed, err = update(c, o, v)

	if err != nil {
		return nil, err
	}

	// Object changed.
	if changed {
		if err = NotifyEmail(cfg, o, r); err != nil {
			fmt.Fprintln(os.Stderr, "[smtp] error sending email:", err)
		}

		return r, nil
	}

	return nil, nil
}
