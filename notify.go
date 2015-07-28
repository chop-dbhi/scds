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
	notifyTemplate = template.New("notify")

	notifyNewObjectEmailBody = `Key: {{.Key}}
Version: {{.Version}}
Time: {{.Time}}
URL: {{.URL}}

# Object

{{.Object}}
`

	notifyChangedObjectEmailBody = `Key: {{.Key}}
Version: {{.Version}}
Time: {{.Time}}
URL: {{.URL}}
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
	// Initialize the email body.
	template.Must(notifyTemplate.New("new-object").Parse(notifyNewObjectEmailBody))
	template.Must(notifyTemplate.New("changed-object").Parse(notifyChangedObjectEmailBody))
}

type context struct {
	Time      time.Time
	Key       string
	Version   int
	URL       string
	Object    string
	Additions string
	Removals  string
	Changes   string
}

func newObjectEmail(cfg *Config, o *Object) (*email.Email, error) {
	var (
		err  error
		byt  []byte
		buff bytes.Buffer
	)

	cxt := context{
		Time:    time.Unix(o.Time, 0).Local(),
		Key:     o.Key,
		Version: o.Version,
		URL:     fmt.Sprintf("http://%s/objects/%s", cfg.HTTP.Addr(), o.Key),
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

	cxt := context{
		Time:    time.Unix(r.Time, 0).Local(),
		Key:     o.Key,
		Version: r.Version,
		URL:     fmt.Sprintf("http://%s/objects/%s", cfg.HTTP.Addr(), o.Key),
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
