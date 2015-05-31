package main

import (
	"github.com/iamdeuterium/go-toggl/toggl"
	"fmt"
	"os"
	"github.com/codegangsta/cli"
	"strings"
)

var config *Configuration = new(Configuration)

func main() {
	if config.Load() == false && len(os.Args) > 2 && os.Args[1] != "config" {
		fmt.Println("Run: toggl config set api_token YOUR_API_TOKEN")

		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: ", r)
		}
	}()

	app := cli.NewApp()
	app.Name = "toggl.com console client"
	app.Usage = "be happy!"
	app.Action = CommandDefault

	app.Commands = []cli.Command{
		{
			Name:      	"start",
			Aliases:   	[]string{"s"},
			Usage:     	"start timer",
			Action: 	CommandStart,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "project, p",
					Usage: "project for task",
				},
				cli.StringFlag{
					Name: "workspace, w",
					Usage: "workspace for task",
				},
			},
		},
		{
			Name:      	"stop",
			Usage:     	"stop timer",
			Action: 	CommandStop,
		},
		{
			Name:      	"config",
			Usage:     	"show config",
			Action: 	CommandConfig,
			Subcommands: []cli.Command{
				{
					Name:  "set",
					Usage: "set",
					Action: CommandConfigSet,
				},
			},
		},
	}

	app.Run(os.Args)
}

func GetApiClient() toggl.ApiClient {
	return *toggl.NewClient(config.ApiToken, nil)
}

func CommandDefault(c *cli.Context) {
	fmt.Println("Nothing...")
}

func CommandStart(c *cli.Context) {
	apiClient := GetApiClient()

	if len(c.Args()) == 1 && c.Args().First() == "last" {
		entries := apiClient.TimeEntries.All()

		if len(*entries) > 0 {
			entry := (*entries)[0]
			apiClient.TimeEntries.Start(entry)
			fmt.Printf("started: %s", entry.Description)
		} else {
			fmt.Println("Nothing to start.")
		}

		return
	}

	if len(c.Args()) == 0 {
		entries := apiClient.TimeEntries.All()

		if len(*entries) == 0 {
			fmt.Println("Nothing to start.")
			return
		}

		for i := 0; i < len(*entries); i++ {
			fmt.Printf("[%d]\t%s\n", i + 1, (*entries)[i].Description)
		}

		var n int

		for {
			fmt.Print("Select last entry: ")
			fmt.Scanf("%d", &n)

			if n > 0 && n <= len(*entries) {
				entry := (*entries)[n - 1]
				apiClient.TimeEntries.Start(entry)

				fmt.Printf("Started selected: %s", entry.Description)

				break
			} else {
				fmt.Println("Mmm?")
				break
			}
		}

		return
	}

	entry := new(toggl.TimeEntry)
	entry.Description = strings.Join(c.Args(), " ")

	workspaceName := c.String("workspace")

	if len(workspaceName) > 0 {
		workspace, ok := apiClient.Workspaces.GetByName(workspaceName)

		if ok {
			entry.WorkspaceID = workspace.ID
		}
	}

	projectName := c.String("project")

	if len(projectName) > 0 {
		if entry.WorkspaceID == 0 {
			entry.WorkspaceID = apiClient.Users.Current().DefaultWorkspaceID
		}

		project, ok := apiClient.Projects.GetByName(projectName, entry.WorkspaceID)

		if ok {
			entry.ProjectID = project.ID
		}
	}

	apiClient.TimeEntries.Start(*entry)

	fmt.Printf("Started: %s", entry.Description)
}

func CommandStop(c *cli.Context) {
	apiClient := GetApiClient()

	entry := apiClient.TimeEntries.Current()

	if entry.IsPersisted() {
		apiClient.TimeEntries.Stop(entry)

		fmt.Printf("Task '%s' was stopped.", entry.Description)
	} else {
		fmt.Println("Nothing to stop.")
	}
}

func CommandConfig(c *cli.Context) {

}

func CommandConfigSet(c *cli.Context) {
	if len(c.Args()) < 2 {
		return
	}

	 if c.Args()[0] == "api_token" {
		config.ApiToken = c.Args()[1]
		config.Save()
	}
}