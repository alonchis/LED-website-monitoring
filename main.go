package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var flag bool

type healthcheck struct {
	Id string `json:"id"`
	status string `json:"status"`
}

type LED struct {
	//red int
	//green int
	green rpio.Pin
	red   rpio.Pin
}

type identifiable interface {
	PinInit(int, int)
	Off()
	flash()
	TestColor()
}

func (l *LED) PinInit(green int, red int) {
	l.green = rpio.Pin(green)
	l.red = rpio.Pin(red)
}

func (l *LED) Off() {
	l.green.Low()
	l.red.Low()
}

func (l *LED) flash() {
	l.red.Low()
	time.Sleep(time.Millisecond * 30)
	l.red.High()
	time.Sleep(time.Millisecond * 100)
	l.red.Low()
	//for i := 0; i < 6; i++ {
	//	pins[i] = rpio.Pin(pinsIndex[i])
	//	pins[i].Low()
	//}
}

func (l *LED) flow() {
	l.red.Low()
	l.green.High()
	time.Sleep(time.Millisecond * 50)
	l.green.Low()
	l.red.High()
}

func (l *LED) ChangeRed() {
	l.green.Low()
	l.red.High()
}

func (l *LED) ChangeGreen() {
	l.red.Low()
	l.green.High()
}

func (l *LED) TestColor() {
	l.Off()
	l.green.High()
	time.Sleep(time.Millisecond * 500) //on for half a second
	l.green.Low()
	l.red.High()
	time.Sleep(time.Millisecond * 500) //on for half a second
	l.Off()
}
var pinsIndex = []int{14, 15, 23, 18, 4, 24} //pattern is red, green, red...
var pins [6]rpio.Pin

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

	//go LightLoop()

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("\nReceived an interrupt, stopping services...\n")
		/** Toggle - infinite loop */
		for i := 0; i < 6; i++ {
			pins[i] = rpio.Pin(pinsIndex[i])
			pins[i].Low()
		}

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
	log.Fatal(http.ListenAndServe(":80", nil))
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
		time.Sleep(time.Minute * 3 )
	}
	return
}

func HandleOff(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "turning all off\n")
	fmt.Println("Endpoint Hit: off")
	flag = false
	Off()
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome to the homepage\n")
	fmt.Println("Endpoint Hit: homepage")
	Blink()
}

func HandleRed(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Turning all LED red\n")
	fmt.Println("Endpoint Hit: RED")
	SolidRed()
}

func HandleGreen(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Turning all LED green\n")
	fmt.Println("Endpoint Hit: GREEN")
	SolidGreen()
}

func Blink() {
		/** flowing LED light in on direction */
		for i := 0; i < len(Leds); i++ {
			Leds[i].flow()
		}

		/** flowing LED light in reverse */
		for i := len(Leds) - 1; i >= 0; i-- {
			Leds[i].flow()
		}

	for i := len(Leds) - 1; i >= 0; i-- {
		Leds[i].Off()
	}
}


//LightLoop does the magic led loop
func LightLoop() {
	//set  leds to red
	for i := 0; i < len(Leds); i++ {
		Leds[i].red.High()
	}

	for {
		/** flowing LED light in on direction */
		for i := 0; i < len(Leds); i++ {
			Leds[i].flow()
		}

		/** flowing LED light in reverse */
		for i := len(Leds) - 1; i >= 0; i-- {
			Leds[i].flow()
		}
	}
}

func SolidGreen() {
	for i := 0; i < len(Leds); i++ {
		Leds[i].red.Low()
		Leds[i].green.High()
	}
}

func SolidRed() {
	for i := 0; i < len(Leds); i++ {
		Leds[i].green.Low()
		Leds[i].red.High()
	}
}


func Off() {
	hr, _, _ := time.Now().Clock()
	fmt.Printf("turning off at %v\n", hr)
	for i := 0; i < len(Leds); i++ {
		Leds[i].Off()
	}
}

func StartChecks() {
	sites := [3]string{}

	//site one
	resp, err := http.Get(sites[0])
	if err != nil {
		// handle error
		Leds[0].ChangeRed()
		log.Printf("site %s unreachable\n", sites[0])
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if resp.StatusCode < 400 {
		log.Printf("site %s returned %v \n", sites[0], resp.StatusCode)
		Leds[0].ChangeGreen()
	}

	//site two
	resp, err = http.Get(sites[1])
	if err != nil {
		log.Printf("site %s unreachable\n", sites[1])
		// handle error
		Leds[1].ChangeRed()
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if resp.StatusCode < 400 {
		log.Printf("site %s returned %v \n", sites[1], resp.StatusCode)
		Leds[1].ChangeGreen()
	}

	//site three
	resp, err = http.Get(sites[2])
	if err != nil {
		log.Printf("site %s unreachable\n", sites[2])
		// handle error
		Leds[2].ChangeRed()
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if resp.StatusCode < 400 {
		log.Printf("site %s returned %v\n", sites[2], resp.StatusCode)
		Leds[2].ChangeGreen()
	} else {
		Leds[2].ChangeRed()
	}
}