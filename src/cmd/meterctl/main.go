package main

import (
	"log"
	"os"

	"GoMeter/src/lib"

	"github.com/jessevdk/go-flags"
	"github.com/simonvetter/modbus"
)

func main() {
	type Options struct {
		DOIndex uint `short:"i" choice:"1" choice:"2" description:"DO Index" required:"true"`
		DOState uint `short:"s" choice:"0" choice:"1" description:"DO State" required:"true" `
	}
	var options Options
	var parser = flags.NewParser(&options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}

	var clientErr error
	var mbclient *modbus.ModbusClient

	mbclient, clientErr = lib.GetClient()
	if clientErr != nil {
		log.Println(clientErr.Error())
		os.Exit(1)
	}

	// state to command map
	cmdstr := "DE_ENERGISE_DO"
	if options.DOState > 0 {
		cmdstr = "ENERGISE_DO"
	}

	var cmderr error
	var cmdStatus uint16

	cmdStatus, cmderr = lib.SendCommand(mbclient, false, cmdstr, []uint16{uint16(options.DOIndex)})
	if cmderr != nil {
		log.Println(cmderr)
	} else {
		log.Println("Command status:", cmdStatus)

	}

	// Print DI and DO state
	var DOreg uint16
	var regerr error
	DOreg, regerr = mbclient.ReadRegister(lib.GetRegisterUint16("DO_STATUS_BITMAP"), modbus.HOLDING_REGISTER)

	if regerr != nil {
		log.Println(regerr)
	} else {
		// 2 DIs on PM5350
		log.Println("DO Status Bitmap:", DOreg&1 > 0, DOreg&2 > 0)
	}
	var DIreg uint16
	DIreg, regerr = mbclient.ReadRegister(lib.GetRegisterUint16("DI_STATUS_BITMAP"), modbus.HOLDING_REGISTER)

	if regerr != nil {
		log.Println(regerr)
	} else {
		// 4 DIs on PM5350
		log.Println("DI Status Bitmap:", DIreg&1 > 0, DIreg&2 > 0, DIreg&3 > 0, DIreg&4 > 0)
	}
	mbclient.Close()
}
