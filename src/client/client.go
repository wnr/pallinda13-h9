package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func Main() {
	server := []string{
		"http://localhost:8080",
		"http://localhost:8081",
		"http://localhost:8082",
	}
	for {
		before := time.Now()
//		res := Get(server[0])
//		res := Read(server[0], time.Second)
		res := MultiRead(server, time.Second)
		after := time.Now()
		fmt.Println("Response:", *res)
		fmt.Println("Time:", after.Sub(before))
		fmt.Println()
		time.Sleep(500 * time.Millisecond)
	}
}

type Response struct {
	Body       string
	StatusCode int
}

// Get makes an HTTP Get request and returns an abbreviated response.
// Status code 200 means that the request was successful.
// The function returns &Response{"", 0} if the request fails
// and it blocks forever if the server doesn't respond.
func Get(url string) *Response {
	res, err := http.Get(url)
	if err != nil {
		return &Response{}
	}
	// res.Body != nil when err == nil
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
	}
	return &Response{string(body), res.StatusCode}
}

//This method calls Get(url) in an go routine that will run until a response from Get is obtained.
//Read will read the response from Get if timeout haven't been exceeded.
//If timeout will be hit, the go routine will still read value from Get but the value
//will not be read on the channel.
//+This fixes the bug that res pointer can be overwritten by old go routines.
//+This also fixes so that the go routines can exit if timeout have been hit.
func Read(url string, timeout time.Duration) (res *Response) {
	done := make(chan *Response, 1) //make buffered so that the go routine can exit. Also make channel accept *Response.
	go func() { //the routine will not modify res since timeout may have been exceeded.
		done <- Get(url) //since its buffered the routine can exit event if no-one reads the channel
	}()
	select {
	case res = <-done:
	case <-time.After(timeout):
		res = &Response{"Gateway timeout\n", 504}
	}
	return
}

// MultiRead makes an HTTP Get request to each url and returns
// the response of the first server to answer with status code 200.
// If none of the servers answer before timeout, the response is
// 503 â€“ Service unavailable.
func MultiRead(urls []string, timeout time.Duration) (res *Response) {
	responseChannel := make(chan *Response, len(urls))

	for _, url := range urls {
		go func() {
			r := Read(url, timeout)
			if (r.StatusCode == 200) {
				responseChannel<- r
			}
		}()
	}

	select {
	case res = <-responseChannel:
	case <-time.After(timeout):
		res = &Response{"Service unavailable\n", 503}
	}
	return
}
