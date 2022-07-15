package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

type Mod struct {
	ID   string
	Path string
}

var App *cli.App

func main() {
	App = cli.NewApp()
	App.Name = "Starbound Workshop Linker"
	App.Usage = ""
	App.Description = ""
	App.Authors = []*cli.Author{
		{
			Name:  "Shane Clark",
			Email: "sclark@frostfire.io",
		},
	}
	App.Version = "1.0.0"
	App.Commands = []*cli.Command{}

	if err := App.Run(os.Args); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
