package docker

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	docker "github.com/docker/docker/client"
)

var client *docker.Client

func ConnectToDocker() error {
	var err error
	client, err = docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())

	if err != nil {
		if docker.IsErrConnectionFailed(err) {
			return ErrConnectionFailed{err: err}
		}

		if docker.IsErrNotFound(err) {
			return ErrNotFound{err: err}
		}
	}

	_, err = client.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "client version") {
			return ErrVersionMismatch{err: err}
		}

		return err
	}

	return nil
}

func GetContainers(ctx context.Context) ([]container.Summary, error) {
	return client.ContainerList(ctx, container.ListOptions{All: true})
}

func GetImageChecksum(ctx context.Context, imageRef string) (string, error) {
	reader, err := client.ImagePull(ctx, imageRef, image.PullOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}
	defer reader.Close()

	_, err = io.Copy(io.Discard, reader)
	if err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}

	inspect, err := client.ImageInspect(ctx, imageRef)
	if err != nil {
		return "", fmt.Errorf("failed to inspect image: %w", err)
	}

	return inspect.ID, nil
}

func UpgradeContainer(ctx context.Context, containerID string, newImageRef string) error {
	inspect, err := client.ContainerInspect(ctx, containerID)
	if err != nil {
		return fmt.Errorf("failed to inspect container: %w", err)
	}

	if err := client.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	if err := client.ContainerRemove(ctx, containerID, container.RemoveOptions{}); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	createConfig := inspect.Config
	createConfig.Image = newImageRef

	resp, err := client.ContainerCreate(ctx, createConfig, inspect.HostConfig, nil, nil, inspect.Name)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	if err := client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}
