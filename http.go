package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

func runHTTP(cfg *Config) {
	app := echo.New()

	app.Use(mw.Logger())
	app.Use(mw.Gzip())
	app.Use(mw.Recover())

	app.SetDebug(cfg.Debug)

	app.Use(func(c *echo.Context) error {
		c.Set("config", cfg)
		return nil
	})

	app.Put("/store/:key", putHandler)
	app.Get("/store/:key", getHandler)
	app.Get("/log/:key", logHandler)

	addr := cfg.HTTP.Addr()
	log.Printf("* [http] Listening on %s", addr)

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

	cfg := c.Get("config").(*Config)

	obj, err := Put(cfg, key, val)

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

	cfg := c.Get("config").(*Config)

	obj, err := Get(cfg, key)

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

	cfg := c.Get("config").(*Config)

	log, err := Log(cfg, key)

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
