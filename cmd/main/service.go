package main

import (
	"AutoDoc/internal/config"
	"AutoDoc/internal/docker"
	"AutoDoc/internal/process"
	"context"
	"fmt"
	"strings"
	"time"
)

func main() {
	config.Load()

	err := docker.ConnectToDocker()
	if err != nil {
		switch err.(type) {
		case docker.ErrConnectionFailed:
			fmt.Println("Connection failed error:", err)
		case docker.ErrNotFound:
			fmt.Println("Not found error:", err)
		case docker.ErrVersionMismatch:
			fmt.Println("Version mismatch error:", err)
		default:
			fmt.Println("Unknown error:", err)
		}
	}

	for {
		containers, err := docker.GetContainers(context.Background())
		if err != nil {
			// TODO: Send notifications? After x attempts?
			fmt.Printf("Failed to get containers:\n%v\n", err)
		}

		for _, container := range containers {
			opts := process.OptsFromLabels(container.Labels)
			if opts.Enable == false {
				continue
			}

			fmt.Printf("Checking %s (%s)\n", strings.Join(container.Names, ", "), container.Image)
			availableChecksum, err := docker.GetImageChecksum(context.Background(), container.Image)
			if err != nil {
				fmt.Printf("Failed to get image checksum for %s:\n%v\n", container.Image, err)
				continue
			}

			if availableChecksum == container.ImageID {
				fmt.Printf("Container %s is up to date\n", strings.Join(container.Names, ", "))
				continue
			}

			fmt.Printf("Container %s needs to be updated\n", strings.Join(container.Names, ", "))
			switch opts.Action {
			case process.ActionNotify:
				fmt.Printf("Sending notification to %s\n", strings.Join(container.Names, ", "))
			case process.ActionUpgrade:
				fmt.Printf("Upgrading container %s\n", strings.Join(container.Names, ", "))
				if err := docker.UpgradeContainer(context.Background(), container.ID, container.Image); err != nil {
					fmt.Printf("Failed to upgrade container %s:\n%v\n", strings.Join(container.Names, ", "), err)
					continue
				}
				fmt.Printf("Upgraded container %s\n", strings.Join(container.Names, ", "))
			}
		}

		time.Sleep(config.Interval)
	}
}
