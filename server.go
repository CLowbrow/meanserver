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

func calculateTemp() float64 {
	//Sin wave of temperatures
	return math.Ceil(math.Sin(time.Since(start).Seconds()/100) * 100)
}

func getTemp(res http.ResponseWriter, req *http.Request) {

	if _, ok := lastRequests[req.RemoteAddr]; !ok ||
		time.Since(lastRequests[req.RemoteAddr]).Seconds() > 2 {
		fmt.Fprintln(res, time.Since(lastRequests[req.RemoteAddr]).Seconds())
		lastRequests[req.RemoteAddr] = time.Now()
		fmt.Fprintln(res, calculateTemp())

	} else {
		lastRequests[req.RemoteAddr] = time.Now()
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
