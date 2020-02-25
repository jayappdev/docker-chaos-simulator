package cmdexec

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
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

func (c *ChaosCommand) Execute() (error, io.Reader) {

	err := c.Validate()
	if err != nil {
		ChaosCommandExecutorLogger.Errorf("Validation errors : %s", err)
	}

	cmdString := c.getFullCommand(ChaosCommandOSExecutor)

	fmt.Printf("Command being executed : %s\n", cmdString)
	ChaosCommandExecutorLogger.Infof("Command being executed : %s", cmdString)

	cmdArgs := strings.Fields(cmdString)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)

	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	oneByte := make([]byte, 100)
	num := 1
	for {
		_, err := stdout.Read(oneByte)
		if err != nil {
			fmt.Printf(err.Error())
			break
		}

		r := bufio.NewReader(stdout)
		line, _, _ := r.ReadLine()
		fmt.Println(string(line))
		num = num + 1
		if num > 3 {
			os.Exit(0)
		}
	}

	cmd.Wait()

	return nil, nil
}

// func main() {

// 	loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stdout))

// 	c := ChaosCommand{}
// 	c.operation = "ping"
// 	c.args = []Arg{
// 		ArgInit("arg1", "value1"),
// 		ArgInit("arg2", "value2"),
// 	}

// 	c.Execute()

// 	path, err := exec.LookPath("fortune")
// 	if err != nil {
// 		log.Fatal("installing fortune is in your future")
// 	}
// 	fmt.Printf("fortune is available at %s\n", path)

// }
