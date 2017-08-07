package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	if len(os.Args) < 4 {
		fmt.Println("Stresser usage: *url endpoint* *number of connections* *how long to run for in seconds*")
		os.Exit(1)
	}
	url := os.Args[1]
	threads, err := strconv.Atoi(os.Args[2])
	check(err)
	length, err := strconv.Atoi(os.Args[3])
	check(err)
	numRun := make(chan int)
	numErr := make(chan int)
	totalCon := 0
	failedCon := 0

	//Create a go routine for every thread asked for from command line
	for i := 0; i < threads; i++ {
		go stress(url, numRun, numErr)
	}
	//Set a timer to finish for the number of seconds passed in from command line
	timer := time.After(time.Duration(time.Second * time.Duration(length)))

	//Continually loop through this code receiving either successful connections or errors
	//Also catches when the timer finishes and then returns out of the entire application
	for {
		select {
		case count := <-numRun:
			totalCon = count + totalCon
		case count := <-numErr:
			failedCon = count + failedCon
		case <-timer:
			log.Println("Finished stressing")
			log.Println("Total connections:", totalCon)
			log.Println("Total failures:", failedCon)
			return
		}
	}
}

//check is a simple error checking helper function that makes it so we don't have the same code repeated throughout the app.
func check(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

//stress infinitely loops while calling the url passed in.
//It also pushes a 1 to either the numRun channel or the numErr channel based on a successful query or not.
func stress(url string, numRun chan int, numErr chan int) {
	//Infinitely call the url passed in
	for {
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err.Error())
			numErr <- 1
			continue
		}
		numRun <- 1
		resp.Body.Close()
	}
}
