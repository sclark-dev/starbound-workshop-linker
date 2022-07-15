package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/nmrshll/go-cp"
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
	App.Commands = []*cli.Command{
		{
			Name:    "symlink",
			Aliases: []string{"s"},
			Usage:   "Links workshop files using symlinks. For running on the same machine as you.",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "workshop", Aliases: []string{"w"}, Usage: "Provides the path to look for mods in. Should end in workshop/content/211820", Required: true},
				&cli.StringFlag{Name: "server", Aliases: []string{"s"}, Usage: "Provides the path to place mods in. Should end in /mods", Required: true},
			},
			Action: func(ctx *cli.Context) error {
				paks, err := getPaks(ctx.String("workshop"))
				if err != nil {
					return err
				}

				paths, err := linkPaks(ctx.String("server"), paks)
				if err != nil {
					return err
				}

				fmt.Printf("Successfully linked %d mods.\n", len(paths))

				return nil
			},
		},
		{
			Name:    "copy",
			Aliases: []string{"c"},
			Usage:   "Copies workshop files. For running on a separate dedicated machine.",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "workshop", Aliases: []string{"w"}, Usage: "Provides the path to look for mods in. Should end in workshop/content/211820", Required: true},
				&cli.StringFlag{Name: "server", Aliases: []string{"s"}, Usage: "Provides the path to place mods in. Should end in /mods", Required: true},
			},
			Action: func(ctx *cli.Context) error {
				return nil
			},
		},
	}

	if err := App.Run(os.Args); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}

func getPaks(dir string) ([]Mod, error) {
	returnString := []Mod{}
	if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || filepath.Ext(path) != ".pak" {
			return nil
		}

		returnString = append(returnString, Mod{ID: filepath.Base(filepath.Dir(path)), Path: path})
		return nil
	}); err != nil {
		return nil, err
	}

	return returnString, nil
}

func linkPaks(dir string, paks []Mod) ([]string, error) {
	returnString := []string{}
	for _, pak := range paks {
		newPath := fmt.Sprintf("%s%c%s.pak", dir, os.PathSeparator, pak.ID)
		if err := os.Symlink(pak.Path, newPath); err != nil {
			return nil, err
		}
		returnString = append(returnString, newPath)
	}

	return returnString, nil
}

func copyPaks(dir string, paks []Mod) ([]string, error) {
	returnString := []string{}
	for _, pak := range paks {
		newPath := fmt.Sprintf("%s%c%s.pak", dir, os.PathSeparator, pak.ID)
		cp.CopyFile(pak.Path, newPath)
		returnString = append(returnString, newPath)
	}

	return returnString, nil
}
