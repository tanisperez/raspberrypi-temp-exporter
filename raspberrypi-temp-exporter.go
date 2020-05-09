package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var temperatureFilePath string
var temperatureReadInterval int64
var prometheusPortExporter string

var (
	cpuTemp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "raspberry_pi_temperature_celsius",
		Help: "Raspberry Pi Temperature",
	})
)

// Loads environment variables
func loadEnvironmentVariables() {
	temperatureFilePath = os.Getenv("TEMPERATURE_FILE_PATH")
	if len(temperatureFilePath) == 0 {
		log.Fatal("TEMPERATURE_FILE_PATH environment variable not set")
	}

	var readInterval = os.Getenv("TEMPERATURE_READ_INTERVAL")
	if len(readInterval) == 0 {
		log.Fatal("TEMPERATURE_READ_INTERVAL environment variable not set")
	}

	temperatureInterval, err := strconv.ParseInt(readInterval, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	temperatureReadInterval = temperatureInterval

	prometheusPortExporter = os.Getenv("PROMETHEUS_EXPORTER_PORT")
	if len(prometheusPortExporter) == 0 {
		log.Fatal("PROMETHEUS_EXPORTER_PORT environment variable not set")
	}
}

// Record the Raspberry Pi Temperature
func recordTemperatureMetrics() {
	go func() {
		for {
			var temperature float64 = getRaspberryPiTemperature()
			cpuTemp.Set(temperature)

			time.Sleep(time.Duration(temperatureReadInterval) * time.Second)
		}
	}()
}

// Read from a system file the current temperature
func getRaspberryPiTemperature() float64 {
	data, err := ioutil.ReadFile(temperatureFilePath)
	if err != nil {
		log.Fatal(err)
	}

	temperature, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		log.Fatal(err)
	}
	var temperatureCelsius float64 = temperature / 1000
	return temperatureCelsius
}

// Main func
func main() {
	loadEnvironmentVariables()
	recordTemperatureMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+prometheusPortExporter, nil)
}
