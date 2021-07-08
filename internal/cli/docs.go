package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/mitchellh/cli"
)

type DocsCommand struct {
	*baseCommand

	commands map[string]cli.CommandFactory
	aliases  map[string]string
}

func (c *DocsCommand) Run(args []string) int {
	os.MkdirAll("./website/content/commands", 0755)
	os.MkdirAll("./website/content/partials/commands", 0755)

	commands := map[string]string{}

	var keys []string

	for k, fact := range c.commands {
		cmd, err := fact()
		if err != nil {
			c.Log.Error("error creating command", "error", err, "command", k)
			return 1
		}

		if _, ok := cmd.(*helpCommand); ok {
			continue
		}

		err = c.genDocs(k, cmd)
		if err != nil {
			c.Log.Error("error generating docs", "error", err, "command", k)
			return 1
		}

		commands[k] = cmd.Synopsis()
		keys = append(keys, k)
	}

	sort.Strings(keys)

	w, err := os.Create("./website/content/partials/commands/command-list.mdx")
	if err != nil {
		c.Log.Error("error creating index page", "error", err)
		return 1
	}

	defer w.Close()

	return 0
}

type HasFlags interface {
	Flags() *flag.Sets
}

func (c *DocsCommand) genDocs(name string, cmd cli.Command) error {
	if name == "cli-docs" {
		return nil
	}

	fmt.Printf("=> %s\n", name)
	goodName := strings.ReplaceAll(name, " ", "-")
	path := filepath.Join("./website", "content", "commands", goodName) + ".mdx"

	w, err := os.Create(path)
	if err != nil {
		return err
	}

	defer w.Close()

	capital := strings.ToUpper(string(name[0])) + name[1:]

	fmt.Fprintf(w, `---
layout: commands
page_title: "Commands: %s"
sidebar_title: "%s"
description: "%s"
---

`, capital, name, cmd.Synopsis())

	fmt.Fprintf(w, "# Waypoint %s\n\nCommand: `waypoint %s`\n\n%s\n\n", capital, name, cmd.Synopsis())

	descFile := goodName + "_desc.mdx"

	fmt.Fprintf(w, "@include \"commands/%s\"\n\n", descFile)

	err = c.touch("./website/content/partials/commands/" + descFile)
	if err != nil {
		return err
	}

	if hf, ok := cmd.(HasFlags); ok {
		flags := hf.Flags()

		fmt.Fprintf(w, "## Usage\n\nUsage: `waypoint %s [options]`\n", name)

		flags.VisitSets(func(name string, set *flag.Set) {
			// Only print a set if it contains vars
			numVars := 0
			set.VisitVars(func(f *flag.VarFlag) { numVars++ })
			if numVars == 0 {
				return
			}

			fmt.Fprintf(w, "\n#### %s\n\n", name)

			set.VisitVars(func(f *flag.VarFlag) {
				if h, ok := f.Value.(flag.FlagVisibility); ok && h.Hidden() {
					return
				}

				name := f.Name
				if t, ok := f.Value.(flag.FlagExample); ok {
					example := t.Example()
					if example != "" {
						name += "=<" + example + ">"
					}
				}

				if len(f.Aliases) > 0 {
					aliases := strings.Join(f.Aliases, "`, `-")

					fmt.Fprintf(w, "- `-%s` (`-%s`) - %s\n", name, aliases, f.Usage)
				} else {
					fmt.Fprintf(w, "- `-%s` - %s\n", name, f.Usage)
				}
			})
		})
	} else {
		fmt.Printf("  ! has no flags\n")
	}

	moreFile := goodName + "_more.mdx"

	fmt.Fprintf(w, "\n@include \"commands/%s\"\n", moreFile)

	return c.touch("./website/content/partials/commands/" + moreFile)
}

func (c *DocsCommand) touch(name string) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	f.Close()

	return nil
}

func (c *DocsCommand) Help() string {
	return "Generate docs"
}

func (c *DocsCommand) Synopsis() string {
	return "Generate docs"
}
