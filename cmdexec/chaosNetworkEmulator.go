package cmdexec

import (
	"encoding/json"
	"io"
	"log"
	"strings"
)

const CONTAINER_NETWORK_EMULATOR = "netem"

type ChaosNetworkEmulatorSimulator struct {
	*ChaosCommand
}

type networkEmulatorStruct struct {
	Command     string
	CommandArgs []Arg
}

func fromJson(jsonStream io.Reader) networkEmulatorStruct {
	var args []networkEmulatorStruct

	dec := json.NewDecoder(jsonStream)
	for {
		var m networkEmulatorStruct
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		args = append(args, m)
	}

	return args[0]
}

func CreateChaosNetworkEmulatorSimulator(container string, jsonStream io.Reader) ChaosNetworkEmulatorSimulator {

	networkEmulatorStructValue := fromJson(jsonStream)
	networkEmulatorOptions := extractParameters(networkEmulatorStructValue.CommandArgs, "--duration", "-d", "--interface", "-i")
	networkEmulatorSubcommandOptions := excludeParameters(networkEmulatorStructValue.CommandArgs, "--duration", "-d", "--interface", "-i")

	containerNetworkEmulator := strings.Join([]string{CONTAINER_NETWORK_EMULATOR, concatenateArgs(networkEmulatorOptions...), networkEmulatorStructValue.Command}, " ")

	return ChaosNetworkEmulatorSimulator{createChaosWithOperationAndArgs(container, containerNetworkEmulator, networkEmulatorSubcommandOptions)}
}
