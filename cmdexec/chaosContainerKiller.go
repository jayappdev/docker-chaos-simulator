package cmdexec

import "io"

const CONTAINER_KILLER = "kill"

type ChaosKillerSimulator struct {
	*ChaosCommand
}

func CreateChaosKillSimulator(container string, jsonStream io.Reader) ChaosKillerSimulator {
	return ChaosKillerSimulator{createChaosWithOperationAndArgsJson(container, CONTAINER_KILLER, jsonStream)}
}
