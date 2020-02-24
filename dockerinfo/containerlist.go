package dockerinfo

import (
	"fmt"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"github.com/juju/loggo"
	"golang.org/x/net/context"
)

var listContainerLogger = loggo.GetLogger("listContainer")

func ListContainer(callContext context.Context) ([]types.Container, error) {
	cli, err := dockerclient.NewEnvClient()
	if err != nil {
		listContainerLogger.Errorf("Error while creating docker Client, Error : %v", err)
		return nil, err
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		listContainerLogger.Errorf("Error while querying for container list, Error : %v", err)
		return nil, err
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}

	return containers, nil

}

// func main() {

// 	containers, err := listContainer(context.Background())
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, container := range containers {
// 		fmt.Printf("%v\n", container)
// 	}
// }
