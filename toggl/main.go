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
	app.Name = "toggl"
	app.Usage = "toggl.com console client"
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
			Name:      	"status",
			Usage:     	"show status",
			Action: 	CommandStatus,
		},
		{
			Name:      	"token",
			Usage:     	"show or set api token",
			Action: 	CommandToken,
		},
	}

	app.Run(os.Args)

	fmt.Println()
}

func GetApiClient() toggl.ApiClient {
	return *toggl.NewClient(config.ApiToken, nil)
}

func CommandDefault(c *cli.Context) {
	fmt.Println("Nothing...")
}

func CommandStatus(c *cli.Context) {
	apiClient := GetApiClient()

	currentEntry := apiClient.TimeEntries.Current()

	fmt.Print("Current active: ")

	if currentEntry.IsPersisted() {
		fmt.Print(currentEntry.Description)
	} else {
		fmt.Println("nothing")
	}
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

		projects, _ := apiClient.Projects.GetByNamePrefix(projectName, entry.WorkspaceID)

		if len(projects) > 0 {
			if len(projects) == 1 {
				entry.ProjectID = projects[0].ID
			} else {
				clients := apiClient.Clients.All()

				for i := 0; i < len(projects); i++ {
					client := "No client"
					for j := 0; j < len(*clients); j++ {
						if projects[i].ClientID == (*clients)[j].ID {
							client = (*clients)[j].Name
						}
					}

					fmt.Printf("[%d]\t%s\t[%s]\n", i + 1, projects[i].Name, client)
				}

				for {
					var n int

					fmt.Print("Select project: ")
					fmt.Scanf("%d", &n)

					if n > 0 && n <= len(projects) {
						entry.ProjectID = projects[n - 1].ID

						break
					} else {
						fmt.Println("Mmm?")
						break
					}
				}
			}
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

func CommandToken(c *cli.Context) {
	if len(c.Args()) == 1 {
		config.ApiToken = c.Args().First()
		config.Save()
	}

	fmt.Printf("Your token: %s", config.ApiToken)
}