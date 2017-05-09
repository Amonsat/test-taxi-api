package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"test-taxi-api-3/tools"
	"time"
)

const maxRequests = 50

var requests [maxRequests]string
var statistics map[string]int

var cCmd chan string
var cOut chan string

func main() {
	statistics = make(map[string]int)
	generateBaseReqeusts()

	cCmd = make(chan string)
	cOut = make(chan string)

	go worker()

	go tools.DoEvery(200*time.Millisecond, func(t time.Time) {
		cCmd <- "change"
	})

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func worker() {
	for {
		command := <-cCmd
		switch command {
		case "change":
			randInt := rand.Intn(maxRequests)
			name := tools.RandStringRunes(2)
			requests[randInt] = name
		case "getReq":
			randInt := rand.Intn(maxRequests)
			rName := requests[randInt]
			cOut <- rName
			statistics[rName] = statistics[rName] + 1
		case "getStat":
			var stat string
			for i, v := range statistics {
				stat += i + ": " + strconv.Itoa(v) + "\n"
			}
			cOut <- stat
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	switch r.URL.Path {
	case "/get1":
		printRequest(w)
	case "/get2":
		printAllRequests(w)
	}
}

func printRequest(w http.ResponseWriter) {
	cCmd <- "getReq"
	req := <-cOut
	fmt.Fprintf(w, "Get random reqest: %s", req)
}

func printAllRequests(w http.ResponseWriter) {
	cCmd <- "getStat"
	text := <-cOut
	fmt.Fprintf(w, "%s", text)
}

func generateBaseReqeusts() {
	for i := 0; i < len(requests); i++ {
		name := tools.RandStringRunes(2)
		requests[i] = name
	}
}
