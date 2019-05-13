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

//pins are using BCM layout
const TRIGGER     = 18
const ECHO        = 24
const MIN_DIST    = 200 //minimum obstacle distance in mm
const SOUND_SPEED float64 = 0.00034 //millimeters per nanosecond
const STOP_URL    = "http://192.168.8.223:5000/move"

var (
    trigger_pin   = rpio.Pin(TRIGGER)
    echo_pin      = rpio.Pin(ECHO)
    previous_dist float64 = 20000 //initialze with some big value for first loop
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
        if dist < MIN_DIST && dist < previous_dist {
            stop_the_car()
            fmt.Println("Distance in mm:")
            fmt.Println(dist)
        }
        previous_dist = dist
        time.Sleep(time.Second / 10)
    }
}

func pin_setup() {
    trigger_pin.Output()
    echo_pin.Input()
    trigger_pin.Low()
}

func stop_the_car() {
    jsonData := map[string]string{"direction": "stop", "speed": "0"}
    jsonValue, _ := json.Marshal(jsonData)
    response, err := http.Post(STOP_URL, "application/json", bytes.NewBuffer(jsonValue))

    if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
    } else {
        data, _ := ioutil.ReadAll(response.Body)
        fmt.Println(string(data))
    }
}

func distance() (distance float64) {
    trigger_pin.High()
    time.Sleep(0.00001)
    trigger_pin.Low()
    pulse_start_time := time.Now().UnixNano()
    for echo_pin != 1 {
        time.Sleep(1)
    }
    pulse_stop_time := time.Now().UnixNano()
    return (float64(pulse_stop_time) - float64(pulse_start_time)) * SOUND_SPEED / 2
}
