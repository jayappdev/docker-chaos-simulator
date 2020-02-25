package cmdexec

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"github.com/juju/loggo"
)

type Arg struct {
	Parameter string
	Value     string
}

type ChaosCommand struct {
	operation string
	container string
	args      []Arg
}

type ChaosCommandValidator interface {
	Validate() error
}

type ChaosCommandExecutor interface {
	Execute() (error, io.Reader)
}

const ChaosCommandOSExecutor = "pumba --log-level info"

var ChaosCommandExecutorLogger = loggo.GetLogger("listContainer")

func ArgInit(parameter1 string, value1 string) Arg {
	return Arg{
		parameter1,
		value1,
	}
}

func (a Arg) String() string {
	fmt.Println("Argument - ", a.Parameter, a.Value)
	return strings.Join([]string{a.Parameter, a.Value}, " ")
}

func createChaosWithOperationOnly(container, operation string) *ChaosCommand {
	return &ChaosCommand{
		operation: operation,
		container: container,
	}
}

func convertToMap(args []Arg) map[string]Arg {
	argLookup := make(map[string]Arg)

	for _, arg := range args {
		argLookup[arg.Parameter] = arg
	}

	return argLookup
}

func extractParameters(args []Arg, parameters ...string) []Arg {
	argLookup := convertToMap(args)

	var returnArgs []Arg

	for _, parameter := range parameters {
		if val, ok := argLookup[parameter]; ok {
			returnArgs = append(returnArgs, val)
		}
	}

	return returnArgs
}

func excludeParameters(args []Arg, parameters ...string) []Arg {
	argLookup := convertToMap(args)

	var returnArgs []Arg

	for _, parameter := range parameters {
		delete(argLookup, parameter)
	}

	for _, v := range argLookup {
		returnArgs = append(returnArgs, v)
	}

	return returnArgs
}

func FromJson(jsonStream io.Reader) []Arg {
	var args []Arg

	dec := json.NewDecoder(jsonStream)
	for {
		var m Arg
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		args = append(args, m)
	}

	return args
}

func createChaosWithOperationAndArgsJson(container, operation string, jsonStream io.Reader) *ChaosCommand {
	c := createChaosWithOperationOnly(container, operation)
	c.args = append(c.args, FromJson(jsonStream)...)
	return c
}

func createChaosWithOperationAndArgs(container, operation string, args []Arg) *ChaosCommand {
	c := createChaosWithOperationOnly(container, operation)
	c.args = append(c.args, args...)
	return c
}

func concatenateArgs(args ...Arg) string {
	var commandsToExecute []string

	for _, arg := range args {
		commandsToExecute = append(commandsToExecute, arg.String())
	}

	return strings.Join(commandsToExecute, " ")
}

func (c ChaosCommand) Validate() error {
	return nil
}

func (c *ChaosCommand) getFullCommand(chaosCommand string) string {
	return strings.Join([]string{chaosCommand, c.operation, concatenateArgs(c.args...), c.container}, " ")
}

func (c *ChaosCommand) Execute() (error, io.Reader, io.Reader) {

	err := c.Validate()
	if err != nil {
		ChaosCommandExecutorLogger.Errorf("Validation errors : %s", err)
	}

	cmdString := c.getFullCommand(ChaosCommandOSExecutor)

	fmt.Printf("Command being executed : %s\n", cmdString)
	ChaosCommandExecutorLogger.Infof("Command being executed : %s", cmdString)

	cmdArgs := strings.Fields(cmdString)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		ChaosCommandExecutorLogger.Errorf("Error while executing output pipe %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		ChaosCommandExecutorLogger.Errorf("Error while executing Error pipe %v", err)
	}

	err = cmd.Start()
	if err != nil {
		ChaosCommandExecutorLogger.Errorf("Error while starting command %v", err)
	}

	lines, _ := ioutil.ReadAll(stdout)
	fmt.Println("Output : " + string(lines))

	errorLines, _ := ioutil.ReadAll(stderr)
	fmt.Println("Error : " + string(errorLines))

	cmd.Wait()
	return nil, strings.NewReader(string(lines)), strings.NewReader(string(errorLines))
}
