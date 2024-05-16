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
	port := flag.String("port", ":8091", "port for server to listen on")

	currentTime := time.Now()

	mux := http.NewServeMux()
	mux.HandleFunc("/", respHome)
	mux.HandleFunc("/cowsay", respCowsay)
	mux.HandleFunc("/fortune", respFortune)
	mux.HandleFunc("/listCows", respListCowfiles)

	fmt.Println(currentTime.String())
	fmt.Println(`
 _________________________________ 
< Starting Cowsay HTTP API Server >
 --------------------------------- 
       \   ,__,
        \  (oo)____
           (__)    )\
              ||--|| *`)
	fmt.Println("Listening on port:", *port)
	log.Fatal(http.ListenAndServe(*port, mux))
}

func respHome(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%s\n", `
 ________________________________ 
< Hi, Welcome to Cowsay HTTP API >
 -------------------------------- 
       \   ,__,
        \  (oo)____
           (__)    )\
              ||--|| *

GET / -- Returns this page

GET /cowsay
  URL PARAMS
    randomCow bool -- Toggle random cowfile
    cowfile string -- Specify a cowfile
    s string -- Thing to say

GET /fortune -- Returns a fortune with an optional pipe to cowsay
  URL PARAMS
    cowsay bool -- Toggle cowsay
      randomCow bool -- Toggle random cowfile
      cowfile string -- Specify a cowfile
      borg bool
      dead bool
      greedy bool
      paranoia bool
      stoned bool
      tired bool
      wired bool
      youthful bool
    time bool -- Print time in response

GET /listCows -- Returns a list of available cows

  `)
}

func respCowsay(w http.ResponseWriter, req *http.Request) {
	parsedUrl, err := url.Parse(req.URL.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing URL string: ", err)
	}
	params, err := url.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing query parameters: ", err)
	}

	cowsayString := params.Get("s")
	// TODO (refactor): allow url encode spaces in a secure way
	if !validCommand(cowsayString) {
		w.WriteHeader(400)
		w.Write([]byte("400 Error - Bad input for s. Parameter must be alphanumeric!\n"))
		return
	}

	cowfile := "default"
	if !validCommand(cowfile) {
		w.WriteHeader(400)
		w.Write([]byte("400 Error - Bad input for cowfile. Parameter must be alphanumeric!\n"))
		return
	}
	if _, ok := params["cowfile"]; ok {
		cowfile = params.Get("cowfile")
	}
	if _, ok := params["randomCow"]; ok {
		cowfile = getRandomCowfile()
	}

	fmt.Fprintf(w, "%s\n", execCowsay(cowsayString, cowfile))
}

func respFortune(w http.ResponseWriter, req *http.Request) {
	parsedUrl, err := url.Parse(req.URL.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing URL string: ", err)
	}
	params, err := url.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing query parameters: ", err)
	}

	if _, ok := params["time"]; ok {
		fmt.Fprintf(w, "%s\n\n", time.Now().String())
	}

	// TODO (refactor): move these to cowsay function
	cowStatusOps := map[string]string{
		"borg":     "-b",
		"dead":     "-d",
		"greedy":   "-g",
		"paranoia": "-p",
		"stoned":   "-s",
		"tired":    "-t",
		"wired":    "-w",
		"youthful": "-y",
	}
	var cowOpts string

	for status, opt := range cowStatusOps {
		if params.Has(status) {
			cowOpts = opt + " "
		}
	}

	// TODO (refactor): move these to cowsay function as well
	if _, ok := params["cowsay"]; ok {
		if _, ok := params["randomCow"]; ok {
			fmt.Fprintf(w, "%s\n", execFortune(true, getRandomCowfile(), cowOpts))
		} else {
			cowfile := params.Get("cowfile")
			if !validCommand(cowfile) {
				w.WriteHeader(400)
				w.Write([]byte("400 Error - Bad input for cowfile. Parameter must be alphanumeric!\n"))
				return
			}

			if !checkCowfile(cowfile) {
				w.WriteHeader(404)
				w.Write([]byte("404 Error - Cowfile not found! Check /listCows for a list of cowfiles!\n"))
				return
			}

			fmt.Fprintf(w, "%s\n", execFortune(true, params.Get("cowfile"), cowOpts))
		}
	} else {
		fmt.Fprintf(w, string(execFortune(false, "default", cowOpts)))
	}
}

func respListCowfiles(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%s\n\n", "avaliable cowfiles:")
	fmt.Fprintf(w, "%v\n", getCowfiles())
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

func execFortune(doCowsay bool, cowfile string, cowOpts string) string {
	fortuneCmd := exec.Command("fortune")

	if !doCowsay {
		fortuneOutput, err := fortuneCmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		return string(fortuneOutput)
	} else {
		// TODO: pass fortune string to cowsay function within go instead of pipes.
		cowsayCmd := exec.Command("cowsay", cowOpts)
		if cowfile != "" {
			cowsayCmd = exec.Command("cowsay", "-f", cowfile, cowOpts)
		}

		fortuneOut, err := fortuneCmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err := fortuneCmd.Start(); err != nil {
			fmt.Println("error starting fortune", err)
		}
		cowsayCmd.Stdin = fortuneOut
		cowsayOut, err := cowsayCmd.StdoutPipe()
		if err != nil {
			fmt.Println("error piping fortune to cowsay", err)
		}
		if err := cowsayCmd.Start(); err != nil {
			fmt.Println("error starting cowsay", err)
		}

		defer fortuneCmd.Wait()
		defer cowsayCmd.Wait()

		cowsayResult, err := io.ReadAll(cowsayOut)
		if err != nil {
			fmt.Printf("Error reading command output: %v\n", err)
		}

		return string(cowsayResult)
	}
}

func execCowsay(say string, cowfile string) string {
	out, err := exec.Command("cowsay", "-f", cowfile, say).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func validCommand(input string) bool {
	// Only allow alphanumeric commands
	allowedChars := regexp.MustCompile(`^[a-zA-Z0-9\\s\-\_]*$`)
	return allowedChars.MatchString(input)
}
