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
	App.Commands = []*cli.Command{}

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
