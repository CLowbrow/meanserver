package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"time"
)

var start = time.Now()
var lastRequests = make(map[string]time.Time)

func getUserIp(req *http.Request) string {
	//deal with heroku being weird
	forwardedIp := req.Header.Get("X-Forwarded-For")
	if forwardedIp == "" {
		return req.RemoteAddr
	} else {
		return forwardedIp
	}
}

func calculateTemp() float64 {
	//Sin wave of temperatures
	return math.Ceil(math.Sin(time.Since(start).Seconds()/100) * 100)
}

func getTemp(res http.ResponseWriter, req *http.Request) {
	ip := getUserIp(req)
	if _, ok := lastRequests[ip]; !ok ||
		time.Since(lastRequests[ip]).Seconds() > 2 {
		fmt.Fprintln(res, time.Since(lastRequests[ip]).Seconds())
		lastRequests[ip] = time.Now()
		fmt.Fprintln(res, calculateTemp())

	} else {
		lastRequests[ip] = time.Now()
		fmt.Fprintln(res, "you lose")
	}

}

func main() {
	http.HandleFunc("/", getTemp)
	fmt.Println("listening...")
	port := ""
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else {
		port = "5000"
	}
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
