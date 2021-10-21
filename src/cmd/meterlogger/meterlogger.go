package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/simonvetter/modbus"

	"github.com/ArrEssJay/PM5350_Modbus/src/lib"

	"github.com/DataDog/datadog-go/v5/statsd"
)

const (
	SAMPLING_INTERVAL uint = 1 //seconds
)

var mbclient *modbus.ModbusClient

func cleanup() {
	fmt.Println("Exiting")
	mbclient.Close()
}

func main() {
	// Handle exit
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()

	var clientErr error

	mbclient, clientErr = lib.GetClient()
	if clientErr != nil {
		log.Fatal(clientErr.Error())
		os.Exit(1)
	}

	statsd, statsderr := statsd.New("unix:///var/run/datadog/dsd.socket", statsd.WithoutClientSideAggregation())
	if statsderr != nil {
		log.Fatal(statsderr)
		os.Exit(1)
	}

	log.Println("Logging Meter Data to Datadog....")

	for {
		for key, value := range lib.ModbusRegistersFloat32 {

			var regVal float32
			var readErr error
			regVal, readErr = mbclient.ReadFloat32(value, modbus.HOLDING_REGISTER)
			if readErr == nil {
				// log.Println(key, "=>", ":", float64(regVal))
				statsd.Distribution(key, float64(regVal), []string{"environment:dev"}, 1)
			} else {
				log.Println("Register:", key, "=>", "ERROR", readErr)
			}

		}
		time.Sleep(time.Duration(SAMPLING_INTERVAL) * time.Second)
	}
}
