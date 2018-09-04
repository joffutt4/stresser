package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var baseport string

func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func stressHandler(c *gin.Context) {
	totalCon := 0
	failedCon := 0
	numRun := make(chan int)
	numErr := make(chan int)

	//Create a go routine for every thread asked for from command line
	// for i := 0; i < threads; i++ {
	// 	go stress(url, numRun, numErr)
	// }
	//Set a timer to finish for the number of seconds passed in from command line
	timer := time.After(time.Duration(time.Second * time.Duration(1)))

	c.HTML(http.StatusOK, "stress.html", gin.H{})

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

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Stresser usage: *port*")
		os.Exit(1)
	}
	baseport := ":" + os.Args[1]
	fmt.Println("To start testing, open http:localhost" + baseport)
	router := gin.Default()
	// Starts a new session
	store := sessions.NewCookieStore([]byte("secret"))
	store.Options(sessions.Options{
		Path: "/",
		// Setting max age to 12 hours
		MaxAge: 43200,
	})
	router.Use(sessions.Sessions("stresser", store))
	router.Static("/css", "./static/css")
	router.Static("/img", "./static/img")
	// router.Static("/header", "./templates/header.tmpl")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", homeHandler)
	router.GET("/stress", stressHandler)

	// providing only the port here makes Gin run on 0.0.0.0:port which allows it to function on any server
	router.Run(baseport)
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
