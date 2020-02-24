package cmdexec

import "io"

const CONTAINER_PAUSER = "pause"

type ChaosPauseSimulator struct {
	*ChaosCommand
}

func CreateChaosPauseSimulator(container string, jsonStream io.Reader) ChaosPauseSimulator {
	return ChaosPauseSimulator{createChaosWithOperationAndArgs(container, CONTAINER_PAUSER, jsonStream)}
}
