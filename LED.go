package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//LED has two rpio.pin types, green and blue, respectively
type LED struct {
	green rpio.Pin
	red   rpio.Pin
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
	//add websites to monitor here
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
