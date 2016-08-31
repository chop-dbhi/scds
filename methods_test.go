package main

import (
	"strconv"
	"testing"
)

var cfg *Config

func init() {
	InitConfig()
	cfg = GetConfig()

	safeMode.WMode = ""
}

func resetDB() {
	cfg = GetConfig()
	cfg.Mongo.Session().DB("").DropDatabase()
}

func TestMethods(t *testing.T) {
	defer cfg.Mongo.Close()
	resetDB()

	// Does not exist.
	o, err := Get(cfg, "bob")

	if o != nil {
		t.Error("get1: object should be nil", o)
	}

	if err != nil {
		t.Error("get1: error should be nil", err)
	}

	// Put the state of bob.
	_, err = Put(cfg, "bob", map[string]interface{}{
		"name": "Bob",
	})

	if err != nil {
		t.Error("put1: error should be nil", err)
	}

	// Get the state of bob.
	o, err = Get(cfg, "bob")

	if o == nil {
		t.Error("get2: object should not be nil")
	}

	if err != nil {
		t.Error("get2: error should be nil", err)
	}

	// Put the state of bob (again).
	_, err = Put(cfg, "bob", map[string]interface{}{
		"name":  "Bob",
		"email": "bob@smith.net",
	})

	if err != nil {
		t.Error("put1: error should be nil", err)
	}

	// Get the revisions.
	h, err := Log(cfg, "bob")

	if len(h) != 2 {
		t.Error("log1: expected 2 revisions, got %d", len(h))
	}
}

func BenchmarkPutInsert(b *testing.B) {
	defer cfg.Mongo.Close()
	resetDB()

	v := map[string]interface{}{
		"name":   "Test",
		"sku":    "1281-291-320191",
		"qty":    31,
		"active": true,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Put(cfg, strconv.Itoa(i), v)
	}
}

func BenchmarkPutUpdate(b *testing.B) {
	defer cfg.Mongo.Close()
	resetDB()

	k := "item"

	v := map[string]interface{}{
		"name":   "Test",
		"sku":    "1281-291-320191",
		"qty":    0,
		"active": true,
	}

	// Initial revision.
	Put(cfg, k, v)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v["qt"] = i + 1
		Put(cfg, k, v)
	}
}

func BenchmarkGet(b *testing.B) {
	defer cfg.Mongo.Close()
	resetDB()

	k := "item"

	v := map[string]interface{}{
		"name":   "Test",
		"sku":    "1281-291-320191",
		"qty":    31,
		"active": true,
	}

	Put(cfg, k, v)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Get(cfg, k)
	}
}
