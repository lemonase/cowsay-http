package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Switches that change the appearence of the cow is some way
var cowsaySwitches = []string{
	"b", // bored
	"d", // dead
	"g", // greedy
	"p", // paranoia
	"s", // st0ned
	"t", // tired
	"w", // wired (not tired)
	"y", // youthful
}

type cowsayOpts struct {
	flags   string
	cowfile string
	say     string
}

func respHome(w http.ResponseWriter, req *http.Request) {
	homeHelpMsg := `Welcome to the cowsay HTTP API!

GET / -- This page (you are here)

GET /cowsay -- Does 'fortune | cowsay' by default (customize with URL parameters)

  URL PARAMS
    say                 string  // Thing to say (defaults to fortune command)
    cowfile,cow,cf      string  // Specify a cowfile (add listCows param to list available cowfiles)
    randomCow,random,r  bool    // Pick a random cowfile
    listCows,list       bool    // List all cowfiles available
    allCows,all					bool    // Get message with all cowfiles available

    // Additional cows flags
    b bool  // Cow appears borg mode
    d bool  // Cow appears dead
    g bool  // Cow appears greedy
    p bool  // Cow appears paranoia
    s bool  // Cow appears st0ned
    t bool  // Cow appears tired
    w bool  // Cow appears wired (not tired)
    y bool  // Cow appears youthful

ALIASES for /cowsay path:
  /say
  /cow
  /cs

---

EXAMPLES:
  # random fortune + cowsay (classic)
  cows.rest/cowsay
  cows.rest/cs

  # random fortune + random cow
  cows.rest/cowsay?random
  cows.rest/cs?r

  # using misc query parameters
  cows.rest/cowsay?d&say=0xDEADBEEF
  cows.rest/cs?d&say=0xDEADBEEF
  cows.rest/cow?say=moo%20world

  # get all cows
  cows.rest/cs?all

TIP:
  # URL escape strings with perl or python:
  perl -nE 'use URI::Escape; chomp $_; print(uri_escape($_))' <<< "some long random text"
  python -c 'import urllib.parse; print(urllib.parse.quote(input()))' <<< "some long random text"

  curl "cows.rest/cowsay?random&say=some+long+random+text"

GITHUB:
https://github.com/lemonase/cowsay-http
  `
	fmt.Fprintf(w, "%s\n", homeHelpMsg)
}

func cowsayRes(w http.ResponseWriter, req *http.Request) {
	var csOpts cowsayOpts
	listCows := false
	allCows := false
	isRandomCowfile := false
	isParamSay := false

	// inital url parsing
	parsedUrl, err := url.Parse(req.URL.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing URL string: ", err)
	}
	params, err := url.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing query parameters: ", err)
	}

	// list cowfiles
	if _, ok := params["listCows"]; ok {
		listCows = true
	}
	if _, ok := params["list"]; ok {
		listCows = true
	}
	if listCows {
		fmt.Fprintf(w, "%s\n\n", "avaliable cowfiles:")
		for index, file := range getCowfiles() {
			fmt.Fprintf(w, "%d: %s\n", index, file)
		}
		return
	}

	// show all cows
	if _, ok := params["all"]; ok {
		allCows = true
	}
	if _, ok := params["allCows"]; ok {
		allCows = true
	}

	// handle generic switches and options
	for _, opt := range cowsaySwitches {
		if _, ok := params[opt]; ok {
			csOpts.flags = csOpts.flags + "-" + opt + " "
		}
	}

	// handle cowfile
	csOpts.cowfile = "default"
	if _, ok := params["cowfile"]; ok {
		csOpts.cowfile = params.Get("cowfile")
	}
	if _, ok := params["cow"]; ok {
		csOpts.cowfile = params.Get("cow")
	}
	if _, ok := params["cf"]; ok {
		csOpts.cowfile = params.Get("cf")
	}
	if !checkCowfile(csOpts.cowfile) {
		http.Error(w, "404 Error - Cowfile not found!\n", http.StatusNotFound)
		return
	}

	// random cowfile
	if _, ok := params["randomCow"]; ok {
		isRandomCowfile = true
	}
	if _, ok := params["random"]; ok {
		isRandomCowfile = true
	}
	if _, ok := params["r"]; ok {
		isRandomCowfile = true
	}
	if isRandomCowfile {
		csOpts.cowfile = getRandomCowfile()
	}
	csOpts.cowfile = sanitizeText(csOpts.cowfile)

	// handle say string
	if _, ok := params["say"]; ok {
		isParamSay = true
	} else {
		isParamSay = false
	}

	if isParamSay {
		sayParam := url.QueryEscape(params.Get("say"))
		sayParam, err := url.QueryUnescape(sayParam)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error decoding query for say param", err)
		}
		csOpts.say = sanitizeText(sayParam)
	} else {
		fortuneOut, err := exec.Command("fortune").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error running fortune command %s\n", err)
		}
		csOpts.say = string(fortuneOut)
	}

	// exec cowsay
	if allCows {
		for index, file := range getCowfiles() {
			csOpts.cowfile = file
			fmt.Fprintf(w, "%d: %s\n%s\n", index, file, execCowsay(csOpts))
		}
	} else {
		fmt.Fprintf(w, "%s\n", execCowsay(csOpts))
	}

}

func execCowsay(csOpts cowsayOpts) []byte {
	if csOpts.flags != "" {
		cowsayOut, err := exec.Command("cowsay", csOpts.flags, csOpts.say).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error running cowsay %s\n", err)
		}
		return cowsayOut
	} else {
		cowsayOut, err := exec.Command("cowsay", "-f", csOpts.cowfile, csOpts.say).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error running cowsay %s\n", err)
		}
		return cowsayOut
	}
}

func checkCowfile(input string) bool {
	for _, file := range getCowfiles() {
		if input == file {
			return true
		}
	}
	return false
}

func getRandomCowfile() string {
	cowfiles := getCowfiles()
	index := rand.Intn(len(cowfiles))
	return cowfiles[index]
}

func getCowfiles() []string {
	cfListCmd := exec.Command("cowsay", "-l")
	grepCmd := exec.Command("grep", "-v", "file")
	xargsCmd := exec.Command("xargs")

	cfListPipe, err := cfListCmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	grepCmd.Stdin = cfListPipe
	grepPipe, err := grepCmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	xargsCmd.Stdin = grepPipe
	xargsPipe, err := xargsCmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cfListCmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := grepCmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := xargsCmd.Start(); err != nil {
		log.Fatal(err)
	}
	defer cfListCmd.Wait()
	defer grepCmd.Wait()
	defer xargsCmd.Wait()

	cmdResult, err := io.ReadAll(xargsPipe)
	if err != nil {
		log.Fatal(err)
	}
	allFiles := string(cmdResult)

	return strings.Split(allFiles, " ")
}

func sanitizeText(input string) string {
	// Do not allow these characters for security purposes (we run shell commands)
	badChars := regexp.MustCompile(`\&|\||\;`)

	return badChars.ReplaceAllString(input, "")
}

func main() {
	defaultPort := "8091"
	defaultIp := ""
	startupMsg := string(execCowsay(cowsayOpts{"", "small", "starting cowsay-http api server"}))

	if ipEnv, ok := os.LookupEnv("IP"); ok {
		defaultIp = ipEnv
	}
	if portEnv, ok := os.LookupEnv("PORT"); ok {
		defaultPort = portEnv
	}
	ip := flag.String("ip", defaultIp, "ip for server to listen on (default is all)")
	port := flag.String("port", defaultPort, "port for server to listen on")

	flag.Parse()

	currentTime := time.Now()

	mux := http.NewServeMux()
	mux.HandleFunc("/", respHome)
	mux.HandleFunc("/cowsay", cowsayRes)
	mux.HandleFunc("/cow", cowsayRes)
	mux.HandleFunc("/say", cowsayRes)
	mux.HandleFunc("/cs", cowsayRes)

	fmt.Println(startupMsg)
	fmt.Println(currentTime.String())
	fmt.Println("Listening on port", *ip+":"+*port)
	log.Fatal(http.ListenAndServe(*ip+":"+*port, mux))
}
