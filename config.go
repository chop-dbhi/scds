package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
)

// HTTPConfig defines configuration fields running the HTTP service.
type HTTPConfig struct {
	Host string
	Port int
}

// Addr returns the HTTP address of the SCDS service.
func (s *HTTPConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// MongoConfig defines configuration fields for connecting to a MongoDB server.
type MongoConfig struct {
	URI string

	mongoSession *mgo.Session
}

// Session returns an initialized MongoDB session.
func (c *MongoConfig) Session() *mgo.Session {
	if c.mongoSession == nil {
		log.Printf("* [mongo] Connecting to %s", c.URI)

		session, err := mgo.Dial(c.URI)

		if err != nil {
			log.Fatal(err)
		}

		if err = session.Ping(); err != nil {
			log.Fatal(err)
		}

		if err = session.DB("").C(colName).EnsureIndexKey("key"); err != nil {
			log.Fatal(err)
		}

		c.mongoSession = session
	}

	return c.mongoSession
}

// Close closes an open MongoDB session.
func (c *MongoConfig) Close() {
	if c.mongoSession != nil {
		log.Print("* [mongo] Closing")
		c.mongoSession.Close()
	}
}

// Config contains all configuration options.
type Config struct {
	Debug bool
	Mongo MongoConfig
	HTTP  HTTPConfig
}
