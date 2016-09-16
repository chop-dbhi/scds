package main

import (
	"io/ioutil"
	"reflect"
	"sort"
	"testing"
)

func TestValidate(t *testing.T) {
	generalSchemaContents := `
		{
			"type": "object",
			"required": ["type"],
			"properties": {
				"type": {
					"enum": ["user", "book"]
				}
			}
		}
	`

	userSchemaContents := `
		{
			"type": "object",
			"required": ["firstName", "lastName"],
			"properties": {
				"firstName": {
					"type": "string"
				},
				"lastName": {
					"type": "string"
				},
				"admin": {
					"type": "boolean"
				}
			}
		}
	`

	bookSchemaContents := `
		{
			"type": "object",
			"required": ["title", "authors"],
			"properties": {
				"title": {
					"type": "string"
				},
				"authors": {
					"type": "array",
					"items": {
						"minItems": 1,
						"type": "string"
					}
				}
			}
		}
	`

	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	generalSchemaPath := f.Name()
	f.Write([]byte(generalSchemaContents))
	f.Close()

	f, err = ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	userSchemaPath := f.Name()
	f.Write([]byte(userSchemaContents))
	f.Close()

	f, err = ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	bookSchemaPath := f.Name()
	f.Write([]byte(bookSchemaContents))
	f.Close()

	generalSchema := Schema{
		Name:  "general",
		Scope: "value",
		Field: "type",
		File:  generalSchemaPath,
	}

	userSchema := Schema{
		Name:    "user",
		Scope:   "object",
		Pattern: "^users:.*",
		File:    userSchemaPath,
	}

	bookSchema := Schema{
		Name:    "book",
		Scope:   "value",
		Field:   "type",
		Pattern: "book",
		File:    bookSchemaPath,
	}

	adminSchema := Schema{
		Name:    "admin",
		Scope:   "value",
		Field:   "id",
		Pattern: `1\d{3}`,
		File:    userSchemaPath,
	}

	schemas := []*Schema{
		&generalSchema,
		&bookSchema,
		&userSchema,
		&adminSchema,
	}

	for _, s := range schemas {
		if err := s.Load(); err != nil {
			t.Fatal(err)
		}
	}

	tests := map[string]struct {
		Valid   bool
		Matches []string
		Key     string
		Value   map[string]interface{}
	}{
		"other (no schema)": {
			Valid:   true,
			Matches: nil,
			Key:     "1",
			Value: map[string]interface{}{
				"foo": 1,
			},
		},

		"other (bad type)": {
			Valid:   false,
			Matches: []string{"general"},
			Key:     "1",
			Value: map[string]interface{}{
				"type": "other",
				"foo":  1,
			},
		},

		"book": {
			Valid: true,
			Matches: []string{
				"general",
				"book",
			},
			Key: "1",
			Value: map[string]interface{}{
				"type":  "book",
				"title": "The Go Programming Language",
				"authors": []string{
					"Brian W. Kernighan",
					"Alan Donovan",
				},
			},
		},

		"book (no authors)": {
			Valid: false,
			Matches: []string{
				"general",
				"book",
			},
			Key: "1",
			Value: map[string]interface{}{
				"type":  "book",
				"title": "The Go Programming Language",
			},
		},

		"user": {
			Valid: true,
			Matches: []string{
				"general",
				"user",
			},
			Key: "users:jdoe",
			Value: map[string]interface{}{
				"type":      "user",
				"firstName": "John",
				"lastName":  "Doe",
			},
		},

		"user (admin)": {
			Valid: true,
			Matches: []string{
				"general",
				"user",
				"admin",
			},
			Key: "users:admin",
			Value: map[string]interface{}{
				"type":      "user",
				"id":        1001,
				"firstName": "Super",
				"lastName":  "User",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := Validate(test.Key, test.Value, schemas...)
			if err != nil {
				t.Fatal(err)
			}

			if !matchesEqual(res.Matches(), test.Matches) {
				t.Errorf("expected %v matches, got %v", test.Matches, res.Matches())
			}

			if res.Valid() && !test.Valid {
				t.Error("expected invalid result, got a valid one")
			} else if !res.Valid() && test.Valid {
				t.Error("expected valid result, got invalid one")
				t.Log(res.Errors())
			}
		})
	}
}

func matchesEqual(a, b []string) bool {
	sort.Strings(a)
	sort.Strings(b)
	return reflect.DeepEqual(a, b)
}
