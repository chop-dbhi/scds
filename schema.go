package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type Schema struct {
	Name    string
	Scope   string
	Field   string
	Pattern string
	File    string

	loader gojsonschema.JSONLoader
}

func (s *Schema) Load() error {
	p, err := filepath.Abs(s.File)
	if err != nil {
		return err
	}

	if _, err := os.Stat(p); err != nil {
		return err
	}

	p = fmt.Sprintf("file://%s", p)
	s.loader = gojsonschema.NewReferenceLoader(p)
	return nil
}

func (s *Schema) Validate(v map[string]interface{}) (*gojsonschema.Result, error) {
	dl := gojsonschema.NewGoLoader(v)
	return gojsonschema.Validate(s.loader, dl)
}

type ResultErrors map[string][]gojsonschema.ResultError

func (r ResultErrors) Error() string {
	var b bytes.Buffer

	for k, errs := range r {
		fmt.Fprintln(&b, k)

		for _, err := range errs {
			fmt.Fprintf(&b, "- %s\n", err)
		}
	}

	return b.String()
}

func (r ResultErrors) MarshalJSON() ([]byte, error) {
	if len(r) == 0 {
		return nil, nil
	}

	var m []map[string]interface{}

	for k, errs := range r {
		var v []map[string]interface{}

		for _, err := range errs {
			v = append(v, map[string]interface{}{
				"type":        err.Type(),
				"description": err.Description(),
				"details":     err.Details(),
			})
		}

		m = append(m, map[string]interface{}{
			"schema": k,
			"errors": v,
		})
	}

	return json.Marshal(m)
}

type Result map[string]*gojsonschema.Result

func (r Result) Matches() []string {
	var m []string

	for k := range r {
		m = append(m, k)
	}

	sort.Strings(m)

	return m
}

func (r Result) Errors() ResultErrors {
	m := make(ResultErrors)

	for k, v := range r {
		errs := v.Errors()

		if len(errs) > 0 {
			m[k] = errs
		}
	}

	return m
}

func (r Result) Valid() bool {
	for _, res := range r {
		if !res.Valid() {
			return false
		}
	}

	return true
}

// Validate validates and object against a set of schemas.
func Validate(k string, m map[string]interface{}, schemas ...*Schema) (Result, error) {
	result := make(Result)

	for _, s := range schemas {
		switch s.Scope {
		// Empty scope matches all documents.
		case "":

		case "object":
			ok, err := compileAndTest(s.Pattern, k)
			if err != nil {
				return nil, err
			}

			if !ok {
				continue
			}

		case "value":
			v, ok := m[s.Field]
			if !ok {
				continue
			}

			// If the pattern is empty, then the presense of the field
			// is all that is required.
			if s.Pattern != "" {
				// Coerce to string for comparison.
				// There are certainly edge cases with this approach, such as nil
				// but the fields being matched on are expected to be simple.
				x := fmt.Sprint(v)

				ok, err := compileAndTest(s.Pattern, x)
				if err != nil {
					return nil, err
				}

				if !ok {
					continue
				}
			}

		// Invalid scope.
		default:
			return nil, fmt.Errorf("invalid schema scope: %s", s.Scope)
		}

		r, err := s.Validate(m)
		if err != nil {
			return nil, err
		}

		result[s.Name] = r
	}

	return result, nil
}

func compileAndTest(expr, str string) (bool, error) {
	if !strings.HasPrefix(expr, "^") {
		expr = "^" + expr
	}
	if !strings.HasSuffix(expr, "$") {
		expr = expr + "$"
	}

	return regexp.MatchString(expr, str)
}
