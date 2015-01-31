package main

import (
	"testing"
)

func TestMessageColor(t *testing.T) {
	origin := "sample message"
	tested := colorMsg(origin, cRed, cGreen)

	if tested != "\x034,3sample message\x0f" {
		t.Fatalf("Unexpected converted message: %+v\n", tested)
	}
}

func TestMessageColorFrontOnly(t *testing.T) {
	origin := "sample message"
	tested := colorMsg(origin, cYellow, cNone)

	if tested != "\x038sample message\x0f" {
		t.Fatalf("Unexpected converted message: %+v\n", tested)
	}
}

func TestMessageMode(t *testing.T) {
	origin := "sample message"
	tested := setIrcMode(ircBold) + origin + setIrcMode(ircCReset)

	if tested != "\x02sample message\x0f" {
		t.Fatalf("Unexpected converted message: %+v\n", tested)
	}
}
