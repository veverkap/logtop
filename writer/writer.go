package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"math/rand" // crypto/rand would be preferred for more secure implementations (https://github.com/golang/go/wiki/CodeReviewComments#crypto-rand
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Returns an int >= min, < max
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

var httpVerbs = [5]string{"GET", "POST", "PUT", "PATCH", "DELETE"}
var users = [5]string{"james", "jill", "frank", "patrick", "lucy"}
var sections = [5]string{"api", "admin", "account", "user", "config"}
var subsections = [5]string{"", "/user", "/widget", "/search", "/update"}
var statusCodes = [7]int{200, 200, 201, 401, 403, 500, 503}
var perSecondRate int
var logFileLocation string

func generateLine() string {
	rand.Seed(time.Now().UnixNano())
	verb := httpVerbs[randomInt(0, 5)]

	rand.Seed(time.Now().UnixNano())
	user := users[randomInt(0, 5)]

	rand.Seed(time.Now().UnixNano())
	section := sections[randomInt(0, 5)]

	rand.Seed(time.Now().UnixNano())
	subsection := subsections[randomInt(0, 5)]

	rand.Seed(time.Now().UnixNano())
	statusCode := statusCodes[randomInt(0, 5)]

	rand.Seed(time.Now().UnixNano())
	byteSize := randomInt(100, 500)
	t := time.Now().UTC()

	return fmt.Sprintf("127.0.0.1 - %s [%02d/%s/%d:%02d:%02d:%02d +0000] \"%s /%s%s HTTP/1.0\" %d %d", user, t.Day(), t.Month().String()[:3], t.Year(), t.Hour(), t.Minute(), t.Second(), verb, section, subsection, statusCode, byteSize)
}

func main() {
	flag.IntVar(&perSecondRate, "rate", 10, "Number of requests per second to write")
	flag.StringVar(&logFileLocation, "file", "/tmp/access.log", "Location of log file")
	flag.Parse()

	ticker := time.NewTicker(time.Second).C

	for {
		select {
		case <-ticker:
			fmt.Printf("Writing events at %d/sec to %s\n", perSecondRate, logFileLocation)
			f, err := os.OpenFile(logFileLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("Could not open file %s", logFileLocation)
			}
			for index := 0; index < perSecondRate; index++ {
				fmt.Fprintln(f, generateLine())
			}
			f.Close()
		}
	}
}
