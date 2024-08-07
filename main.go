package main

import (
	"flag"
	"fmt"
	"html/template"
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

var homeHelpMsg string = `Welcome to the cowsay HTTP API!

GET /api -- This page (you are here)

GET /api/cowsay -- Does 'fortune | cowsay' by default (customize with URL parameters)

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

ALIASES for /api/cowsay path:
  /api/say
  /api/cow
  /api/cs

---

EXAMPLES:
  # random fortune + cowsay (classic)
  cows.rest/api/cowsay
  cows.rest/api/cs

  # random fortune + random cow
  cows.rest/api/cowsay?random
  cows.rest/api/cs?r

  # using misc query parameters
  cows.rest/api/cowsay?d&say=0xDEADBEEF
  cows.rest/api/cs?d&say=0xDEADBEEF
  cows.rest/api/cow?say=moo%20world

  # get all cows
  cows.rest/api/cs?all

TIP:
  # URL escape strings with perl or python:
  perl -nE 'use URI::Escape; chomp $_; print(uri_escape($_))' <<< "some long random text"
  python -c 'import urllib.parse; print(urllib.parse.quote(input()))' <<< "some long random text"

  curl "cows.rest/api/cowsay?random&say=some+long+random+text"

GITHUB:
https://github.com/lemonase/cowsay-http
  `

func respHomeApi(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%s\n", homeHelpMsg)
}

func cowsayRes(w http.ResponseWriter, req *http.Request) {
	var csOpts cowsayOpts
	listCows := false
	allCows := false
	isRandomCowfile := false
	isParamSay := false
	csOpts.cowfile = "default"

	// inital url parsing
	parsedUrl, err := url.Parse(req.URL.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing URL string: ", err)
	}
	params, err := url.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing query parameters: ", err)
	}

	// processing URL parameters
	// params that handle generic switches and options
	for _, opt := range cowsaySwitches {
		if _, ok := params[opt]; ok {
			csOpts.flags = csOpts.flags + "-" + opt + " "
		}
	}
	// param to list all cowfiles
	if params.Has("listCows") || params.Has("list") {
		listCows = true
	}
	// param to get random cowfile
	if params.Has("randomCow") || params.Has("random") || params.Has("r") {
		isRandomCowfile = true
	}
	// param for say string
	if params.Has("say") {
		isParamSay = true
	}
	// param to run with all cowfiles
	if params.Has("all") || params.Has("allCows") {
		allCows = true
	}
	// param for cowfile
	if params.Has("cowfile") {
		csOpts.cowfile = params.Get("cowfile")
	}
	if params.Has("cow") {
		csOpts.cowfile = params.Get("cow")
	}
	if params.Has("cf") {
		csOpts.cowfile = params.Get("cf")
	}
	if !checkCowfile(csOpts.cowfile) {
		http.Error(w, "404 Error - Cowfile not found!\n", http.StatusNotFound)
		return
	}

	// execute based on params
	if listCows {
		fmt.Fprintf(w, "%s\n\n", "avaliable cowfiles:")
		for index, file := range getCowfiles() {
			fmt.Fprintf(w, "%d: %s\n", index, file)
		}
		return
	}
	if isRandomCowfile {
		csOpts.cowfile = getRandomCowfile()
	}
	csOpts.cowfile = sanitizeText(csOpts.cowfile)

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
	currentTime := time.Now()

	if ipEnv, ok := os.LookupEnv("IP"); ok {
		defaultIp = ipEnv
	}
	if portEnv, ok := os.LookupEnv("PORT"); ok {
		defaultPort = portEnv
	}
	ip := flag.String("ip", defaultIp, "ip for server to listen on (default is all)")
	port := flag.String("port", defaultPort, "port for server to listen on")
	flag.Parse()

	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api", respHomeApi)
	mux.HandleFunc("/api/cowsay", cowsayRes)
	mux.HandleFunc("/api/cow", cowsayRes)
	mux.HandleFunc("/api/say", cowsayRes)
	mux.HandleFunc("/api/cs", cowsayRes)

	// Templates
	tmpl := template.Must(template.ParseFiles("pages/index.html"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, "")
	})

	fmt.Println(startupMsg)
	fmt.Println(currentTime.String())
	fmt.Println("Listening on port", *ip+":"+*port)
	log.Fatal(http.ListenAndServe(*ip+":"+*port, mux))
}
