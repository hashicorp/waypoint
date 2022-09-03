package hclgen

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"

	embedJson "github.com/hashicorp/waypoint/embedJson"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/plugin"
	fmtpkg "github.com/hashicorp/waypoint/pkg/config"
)

type PlugDocs struct {
	PlugDocs []struct {
		PlugSubDocs []struct {
			Field    string `json:"Field"`
			Type     string `json:"Type"`
			Synopsis string `json:"Synopsis"`
			Optional bool   `json:"Optional"`
			Category bool   `json:"Category"`
			EnvVar   string `json:"EnvVar"`
		} `json:"SubFields"`
		Field    string `json:"Field"`
		Type     string `json:"Type"`
		Synopsis string `json:"Synopsis"`
		Optional bool   `json:"Optional"`
		Category bool   `json:"Category"`
		EnvVar   string `json:"EnvVar"`
	} `json:"requiredFields"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type fieldInfo struct {
	contents string
	isParent bool
	children map[string]string
}

func HclGen(ui terminal.UI) bool {
	brackets := 0
	hclFile, err := os.Create("waypoint.hcl")
	var hclFileByte []byte
	if err != nil {
		ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	}
	ui.Output("Initial waypoint.hcl created!", terminal.WithStyle(terminal.SuccessBoldStyle))
	ui.Output("Type \"exit\" at any point to exit the generator")
	ui.Output("Name your project", terminal.WithHeaderStyle())
	projName, err, close := getName("project", ui)
	if err != nil {
		ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	} else if close {
		exitSafe(hclFile, brackets, ui, hclFileByte)
		return false
	}
	hclFileByte = append(hclFileByte, []byte(fmt.Sprintf("project = \"%s\"\n", projName))...)
	ui.Output("Name your app", terminal.WithHeaderStyle())
	appName, err, close := getName("app", ui)
	if err != nil {
		ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	} else if close {
		exitSafe(hclFile, brackets, ui, hclFileByte)
		return false
	}
	hclFileByte = append(hclFileByte, []byte(fmt.Sprintf("app \"%s\" {\n", appName))...)
	brackets++

	var pluginNames []string
	for pluginName := range plugin.Builtins {
		pluginNames = append(pluginNames, pluginName)
	}
	var fList []string
	dirList, _ := embedJson.Files.ReadDir("gen")
	for _, dirE := range dirList {
		fList = append(fList, dirE.Name())
	}

	ui.Output(
		"Choose build, registry, deployment platform, and releaser plugins",
		terminal.WithHeaderStyle(),
	)

	ui.Output("Configure builder", terminal.WithHeaderStyle())
	// Select a builder
	plug, close, err := selectPlugin("builder", fList, embedJson.Files, ui)
	if err != nil {
		ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	} else if close {
		exitSafe(hclFile, brackets, ui, hclFileByte)
		return false
	}
	hclFileByte = append(hclFileByte, []byte(fmt.Sprintf(genIndent(brackets)+"build {\n"))...)
	brackets++
	if plug.Name != "" {
		hclFileByte = append(hclFileByte, []byte(fmt.Sprintf(genIndent(brackets)+"use \"%s\" {\n", plug.Name))...)
		brackets++
		fieldMap, close, err := populatePlugins(plug, ui)
		if err != nil {
			ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return false
		} else if close {
			exitSafe(hclFile, brackets, ui, hclFileByte)
			return false
		}
		hclFileByte, err = writeFields(hclFileByte, fieldMap, brackets, ui)
		if err != nil {
			ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return false
		}
		ui.Output(
			"Step complete: builder configuration",
			terminal.WithSuccessStyle(),
		)
	}

	// Here we want to close a bracket so that the registry does not appear in the "use" stanza
	hclFileByte = closeBrackets(hclFileByte, 1, brackets)
	brackets--

	ui.Output("Configure registry", terminal.WithHeaderStyle())
	// Select a registry
	plug, close, err = selectPlugin("registry", fList, embedJson.Files, ui)
	if err != nil {
		ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	} else if close {
		exitSafe(hclFile, brackets, ui, hclFileByte)
		return false
	}
	// A registry stanza will only appear in the file if one is chosen
	if plug.Name != "" {
		hclFileByte = append(hclFileByte, []byte(fmt.Sprintf(genIndent(brackets)+"registry {\n"))...)
		brackets++
		hclFileByte = append(hclFileByte, []byte(fmt.Sprintf(genIndent(brackets)+"use \"%s\" {\n", plug.Name))...)
		brackets++
		fieldMap, close, err := populatePlugins(plug, ui)
		if err != nil {
			ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return false
		} else if close {
			exitSafe(hclFile, brackets, ui, hclFileByte)
			return false
		}
		hclFileByte, err = writeFields(hclFileByte, fieldMap, brackets, ui)
		if err != nil {
			ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		}
		ui.Output(
			"Step complete: registry configuration",
			terminal.WithSuccessStyle(),
		)
	}

	// After the registry stanza we want to close the brackets on the build
	// and registry (if it exists) stanzas
	hclFileByte = closeBrackets(hclFileByte, brackets-1, brackets)
	brackets = 1

	ui.Output("Configure deployment platform", terminal.WithHeaderStyle())
	// Select a deployer
	plug, close, err = selectPlugin("deployment platform", fList, embedJson.Files, ui)
	if err != nil {
		ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	} else if close {
		exitSafe(hclFile, brackets, ui, hclFileByte)
		return false
	}

	// A deployer stanza will only appear in the file if one is chosen
	if plug.Name != "" {
		hclFileByte = append(hclFileByte, []byte(fmt.Sprintf(genIndent(brackets)+"deploy {\n"))...)
		brackets++
		hclFileByte = append(hclFileByte, []byte(fmt.Sprintf(genIndent(brackets)+"use \"%s\" {\n", plug.Name))...)
		brackets++
		fieldMap, close, err := populatePlugins(plug, ui)
		if err != nil {
			ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return false
		} else if close {
			exitSafe(hclFile, brackets, ui, hclFileByte)
			return false
		}
		hclFileByte, err = writeFields(hclFileByte, fieldMap, brackets, ui)
		if err != nil {
			ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		}
		ui.Output(
			"Step complete: deployment platform configuration",
			terminal.WithSuccessStyle(),
		)
	}
	// After the deployer stanza we want to close the brackets on the deployer stanza
	hclFileByte = closeBrackets(hclFileByte, brackets-1, brackets)
	brackets = 1

	ui.Output("Configure releaser", terminal.WithHeaderStyle())
	// Select a releaser
	plug, close, err = selectPlugin("releaser", fList, embedJson.Files, ui)
	if err != nil {
		ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	} else if close {
		exitSafe(hclFile, brackets, ui, hclFileByte)
		return false
	}

	// A releaser stanza will only appear in the file if one is chosen
	if plug.Name != "" {
		hclFileByte = append(hclFileByte, []byte(fmt.Sprintf(genIndent(brackets)+"release {\n"))...)
		brackets++
		hclFileByte = append(hclFileByte, []byte(fmt.Sprintf(genIndent(brackets)+"use \"%s\" {\n", plug.Name))...)
		brackets++
		fieldMap, close, err := populatePlugins(plug, ui)
		if err != nil {
			ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return false
		} else if close {
			exitSafe(hclFile, brackets, ui, hclFileByte)
			return false
		}
		hclFileByte, err = writeFields(hclFileByte, fieldMap, brackets, ui)
		if err != nil {
			ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		}
		ui.Output("Step complete: releaser configuration", terminal.WithSuccessStyle())
	}
	// After the releaser stanza we want to close all the brackets
	hclFileByte = closeBrackets(hclFileByte, brackets, brackets)
	hclFile.Write(hclFileByte)
	hclFile.Close()
	ui.Output("\nAll plugin configuration complete", terminal.WithSuccessStyle())
	ui.Output("\nwaypoint.hcl saved!", terminal.WithStyle(terminal.SuccessBoldStyle))
	ui.Output(
		"\nIf you skipped any steps, open your waypoint.hcl file to add missing plugins or fields before continuing. (See https://www.waypointproject.io/plugins)",
	)
	ui.Output("Otherwise, run \"waypoint init\" again to start using Waypoint!\n")
	ui.Output("Now attempting to format the HCL file:\n")
	out, err := fmtpkg.Format(hclFileByte, "waypoint.hcl")
	if err != nil {
		ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	}
	if err := ioutil.WriteFile("waypoint.hcl", out, 0644); err != nil {
		ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	}
	ui.Output("\nFormatting successful!", terminal.WithSuccessStyle())
	return true
}

func writeFields(byteS []byte, fieldMap map[string]fieldInfo, brackets int, ui terminal.UI) ([]byte, error) {
	for key, elem := range fieldMap {
		if elem.isParent {
			byteS = append(byteS, []byte(fmt.Sprintf(genIndent(brackets)+"%s {\n", key))...)
			brackets++
			for name, cont := range elem.children {
				if cont == "" {
					byteS = append(byteS, []byte(fmt.Sprintf(genIndent(brackets)+"//The following field was skipped during file generation\n"))...)
				}
				byteS = append(byteS, []byte(fmt.Sprintf(genIndent(brackets)+"%s = \"%s\"\n", name, cont))...)
			}
			byteS = closeBrackets(byteS, 1, brackets)
			brackets--
		} else {
			if elem.contents == "" {
				byteS = append(byteS, []byte(fmt.Sprintf(genIndent(brackets)+"//The following field was skipped during file generation\n"))...)
			}
			byteS = append(byteS, []byte(fmt.Sprintf(genIndent(brackets)+"%s = \"%s\"\n", key, elem.contents))...)
		}
	}
	return byteS, nil
}

func exitSafe(file *os.File, outstanding int, ui terminal.UI, byteS []byte) error {
	byteS = closeBrackets(byteS, outstanding, outstanding)
	file.Write(byteS)
	file.Close()
	ui.Output("Generator exited. Any information you added before exiting has been included in your waypoint.hcl file. Edit this file manually before using Waypoint.")
	return nil
}

func closeBrackets(byteS []byte, toClose int, outstanding int) []byte {
	extra := outstanding - toClose
	toPrint := ""
	for i := toClose; i > 0; i-- {
		for k := extra + i; k > 1; k-- {
			toPrint = toPrint + "    "
		}
		toPrint = toPrint + "}\n"
		byteS = append(byteS, []byte(fmt.Sprintf(toPrint))...)
		toPrint = ""
	}
	return byteS
}

func genIndent(outstanding int) string {
	spaces := ""
	for i := outstanding; i > 0; i-- {
		spaces = spaces + "    "
	}
	return spaces
}

// populatePlugins returns a map of field name strings as keys mapped to fieldInfo structs containing the
// user generated value for the field and information about if the field is a subfield and its children if so.
// It also returns a boolean indicating if we want to exit the generator completely and an error if one exists
func populatePlugins(plug PlugDocs, ui terminal.UI) (map[string]fieldInfo, bool, error) {
	m := make(map[string]fieldInfo)
	fCount := 0
	for _, f := range plug.PlugDocs {
		if f.Category {
			for _, sf := range f.PlugSubDocs {
				if !sf.Optional {
					fCount++
				}
			}
		} else {
			fCount++
		}
	}
	if plug.PlugDocs == nil || fCount == 0 {
		ui.Output(
			"There are no required fields for this %s plugin, but there may be optional fields you can add to your .hcl file later. See the Waypoint plugin documentation for more information.",
			plug.Type,
		)
	} else {
		if fCount == 1 {
			ui.Output(
				"Please complete the following %d required field for %s, or hit \"return\" to skip.",
				fCount, plug.Name,
				terminal.WithHeaderStyle(),
			)
		} else {
			ui.Output(
				"Please complete the following %d required fields for %s, or hit \"return\" to skip.",
				fCount, plug.Name,
				terminal.WithHeaderStyle(),
			)
		}
		for _, field := range plug.PlugDocs {
			if field.Category {
				// Subfield handling
				m[field.Field] = fieldInfo{isParent: true, children: make(map[string]string)}
				for _, sfield := range field.PlugSubDocs {
					if !sfield.Optional {
						cont, close, err := populateField(sfield.Field, sfield.Type, ui)
						if err != nil {
							return m, false, err
						} else if close {
							return m, true, nil
						}
						m[field.Field].children[sfield.Field] = cont
					}
				}
			} else {
				cont, close, err := populateField(field.Field, field.Type, ui)
				if err != nil {
					return m, false, err
				} else if close {
					return m, true, nil
				}
				m[field.Field] = fieldInfo{contents: cont, isParent: false}
			}
		}
	}
	return m, false, nil
}

// populateField returns the string value inputted by the user for a field, a boolean indicating if we
// want to exit, and an error if one exists
func populateField(name string, fType string, ui terminal.UI) (string, bool, error) {
	getField := true
	typeString := fType
	if typeString == "" {
		typeString = "No_Type_Specified"
	}
	for getField {
		fieldVal, err := ui.Input(&terminal.Input{
			Prompt: fmt.Sprintf("%s <%s>: ", name, typeString),
			Style:  "",
			Secret: false,
		})
		if err != nil {
			ui.Output(
				"Error getting input: %s",
				clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return "", false, err
		} else if strings.ToLower(fieldVal) == "exit" {
			return "", true, nil
		} else if strings.ToLower(fieldVal) == "" {
			ui.Output(fmt.Sprintf("You have selected to skip the %s field.", name))
			pNameConfirm, err := ui.Input(&terminal.Input{
				Prompt: fmt.Sprintf("Do you really want to skip the %s field? (y/N): ", name),
				Style:  "",
				Secret: false,
			})
			if err != nil {
				ui.Output(
					"Error getting input: %s",
					clierrors.Humanize(err),
					terminal.WithErrorStyle(),
				)
				return "", false, err
			} else if strings.ToLower(pNameConfirm) == "exit" {
				return "", true, nil
			} else if strings.ToLower(pNameConfirm) == "yes" || strings.ToLower(pNameConfirm) == "y" {
				ui.Output("%s skipped\n", strings.Title(name), terminal.WithWarningStyle())
				return "", false, nil
			} else {
				ui.Output("Skip cancelled\n")
			}
		} else {
			// TODO: field input type checking
			ui.Output("You inputted \"%s\"\n", fieldVal)
			fieldConfirm, err := ui.Input(&terminal.Input{
				Prompt: "Is this correct? (y/N): ",
				Style:  "",
				Secret: false,
			})
			if err != nil {
				ui.Output(
					"Error getting input: %s",
					clierrors.Humanize(err),
					terminal.WithErrorStyle(),
				)
				return "", false, err
			} else if strings.ToLower(fieldConfirm) == "exit" {
				return "", true, nil
			} else if strings.ToLower(fieldConfirm) == "yes" || strings.ToLower(fieldConfirm) == "y" {
				ui.Output("%s confirmed\n", strings.Title(name), terminal.WithSuccessStyle())
				return fieldVal, false, nil
			} else {
				ui.Output("%s rejected\n", strings.Title(name))
			}
		}
	}
	return "", true, nil
}

// <plug> indicates the plugin that the user needs to select.
// selectPlugin returns a PlugDocs struct filled with the relevant information for the user selected plugin
// selectPlugin also returns a boolean indicating if the program should exit early and an error if one exists
func selectPlugin(plug string, fList []string, fSystem embed.FS, ui terminal.UI) (PlugDocs, bool, error) {
	var plugDocs PlugDocs
	var plugList []string
	for _, file := range fList {
		if filepath.Ext(file) == ".json" {
			switch plug {
			case "builder":
				if strings.HasPrefix(file, "builder") {
					plugList = append(plugList, file)
				}
			case "deployment platform":
				if strings.HasPrefix(file, "platform") {
					plugList = append(plugList, file)
				}
			case "releaser":
				if strings.HasPrefix(file, "release") {
					plugList = append(plugList, file)
				}
			case "registry":
				if strings.HasPrefix(file, "registry") {
					plugList = append(plugList, file)
				}
			}
		}
	}
	sort.Strings(plugList)
	ui.Output(fmt.Sprintf("Select a %s: learn more at https://www.waypointproject.io/plugins. To use a %s that’s not shown here hit return, then edit the .hcl file after it’s been generated.\n",
		plug, plug))
	jMap := make(map[string]interface{})
	var selList []string
	var nameSelList []string
	count := 1
	for _, f := range plugList {
		byteValue, err := fSystem.ReadFile(fmt.Sprintf("gen/%s", f))
		if err != nil {
			return plugDocs, false, err
		}
		json.Unmarshal(byteValue, &jMap)

		// There is an assumption here that all valid plugins will have a description,
		// we have to unmarshal all the plugins for a given stage to get an accurate name
		// and ensure that they exist
		if _, ok := jMap["description"]; ok {
			ui.Output(fmt.Sprintf("%d: %s", count, jMap["name"]), terminal.WithInfoStyle())
			count++
			selList = append(selList, f)
			nameSelList = append(nameSelList, fmt.Sprintf("%s", jMap["name"]))
		}
		for k := range jMap {
			delete(jMap, k)
		}
	}
	// This generates a newline after the list of plugins
	ui.Output("")
	selFileName := ""
	getSelect := true
	for getSelect {
		num, err := ui.Input(&terminal.Input{
			Prompt: fmt.Sprintf(
				"Please select a plugin by typing its corresponding number or hit \"return\" to skip this step (1-%d): ",
				count-1,
			),
			Style:  "",
			Secret: false,
		})
		if err != nil {
			ui.Output(
				"Error getting input: %s",
				clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return plugDocs, false, err
		} else if strings.ToLower(num) == "exit" {
			return plugDocs, true, nil
		} else if val, err := strconv.Atoi(num); err == nil && (0 < val && val < count) {
			ui.Output(fmt.Sprintf("You have selected the %s plugin.", nameSelList[val-1]))
			pNameConfirm, err := ui.Input(&terminal.Input{
				Prompt: fmt.Sprintf("Is this %s plugin correct? (y/N): ", plug),
				Style:  "",
				Secret: false,
			})
			if err != nil {
				ui.Output(
					"Error getting input: %s",
					clierrors.Humanize(err),
					terminal.WithErrorStyle(),
				)
				return plugDocs, false, err
			} else if strings.ToLower(pNameConfirm) == "exit" {
				return plugDocs, true, nil
			} else if strings.ToLower(pNameConfirm) == "yes" || strings.ToLower(pNameConfirm) == "y" {
				ui.Output("%s plugin confirmed\n", strings.Title(plug), terminal.WithSuccessStyle())
				selFileName = selList[val-1]
				getSelect = false
			} else {
				ui.Output("%s plugin rejected\n", strings.Title(plug))
			}
		} else if num == "" {
			ui.Output(fmt.Sprintf("You have selected to skip the %s stage.", plug))
			pNameConfirm, err := ui.Input(&terminal.Input{
				Prompt: fmt.Sprintf("Do you really want to skip the %s stage? (y/N): ", plug),
				Style:  "",
				Secret: false,
			})
			if err != nil {
				ui.Output(
					"Error getting input: %s",
					clierrors.Humanize(err),
					terminal.WithErrorStyle(),
				)
				return plugDocs, false, err
			} else if strings.ToLower(pNameConfirm) == "exit" {
				return plugDocs, true, nil
			} else if strings.ToLower(pNameConfirm) == "yes" || strings.ToLower(pNameConfirm) == "y" {
				ui.Output("Step complete: %s stage skipped", strings.Title(plug), terminal.WithWarningStyle())
				plugDocs.Name = ""
				return plugDocs, false, nil
			} else {
				ui.Output("Skip cancelled\n")
			}
		} else {
			ui.Output("Please select a numbered entry or type nothing to skip.\n")
		}
	}
	// We again unmarshal the JSON file corresponding to the file the user has selected
	if selFileName != "" {
		byteValue, err := fSystem.ReadFile(fmt.Sprintf("gen/%s", selFileName))
		if err != nil {
			return plugDocs, false, err
		}
		json.Unmarshal(byteValue, &plugDocs)
		if plugDocs.Name != "" {
			ui.Output(fmt.Sprintf("You have selected the %s %s plugin.", plugDocs.Name, plug))

		} else {
			//TODO: better error here, do we need to check again here?
			return plugDocs, true, nil
		}
		return plugDocs, false, nil
	}
	return plugDocs, false, nil
}

// Gets either a project or app name for an HCL file, pa should be either "project" or "app"
func getName(pa string, ui terminal.UI) (string, error, bool) {
	if pa == "project" {
		ui.Output(
			"A project contains your app and typically maps 1:1 to a VCS repository. This name must be unique for your Waypoint server. If you're running in local mode, this must be unique to your machine.\n",
		)
	}
	prompt := ""
	if pa == "project" {
		prompt = "Please enter a project name"
	} else {
		prompt = "Please enter an app name"
	}
	getName := true
	name := ""
	for getName {
		paName, err := ui.Input(&terminal.Input{
			Prompt: prompt + ": ",
			Style:  "",
			Secret: false,
		})
		if err != nil {
			ui.Output(
				"Error getting input: %s",
				clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return "", err, false
		} else if strings.ToLower(paName) == "exit" {
			return "", nil, true
		} else if strings.ToLower(paName) == "" {
			ui.Output(prompt + ".\n")
		} else {
			ui.Output("You inputted \"%s\"\n", paName)
			pNameConfirm, err := ui.Input(&terminal.Input{
				Prompt: fmt.Sprintf("Is this %s name correct? (y/N): ", pa),
				Style:  "",
				Secret: false,
			})
			if err != nil {
				ui.Output(
					"Error getting input: %s",
					clierrors.Humanize(err),
					terminal.WithErrorStyle(),
				)
				return "", err, false
			} else if strings.ToLower(pNameConfirm) == "exit" {
				return "", nil, true
			} else if strings.ToLower(pNameConfirm) == "yes" || strings.ToLower(pNameConfirm) == "y" {
				ui.Output("%s name confirmed", strings.Title(pa), terminal.WithSuccessStyle())
				name = paName
				getName = false
			} else {
				ui.Output("%s name rejected", strings.Title(pa))
			}
		}
	}
	return name, nil, false
}
