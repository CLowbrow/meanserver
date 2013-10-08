package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Weather struct {
	Temperature float64
	Conditions  string
}

var start = time.Now()
var lastRequests = make(map[string]time.Time)
var conditions = [...]string{"cloudy", "sunny", "foggy"}

func getUserIp(req *http.Request) string {
	//normalize for running locally or on heroku
	forwardedIp := req.Header.Get("X-Forwarded-For")
	if forwardedIp == "" {
		return req.RemoteAddr
	} else {
		return forwardedIp
	}
}

func calculateTemp() float64 {
	//Sin wave of temperatures
	return math.Ceil(math.Sin(time.Since(start).Seconds()/300) * 100)
}

func getWeather() Weather {
	return Weather{calculateTemp(), conditions[rand.Intn(3)]}
}

func generateRes(res http.ResponseWriter) {
	// Randomly give back a good response or random garbage :)
	dice := rand.Intn(6)

	switch dice {
	default:
		res.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(getWeather())
		fmt.Fprintln(res, string(response))
	case 1:
		res.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(res, "{ weather: malformed json here ++_--_(*&^$#$^&*(")
	case 2:
		res.WriteHeader(http.StatusTeapot)
		fmt.Fprintln(res, "I'm A Teapot!")
	}

}

func getTemp(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	ip := getUserIp(req)
	//Reject if this IP has made a request in the last two seconds
	if _, ok := lastRequests[ip]; !ok ||
		time.Since(lastRequests[ip]).Seconds() > 2 {
		lastRequests[ip] = time.Now()
		generateRes(res)
	} else {
		lastRequests[ip] = time.Now()
		res.WriteHeader(429)
		fmt.Fprintln(res, "Too many requests")
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
