package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TODO: test http handlers

// Cowsay Tests

func TestExecCowsay(t *testing.T) {
	got := string(execCowsay(cowsayOpts{"", "default", "test"}))
	want, err := exec.Command("cowsay", "test").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cowsay command has failed")
	}
	if got != string(want) {
		t.Errorf("bad cowsay output: got %q want %q", got, want)
	}
}

// Cowfile Tests

func TestNoCowfile(t *testing.T) {
	got := checkCowfile("does-not-exist")
	want := false

	if got != want {
		t.Errorf("no cowfile exists: got %t want %t", got, want)
	}
}

func TestGetRandomCowfile(t *testing.T) {
	got := checkCowfile(getRandomCowfile())
	want := true

	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}

func TestGetCowfiles(t *testing.T) {
	got := strings.TrimSuffix(strings.Join(getCowfiles(), " "), "\n")
	want := "beavis.zen blowfish bong bud-frogs bunny cheese cower daemon default dragon dragon-and-cow elephant elephant-in-snake eyes flaming-sheep ghostbusters head-in hellokitty kiss kitty koala kosh luke-koala meow milk moofasa moose mutilated ren satanic sheep skeleton small stegosaurus stimpy supermilker surgery three-eyes turkey turtle tux udder vader vader-koala www"
	if got != want {
		t.Errorf("bad cowfiles got %q want %q", got, want)
	}
}
