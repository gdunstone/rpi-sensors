package main

import (
	"github.com/BurntSushi/toml"
	"github.com/shawntoffel/darksky"
	"fmt"
	"os"
	"time"
	"strings"
	"reflect"
	"flag"
	"sort"
)

type Config struct{
	Key string `toml:"key"`
	Latitude float64 `toml:"latitude"`
	Longitude float64 `toml:"longitude"`
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
	configPath := flag.Args()[0]

	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil{
		panic(err)
	}
	client := darksky.New(config.Key)
	request := darksky.ForecastRequest{}

	request.Latitude = darksky.Measurement(config.Latitude)
	request.Longitude = darksky.Measurement(config.Longitude)
	request.Options = darksky.ForecastRequestOptions{Exclude: "hourly,minutely,flags", Units: "si"}

	forecast, err := client.Forecast(request)

	if err != nil{
		panic(err)
	}
	mapData, err := ToMap(forecast.Currently)
	if err != nil{
		panic(err)
	}
	formatOutput("darksky", mapData, time.Now().UnixNano())
	for _,fcast := range forecast.Daily.Data{
		md, err := ToMap(fcast)
		if err != nil{
			fmt.Errorf("%s", err)
			continue
		}

		formatOutput("darksky-forecast", md, int64(fcast.Time) * 1000000000)

	}
	os.Exit(0)
}

// ToMap converts a struct to a map using the struct's tags.
//
// ToMap uses tags on struct fields to decide which fields to add to the
// returned map.
func ToMap(in interface{}) (map[string]interface{}, error){
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		//fi := typ.Field(i)
		name := typ.Field(i).Name
		out[strings.ToLower(name)] = v.Field(i).Interface()

	}
	return out, nil
}