package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine"

	"gopkg.in/labstack/echo.v2/engine/standard"
	mw "gopkg.in/labstack/echo.v2/middleware"
	"gopkg.in/mgo.v2/bson"
)

const StatusUnprocessableEntity = 422

func runHTTP(cfg *Config) {
	app := echo.New()

	app.SetDebug(cfg.Debug)

	app.Pre(mw.RemoveTrailingSlash())
	app.Use(mw.Logger())
	app.Use(mw.Recover())

	app.Use(mw.GzipWithConfig(mw.GzipConfig{
		Level: 5,
	}))

	if cfg.HTTP.CORS {
		app.Use(mw.CORS())
	}

	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("config", cfg)
			return next(c)
		}
	})

	app.Get("/", rootHandler)

	app.Get("/keys", keysHandler)

	app.Get("/subscribers", getSubscribersHandler)
	app.Post("/subscribers", addSubscribersHandler)
	app.Delete("/subscriber/:token", deleteSubscriberHandler)

	app.Put("/objects/:key", putHandler)
	app.Get("/objects/:key", getHandler)
	app.Get("/objects/:key/v/:version", getHandler)
	app.Get("/objects/:key/t/:time", getHandler)

	app.Get("/log/:key", logHandler)

	addr := cfg.HTTP.Addr()
	log.Printf("* [http] Listening on %s", addr)

	ecfg := engine.Config{
		Address: addr,
	}

	if cfg.HTTP.TLSKey != "" {
		ecfg.TLSCertFile = cfg.HTTP.TLSCert
		ecfg.TLSKeyFile = cfg.HTTP.TLSKey
	}

	app.Run(standard.WithConfig(ecfg))
}

func rootHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"name":    "SCDS",
		"version": "1.0.0",
	})
}

func putHandler(c echo.Context) error {
	var val map[string]interface{}

	if err := c.Bind(&val); err != nil {
		return err
	}

	cfg := c.Get("config").(*Config)

	key := c.Param("key")
	obj, err := Put(cfg, key, val)

	if err != nil {
		// Failed validation.
		if errs, ok := err.(ResultErrors); ok {
			return c.JSON(StatusUnprocessableEntity, errs)
		}

		return err
	}

	// No change.
	if obj == nil {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, obj)
}

func keysHandler(c echo.Context) error {
	cfg := c.Get("config").(*Config)

	keys, err := Keys(cfg)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, keys)
}

func getHandler(c echo.Context) error {
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

	return c.JSON(http.StatusOK, obj)
}

func logHandler(c echo.Context) error {
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

	return c.JSON(http.StatusOK, log)
}

func getSubscribersHandler(c echo.Context) error {
	cfg := c.Get("config").(*Config)

	subs, err := AllSubscribers(cfg)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, subs)
}

func addSubscribersHandler(c echo.Context) error {
	cfg := c.Get("config").(*Config)

	var (
		err    error
		emails []string
	)

	if err = c.Bind(&emails); err != nil {
		return c.JSON(StatusUnprocessableEntity, map[string]interface{}{
			"message": "problem decoding request body",
			"error":   err,
		})
	}

	subs, err := SubscribeEmail(cfg, emails...)

	if err != nil {
		return c.JSON(StatusUnprocessableEntity, map[string]interface{}{
			"message": "problem subscribing emails",
			"error":   err,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"subscribed": len(subs),
	})
}

func deleteSubscriberHandler(c echo.Context) error {
	token := c.Param("token")

	cfg := c.Get("config").(*Config)

	if !bson.IsObjectIdHex(token) {
		return c.NoContent(http.StatusNotFound)
	}

	ok, err := UnsubscribeID(cfg, bson.ObjectIdHex(token))

	if err != nil {
		return c.JSON(StatusUnprocessableEntity, map[string]interface{}{
			"message": "problem unsubscribing email",
			"error":   err,
		})
	}

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	return c.NoContent(http.StatusOK)
}
