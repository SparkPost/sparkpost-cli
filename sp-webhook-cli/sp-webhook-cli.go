package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"

	sp "github.com/SparkPost/gosparkpost"
)

func main() {

	ValidParameters := []string{
		"timezone", "limit",
	}

	app := cli.NewApp()

	app.Version = "0.0.1"
	app.Name = "sparkpost-message-event-cli"
	app.Usage = "SparkPost Message Event CLI\n\n\tSee https://developers.sparkpost.com/api/webhooks.html"
	app.Flags = []cli.Flag{
		// Core Client Configuration
		cli.StringFlag{
			Name:   "baseurl, u",
			Value:  "https://api.sparkpost.com",
			Usage:  "Optional baseUrl for SparkPost.",
			EnvVar: "SPARKPOST_BASEURL",
		},
		cli.StringFlag{
			Name:   "apikey, k",
			Value:  "",
			Usage:  "Required SparkPost API key",
			EnvVar: "SPARKPOST_API_KEY",
		},
		cli.StringFlag{
			Name:  "username",
			Value: "",
			Usage: "Username this is a special case it is more common to use apikey",
		},
		cli.StringFlag{
			Name:  "password, p",
			Value: "",
			Usage: "Username this is a special it is more common to use apikey",
		},
		cli.StringFlag{
			Name:  "verbose",
			Value: "false",
			Usage: "Dumps additional information to console",
		},

		// Webhook Parameters
		cli.StringFlag{
			Name:  "command, c",
			Value: "list",
			Usage: "Optional one of list, query, status. Default is \"list\"",
		},
		cli.StringFlag{
			Name:  "timezone, tz",
			Value: "",
			Usage: "Optional Standard timezone identification string, defaults to UTC Example: America/New_York.",
		},
		cli.StringFlag{
			Name:  "id",
			Value: "",
			Usage: "Optional UUID identifying a webhook Example: 12affc24-f183-11e3-9234-3c15c2c818c2.",
		},
		cli.StringFlag{
			Name:  "limit",
			Value: "",
			Usage: "Optional Maximum number of results to return. Defaults to 1000. Example: 1000.",
		},
	}
	app.Action = func(c *cli.Context) {

		if c.String("baseurl") == "" {
			log.Fatalf("Error: SparkPost BaseUrl must be set\n")
			return
		}

		if c.String("apikey") == "" && c.String("username") == "" && c.String("password") == "" {
			log.Fatalf("Error: SparkPost API key must be set\n")
			return
		}

		isVerbose := false

		if c.String("verbose") == "true" {
			isVerbose = true
		}

		//println("SparkPost baseUrl: ", c.String("baseurl"))

		cfg := &sp.Config{
			BaseUrl:    c.String("baseurl"),
			ApiKey:     c.String("apikey"),
			Username:   c.String("username"),
			Password:   c.String("password"),
			ApiVersion: 1,
			Verbose:    isVerbose,
		}

		var client sp.Client
		err := client.Init(cfg)
		if err != nil {
			log.Fatalf("SparkPost client init failed: %s\n", err)
		}

		parameters := make(map[string]string)

		for i, val := range ValidParameters {

			if c.String(ValidParameters[i]) != "" {
				parameters[val] = c.String(val)
			}
		}

		switch c.String("command") {
		case "list":
			doList(&client, parameters)
		case "query":
			doItemQuery(&client, parameters, c.String("id"))
		case "status":
			doStatus(&client, parameters, c.String("id"))

		default:
			log.Fatalf("ERROR: Unknown \"command\" [%s]. Try --help for a list of available commands.\n", c.String("command"))
		}
	}
	app.Run(os.Args)

}

func doStatus(client *sp.Client, parameters map[string]string, id string) {

	statusWrapper := &sp.WebhookStatusWrapper{}
	statusWrapper.Params = parameters
	statusWrapper.ID = id

	e, err := client.WebhookStatus(statusWrapper)

	if err != nil {
		log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
		return
	} else if e.Errors != nil {
		log.Fatalf("ERROR: %v.\n\nFor additional information try using `--verbose true`\n\n\n", e.Errors)
		return
	} else if statusWrapper.Errors != nil {
		log.Fatalf("ERROR: %v.\n\nFor additional information try using `--verbose true`\n\n\n", statusWrapper.Errors)
		return
	}

	for _, element := range statusWrapper.Results {
		webhookStatusPrinter(element)
	}
}

func doItemQuery(client *sp.Client, parameters map[string]string, id string) {
	queryWrapper := &sp.WebhookQueryWrapper{}
	queryWrapper.Params = parameters
	queryWrapper.ID = id

	e, err := client.QueryWebhook(queryWrapper)

	if err != nil {
		log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
		return
	} else if e.Errors != nil {
		log.Fatalf("ERROR: %v.\n\nFor additional information try using `--verbose true`\n\n\n", e.Errors)
		return
	} else if queryWrapper.Errors != nil {
		log.Fatalf("ERROR: %v.\n\nFor additional information try using `--verbose true`\n\n\n", queryWrapper.Errors)
		return
	}

	// for _, element := range e.Results {
	webhookDetailPrinter(queryWrapper.Results)
	// }
}

func doList(client *sp.Client, parameters map[string]string) {
	listWrapper := &sp.WebhookListWrapper{}
	listWrapper.Params = parameters

	e, err := client.Webhooks(listWrapper)

	if err != nil {
		log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
		return
	} else if e.Errors != nil {
		log.Fatalf("ERROR: %v.\n\nFor additional information try using `--verbose true`\n\n\n", e.Errors)
		return
	} else if listWrapper.Errors != nil {
		log.Fatalf("ERROR: %v.\n\nFor additional information try using `--verbose true`\n\n\n", listWrapper.Errors)
		return
	}

	for _, element := range listWrapper.Results {
		listSummaryPrinter(element)
	}
}

func listSummaryPrinter(event *sp.WebhookItem) {
	row := ""

	row = fmt.Sprintf("Name: \"%s\"\n", event.Name)
	row = fmt.Sprintf("%s\thook ID:   %s\n", row, event.ID)
	row = fmt.Sprintf("%s\tTarget:    %s\n", row, event.Target)
	row = fmt.Sprintf("%s\tSuccess:   %s\n", row, event.LastSuccessful)
	row = fmt.Sprintf("%s\tFail:      %s\n", row, event.LastFailure)
	row = fmt.Sprintf("%s\tAuthType:  %s\n", row, event.AuthType)

	fmt.Println(row)
}

func listHeaderPrinter(fields []string) {
	row := "domain, "
	for i := range fields {
		row = fmt.Sprintf("%s%s, ", row, fields[i])
	}

	fmt.Println(row)
}

func webhookDetailPrinter(event *sp.WebhookItem) {
	row := ""
	row = fmt.Sprintf("Name: \"%s\"\n", event.Name)
	row = fmt.Sprintf("%s\thook ID:   %s\n", row, event.ID)
	row = fmt.Sprintf("%s\tTarget:    %s\n", row, event.Target)
	row = fmt.Sprintf("%s\tSuccess:   %s\n", row, event.LastSuccessful)
	row = fmt.Sprintf("%s\tFail:      %s\n", row, event.LastFailure)
	row = fmt.Sprintf("%s\tAuthType:  %s\n", row, event.AuthType)
	if event.Events != nil {
		row = fmt.Sprintf("%s\tEvents:\n", row)
		for i := range event.Events {
			row = fmt.Sprintf("%s\t\t%s\n", row, event.Events[i])
		}
	}

	fmt.Println(row)
}

func webhookStatusPrinter(event *sp.WebhookStatus) {
	row := ""
	row = fmt.Sprintf("BatchId: \"%s\"\n", event.BatchID)
	row = fmt.Sprintf("%s\tTime:       %s\n", row, event.Timestamp)
	row = fmt.Sprintf("%s\tAttempts:   %d\n", row, event.Attempts)
	row = fmt.Sprintf("%s\tRespCode:   %s\n", row, event.ResponseCode)

	fmt.Println(row)
}
