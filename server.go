package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Weather struct {
	Temperature float64
	Conditions  string
}

var start = time.Now()
var lastRequests = struct {
	sync.RWMutex
	m map[string]time.Time
}{m: make(map[string]time.Time)}

var conditions = [...]string{"cloudy", "sunny", "foggy"}

func getUserIp(req *http.Request) string {
	//normalize for running locally or on heroku
	forwardedIp := req.Header.Get("X-Forwarded-For")
	if forwardedIp == "" {
		parts := strings.SplitN(req.RemoteAddr, ":", 2)
		if len(parts) > 0 {
			return parts[0]
		}
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
	dice := rand.Intn(10)

	switch dice {
	default:
		res.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(getWeather())
		fmt.Fprintln(res, string(response))
	case 1:
		res.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(res, "{ weather: ++_--_(*&^$#$^&*(")
	case 2:
		res.WriteHeader(http.StatusTeapot)
		fmt.Fprintln(res, "I'm A Teapot!")
	case 3:
		res.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(res, "{\"Server Tired\": \"ZzZzZzZzZzZzZzZzZ\" }")
	case 4:
		res.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(res, "{\"Temperature\": \"<script>window.location = 'http://www.google.com'</script>\", \"Conditions\":\"<script>window.location = 'http://www.google.com'</script>\" }")
	}

}

func getTemp(res http.ResponseWriter, req *http.Request) {
	//Allows cross-domain requests in modern browsers
	res.Header().Set("Access-Control-Allow-Origin", "*")
	ip := getUserIp(req)
	//Reject if this IP has made a request in the last second
	lastRequests.RLock()
	lr, ok := lastRequests.m[ip]
	lastRequests.RUnlock()
	maxRate := time.Duration(1)
	if !ok || time.Since(lr) > maxRate*time.Second {
		lastRequests.Lock()
		lastRequests.m[ip] = time.Now()
		lastRequests.Unlock()
		generateRes(res)
	} else {
		// Add a 3 second penalty
		penalty := time.Duration(3)
		lastRequests.Lock()
		lastRequests.m[ip] = time.Now().Add(penalty * time.Second)
		lastRequests.Unlock()
		if rand.Intn(4) > 1 {
			res.WriteHeader(429)
			fmt.Fprintf(res, "Exceeded one request every %d seconds. Now you have to wait %d seconds!\n", maxRate, penalty+maxRate)
		} else {
			// Thanks, Twitter.
			res.WriteHeader(420)
			fmt.Fprintln(res, "Enhance your calm")
		}
	}
}

func main() {
	rand.Seed(time.Now().Unix())

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
