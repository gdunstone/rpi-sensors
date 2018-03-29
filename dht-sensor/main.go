package main

import (
	"flag"
	"fmt"
	"github.com/d2r2/go-dht"
	"math"
	"os"
	"strings"
	"time"
)

var (
	pin           int
	stype         string
	boostPerfFlag bool
)

func init() {
	flag.IntVar(&pin, "pin", 4, "pin")
	flag.StringVar(&stype, "sensor-type", "dht22", "sensor type (dht22, dht11)")
	flag.BoolVar(&boostPerfFlag, "boost", false, "boost performance")
}

func formatOutput(sensorType string, values map[string]interface{}) {

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
	str := fmt.Sprintf("%s %s", sensorType, csv)
	// add timestamp
	str = fmt.Sprintf("%s %d", str, time.Now().UnixNano())
	fmt.Fprintln(os.Stdout, str)
	os.Exit(0)
}

func main() {
	flag.Parse()
	var sensorType dht.SensorType

	if stype == "dht22" || stype == "am2302" {
		sensorType = dht.DHT22
	} else if stype == "dht11" {
		sensorType = dht.DHT11
	}
	values := make(map[string]interface{})

	temperature, humidity, retried, err :=
		dht.ReadDHTxxWithRetry(sensorType, pin, boostPerfFlag, 10)
	if err != nil {
		panic(err)
	}

	values["retried"] = retried
	values["temperature"] = temperature
	values["humidity"] = humidity

	// calculate vpd
	// J. Win. (https://physics.stackexchange.com/users/1680/j-win),
	// How can I calculate Vapor Pressure Deficit from Temperature and Relative Humidity?,
	// URL (version: 2011-02-03): https://physics.stackexchange.com/q/4553
	temperature64 := float64(temperature)

	humidity64 := float64(humidity)

	es := 0.6108 * math.Exp(17.27*temperature64/(temperature64+237.3))
	ea := humidity64 / 100 * es

	// this equation returns a negative value (in kPa), which while technically correct,
	// is invalid in this case because we are talking about a deficit.
	vpd := (ea - es) * -1
	values["vpd"] = vpd
	formatOutput(stype, values)
}
