package main

import (
	"flag"
	"fmt"
	"github.com/aquarat/embd"
	_ "github.com/aquarat/embd/host/all"
	"github.com/aquarat/embd/sensor/bh1750fvi"
	"github.com/aquarat/embd/sensor/bme280"
	"github.com/aquarat/embd/sensor/bmp085"
	"github.com/aquarat/embd/sensor/bmp180"
	"github.com/aquarat/embd/sensor/l3gd20"
	"github.com/aquarat/embd/sensor/lsm303"
	"math"
	"os"
	"strings"
	"time"
)

var (
	stype  string
	pullUp bool
)

func init() {
	flag.StringVar(&stype, "sensor-type", "bmp180", "sensor type (bh1750fvi, bme280, bmp085, bmp180, l3gd20, lsm303)")
	flag.BoolVar(&pullUp, "pull-up", false, "use pull-up address, for when SDO is pulled up (connected to vddio)")
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

	bus := embd.NewI2CBus(1)

	values := make(map[string]interface{})
	switch stype {
	case "bmp180":
		sensor := bmp180.New(bus)
		temperature, err := sensor.Temperature()
		if err != nil {
			panic(err)
		}
		values["temperature"] = temperature
		pressure, err := sensor.Pressure()
		if err != nil {
			panic(err)
		}
		values["pressure"] = pressure
		altitude, err := sensor.Altitude()
		if err != nil {
			panic(err)
		}
		values["altitude"] = altitude
	case "bme280":
		// by default pull SDO down.
		// if using GY-BME280, connect to 3v3 ONLY! not 5v.
		var addr byte
		addr = 0x76
		if pullUp {
			addr = 0x77
		}

		sensor, err := bme280.New(bus, addr)
		if err != nil {
			panic(err)
		}
		measurements, err := sensor.Measurements()
		if err != nil {
			panic(err)
		}

		temperature64 := sensor.Temperature(measurements)
		humidity64 := sensor.Humidity(measurements)

		values["temperature"] = temperature64
		values["humidity"] = humidity64
		values["pressure"] = sensor.Pressure(measurements)

		es := 0.6108 * math.Exp(17.27*temperature64/(temperature64+237.3))
		ea := humidity64 / 100 * es
		// this equation returns a negative value (in kPa), which while technically correct,
		// is invalid in this case because we are talking about a deficit.
		values["vpd"] = (ea - es) * -1
	case "bmp085":
		sensor := bmp085.New(bus)

		pressure, err := sensor.Pressure()
		if err != nil {
			panic(err)
		}
		values["pressure"] = pressure
		altitude, err := sensor.Altitude()
		if err != nil {
			panic(err)
		}
		values["altitude"] = altitude
	case "lsm303":
		sensor := lsm303.New(bus)
		heading, err := sensor.Heading()
		if err != nil {
			panic(err)
		}
		values["heading"] = heading
	case "bh1750fvi":
		sensor := bh1750fvi.New("H2", bus)
		lux, err := sensor.Lighting()
		if err != nil {
			panic(err)
		}
		values["lux"] = lux
	case "l3gd20":
		sensor := l3gd20.New(bus, l3gd20.R2000DPS)
		temperature, err := sensor.Temperature()
		if err != nil {
			panic(err)
		}
		values["temperature"] = temperature
	}
	formatOutput(stype, values)
}
