package docker

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerCommandRunner struct {
	client *client.Client
}

func NewDockerCommandRunner() *DockerCommandRunner {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	log.Println("Docker client is created")
	return &DockerCommandRunner{
		client: cli,
	}
}

func (d *DockerCommandRunner) ExecuteCommand(ctx context.Context, containerName string, cmd ...string) ([]byte, error) {
	log.Println("Running cmd:", cmd)
	execConfig := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}

	resp, err := d.client.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		return nil, fmt.Errorf("can't create a task: %w", err)
	}

	execResponse, err := d.client.ContainerExecAttach(ctx, resp.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		return nil, fmt.Errorf("can't attach to container: %w", err)
	}
	defer execResponse.Close()

	output, err := ioutil.ReadAll(execResponse.Reader)
	if err != nil {
		return nil, fmt.Errorf("can't read the output: %w", err)
	}

	log.Println("Executing was successful!")
	return output, nil
}
