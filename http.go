package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

	app.Get("/keys", keysHandler)

	app.Put("/objects/:key", putHandler)
	app.Get("/objects/:key", getHandler)
	app.Get("/objects/:key/v/:version", getHandler)
	app.Get("/objects/:key/t/:time", getHandler)

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

func keysHandler(c *echo.Context) error {
	cfg := c.Get("config").(*Config)

	keys, err := Keys(cfg)

	if err != nil {
		return err
	}

	resp := c.Response()

	return json.NewEncoder(resp).Encode(keys)
}

func getHandler(c *echo.Context) error {
	key := c.Param("key")

	vs := c.Param("version")
	ts := c.Param("time")

	cfg := c.Get("config").(*Config)

	obj, err := get(cfg, key, true)

	if err != nil {
		return err
	}

	if vs != "" {
		v, err := strconv.Atoi(vs)

		// Invalid parameter for version, treat as a 404.
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}

		// Version is greater than what is available.
		if v > obj.Version || v < 0 {
			return c.NoContent(http.StatusNotFound)
		}

		obj = obj.AtVersion(v)
	} else if ts != "" {
		t, err := ParseTimeString(ts)

		// Invalid parameter for version, treat as a 404.
		if err != nil {
			return c.NoContent(http.StatusNotFound)
		}

		obj = obj.AtTime(t)

		if obj == nil {
			return c.NoContent(http.StatusNotFound)
		}
	}

	// Does not exist.
	if obj == nil {
		return c.NoContent(http.StatusNoContent)
	}

	// Do not include history in output.
	obj.History = nil

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
