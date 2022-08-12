package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-getter"
	"github.com/pkg/errors"
	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"

	"github.com/hashicorp/waypoint/internal/cli/datagen"
	clientpkg "github.com/hashicorp/waypoint/internal/client"
	"github.com/hashicorp/waypoint/internal/clierrors"
	configpkg "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/datasource"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

type InitCommand struct {
	*baseCommand

	fromProject string
	into        string
	update      bool
	from        string

	project *clientpkg.Project
	cfg     *configpkg.Config
}

func (c *InitCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoConfig(),
		WithNoClient(),
	); err != nil {
		return 1
	}

	if c.fromProject != "" {
		if c.into == "" {
			if u, err := url.Parse(c.fromProject); err == nil {
				c.into = filepath.Base(u.Path)
			} else {
				c.into = filepath.Base(c.fromProject)
			}

			ext := filepath.Ext(c.into)
			if ext != "" {
				c.into = c.into[:len(c.into)-len(ext)]
			}
		}

		var dir string

		if filepath.IsAbs(c.into) {
			dir = c.into
		} else {
			dir = "./" + c.into
		}

		if _, err := os.Stat(dir); err == nil {
			c.ui.Output("Cannot perform a remote initialization", terminal.WithStyle(terminal.ErrorBoldStyle))
			c.ui.Output("")
			c.ui.Output(
				"Waypoint has detected an existing directory '"+dir+"' and will not \n"+
					"overwrite an existing application with a remote one.",
				terminal.WithErrorStyle(),
			)

			return 1
		}

		c.ui.Output("Initializing local application from remote location:")
		c.ui.NamedValues([]terminal.NamedValue{
			{
				Name:  "Location",
				Value: c.fromProject,
			},
			{
				Name:  "Directory",
				Value: dir,
			},
		}, terminal.WithInfoStyle())

		pwd, err := os.Getwd()
		if err != nil {
			c.ui.Output("")
			c.ui.Output("Project had errors during unpacking.", terminal.WithStyle(terminal.ErrorBoldStyle))
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())

			return 1
		}

		client := &getter.Client{
			Src: c.fromProject,
			Dst: dir,
			Pwd: pwd,
			Dir: true,
			Getters: map[string]getter.Getter{
				"http": &getter.HttpGetter{
					Netrc:            false,
					HeadFirstTimeout: 10 * time.Second,
					ReadTimeout:      30 * time.Second,
					MaxBytes:         500000000, // 500 MB
				},
				"file": &getter.FileGetter{
					Copy: true,
				},
				"git": &getter.GitGetter{
					Timeout: 30 * time.Second,
				},
			},
			DisableSymlinks: true,
		}

		err = client.Get()
		if err != nil {
			c.ui.Output("")
			c.ui.Output("Project had errors during unpacking.", terminal.WithStyle(terminal.ErrorBoldStyle))
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())

			return 1
		}

		err = os.Chdir(dir)
		if err != nil {
			c.ui.Output("")
			c.ui.Output("Project had errors during unpacking.", terminal.WithStyle(terminal.ErrorBoldStyle))
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())

			return 1
		}

		c.ui.Output("Project fetched into '%s'", dir)
		return 0
	}

	path, err := c.initConfigPath(c.fromProject)
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return 1
	}

	// If we have no config, initialize a new one.
	if path == "" {
		proceed, err := c.ui.Input(&terminal.Input{
			Prompt: "Do you want help generating a waypoint.hcl file? Type 'yes' to initialize the interactive generator or 'no' to generate a template waypoint.hcl file: ",
			Style:  "",
			Secret: false,
		})
		if err != nil {
			c.ui.Output(
				"Error getting input: %s",
				clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
		} else if strings.ToLower(proceed) == "yes" || strings.ToLower(proceed) == "y" {
			c.ui.Output("Starting interactive .hcl generator.\n")
			if !c.hclGen() {
				return 1
			}
		} else if strings.ToLower(proceed) == "no" {
			c.ui.Output("Generating template file.\n")
			if !c.initNew() {
				return 1
			}
		} else {
			c.ui.Output("Input did not match any option, generating template file\n")
			if !c.initNew() {
				return 1
			}
		}

		return 0
	}

	// Steps to run
	steps := []func() bool{
		c.validateConfig,
		c.validateServer,
		c.validateProject,
		// NOTE(mitchellh): this is disabled as of 0.6 since we can't load
		// plugins in the CLI safely due to the usage of input variables +
		// remote runners. This will be fixed in the future by migrating the
		// init task to the InitOp remote operation. We're keeping this code
		// around so we can migrate it then.
		// c.validatePlugins,

		// NOTE(mitchellh): this is disabled as of 0.2 since we can't load
		// config anymore. We're keeping the code around so that we can migrate
		// it in the future.
		// c.validateAuth,
	}
	for _, step := range steps {
		if !step() {
			c.ui.Output("Project had errors during initialization.\n"+
				"Waypoint experienced some errors during project initialization. The output\n"+
				"above should contain the failure messages. Please correct these errors and\n"+
				"run 'waypoint init' again.",
				terminal.WithStyle(terminal.ErrorBoldStyle),
			)

			return 1
		}
	}

	c.ui.Output("")
	c.ui.Output("Project initialized!", terminal.WithStyle(terminal.SuccessBoldStyle))
	c.ui.Output("")
	c.ui.Output(
		"You may now call 'waypoint up' to deploy your project or\n" +
			"commands such as 'waypoint build' to perform steps individually.",
	)

	return 0
}

func (c *InitCommand) initNew() bool {
	data, err := datagen.Asset("init.tpl.hcl")
	if err != nil {
		// Should never happen because it is embedded.
		panic(err)
	}

	if err := ioutil.WriteFile(configpkg.Filename, data, 0644); err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	}

	c.ui.Output("Initial Waypoint configuration created!", terminal.WithStyle(terminal.SuccessBoldStyle))
	c.ui.Output(strings.TrimSpace(`
No Waypoint configuration was found in this directory.

A sample configuration has been created in the file "waypoint.hcl". This
file is heavily commented to help you get started.

Once you've setup your initial configuration, run "waypoint init" again to
validate the configuration and initialize your project.
`),
		terminal.WithSuccessStyle(),
	)

	return true
}

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

func (c *InitCommand) hclGen() bool {
	brackets := 0
	hclFile, err := os.Create("waypoint.hcl")
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	}
	defer hclFile.Close()
	c.ui.Output("Initial waypoint.hcl created!", terminal.WithStyle(terminal.SuccessBoldStyle))
	c.ui.Output("Type \"exit\" at any point to exit the generator")
	c.ui.Output("Name your project", terminal.WithHeaderStyle())
	projName, err, close := c.getName("project")
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	} else if close == true {
		c.exitSafe(hclFile, brackets)
		return false
	}
	hclFile.Write([]byte(fmt.Sprintf("project = \"%s\"\n", projName)))
	c.ui.Output("Name your app", terminal.WithHeaderStyle())
	appName, err, close := c.getName("app")
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	} else if close == true {
		c.exitSafe(hclFile, brackets)
		return false
	}
	hclFile.Write([]byte(fmt.Sprintf("app \"%s\" {\n", appName)))
	brackets++

	// TODO: this is a placeholder, the real implementation will use the JSON files as they are included in the waypoint binary
	// Not a final implemenation so hardcoded with a relative path
	fPath := "./docs/gen"
	file, err := os.Open(fPath)
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	}
	defer file.Close()
	fList, err := file.Readdirnames(0)
	//TODO: replace all above file code with this line
	//fList, err := fs.Glob(embedJson.files, "*.json")
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	}

	c.ui.Output("Choose build, registry, deployment platform, and releaser plugins", terminal.WithHeaderStyle())

	c.ui.Output("Configure builder", terminal.WithHeaderStyle())
	// Select a builder
	plug, err, close := c.selectPlugin(1, fList, fPath)
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	} else if close == true {
		c.exitSafe(hclFile, brackets)
		return false
	}
	hclFile.Write([]byte(fmt.Sprintf(c.genIndent(brackets) + "build {\n")))
	brackets++
	if plug.Name != "" {
		hclFile.Write([]byte(fmt.Sprintf(c.genIndent(brackets)+"use \"%s\" {\n", plug.Name)))
		brackets++
		fieldMap, err, close := c.populatePlugins(plug)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return false
		} else if close == true {
			c.exitSafe(hclFile, brackets)
			return false
		}
		for key, elem := range fieldMap {
			hclFile.Write([]byte(fmt.Sprintf(c.genIndent(brackets)+"%s = \"%s\"\n", key, elem)))
		}
		c.ui.Output("Step complete: builder configuration complete", terminal.WithSuccessStyle())
	}

	c.ui.Output("Configure registry", terminal.WithHeaderStyle())
	// Select a registry
	plug, err, close = c.selectPlugin(4, fList, fPath)
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	} else if close == true {
		c.exitSafe(hclFile, brackets)
		return false
	}
	// A registry stanza will only appear in the file if one is chosen
	if plug.Name != "" {
		hclFile.Write([]byte(fmt.Sprintf(c.genIndent(brackets) + "registry {\n")))
		brackets++
		hclFile.Write([]byte(fmt.Sprintf(c.genIndent(brackets)+"use \"%s\" {\n", plug.Name)))
		brackets++
		fieldMap, err, close := c.populatePlugins(plug)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return false
		} else if close == true {
			c.exitSafe(hclFile, brackets)
			return false
		}
		for key, elem := range fieldMap {
			hclFile.Write([]byte(fmt.Sprintf(c.genIndent(brackets)+"%s = \"%s\"\n", key, elem)))
		}
		c.ui.Output("Step complete: registry configuration complete", terminal.WithSuccessStyle())
	}

	// After the registry stanza we want to close the brackets on the build and registry (if it exists) stanzas
	err = c.closeBrackets(hclFile, brackets-1, brackets)
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	}
	brackets = 1

	c.ui.Output("Configure deployment platform", terminal.WithHeaderStyle())
	// Select a deployer
	plug, err, close = c.selectPlugin(2, fList, fPath)
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	} else if close == true {
		c.exitSafe(hclFile, brackets)
		return false
	}

	// A deployer stanza will only appear in the file if one is chosen
	if plug.Name != "" {
		hclFile.Write([]byte(fmt.Sprintf(c.genIndent(brackets) + "deploy {\n")))
		brackets++
		hclFile.Write([]byte(fmt.Sprintf(c.genIndent(brackets)+"use \"%s\" {\n", plug.Name)))
		brackets++
		fieldMap, err, close := c.populatePlugins(plug)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return false
		} else if close == true {
			c.exitSafe(hclFile, brackets)
			return false
		}
		for key, elem := range fieldMap {
			hclFile.Write([]byte(fmt.Sprintf(c.genIndent(brackets)+"%s = \"%s\"\n", key, elem)))
		}
		c.ui.Output("Step complete: deployment platform configuration complete", terminal.WithSuccessStyle())
	}
	// After the deployer stanza we want to close the brackets on the deployer stanza
	err = c.closeBrackets(hclFile, brackets-1, brackets)
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	}
	brackets = 1

	c.ui.Output("Configure releaser", terminal.WithHeaderStyle())
	// Select a releaser
	plug, err, close = c.selectPlugin(3, fList, fPath)
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	} else if close == true {
		c.exitSafe(hclFile, brackets)
		return false
	}

	// A releaser stanza will only appear in the file if one is chosen
	if plug.Name != "" {
		hclFile.Write([]byte(fmt.Sprintf(c.genIndent(brackets) + "release {\n")))
		brackets++
		hclFile.Write([]byte(fmt.Sprintf(c.genIndent(brackets)+"use \"%s\" {\n", plug.Name)))
		brackets++
		fieldMap, err, close := c.populatePlugins(plug)
		if err != nil {
			c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
			return false
		} else if close == true {
			c.exitSafe(hclFile, brackets)
			return false
		}
		for key, elem := range fieldMap {
			hclFile.Write([]byte(fmt.Sprintf(c.genIndent(brackets)+"%s = \"%s\"\n", key, elem)))
		}
		c.ui.Output("Step complete: releaser configuration complete", terminal.WithSuccessStyle())
	}
	// After the releaser stanza we want to close all the brackets
	err = c.closeBrackets(hclFile, brackets, brackets)
	if err != nil {
		c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
		return false
	}
	c.ui.Output("All plugin configuration complete", terminal.WithSuccessStyle())
	c.ui.Output("waypoint.hcl saved!", terminal.WithStyle(terminal.SuccessBoldStyle))
	c.ui.Output("If you skipped any steps, open your waypoint.hcl file to add missing plugins or fields before continuing. (See https://www.waypointproject.io/plugins)")
	c.ui.Output("Otherwise, run \"waypoint init\" again to start using Waypoint!")
	return true
}

func (c *InitCommand) exitSafe(file *os.File, outstanding int) error {
	c.closeBrackets(file, outstanding, outstanding)
	c.ui.Output("Generator exited")
	return nil
}

func (c *InitCommand) closeBrackets(file *os.File, toClose int, outstanding int) error {
	extra := outstanding - toClose
	toPrint := ""
	for i := toClose; i > 0; i-- {
		for k := extra + i; k > 1; k-- {
			toPrint = toPrint + "    "
		}
		toPrint = toPrint + "}\n"
		file.Write([]byte(fmt.Sprintf(toPrint)))
		toPrint = ""
	}
	return nil
}

func (c *InitCommand) genIndent(outstanding int) string {
	spaces := ""
	for i := outstanding; i > 0; i-- {
		spaces = spaces + "    "
	}
	return spaces
}

func (c *InitCommand) populatePlugins(plug PlugDocs) (map[string]string, error, bool) {
	m := make(map[string]string)
	if plug.PlugDocs == nil {
		c.ui.Output("There are no required fields for this %s plugin, but there may be optional fields you can add to your .hcl file later. See the Waypoint plugin documentation for more information.", plug.Type)
	} else {
		fCount := 0
		for _, f := range plug.PlugDocs {
			if f.Category == true {
				for _, sf := range f.PlugSubDocs {
					if sf.Optional == false {
						fCount++
					}
				}
			} else {
				fCount++
			}
		}
		if fCount == 1 {
			c.ui.Output("Please complete the following %d required field for %s, or hit \"return\" to skip.", fCount, plug.Name, terminal.WithHeaderStyle())
		} else {
			c.ui.Output("Please complete the following %d required fields for %s, or hit \"return\" to skip.", fCount, plug.Name, terminal.WithHeaderStyle())
		}
		fCount = 0
		for _, field := range plug.PlugDocs {
			if field.Category == true {
				// Subfield handling
				for _, sfield := range field.PlugSubDocs {
					if sfield.Optional == false {
						cont, err, close := c.populateField(sfield.Field, sfield.Type, fCount)
						fCount++
						if err != nil {
							return m, err, false
						} else if close == true {
							return m, nil, true
						}
						m[sfield.Field] = cont
					}
				}

			} else {
				cont, err, close := c.populateField(field.Field, field.Type, fCount)
				fCount++
				if err != nil {
					return m, err, false
				} else if close == true {
					return m, nil, true
				}
				m[field.Field] = cont
			}
		}
	}
	return m, nil, false
}

func (c *InitCommand) populateField(name string, fType string, count int) (string, error, bool) {
	getField := true
	typeString := fType
	if typeString == "" {
		typeString = "No_Type_Specified"
	}
	for getField {
		fieldVal, err := c.ui.Input(&terminal.Input{
			Prompt: fmt.Sprintf("%s <%s>: ", strings.Title(name), typeString),
			Style:  "",
			Secret: false,
		})
		if err != nil {
			c.ui.Output(
				"Error getting input: %s",
				clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return "", err, false
		} else if strings.ToLower(fieldVal) == "exit" {
			return "", nil, true
		} else if strings.ToLower(fieldVal) == "" {
			c.ui.Output(fmt.Sprintf("You have selected to skip the %s field.", name))
			pNameConfirm, err := c.ui.Input(&terminal.Input{
				Prompt: fmt.Sprintf("Do you really want to skip the %s field? (y/N): ", name),
				Style:  "",
				Secret: false,
			})
			if err != nil {
				c.ui.Output(
					"Error getting input: %s",
					clierrors.Humanize(err),
					terminal.WithErrorStyle(),
				)
				return "", err, false
			} else if strings.ToLower(pNameConfirm) == "exit" {
				return "", nil, true
			} else if strings.ToLower(pNameConfirm) == "yes" || strings.ToLower(pNameConfirm) == "y" {
				c.ui.Output("%s skipped\n", strings.Title(name), terminal.WithWarningStyle())
				getField = false
			} else {
				c.ui.Output("Skip cancelled\n")
			}
		} else {
			// TODO: field input type checking
			c.ui.Output("You inputted \"%s\"\n", fieldVal)
			fieldConfirm, err := c.ui.Input(&terminal.Input{
				Prompt: fmt.Sprintf("Is this correct? (y/N): "),
				Style:  "",
				Secret: false,
			})
			if err != nil {
				c.ui.Output(
					"Error getting input: %s",
					clierrors.Humanize(err),
					terminal.WithErrorStyle(),
				)
				return "", err, false
			} else if strings.ToLower(fieldConfirm) == "exit" {
				return "", nil, true
			} else if strings.ToLower(fieldConfirm) == "yes" || strings.ToLower(fieldConfirm) == "y" {
				c.ui.Output("%s confirmed\n", strings.Title(name), terminal.WithSuccessStyle())
				return fieldVal, nil, false
			} else {
				c.ui.Output("%s rejected\n", strings.Title(name))
			}
		}
	}
	return "", nil, true
}

// plug indicates the plugin that the user needs to select. 1: Builder, 2: Deployer/Platform, 3: Releaser, 4: Registry
func (c *InitCommand) selectPlugin(plug int, fList []string, fPath string) (PlugDocs, error, bool) {
	var plugType string
	var plugDocs PlugDocs
	switch plug {
	case 1:
		plugType = "builder"
	case 2:
		plugType = "deployment platform"
	case 3:
		plugType = "releaser"
	case 4:
		plugType = "registry"
	}
	var plugList []string
	for _, file := range fList {
		if filepath.Ext(file) == ".json" {
			switch plug {
			case 1:
				if strings.HasPrefix(file, "builder") {
					plugList = append(plugList, file)
				}
			case 2:
				if strings.HasPrefix(file, "platform") {
					plugList = append(plugList, file)
				}
			case 3:
				if strings.HasPrefix(file, "release") {
					plugList = append(plugList, file)
				}
			case 4:
				if strings.HasPrefix(file, "registry") {
					plugList = append(plugList, file)
				}
			}
		}
	}
	sort.Strings(plugList)
	c.ui.Output(fmt.Sprintf("Select a %s: learn more at https://www.waypointproject.io/plugins. To use a %s that’s not shown here enter nothing, then edit the .hcl file after it’s been generated.\n", plugType, plugType))
	jMap := make(map[string]interface{})
	var selList []string
	var nameSelList []string
	count := 1
	for _, f := range plugList {
		jsonFile, err := os.Open(fmt.Sprintf("%s/%s", fPath, f))
		if err != nil {
			return plugDocs, err, false
		}
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &jMap)

		//TODO: REMOVE TEST CODE
		json.Unmarshal(byteValue, &plugDocs)

		// There is an assumption here that all plugins will have a description, we have to unmarshal all the plugins
		// for a given plugin to get an accurate name and ensure that they exist
		if _, ok := jMap["description"]; ok {
			c.ui.Output(fmt.Sprintf("%d: %s", count, jMap["name"]), terminal.WithInfoStyle())
			count++
			selList = append(selList, f)
			nameSelList = append(nameSelList, fmt.Sprintf("%s", jMap["name"]))
		}
		for k := range jMap {
			delete(jMap, k)
		}
	}
	// This generates a newline after the list of plugins
	c.ui.Output("")
	selFileName := ""
	getSelect := true
	for getSelect {
		num, err := c.ui.Input(&terminal.Input{
			Prompt: fmt.Sprintf("Please select a plugin by typing its corresponding number or hit \"return\" to skip this step (1-%d): ", count-1),
			Style:  "",
			Secret: false,
		})
		if err != nil {
			c.ui.Output(
				"Error getting input: %s",
				clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return plugDocs, err, false
		} else if strings.ToLower(num) == "exit" {
			return plugDocs, nil, true
		} else if val, err := strconv.Atoi(num); err == nil && (0 < val && val < count) {
			c.ui.Output(fmt.Sprintf("You have selected the %s plugin.", nameSelList[val-1]))
			pNameConfirm, err := c.ui.Input(&terminal.Input{
				Prompt: fmt.Sprintf("Is this %s plugin correct? (y/N): ", plugType),
				Style:  "",
				Secret: false,
			})
			if err != nil {
				c.ui.Output(
					"Error getting input: %s",
					clierrors.Humanize(err),
					terminal.WithErrorStyle(),
				)
				return plugDocs, err, false
			} else if strings.ToLower(pNameConfirm) == "exit" {
				return plugDocs, nil, true
			} else if strings.ToLower(pNameConfirm) == "yes" || strings.ToLower(pNameConfirm) == "y" {
				c.ui.Output("%s plugin confirmed\n", strings.Title(plugType))
				selFileName = selList[val-1]
				getSelect = false
			} else {
				c.ui.Output("%s plugin rejected\n", strings.Title(plugType))
			}
		} else if num == "" {
			c.ui.Output(fmt.Sprintf("You have selected to skip the %s stage.", plugType))
			pNameConfirm, err := c.ui.Input(&terminal.Input{
				Prompt: fmt.Sprintf("Do you really want to skip the %s stage? (y/N): ", plugType),
				Style:  "",
				Secret: false,
			})
			if err != nil {
				c.ui.Output(
					"Error getting input: %s",
					clierrors.Humanize(err),
					terminal.WithErrorStyle(),
				)
				return plugDocs, err, false
			} else if strings.ToLower(pNameConfirm) == "exit" {
				return plugDocs, nil, true
			} else if strings.ToLower(pNameConfirm) == "yes" || strings.ToLower(pNameConfirm) == "y" {
				c.ui.Output("Step complete: %s stage skipped", strings.Title(plugType), terminal.WithWarningStyle())
				plugDocs.Name = ""
				return plugDocs, nil, false
			} else {
				c.ui.Output("Skip cancelled\n")
			}
		} else {
			c.ui.Output("Please select a numbered entry or type nothing to skip.\n")
		}
	}
	// We again unmarshal the JSON file corresponding to the file the user has selected
	if selFileName != "" {

		jsonFile, err := os.Open(fmt.Sprintf("%s/%s", fPath, selFileName))
		if err != nil {
			return plugDocs, err, false
		}
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &plugDocs)
		if plugDocs.Name != "" {
			c.ui.Output(fmt.Sprintf("You have selected the %s %s plugin.", plugDocs.Name, plugType))

		} else {
			//TODO: better error here, do we need to check again here?
			return plugDocs, nil, true
		}
		return plugDocs, nil, false
	}
	return plugDocs, nil, false
}

// Gets either a project or app name for an HCL file, pa should be either "project" or "app"
func (c *InitCommand) getName(pa string) (string, error, bool) {
	if pa == "project" {
		c.ui.Output("Please enter the name of your project. A project typically maps 1:1 to a VCS repository. This name must be unique for your Waypoint server. If you're running in local mode, this must be unique to your machine.\n")
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
		paName, err := c.ui.Input(&terminal.Input{
			Prompt: prompt + ": ",
			Style:  "",
			Secret: false,
		})
		if err != nil {
			c.ui.Output(
				"Error getting input: %s",
				clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return "", err, false
		} else if strings.ToLower(paName) == "exit" {
			return "", nil, true
		} else if strings.ToLower(paName) == "" {
			c.ui.Output(prompt + ".\n")
		} else {
			c.ui.Output("You inputted \"%s\"\n", paName)
			pNameConfirm, err := c.ui.Input(&terminal.Input{
				Prompt: fmt.Sprintf("Is this %s name correct? (y/N): ", pa),
				Style:  "",
				Secret: false,
			})
			if err != nil {
				c.ui.Output(
					"Error getting input: %s",
					clierrors.Humanize(err),
					terminal.WithErrorStyle(),
				)
				return "", err, false
			} else if strings.ToLower(pNameConfirm) == "exit" {
				return "", nil, true
			} else if strings.ToLower(pNameConfirm) == "yes" || strings.ToLower(pNameConfirm) == "y" {
				c.ui.Output("%s name confirmed", strings.Title(pa), terminal.WithSuccessStyle())
				name = paName
				getName = false
			} else {
				c.ui.Output("%s name rejected", strings.Title(pa))
			}
		}
	}
	return name, nil, false
}

func (c *InitCommand) validateConfig() bool {
	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Validating configuration file...")
	cfg, _, err := c.initConfig(c.fromProject)
	if err != nil {
		c.stepError(s, initStepConfig, err)
		return false
	}

	if cfg == nil {
		// This should never happen, because if there is no config, init should have created
		// it and exited earlier.
		err = errors.New("No configuration file found")
		c.stepError(s, initStepConfig, err)
		return false
	}

	c.cfg = cfg
	c.refProject = &pb.Ref_Project{Project: cfg.Project}

	s.Update("Configuration file appears valid")
	s.Status(terminal.StatusOK)
	s.Done()

	return true
}

func (c *InitCommand) validateServer() bool {
	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Validating server credentials...")
	client, err := c.initClient(nil)
	if err != nil {
		c.stepError(s, initStepConnect, err)
		return false
	}
	c.project = client

	if c.project.Local() {
		s.Update("Local mode initialized successfully")
	} else {
		s.Update("Connection to Waypoint server was successful")
	}

	s.Status(terminal.StatusOK)
	s.Done()
	return true
}

func (c *InitCommand) validateProject() bool {
	sg := c.ui.StepGroup()
	defer sg.Wait()

	ref := c.project.Ref()

	s := sg.Add("Checking if project %q is registered...", ref.Project)

	client := c.project.Client()
	resp, err := client.GetProject(c.Ctx, &pb.GetProjectRequest{Project: ref})
	if status.Code(err) == codes.NotFound {
		err = nil
		resp = nil
	}
	if err != nil {
		c.stepError(s, initStepProject, err)
		return false
	}

	var project *pb.Project
	if resp != nil {
		project = resp.Project
	}

	// If the project itself is missing, then register that.
	if project == nil || c.update {
		if project == nil {
			s.Status(terminal.StatusWarn)
			s.Update("Project %q is not registered with the server. Registering...", ref.Project)
		} else {
			s.Update("Updating project %q...", ref.Project)
		}

		// We need to load the data source configuration if we have it
		var ds *pb.Job_DataSource
		if dscfg := c.cfg.Runner.DataSource; dscfg != nil {
			factory, ok := datasource.FromString[dscfg.Type]
			if !ok {
				c.stepError(s, initStepProject, fmt.Errorf(
					"runner data source type %q unknown", dscfg.Type))
				return false
			}

			source := factory()
			ds, err = source.ProjectSource(dscfg.Body, c.cfg.HCLContext())
			if err != nil {
				c.stepError(s, initStepProject, err)
				return false
			}
		}

		var poll *pb.Project_Poll
		if v := c.cfg.Runner.Poll; v != nil {
			poll = &pb.Project_Poll{
				Enabled:  v.Enabled,
				Interval: v.Interval,
			}
		}

		resp, err := client.UpsertProject(c.Ctx, &pb.UpsertProjectRequest{
			Project: &pb.Project{
				Name:           ref.Project,
				RemoteEnabled:  c.cfg.Runner.Enabled,
				DataSource:     ds,
				DataSourcePoll: poll,
			},
		})
		if err != nil {
			c.stepError(s, initStepProject, err)
			return false
		}
		s.Status(terminal.StatusOK)

		project = resp.Project
	}

	pt := &serverptypes.Project{Project: project}
	for _, name := range c.cfg.Apps() {
		if pt.App(name) >= 0 {
			continue
		}

		// Missing an application, register it.
		s.Status(terminal.StatusWarn)
		s.Update("Application %q is not registered with the server. Registering...", name)

		_, err := client.UpsertApplication(c.Ctx, &pb.UpsertApplicationRequest{
			Project: ref,
			Name:    name,
		})
		if err != nil {
			c.stepError(s, initStepProject, err)
			return false
		}
		s.Status(terminal.StatusOK)
	}

	s.Update("Project %q and all apps are registered with the server.", ref.Project)
	s.Status(terminal.StatusOK)
	s.Done()
	return true
}

func (c *InitCommand) validatePlugins() bool {
	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Validating required plugins...")

	_, err := c.project.Validate(c.Ctx, &pb.Job_ValidateOp{})
	if err != nil {
		c.stepError(s, initStepPluginConfig, err)
		return false
	}

	s.Update("Plugins loaded and configured successfully")
	s.Status(terminal.StatusOK)
	s.Done()
	return true
}

func (c *InitCommand) validateAuth() bool {
	sg := c.ui.StepGroup()
	defer func() { sg.Wait() }() // defer a func so we can overwrite sg

	s := sg.Add("Checking auth for the configured components...")

	failures := false
	for _, name := range c.cfg.Apps() {
		app := c.project.App(name)

		ref := app.Ref()
		s.Update("Checking auth for app: %q", ref.Application)

		result, err := app.Auth(c.Ctx, &pb.Job_AuthOp{
			CheckOnly: true,
		})
		if err != nil {
			c.stepError(s, initStepAuth, err)
			return false
		}

		var requiresAuth []*pb.Component
		for _, r := range result.Results {
			if r.CheckResult {
				continue
			}

			requiresAuth = append(requiresAuth, r.Component)
		}

		if len(requiresAuth) == 0 {
			continue
		}
		failures = true

		// Update the status and end the step so we can output normal text
		s.Status(terminal.StatusWarn)
		s.Update("%q has plugins that require authentication:", ref.Application)
		s.Done()
		sg.Wait()

		for _, comp := range requiresAuth {
			c.ui.Output("- %s %q",
				strings.Title(strings.ToLower(comp.Type.String())),
				comp.Name,
				terminal.WithStyle(terminal.WarningStyle))
		}

		if c.ui.Interactive() {
			c.ui.Output("")
			c.ui.Output(
				strings.TrimSpace(initStepStrings[initStepAuth].Other["guide"])+"\n",
				terminal.WithStyle(terminal.WarningBoldStyle),
			)

			auth, err := c.inputContinue(terminal.WarningBoldStyle)
			if err != nil {
				c.stepError(s, initStepAuth, err)
				return false
			}
			if !auth {
				return false
			}

			// Mark failures as false since the user is trying to auth!
			failures = false

			for i, comp := range requiresAuth {
				c.ui.Output("Authenticating %s %q",
					strings.Title(strings.ToLower(comp.Type.String())),
					comp.Name,
					terminal.WithStyle(terminal.HeaderStyle),
				)

				resultRaw, err := app.Auth(c.Ctx, &pb.Job_AuthOp{
					Component: &pb.Ref_Component{
						Type: comp.Type,
						Name: comp.Name,
					},
				})
				if err != nil {
					c.stepError(s, initStepAuth, err)
					return false
				}

				// This should always be exactly one...
				if len(resultRaw.Results) != 1 {
					c.stepError(s, initStepAuth, fmt.Errorf(
						"unexpected result from server on auth: %#v",
						resultRaw))
					return false
				}
				result := resultRaw.Results[0]

				// Check the results
				if !result.AuthCompleted {
					// If we didn't authenticate at all, we still have failures.
					failures = true
				} else if !result.CheckResult {
					// If auth failed, then we still have failures but we also
					// should tell the user.
					failures = true

					c.ui.Output(
						strings.TrimSpace(initStepStrings[initStepAuth].Other["auth-failure"]),
						status.FromProto(result.CheckError).Message(),
						terminal.WithStyle(terminal.WarningBoldStyle),
					)
				} else {
					sg = c.ui.StepGroup()
					s = sg.Add("%s %q authenticated successfully.",
						strings.Title(strings.ToLower(comp.Type.String())),
						comp.Name,
					)
					s.Done()
					sg.Wait()
				}

				if i+1 < len(requiresAuth) {
					auth, err := c.inputContinue(terminal.WarningBoldStyle)
					if err != nil {
						c.stepError(s, initStepAuth, err)
						return false
					}
					if !auth {
						return false
					}
				}
			}
		}

		// Initialize a new step group for remaining apps
		sg = c.ui.StepGroup()
		s = sg.Add("")
	}

	if !failures {
		s.Update("Authentication requirements appear satisfied.")
		s.Status(terminal.StatusOK)
	} else {
		s.Update("Authentication checks had failures.")
		s.Status(terminal.StatusError)
	}

	// If we aren't interactive with failures, then we want to report as
	// an error since the user couldn't have corrected them.
	if !c.ui.Interactive() && failures {
		c.stepError(s, initStepAuth, fmt.Errorf(
			"The plugins above reported that they aren't authenticated."))
		return false
	}

	s.Done()
	return !failures
}

func (c *InitCommand) stepError(s terminal.Step, step initStepType, err error) {
	stepStrings := initStepStrings[step]

	s.Status(terminal.StatusError)
	s.Update(stepStrings.Error)
	s.Done()
	c.ui.Output("")
	if v := stepStrings.ErrorDetails; v != "" {
		c.ui.Output(strings.TrimSpace(v), terminal.WithErrorStyle())
		c.ui.Output("")
	}
	c.ui.Output(clierrors.Humanize(err), terminal.WithErrorStyle())
}

func (c *InitCommand) inputContinue(style string) (bool, error) {
	for {
		result, err := c.ui.Input(&terminal.Input{
			Prompt: "Continue? [y/n]",
			Style:  style,
		})
		if err != nil {
			return false, err
		}
		if result == "y" || result == "n" {
			return result == "y", nil
		}
	}
}

func (c *InitCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:    "from-project",
			Target:  &c.fromProject,
			Default: "",
			Usage: "Create a new application by fetching the given application from " +
				"a remote source or from a local project folder or file on disk.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "into",
			Target:  &c.into,
			Default: "",
			Usage:   "Where to write the application fetched via -from-project",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "update",
			Target:  &c.update,
			Default: false,
			Usage: "Update the project configuration if it already exists. This can be used " +
				"to update settings such as the remote runner data source.",
		})
	})
}

func (c *InitCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *InitCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *InitCommand) Synopsis() string {
	return "Initialize and validate a project"
}

func (c *InitCommand) Help() string {
	return formatHelp(`
Usage: waypoint init [options]

  Initialize and validate a project.

  This is the first command that should be run for any new or existing
  Waypoint project per machine. This sets up the project if required and
  also validates that operations such as "up" will most likely work.

  This command is always safe to run multiple times. This command will never
  delete your configuration or any data in the server.

` + c.Flags().Help())
}

type initStepType uint

const (
	initStepInvalid initStepType = iota
	initStepConfig
	initStepConnect
	initStepPluginConfig
	initStepProject
	initStepAuth
)

var initStepStrings = map[initStepType]struct {
	Error        string
	ErrorDetails string
	Other        map[string]string
}{
	initStepConfig: {
		Error: "Error loading configuration!",
	},

	initStepConnect: {
		Error: "Failed to initialize client for Waypoint server.",
		ErrorDetails: `
The Waypoint client validation step validates that we can connect to the
configured Waypoint server. If this is a local-only operation (no Waypoint
server is configured), then we validate that we can initialize local writes.
The error for this failure is shown below.
			`,
	},

	initStepPluginConfig: {
		Error: "Failed to load and validate plugins!",
		ErrorDetails: `
This validation check ensures that you have all the required plugins available
and the configuration for each plugin (if it exists) is valid. The error message
below should tell you which plugin(s) failed.
		`,
	},

	initStepProject: {
		Error: "Error while checking for project registration.",
		ErrorDetails: `
There was an error while the checking if the project and applications
are registered with the Waypoint server. This error may be temporary and
you may retry to init. See the error message below.
		`,

		Other: map[string]string{
			"unregistered-desc": `
The project and apps must be registered prior to performing any operations.
This creates some metadata with the server. We require registration as a
verification that the project/app names are correct and that you're targeting
the correct server.
			`,
		},
	},

	initStepAuth: {
		Error: "Failed to check authentication requirements!",
		ErrorDetails: `
This step verifies that Waypoint has access to the configured systems.
This is a best-effort check, since not all plugins support this check
and the check can often only check that any known credentials work at
a minimal level.

There was an error during this step and it is shown below.
		`,

		Other: map[string]string{
			"guide": `
Waypoint will guide you through the authentication process one plugin
at a time. Plugins may interactively attempt to authenticate or they may
just output help text to guide you there. You can use Ctrl-C at any point
to cancel and run "waypoint init" again later.
			`,

			"auth-failure": `
Authentication failed with error: %s
			`,
		},
	},
}
