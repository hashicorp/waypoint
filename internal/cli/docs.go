package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

		// Generate the Usage headers based on the cmd Help text
		helpText := strings.Split(cmd.Help(), "\n")
		usage := helpText[0]

		var optionalAlias string
		if len(helpText) > 1 {
			optionalAlias = helpText[1]
		}

		reUsage := regexp.MustCompile(`waypoint (?P<cmd>.*)$`)
		reAlias := regexp.MustCompile(`Alias: `)

		matches := reUsage.FindStringSubmatch(usage)

		if len(matches) > 0 {
			fmt.Fprintf(w, fmt.Sprintf("## Usage\n\nUsage: `waypoint %s`\n", matches[1]))

			hasAlias := false
			if optionalAlias != "" {
				matchAlias := reAlias.FindStringSubmatch(optionalAlias)
				if len(matchAlias) > 0 {
					hasAlias = true
					aliasMatch := reUsage.FindStringSubmatch(optionalAlias)
					fmt.Fprintf(w, fmt.Sprintf("\nAlias: `waypoint %s`\n", aliasMatch[1]))
				}
			}

			// Don't include flag options, we do this later. We trim it here because
			// most commands include it in the help text.
			reOptions := regexp.MustCompile(` Options:`)
			optionsIndex := 0
			for i, opt := range helpText {
				optMatch := reOptions.FindStringSubmatch(opt)
				if len(optMatch) > 0 {
					optionsIndex = i
					break
				}
			}

			if optionsIndex > 1 {
				// Assume all commands have at least "Global Options"
				startIndex := 1
				helpMsg := ""

				if hasAlias {
					startIndex = 2
				}

				helpMsg = strings.Join(helpText[startIndex:optionsIndex], "\n")

				// Strip any color formatting
				reAsciColor := regexp.MustCompile(ansi)
				helpMsg = reAsciColor.ReplaceAllString(helpMsg, "")

				// Trim any left leading whitespace, if any. We do this because any
				// chunk of text that's indented in markdown will be rendered as a
				// code block rather than a paragraph of text.
				helpMsg = strings.TrimLeft(helpMsg, " ")
				fmt.Fprintf(w, "\n%s", helpMsg)
			}
		} else {
			// Fail over to simple docs gen. These are for top level commands
			// like `waypoint context` that don't work without a subcommand and fail the regex match.
			fmt.Fprintf(w, "## Usage\n\nUsage: `waypoint %s [options]`\n", name)
		}

		// Generate flag options
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

				str := fmt.Sprintf("- `-%s` - %s", name, f.Usage)
				if len(f.Aliases) > 0 {
					aliases := strings.Join(f.Aliases, "`, `-")

					str = fmt.Sprintf("- `-%s` (`-%s`) - %s", name, aliases, f.Usage)
				}

				// Add a period at the end of the doc sentence if the field didn't add
				// one already.
				if !strings.HasSuffix(str, ".") {
					str += "."
				}

				if f.Default != "" {
					str = fmt.Sprintf("%s The default is %s.", str, f.Default)
				}
				str += "\n"

				fmt.Fprint(w, str)
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

const (
	// NOTE: adapted from https://github.com/acarl005/stripansi
	ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
)
