package main

import (
	"flag"
	"fmt"
	"os"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
	"strconv"
	"strings"
	"time"
)

// this many mm per bucket tip
const mmPerBucketTip = 0.2794 // mm per bucket tip

const stype = "rainmeter"

var pin int

func init() {
	flag.IntVar(&pin, "pin", 17, "pin that the rainmeter is connected to")
}

func formatOutput(sensor_type string, values map[string]interface{}) {

	keyvaluepairs := make([]string, 0)

	for key, val := range values {
		switch v := val.(type) {
		case int:
			keyvaluepairs = append(keyvaluepairs, fmt.Sprintf("%s=%d", key, v))
		case float64:
			keyvaluepairs = append(keyvaluepairs, fmt.Sprintf("%s=%f", key, v))
		case float32:
			keyvaluepairs = append(keyvaluepairs, fmt.Sprintf("%s=%f", key, v))
		case string:
			keyvaluepairs = append(keyvaluepairs, fmt.Sprintf("%s=%s", key, v))
		}
	}
	csv := strings.Join(keyvaluepairs, ",")
	str := fmt.Sprintf("%s %s", sensor_type, csv)
	// add timestamp
	str = fmt.Sprintf("%s %d", str, time.Now().UnixNano())
	fmt.Fprintln(os.Stdout, str)
	os.Exit(0)
}

func main() {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		panic(err)
	}

	// Lookup a pin by its number:
	p := gpioreg.ByName(strconv.Itoa(pin))

	fmt.Printf("%s: %s\n", p, p.Function())

	// Set it as input, with an internal pull down resistor:
	if err := p.In(gpio.PullUp, gpio.BothEdges); err != nil {
		panic(err)
	}

	// Wait for edges as detected by the hardware, and print the value read:

	values := make(map[string]interface{})
	values["mmPerHour"] = 0


	for {
		detected := p.WaitForEdge(time.Duration(10) * time.Second)
		read := p.Read()
		if !detected{
			formatOutput(stype, values)
		}
		if read == gpio.High{
			break
		}
	}
	start := time.Now()

	for {
		detected := p.WaitForEdge(time.Duration(10) * time.Second)
		read := p.Read()
		if !detected{
			formatOutput(stype, values)
		}
		if read == gpio.High{
			break
		}
	}

	dur := time.Since(start)

	values["mmPerHour"] = mmPerBucketTip / (dur.Seconds() / 60)

	formatOutput(stype, values)

}
