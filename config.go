package main

import (
	"fmt"
	"log"
	"net/smtp"

	"gopkg.in/mgo.v2"
)

// SMTPConfig defines configuration fields for communicating with an SMTP server.
// This is used for sending notification emails when changes occur.
type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string

	client *smtp.Client
}

// Addr returns the address of the SMTP host.
func (s *SMTPConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// Auth returns an authorization value for the SMTP server.
func (s *SMTPConfig) Auth() smtp.Auth {
	if s.User == "" {
		return nil
	}

	return smtp.PlainAuth("", s.User, s.Password, s.Addr())
}

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
		c.mongoSession.Close()
	}
}

// Config contains all configuration options.
type Config struct {
	Debug bool
	Mongo MongoConfig
	HTTP  HTTPConfig
	SMTP  SMTPConfig
}
