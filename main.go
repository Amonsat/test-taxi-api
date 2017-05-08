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

var cUser chan string
var cCounter chan string
var cChanger chan int
var cRander chan int
var cStatistic chan string
var cGetStatistic chan bool

func main() {
	statistics = make(map[string]int)
	generateBaseReqeusts()

	cUser = make(chan string)
	cCounter = make(chan string)
	cRander = make(chan int)
	cChanger = make(chan int)
	cStatistic = make(chan string)
	cGetStatistic = make(chan bool)

	go func() {
		for {
			select {
			case randInt := <-cChanger:
				name := tools.RandStringRunes(2)
				println("change name " + requests[randInt] + " to " + name)
				requests[randInt] = name
			case randInt := <-cRander:
				rName := requests[randInt]
				cUser <- rName
				statistics[rName] = statistics[rName] + 1
			case <-cGetStatistic:
				println("get stat")
				var text string
				for i, v := range statistics {
					text += i + ": " + strconv.Itoa(v) + "\n"
				}
				cStatistic <- text
			}
		}
	}()

	go tools.DoEvery(200*time.Millisecond, func(t time.Time) {
		cChanger <- rand.Intn(maxRequests)
	})

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
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
	cRander <- rand.Intn(maxRequests)
	req := <-cUser
	fmt.Fprintf(w, "Get random reqest: %s", req)
}

func printAllRequests(w http.ResponseWriter) {
	cGetStatistic <- true
	text := <-cStatistic
	fmt.Fprintf(w, "%s", text)
}

func generateBaseReqeusts() {
	for i := 0; i < len(requests); i++ {
		name := tools.RandStringRunes(2)
		requests[i] = name
	}
}
