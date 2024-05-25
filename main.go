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

func main() {
	defaultPort := "8091"
	if portEnv, ok := os.LookupEnv("PORT"); ok {
		defaultPort = portEnv
	}
	port := flag.String("port", defaultPort, "port for server to listen on")
	flag.Parse()

	currentTime := time.Now()

	mux := http.NewServeMux()
	mux.HandleFunc("/", respHome)
	mux.HandleFunc("/cs", cowsayRes)
	mux.HandleFunc("/cowsay", cowsayRes)

	fmt.Println(`_________________________________ 
< Starting Cowsay HTTP API Server >
 --------------------------------- 
       \   ,__,
        \  (oo)____
           (__)    )\
              ||--|| *`)
	fmt.Println(currentTime.String())
	fmt.Println("Listening on port", ":"+*port)
	log.Fatal(http.ListenAndServe(":"+*port, mux))
}

func respHome(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%s\n", `
GET / -- Returns this page

GET /cowsay -- Does cowsay (customize with URL parameters)
GET /cs
  URL PARAMS
    s string -- Thing to say (defaults to fortune command)
    cf string -- Specify a cowfile (see /list or add l param to request)
    r bool -- Pick a random cowfile
    l bool -- List all cowfiles available
  `)
}

func cowsayRes(w http.ResponseWriter, req *http.Request) {
	var say string
	var cowsayFlags string
	var cowfile string

	parsedUrl, err := url.Parse(req.URL.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing URL string: ", err)
	}
	params, err := url.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing query parameters: ", err)
	}

	if _, ok := params["l"]; ok {
		fmt.Fprintf(w, "%s\n\n", "avaliable cowfiles:")
		for index, file := range getCowfiles() {
			fmt.Fprintf(w, "%d: %s\n", index, file)
		}
		return
	}

	if _, ok := params["s"]; ok {
		say = url.QueryEscape(params.Get("s"))
	} else {
		fortuneOut, err := exec.Command("fortune").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error running fortune command %s\n", err)
		}
		say = string(fortuneOut)
	}

	if _, ok := params["b"]; ok {
		cowsayFlags += "-b "
	}
	if _, ok := params["d"]; ok {
		cowsayFlags += "-d "
	}
	if _, ok := params["g"]; ok {
		cowsayFlags += "-g "
	}
	if _, ok := params["p"]; ok {
		cowsayFlags += "-p "
	}
	if _, ok := params["st"]; ok {
		cowsayFlags += "-s "
	}
	if _, ok := params["t"]; ok {
		cowsayFlags += "-t "
	}
	if _, ok := params["w"]; ok {
		cowsayFlags += "-w "
	}
	if _, ok := params["y"]; ok {
		cowsayFlags += "-y "
	}

	cowfile = "default"
	if _, ok := params["cf"]; ok {
		cowfile = params.Get("cf")
	}
	if _, ok := params["r"]; ok {
		cowfile = getRandomCowfile()
	}

	if !validText(cowfile) {
		http.Error(w, "400 Error - Bad input for cowfile. Parameter must be alphanumeric!\n", http.StatusBadRequest)
	}
	if !checkCowfile(cowfile) {
		http.Error(w, "404 Error - Cowfile not found!\n", http.StatusNotFound)
	}

	unEscSay, err := url.QueryUnescape(say)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error decoding query for say param", err)
	}
	if !validText(unEscSay) {
		http.Error(w, "400 Error - Bad input for say param. Parameter must be alphanumeric!\n", http.StatusBadRequest)
	}

	var cowsayOut []byte
	if cowsayFlags == "" {
		cowsayOut, err = exec.Command("cowsay", "-f", cowfile, unEscSay).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error running command %s\n", err)
			http.Error(w, fmt.Sprintf("error running command %s", err), http.StatusInternalServerError)
		}
	} else {
		cowsayOut, err = exec.Command("cowsay", "-f", cowfile, cowsayFlags, unEscSay).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error running command %s\n", err)
			http.Error(w, fmt.Sprintf("error running command %s", err), http.StatusInternalServerError)
		}
	}
	fmt.Fprintf(w, "%s\n", cowsayOut)
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

func validText(input string) bool {
	// Only allow these characters for security purposes (we run shell commands)
	allowedChars := regexp.MustCompile(`^[a-zA-Z0-9\\s\-\_\.]*$`)
	return allowedChars.MatchString(input)
}
