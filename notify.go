package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"

	"github.com/jordan-wright/email"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
)

// Template tree of notification templates.
var notifyTemplate = template.New("notify")

func init() {
	// Template of the email body when new a object is created.
	newObject := string(MustAsset("email/new_object_email_body.txt"))

	// Template of the email body when an object has changed.
	changedObject := string(MustAsset("email/changed_object_email_body.txt"))

	// Compile the email bodies.
	template.Must(notifyTemplate.New("new-object").Parse(newObject))
	template.Must(notifyTemplate.New("changed-object").Parse(changedObject))
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

func NotifyEmail(cfg *Config, o *Object, r *Revision) error {
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

	q := bson.M{
		"subscribed": true,
	}

	p := bson.M{
		"email": 1,
	}

	var subs []*Subscriber

	if err = cfg.Mongo.Subscribers().Find(q).Select(p).All(&subs); err != nil {
		return err
	}

	// No subscribers.
	if len(subs) == 0 {
		return nil
	}

	e.From = cfg.SMTP.From
	e.To = make([]string, 1)

	for _, sub := range subs {
		e.To[0] = sub.Email

		if err = e.Send(cfg.SMTP.Addr(), cfg.SMTP.Auth()); err != nil {
			log.Printf("error sending email: %s", err)
		}
	}

	return nil
}

// Subscriber
type Subscriber struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Email string
	Time  time.Time
}

func AllSubscribers(cfg *Config) ([]*Subscriber, error) {
	c := cfg.Mongo.Subscribers()

	var subs []*Subscriber

	if err := c.Find(nil).All(&subs); err != nil {
		return nil, err
	}

	return subs, nil
}

// SubscribeEmail subscribes one or more email addresses to receive notification
// emails when object events occurs. Returned are the new subscribers or an error
// if one occurred.
func SubscribeEmail(cfg *Config, emails ...string) ([]*Subscriber, error) {
	c := cfg.Mongo.Subscribers()

	var (
		q    bson.M
		err  error
		chg  mgo.Change
		info *mgo.ChangeInfo
		sub  *Subscriber
		subs []*Subscriber
	)

	// Upsert the subscribers based on the email address. Email
	// addresses to lowercased for consistency.
	for _, email := range emails {
		sub = &Subscriber{
			Email: strings.ToLower(email),
			Time:  time.Now().UTC(),
		}

		q = bson.M{
			"email": sub.Email,
		}

		chg = mgo.Change{
			Upsert:    true,
			ReturnNew: true,
			Update: bson.M{
				"$setOnInsert": sub,
			},
		}

		if info, err = c.Find(q).Apply(chg, sub); err != nil {
			break
		}

		// If the document was not updated, it was inserted.
		if info.Updated == 0 {
			subs = append(subs, sub)
		}

	}

	return subs, err
}

// UnsubscribeEmail unsubscribes email addresses from receiving notifications.
func UnsubscribeEmail(cfg *Config, emails ...string) (int, error) {
	c := cfg.Mongo.Subscribers()

	var (
		n   int
		err error
		q   bson.M
	)

	for _, email := range emails {
		email = strings.ToLower(email)

		q = bson.M{
			"email": email,
		}

		err = c.Remove(q)

		// No subscriber with the email.
		if err == mgo.ErrNotFound {
			continue
		}

		if err != nil {
			break
		}

		n++
	}

	return n, err
}

func UnsubscribeID(cfg *Config, id bson.ObjectId) (bool, error) {
	c := cfg.Mongo.Subscribers()

	q := bson.M{
		"_id": id,
	}

	err := c.Remove(q)

	// No subscriber with the id.
	if err == mgo.ErrNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
