package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

var start = time.Now()
var lastRequests = make(map[string]time.Time)

func calculateTemp() float64 {
	//Sin wave of temperatures
	return math.Ceil(math.Sin(time.Since(start).Seconds()/100) * 100)
}

func getTemp(res http.ResponseWriter, req *http.Request) {
	ip := strings.Split(req.RemoteAddr, ":")[0]
	if time.Since(lastRequests[ip]).Seconds() > 2 {
		fmt.Fprintln(res, calculateTemp())
	} else {
		fmt.Fprintln(res, "you lose")
	}
	lastRequests[ip] = time.Now()
}

func main() {
	http.HandleFunc("/", getTemp)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
