package main

import (
	"github.com/codegangsta/envy/lib"
	"github.com/fsouza/go-dockerclient"
	"log"
)

func main() {
	host := envy.MustGet("DOCKER_HOST")

	c, err := docker.NewClient(host)
	if err != nil {
		log.Fatal(err)
	}

	containers, err := c.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, container := range containers {

		con, err := c.InspectContainer(container.ID)
		if err != nil {
			log.Fatal(err)
		}

		// copy the config
		cfg := con.Config
		name := con.Name

		cfg.Env = append(cfg.Env, "FOO=bar")
		log.Println(cfg.Env)

		log.Println("Removing", name)
		err = c.RemoveContainer(docker.RemoveContainerOptions{
			ID:    container.ID,
			Force: true,
		})
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Creating new instance of", name)
		// Create the new one with a new env
		con, err = c.CreateContainer(docker.CreateContainerOptions{
			Name:   name,
			Config: cfg,
		})
		if err != nil {
			log.Fatal(err)
		}

		err = c.RestartContainer(con.ID, 10)
		if err != nil {
			log.Fatal(err)
		}
	}
}
