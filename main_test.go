package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestHTTPCowsayDeadbeef(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(cowsayRes))
	resp, err := http.Get(server.URL + "/cowsay?d&say=0xDEADBEEF")
	if err != nil {
		t.Error(err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	expected := ` ____________ 
< 0xDEADBEEF >
 ------------ 
        \   ^__^
         \  (xx)\_______
            (__)\       )\/\
             U  ||----w |
                ||     ||

`

	if string(b) != expected {
		t.Errorf("got cowsay: %s expected cowsay %s", string(b), string(expected))
	}
}

func TestHTTPCowsayMoo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(cowsayRes))
	resp, err := http.Get(server.URL + "/cowsay?say=moo+world")
	if err != nil {
		t.Error(err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	expected := ` ___________ 
< moo world >
 ----------- 
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||

`

	if string(b) != expected {
		t.Errorf("got cowsay: %s expected cowsay %s", string(b), string(expected))
	}

}

func TestHTTPServerOk(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(cowsayRes))
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("did not recieve 200 status: %d", resp.StatusCode)
	}
}

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
