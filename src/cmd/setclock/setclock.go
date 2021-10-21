package main

import (
	"log"
	"os"
	"time"

	"github.com/ArrEssJay/PM5350_Modbus/src/lib"

	"github.com/simonvetter/modbus"
)

func main() {
	var clientErr error
	var mbclient *modbus.ModbusClient

	mbclient, clientErr = lib.GetClient()
	if clientErr != nil {
		log.Println(clientErr.Error())
		os.Exit(1)
	}

	log.Println("Setting meter time to local time")

	var cmderr error
	var cmdStatus uint16

	var meterDate []uint16

	now := time.Now()
	log.Println("Local time:", now)

	// Only setting the meter time to second precision as network latency etc.
	// make millisecond precision unlikely to be useful
	cmdStatus, cmderr = lib.SendCommand(mbclient, true, "SET_DT", []uint16{uint16(now.Year()), uint16(now.Month()), uint16(now.Day()), uint16(now.Hour()), uint16(now.Minute()), uint16(now.Second()), uint16(now.Nanosecond() / 1e6)})
	if cmderr != nil {
		log.Println("Error:", cmderr.Error())

	} else if cmdStatus != 0 {
		log.Println("Command Status !=0:", cmdStatus)
	} else {
		log.Println("Set meter time ok")
	}

	meterDate, cmderr = mbclient.ReadRegisters(lib.GetRegisterUint16("DATE_YEAR"), 7, modbus.HOLDING_REGISTER)
	if cmderr != nil {
		log.Println("Error reading meter date:", cmderr.Error())
	} else {
		var meterTime = time.Date(
			int(meterDate[0]),        //year
			time.Month(meterDate[1]), //month
			int(meterDate[2]),        //day
			int(meterDate[3]),        //hour
			int(meterDate[4]),        //minute
			int(meterDate[5]),        //second
			int(meterDate[6])*1e6,    //nano
			time.UTC)
		log.Println("New meter date:", meterTime, meterDate)
	}
	mbclient.Close()
}
