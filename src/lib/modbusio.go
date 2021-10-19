package lib

import (
	"log"
	"time"

	"github.com/simonvetter/modbus"
)

var ModbusRegistersFloat32 = map[string]uint16{
	"CURRENT_A": 3000 - 1,
	"CURRENT_B": 3002 - 1,
	"CURRENT_C": 3004 - 1,

	"VOLTAGE_A_B": 3020 - 1,
	"VOLTAGE_B_C": 3022 - 1,
	"VOLTAGE_C_A": 3024 - 1,

	"VOLTAGE_A_N": 3028 - 1,
	"VOLTAGE_B_N": 3030 - 1,
	"VOLTAGE_C_N": 3032 - 1,

	"ACTIVE_POWER_A": 3054 - 1,
	"ACTIVE_POWER_B": 3056 - 1,
	"ACTIVE_POWER_C": 3058 - 1,

	"FREQUENCY": 3110 - 1,
}

// Registers are 'published register list - 1'
// for clarity
// Modbus registers are zero indexed
var ModbusRegistersUint16 = map[string]uint16{
	"DATE_YEAR":          1837 - 1,
	"CMD_PRIV_INTERFACE": 5000 - 1,
	"CMD_INTERFACE":      5250 - 1,
	"CMD_PRIV_RESULT":    5126 - 1,
	"CMD_RESULT":         5376 - 1,
	"CMD_SEMAPHORE":      5680 - 1,
	"DO_STATUS_BITMAP":   9667 - 1,
	"DI_STATUS_BITMAP":   8905 - 1,
}

var ModbusCommands = map[string]uint16{
	"SET_DT":         1003,
	"DE_ENERGISE_DO": 6002,
	"ENERGISE_DO":    6003,
}

const (
	URL string = "tcp://10.3.0.3:502"
)

func GetRegisterUint16(register string) uint16 {
	var regUInt uint16 = ModbusRegistersUint16[register]
	log.Println("Register for", register, regUInt) // as unsigned integer
	return regUInt
}

func GetCommand(command string) uint16 {
	var commandUInt uint16 = ModbusCommands[command]
	log.Println("Register for", command, commandUInt) // as unsigned integer
	return commandUInt
}

func GetRegisterFloat32(register string) uint16 {
	var regUInt uint16 = ModbusRegistersFloat32[register]
	log.Println("Register for", register, regUInt) // as unsigned integer
	return regUInt
}

func SendCommand(client *modbus.ModbusClient, usePrivileged bool, command string, params []uint16) (uint16, error) {
	var err error
	var cmdInt uint16
	var cmdAddress uint16 = GetRegisterUint16("CMD_INTERFACE")
	var semaphore uint16 = 0

	// Privileged command interface needs the semaphore value to be retrieved and provided
	// with the command.
	// It will only be issued every 4 minutes. This tool does not cache the value.

	if usePrivileged {
		log.Println("Privileged command: Acquiring semaphore") // as unsigned integer

		cmdAddress = GetRegisterUint16("CMD_PRIV_INTERFACE")
		semaphore, err = client.ReadRegister(GetRegisterUint16("CMD_SEMAPHORE"), modbus.HOLDING_REGISTER)
		if err != nil {
			log.Println("Could not get semaphore") // as unsigned integer
		} else {
			if semaphore == 0 {
				log.Println("Invalid semaphore received - timeout?") // as unsigned integer
			}
			log.Println("Semaphore value:", semaphore) // as unsigned integer
		}
	}

	log.Println("Issuing command:", command, params)
	cmdInt = GetCommand(command)

	cmd := []uint16{cmdInt, semaphore}
	cmd = append(cmd, params...)
	log.Println("Address:", cmdAddress, "Data:", cmd)

	err = client.WriteRegisters(cmdAddress, cmd)
	if err != nil {
		return 0, err
	}
	var cmdStatus uint16

	var resultRegister string = "CMD_RESULT"
	if usePrivileged {
		resultRegister = "CMD_PRIV_RESULT"
	}
	cmdStatus, err = client.ReadRegister(GetRegisterUint16(resultRegister), modbus.HOLDING_REGISTER)
	return cmdStatus, err

}

func GetClient() (*modbus.ModbusClient, error) {
	var client *modbus.ModbusClient
	var err error

	client, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:     URL,
		Timeout: 1 * time.Second,
	})

	if err != nil {
		return client, err
	}

	err = client.Open()

	log.Println("Opened Connection to", URL)

	return client, err
}
