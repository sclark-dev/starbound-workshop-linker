package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"

	"github.com/nmrshll/go-cp"
	"github.com/urfave/cli/v2"
)

type Mod struct {
	ID   string
	Slug string
	Path string
}

type WorkshopResponse struct {
	Response struct {
		Result               int `json:"result"`
		Resultcount          int `json:"resultcount"`
		Publishedfiledetails []struct {
			Publishedfileid       string `json:"publishedfileid"`
			Result                int    `json:"result"`
			Creator               string `json:"creator"`
			CreatorAppID          int    `json:"creator_app_id"`
			ConsumerAppID         int    `json:"consumer_app_id"`
			Filename              string `json:"filename"`
			FileSize              string `json:"file_size"`
			FileURL               string `json:"file_url"`
			HcontentFile          string `json:"hcontent_file"`
			PreviewURL            string `json:"preview_url"`
			HcontentPreview       string `json:"hcontent_preview"`
			Title                 string `json:"title"`
			Description           string `json:"description"`
			TimeCreated           int    `json:"time_created"`
			TimeUpdated           int    `json:"time_updated"`
			Visibility            int    `json:"visibility"`
			Banned                int    `json:"banned"`
			BanReason             string `json:"ban_reason"`
			Subscriptions         int    `json:"subscriptions"`
			Favorited             int    `json:"favorited"`
			LifetimeSubscriptions int    `json:"lifetime_subscriptions"`
			LifetimeFavorited     int    `json:"lifetime_favorited"`
			Views                 int    `json:"views"`
			Tags                  []struct {
				Tag string `json:"tag"`
			} `json:"tags"`
		} `json:"publishedfiledetails"`
	} `json:"response"`
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
			Email: "sclark@jadesoftware.net",
		},
	}
	App.Version = "1.1.0"
	App.Commands = []*cli.Command{
		{
			Name:    "symlink",
			Aliases: []string{"s"},
			Usage:   "Links workshop files using symlinks. For running on the same machine as you.",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "workshop", Aliases: []string{"w"}, Usage: "Provides the path to look for mods in. Should end in workshop/content/211820", Required: true},
				&cli.StringFlag{Name: "server", Aliases: []string{"s"}, Usage: "Provides the path to place mods in. Should end in /mods", Required: true},
				&cli.BoolFlag{Name: "api", Aliases: []string{"a"}, Usage: "Use the Steam Workshop API to fetch and include the mod's name."},
			},
			Action: func(ctx *cli.Context) error {
				paks, err := getPaks(ctx.String("workshop"), ctx.Bool("api"))
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
				&cli.BoolFlag{Name: "api", Aliases: []string{"a"}, Usage: "Use the Steam Workshop API to fetch and include the mod's name."},
			},
			Action: func(ctx *cli.Context) error {
				paks, err := getPaks(ctx.String("workshop"), ctx.Bool("api"))
				if err != nil {
					return err
				}

				paths, err := copyPaks(ctx.String("server"), paks)
				if err != nil {
					return err
				}

				fmt.Printf("Successfully copied %d mods.\n", len(paths))

				return nil
			},
		},
		{
			Name:    "unlink",
			Aliases: []string{"u"},
			Usage:   "Removed all currently linked mods.",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "server", Aliases: []string{"s"}, Usage: "Provides the path to remove mods from. Should end in /mods", Required: true},
			},
			Action: func(ctx *cli.Context) error {
				paks, err := getPaks(ctx.String("server"), false)
				if err != nil {
					return err
				}

				for _, pak := range paks {
					if err := os.Remove(pak.Path); err != nil {
						return err
					}
				}

				fmt.Printf("Successfully removed %d mods.\n", len(paks))

				return nil
			},
		},
	}

	if err := App.Run(os.Args); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}

func getPaks(dir string, api bool) ([]Mod, error) {
	var mods []Mod
	if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || filepath.Ext(path) != ".pak" {
			return nil
		}

		mod := Mod{ID: filepath.Base(filepath.Dir(path)), Path: path}
		if api {
			resp, err := fetchWorkshopData(mod.ID)
			if err != nil {
				fmt.Printf("Failed to fetch workshop data for %s: %v\n", mod.ID, err)
			} else {
				mod.Slug = slug.Make(resp.Response.Publishedfiledetails[0].Title)
			}
		}

		mods = append(mods, mod)
		return nil
	}); err != nil {
		return nil, err
	}

	return mods, nil
}

func linkPaks(dir string, paks []Mod) ([]string, error) {
	var returnString []string
	for _, pak := range paks {
		var newPath string
		if pak.Slug == "" {
			newPath = fmt.Sprintf("%s%c%s.pak", dir, os.PathSeparator, pak.ID)
		} else {
			newPath = fmt.Sprintf("%s%c%s_%s.pak", dir, os.PathSeparator, pak.ID, pak.Slug)
		}
		if err := os.Symlink(pak.Path, newPath); err != nil {
			return nil, err
		}
		returnString = append(returnString, newPath)
	}

	return returnString, nil
}

func copyPaks(dir string, paks []Mod) ([]string, error) {
	var returnString []string
	for _, pak := range paks {
		newPath := fmt.Sprintf("%s%c%s.pak", dir, os.PathSeparator, pak.ID)
		err := cp.CopyFile(pak.Path, newPath)
		if err != nil {
			return nil, err
		}
		returnString = append(returnString, newPath)
	}

	return returnString, nil
}

func fetchWorkshopData(ID string) (WorkshopResponse, error) {
	data := url.Values{}
	data.Set("itemcount", "1")
	data.Set("publishedfileids[0]", ID)

	resp, err := http.Post("https://api.steampowered.com/ISteamRemoteStorage/GetPublishedFileDetails/v1/",
		"application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return WorkshopResponse{}, err
	}

	defer resp.Body.Close()

	var workshopResponse WorkshopResponse
	err = json.NewDecoder(resp.Body).Decode(&workshopResponse)
	if err != nil {
		return WorkshopResponse{}, err
	}

	return workshopResponse, nil
}
