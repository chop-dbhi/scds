package main

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/jordan-wright/email"
	"gopkg.in/yaml.v2"
)

var (
	// Root template tree of notification templates.
	notifyTemplate = template.New("notify")

	// Template of the email body when new a object is created.
	notifyNewObjectEmailBody = `Key: {{.Key}}
Version: {{.Version}}
Time: {{.Time}}
URL: {{.URL}}
Version URL: {{.VersionURL}}

# Object

{{.Object}}
`

	// Template of the email body when an object has changed.
	notifyChangedObjectEmailBody = `Key: {{.Key}}
Version: {{.Version}}
Time: {{.Time}}
URL: {{.URL}}
Version URL: {{.VersionURL}}

{{if .Changes}}
# Changes

{{.Changes}}{{end}}{{if .Additions}}
# Additions

{{.Additions}}{{end}}{{if .Removals}}
# Removals

{{.Removals}}{{end}}
`
)

func init() {
	// Compile the email bodies.
	template.Must(notifyTemplate.New("new-object").Parse(notifyNewObjectEmailBody))
	template.Must(notifyTemplate.New("changed-object").Parse(notifyChangedObjectEmailBody))
}

// EmailContext is the template EmailContext for email bodies.
type EmailContext struct {
	Time       time.Time
	Key        string
	Version    int
	URL        string
	VersionURL string
	Object     string
	Additions  string
	Removals   string
	Changes    string
}

func newObjectEmail(cfg *Config, o *Object) (*email.Email, error) {
	var (
		err  error
		byt  []byte
		buff bytes.Buffer
	)

	cxt := EmailContext{
		Time:       time.Unix(o.Time, 0).Local(),
		Key:        o.Key,
		Version:    o.Version,
		URL:        fmt.Sprintf("http://%s/objects/%s", cfg.HTTP.Addr(), o.Key),
		VersionURL: fmt.Sprintf("http://%s/objects/%s/v/%d", cfg.HTTP.Addr(), o.Key, o.Version),
	}

	byt, _ = yaml.Marshal(o.Value)
	cxt.Object = string(byt)

	if err = notifyTemplate.Lookup("new-object").Execute(&buff, &cxt); err != nil {
		return nil, err
	}

	e := email.NewEmail()

	e.Subject = "[SCDS] New Object"
	e.Text = buff.Bytes()

	return e, nil
}

func changedObjectEmail(cfg *Config, o *Object, r *Revision) (*email.Email, error) {
	var (
		err  error
		byt  []byte
		buff bytes.Buffer
	)

	cxt := EmailContext{
		Time:       time.Unix(r.Time, 0).Local(),
		Key:        o.Key,
		Version:    r.Version,
		URL:        fmt.Sprintf("http://%s/objects/%s", cfg.HTTP.Addr(), o.Key),
		VersionURL: fmt.Sprintf("http://%s/objects/%s/v/%d", cfg.HTTP.Addr(), o.Key, r.Version),
	}

	if r.Changes != nil {
		byt, _ = yaml.Marshal(r.Changes)
		cxt.Changes = string(byt)
	}

	if r.Additions != nil {
		byt, _ = yaml.Marshal(r.Additions)
		cxt.Additions = string(byt)
	}

	if r.Removals != nil {
		byt, _ = yaml.Marshal(r.Removals)
		cxt.Removals = string(byt)
	}

	if err = notifyTemplate.Lookup("changed-object").Execute(&buff, &cxt); err != nil {
		return nil, err
	}

	e := email.NewEmail()

	e.Subject = "[SCDS] Object Changed"
	e.Text = buff.Bytes()

	return e, nil
}

func SendNotifications(cfg *Config, o *Object, r *Revision) error {
	var (
		e   *email.Email
		err error
	)

	// First version.
	if r.Version == 1 {
		e, err = newObjectEmail(cfg, o)
	} else {
		e, err = changedObjectEmail(cfg, o, r)
	}

	if err != nil {
		return err
	}

	e.From = cfg.SMTP.From
	e.To = []string{cfg.SMTP.From}
	e.Bcc = []string{}

	return e.Send(cfg.SMTP.Addr(), cfg.SMTP.Auth())
}
