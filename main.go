package main

import (
			"fmt"
			"os"
			"time"
			"github.com/stianeikeland/go-rpio"
			"bytes"
    		"encoding/json"
		    "io/ioutil"
    		"net/http"
		)

const TRIGGER   = 18
const ECHO      = 24
const RED_LED   = 23
const GREEN_LED = 25

var (
	trigger_pin = rpio.Pin(TRIGGER)
	echo_pin    = rpio.Pin(ECHO)
	red_pin     = rpio.Pin(RED_LED)
	green_pin   = rpio.Pin(GREEN_LED)
	previous_dist float64

)
func main() {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()

	pin_setup()    
    for {
    	dist := distance()
    	fmt.Println("Distance in mm:")
    	fmt.Println(dist)
    	time.Sleep(1 * time.Second)
    }
}

func pin_setup() {
	trigger_pin.Output()
	echo_pin.Input()
	red_pin.Output()
	green_pin.Output()
	
	for x := 0; x < 5; x++ {
		red_pin.Toggle()
		time.Sleep(time.Second / 5)
	}
	for x := 0; x < 5; x++ {
		green_pin.Toggle()
		time.Sleep(time.Second / 5)
	}
	red_pin.Low()
	green_pin.Low()
	trigger_pin.Low()
}

func stop_the_car() {
	jsonData := map[string]string{"direction": "stop", "speed": "0"}
    jsonValue, _ := json.Marshal(jsonData)
    response, err := http.Post("http://192.168.4.1:5000/move", "application/json", bytes.NewBuffer(jsonValue))
    if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
    } else {
        data, _ := ioutil.ReadAll(response.Body)
        fmt.Println(string(data))
    }
}

func distance() (distance float64) {
	trigger_pin.Low()
	time.Sleep(5 * time.Microsecond)
	trigger_pin.High()
	time.Sleep(10 * time.Microsecond)
	trigger_pin.Low()
	start := time.Now()
	stop  := time.Now()

	for echo_pin.Read() == rpio.Low {
		start = time.Now()
	}
	for echo_pin.Read() == rpio.High {
		stop = time.Now()		
	}

	time_elapsed := stop.Sub(start)
	distance = (float64(time_elapsed) / 2) * float64(0.00034)
	if previous_dist==0 { previous_dist = distance }
	if distance < 200 && distance < previous_dist {
    	stop_the_car()
    }
    previous_dist = distance
	return distance
}