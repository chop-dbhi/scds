package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

const (
	mongoObjects     = "objects"
	mongoSubscribers = "subcribers"
)

// Safety mode of the MongoDB instance.
var safeMode = &mgo.Safe{
	WMode:    "majority",
	WTimeout: 2000,
	FSync:    true,
}

func InitConfig() {
	viper.SetConfigName("scds")
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("scds")
	viper.AutomaticEnv()

	// Replaces underscores with periods when mapping environment variables.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set non-zero defaults. Nested options take a lower precedence than
	// dot-delimited ones, so namespaced options are defined here as maps.
	viper.SetDefault("mongo", map[string]interface{}{
		"uri": "localhost/scds",
	})

	viper.SetDefault("http", map[string]interface{}{
		"host": "localhost",
		"port": 5000,
	})

	viper.SetDefault("smtp", map[string]interface{}{
		"host": "localhost",
		"port": 25,
	})

	// Read the default config file from the working directory.
	dir, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	viper.AddConfigPath(dir)
}

func GetConfig() *Config {
	// Load custom config file if explicitly set.
	if path := viper.GetString("config"); path != "" {
		file, err := os.Open(path)

		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		if err = viper.ReadConfig(file); err != nil {
			log.Fatal(err)
		}
	} else {
		viper.ReadInConfig()
	}

	return &Config{
		Debug:  viper.GetBool("debug"),
		Config: viper.GetString("config"),

		Mongo: MongoConfig{
			URI: viper.GetString("mongo.uri"),
		},

		HTTP: HTTPConfig{
			Host:    viper.GetString("http.host"),
			Port:    viper.GetInt("http.port"),
			CORS:    viper.GetBool("http.cors"),
			TLSCert: viper.GetString("http.tlscert"),
			TLSKey:  viper.GetString("http.tlskey"),
		},

		SMTP: SMTPConfig{
			Host:     viper.GetString("smtp.host"),
			Port:     viper.GetInt("smtp.port"),
			User:     viper.GetString("smtp.user"),
			Password: viper.GetString("smtp.password"),
			From:     viper.GetString("smtp.from"),
		},
	}
}

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
	Host    string
	Port    int
	CORS    bool
	TLSCert string
	TLSKey  string
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

		session.SetSafe(safeMode)

		if err = session.DB("").C(mongoObjects).EnsureIndexKey("key"); err != nil {
			log.Fatal(err)
		}

		if err = session.DB("").C(mongoSubscribers).EnsureIndexKey("email"); err != nil {
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

// Objects returns the default objects collection.
func (c *MongoConfig) Objects() *mgo.Collection {
	return c.Session().DB("").C(mongoObjects)
}

// Subscribers returns the subscribers collection.
func (c *MongoConfig) Subscribers() *mgo.Collection {
	return c.Session().DB("").C(mongoSubscribers)
}

// Config contains all configuration options.
type Config struct {
	Debug  bool
	Config string
	Mongo  MongoConfig
	HTTP   HTTPConfig
	SMTP   SMTPConfig
}
