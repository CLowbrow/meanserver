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
		time.Since(lastRequests[req.RemoteAddr]).Seconds() > 1 {
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
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
