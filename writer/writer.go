package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
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
	flag.Parse()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case <-ticker:
			fmt.Printf("Writing events at %d/sec\n", perSecondRate)
			f, err := os.OpenFile("/tmp/access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			check(err)
			for index := 0; index < perSecondRate; index++ {
				fmt.Fprintln(f, generateLine())
			}
			f.Close()
		}
	}
}
