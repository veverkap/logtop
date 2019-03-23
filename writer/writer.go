package main

import (
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

func main() {
	HTTPVerbs := [5]string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	Users := [5]string{"james", "jill", "frank", "patrick", "lucy"}
	Sections := [5]string{"api", "admin", "account", "user", "config"}
	Subsections := [5]string{"", "/user", "/widget", "/search", "/update"}
	StatusCodes := [7]int{200, 200, 201, 401, 403, 500, 503}

	//127.0.0.1 - james [09/May/2018:16:00:39 +0000] "GET /report HTTP/1.0" 200 123

	for true {
		f, err := os.OpenFile("/tmp/access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		check(err)
		// defer f.Close()

		rand.Seed(time.Now().UnixNano())
		verb := HTTPVerbs[randomInt(0, 5)]

		rand.Seed(time.Now().UnixNano())
		user := Users[randomInt(0, 5)]

		rand.Seed(time.Now().UnixNano())
		section := Sections[randomInt(0, 5)]

		rand.Seed(time.Now().UnixNano())
		subsection := Subsections[randomInt(0, 5)]

		rand.Seed(time.Now().UnixNano())
		statusCode := StatusCodes[randomInt(0, 5)]

		rand.Seed(time.Now().UnixNano())
		byteSize := randomInt(100, 500)
		t := time.Now().UTC()

		fmt.Fprintf(f, "127.0.0.1 - %s [%02d/%s/%d:%02d:%02d:%02d +0000] \"%s /%s%s HTTP/1.0\" %d %d\n", user, t.Day(), t.Month().String()[:3], t.Year(), t.Hour(), t.Minute(), t.Second(), verb, section, subsection, statusCode, byteSize)

		f.Close()
		rand.Seed(time.Now().UnixNano())

		Sleeps := [5]int{
			randomInt(10, 200),
			randomInt(200, 400),
			randomInt(400, 500),
			randomInt(500, 1000),
			randomInt(1000, 5000)}

		randomSleep := Sleeps[randomInt(0, 5)]
		fmt.Printf("Sleeping for %d ms\n", randomSleep)
		time.Sleep(time.Duration(randomSleep) * time.Millisecond)
	}
	//

}
