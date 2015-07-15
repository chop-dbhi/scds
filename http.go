package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2"
)

func runHTTP(s *mgo.Session) {
	app := echo.New()

	app.Use(mw.Logger())
	app.Use(mw.Gzip())
	app.Use(mw.Recover())

	app.Use(func(c *echo.Context) error {
		c.Set("session", s)
		return nil
	})

	app.Put("/store/:key", putHandler)
	app.Get("/store/:key", getHandler)
	app.Get("/log/:key", logHandler)

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Listening on %s...", addr)

	app.Run(addr)
}

func putHandler(c *echo.Context) error {
	key := c.Param("key")
	req := c.Request()

	defer req.Body.Close()

	var val map[string]interface{}

	if err := json.NewDecoder(req.Body).Decode(&val); err != nil {
		return err
	}

	sess := c.Get("session").(*mgo.Session)

	obj, err := Put(sess, key, val)

	if err != nil {
		return err
	}

	// No change.
	if obj == nil {
		return c.NoContent(http.StatusNoContent)
	}

	resp := c.Response()

	return json.NewEncoder(resp).Encode(obj)
}

func getHandler(c *echo.Context) error {
	key := c.Param("key")

	sess := c.Get("session").(*mgo.Session)

	obj, err := Get(sess, key)

	if err != nil {
		return err
	}

	// Does not exist.
	if obj == nil {
		return c.NoContent(http.StatusNoContent)
	}

	resp := c.Response()

	return json.NewEncoder(resp).Encode(obj)
}

func logHandler(c *echo.Context) error {
	key := c.Param("key")

	sess := c.Get("session").(*mgo.Session)

	log, err := Log(sess, key)

	if err != nil {
		return err
	}

	// Does not exist.
	if log == nil {
		return c.NoContent(http.StatusNoContent)
	}

	resp := c.Response()

	return json.NewEncoder(resp).Encode(log)
}
