package pkg

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"log"
	"os"
	"runtime"
)

// runAndGetDockerContainerLogStream runs the user container and returns the log stream
func (t *task) runAndGetDockerContainerLogStream() io.ReadCloser {
	c := t.getDockerClientOrDie()
	t.pullDockerImage(c)
	containerID := t.createDockerContainer(c)
	t.startDockerContainer(containerID, c)
	return t.getContainerLogsStream(containerID, c)
}

// getDockerClientOrDie returns the docker client
func (t *task) getDockerClientOrDie() *client.Client {
	var host string
	switch runtime.GOOS {
	case "darwin":
		host = defaultDockerHost
	case "linux":
		log.Println("untested on linux, but it should work on linux too, continuing...")
		host = defaultDockerHost
	default:
		log.Println("Windows (or other) system may have different docker host. Untested GOOS found, exiting.")
		os.Exit(0)
	}

	c, err := client.NewClientWithOpts(client.WithHost(host))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("docker client initialized")
	return c
}

// pullDockerImage pulls users docker image
func (t *task) pullDockerImage(c *client.Client) {
	_, err := c.ImagePull(context.TODO(), t.dockerImage, types.ImagePullOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("docker image pulled successfully or already exists", "dockerImage", t.dockerImage)
}

// createDockerContainer create user's docker container
func (t *task) createDockerContainer(c *client.Client) string {
	containerCreated, err := c.ContainerCreate(context.TODO(), &container.Config{
		Image: t.dockerImage,
		Cmd:   []string{"sh", "-c", t.bashCommand},
	}, &container.HostConfig{}, nil, nil, "")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("container created successfully", "container id", containerCreated.ID, "bash command", t.bashCommand)
	return containerCreated.ID
}

// startDockerContainer start user's docker container
func (t *task) startDockerContainer(containerID string, c *client.Client) {
	err := c.ContainerStart(context.TODO(), containerID, types.ContainerStartOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("container started successfully", "container id", containerID)
}

// getContainerLogsStream returns container log stream
func (t *task) getContainerLogsStream(containerID string, c *client.Client) io.ReadCloser {
	i, err := c.ContainerLogs(context.Background(), containerID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: true,
		Follow:     true,
		Tail:       "40",
	})

	if err != nil {
		log.Fatal(err)
	}
	log.Println("fetching container logs successfully", "container id", containerID)
	return i
}
