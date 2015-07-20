package main

import (
	"os"
	"strconv"
	"testing"

	"gopkg.in/mgo.v2"
)

var session *mgo.Session

func init() {
	// Set the global.
	mongoURI = os.Getenv("MONGO_URI")
	session = initDB()
}

func resetDB() {
	session.DB("").DropDatabase()
}

func TestMethods(t *testing.T) {
	resetDB()

	// Does not exist.
	o, err := Get(session, "bob")

	if o != nil {
		t.Error("get1: object should be nil", o)
	}

	if err != nil {
		t.Error("get1: error should be nil", err)
	}

	// Put the state of bob.
	_, err = Put(session, "bob", map[string]interface{}{
		"name": "Bob",
	})

	if err != nil {
		t.Error("put1: error should be nil", err)
	}

	// Get the state of bob.
	o, err = Get(session, "bob")

	if o == nil {
		t.Error("get2: object should not be nil")
	}

	if err != nil {
		t.Error("get2: error should be nil", err)
	}

	// Put the state of bob (again).
	_, err = Put(session, "bob", map[string]interface{}{
		"name":  "Bob",
		"email": "bob@smith.net",
	})

	if err != nil {
		t.Error("put1: error should be nil", err)
	}

	// Get the revisions.
	h, err := Log(session, "bob")

	if len(h) != 2 {
		t.Error("log1: expected 2 revisions, got %d", len(h))
	}
}

func BenchmarkPutInsert(b *testing.B) {
	resetDB()

	v := map[string]interface{}{
		"name":   "Test",
		"sku":    "1281-291-320191",
		"qty":    31,
		"active": true,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Put(session, strconv.Itoa(i), v)
	}
}

func BenchmarkPutUpdate(b *testing.B) {
	resetDB()

	k := "item"

	v := map[string]interface{}{
		"name":   "Test",
		"sku":    "1281-291-320191",
		"qty":    0,
		"active": true,
	}

	// Initial revision.
	Put(session, k, v)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v["qt"] = i + 1
		Put(session, k, v)
	}
}

func BenchmarkGet(b *testing.B) {
	resetDB()

	k := "item"

	v := map[string]interface{}{
		"name":   "Test",
		"sku":    "1281-291-320191",
		"qty":    31,
		"active": true,
	}

	Put(session, k, v)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Get(session, k)
	}
}
