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
	"log"
	"sort"
	"time"
	"github.com/shawntoffel/darksky"
)

var (
	stype  string
	pullUp bool
)

func init() {
	// if we dont do this, embd will fuck up our output by outputting to stdout
	log.SetOutput(os.Stderr)
	flag.StringVar(&stype, "sensor-type", "bmp180", "sensor type (bh1750fvi, bme280, bmp280, bmp085, bmp180, l3gd20, lsm303)")
	flag.BoolVar(&pullUp, "pull-up", false, "use pull-up address, for when SDO is pulled up (connected to vddio)")
}

func formatOutput(sensorType string, values map[string]interface{}, t int64) {

	keyvaluepairs := make([]string, 0)

	keys := make([]string, 0)
	for k, _ := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := values[key]
		switch v := val.(type) {
		case int:
			keyvaluepairs = append(keyvaluepairs, fmt.Sprintf("%s=%di", key, v))
		case darksky.Measurement:
			keyvaluepairs = append(keyvaluepairs, fmt.Sprintf("%s=%f", key, float64(v)))
		case float64:
			keyvaluepairs = append(keyvaluepairs, fmt.Sprintf("%s=%f", key, v))
		case float32:
			keyvaluepairs = append(keyvaluepairs, fmt.Sprintf("%s=%f", key, v))
		case string:
			if v == ""{
				continue
			}
			keyvaluepairs = append(keyvaluepairs, fmt.Sprintf("%s=\"%s\"", key, v))
		case bool:
			keyvaluepairs = append(keyvaluepairs, fmt.Sprintf("%s=%t", key, v))
		}
	}
	csv := strings.Join(keyvaluepairs, ",")
	str := fmt.Sprintf("%s %s", sensorType, csv)
	// add timestamp
	str = fmt.Sprintf("%s %d", str, t)
	fmt.Fprintln(os.Stdout, str)
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
	case "bme280", "bmp280":
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
		values["temperature"] = temperature64
		// pressure in hectopascals
		pressure64 := sensor.Pressure(measurements)
		values["pressure"] = pressure64
		if stype == "bme280" {
			// bmp280 doesnt have humidity

			humidity64 := sensor.Humidity(measurements)
			values["humidity"] = humidity64
			// saturated vapor pressure
			es := 0.6108 * math.Exp(17.27 * temperature64 / (temperature64 + 237.3))

			// actual vapor pressure
			ea := humidity64 / 100 * es

			// mixing ratio
			//w := 621.97 * ea / ((pressure64/10) - ea)
			// saturated mixing ratio
			//ws := 621.97 * es / ((pressure64/10) - es)
			// absolute humidity (in kg/m³)
			ah := ea / (461.5 * (temperature64 + 273.15))

			// report it as g/m³
			values["absolute_humidity"] = ah*1000
			// this equation returns a negative value (in kPa), which while technically correct,
			// is invalid in this case because we are talking about a deficit.
			values["vpd"] = (ea - es) * -1
		}
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
	formatOutput(stype, values, time.Now().UnixNano())
}
