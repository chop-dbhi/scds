package main

import "testing"

func TestSubscribeEmail(t *testing.T) {
	resetDB()

	email := "test@example.com"

	subs, err := SubscribeEmail(cfg, email)

	if err != nil {
		t.Fatal(err)
	}

	if len(subs) != 1 {
		t.Errorf("expected 1 subscriber, got %d", len(subs))
	}

	n, err := UnsubscribeEmail(cfg, email)

	if err != nil {
		t.Fatal(err)
	}

	if n != 1 {
		t.Errorf("expected 1 subscriber, got %d", n)
	}
}
