package main

import (
	"os"

        "github.com/rs/zerolog"
        "github.com/rs/zerolog/log"

        "github.com/ArrEssJay/PM5350_Modbus/src/lib"

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
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var clientErr error
	var mbclient *modbus.ModbusClient

	mbclient, clientErr = lib.GetClient()
	if clientErr != nil {
		log.Error().Err(clientErr)
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
		log.Error().Err(cmderr)
	} else {
		log.Debug().
        	Uint16("Command Status", cmdStatus)
	}

	// Print DI and DO state
	var DOreg uint16
	var regerr error
	DOreg, regerr = mbclient.ReadRegister(lib.GetRegisterUint16("DO_STATUS_BITMAP"), modbus.HOLDING_REGISTER)

	if regerr != nil {
		log.Error().Err(regerr)
	} else {
		// 2 DIs on PM5350
		log.Info().
		Bool("DO1", DOreg&1 > 0).
                Bool("DO2", DOreg&2 > 0).Msg("DO Status Bitmap")
	}
	var DIreg uint16
	DIreg, regerr = mbclient.ReadRegister(lib.GetRegisterUint16("DI_STATUS_BITMAP"), modbus.HOLDING_REGISTER)

	if regerr != nil {
		log.Error().Err(regerr)
	} else {
		// 4 DIs on PM5350
		log.Info().
                Bool("DI1", DIreg&1 > 0).
                Bool("DI2", DIreg&2 > 0).
                Bool("DI3", DIreg&3 > 0).
                Bool("DI4", DIreg&4 > 0).Msg("DI Status Bitmap")
	}
	mbclient.Close()
}
