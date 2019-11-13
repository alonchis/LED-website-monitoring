package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var flag bool
var pinsIndex = []int{14, 15, 23, 18, 4, 24} //pattern is red, green, red...
var Leds [3]LED //array to hold all leds

func main() {
	/** Open and map memory to access gpio, check for errors */
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	/** Unmap gpio memory when done */
	defer rpio.Close()

	//Init all leds
	// probably could be written better
	counter := 0
	for i := 0; i < len(Leds); i++ {
		Leds[i].PinInit(pinsIndex[counter], pinsIndex[counter+1])
		counter = counter + 2
		Leds[i].Off() //set state to off
	}

	//Test leds
	for i := 0; i < len(Leds); i++ {
		Leds[i].TestColor()
	}

	//start api
	handleRequests()

	//listen for interrupt and teardown (turn off leds)
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("\nReceived an interrupt, stopping services...\n")
		/** turn off leds on exit */
		Off()

		close(cleanupDone)
	}()
	<-cleanupDone

	fmt.Println("this code should be unreachable")
}

func handleRequests() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/red", HandleRed)
	http.HandleFunc("/green", HandleGreen)
	http.HandleFunc("/off", HandleOff)
	http.HandleFunc("/healthcheck", StartCheck)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func StartCheck(writer http.ResponseWriter, request *http.Request) {
	hr, _, _ := time.Now().Clock()

	if hr == 0 || hr == 1 || hr == 2 || hr == 3 || hr == 4 || hr == 5 || hr == 6 {
		return
	}

	fmt.Fprintf(writer, "Starting healthchecks\n")
	fmt.Println("Endpoint Hit: off")
	flag = true
	for flag != false {
		go StartChecks()
		fmt.Println("go func called")
		time.Sleep(time.Minute * 3)
	}
	return
}

// HandleOff handles func to turn all leds OFF
func HandleOff(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "turning all off\n")
	fmt.Println("Endpoint Hit: off")
	flag = false
	Off()
}

// HomePage default route
func HomePage(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./views"))
	fmt.Fprintf(w, "welcome to the homepage\n")
	fmt.Println("Endpoint Hit: homepage")

	Blink()
}

// HandleRed handles func to turn all leds red
func HandleRed(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Turning all LED red\n")
	fmt.Println("Endpoint Hit: RED")
	SolidRed()
}

// HandleGreen handles func to turn all leds green
func HandleGreen(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Turning all LED green\n")
	fmt.Println("Endpoint Hit: GREEN")
	SolidGreen()
}
