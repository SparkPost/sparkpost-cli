package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"

	sp "github.com/SparkPost/gosparkpost"
)

func main() {

	VALID_PARAMETERS := []string{
		"bounce_classes", "campaign_ids", "events", "friendly_froms", "from",
		"message_ids", "page", "per_page", "reason", "recipients", "template_ids",
		"timezone", "to", "transmission_ids", "subaccounts",
	}

	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "sparkpost-message-event-cli"
	app.Usage = "SparkPost Message Event CLI"
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
			Usage: "Username it is more common to use apikey",
		},
		cli.StringFlag{
			Name:  "password, p",
			Value: "",
			Usage: "Username it is more common to use apikey",
		},
		cli.StringFlag{
			Name:  "verbose",
			Value: "false",
			Usage: "Dumps additional information to console",
		},

		// Metrics Parameters
		cli.StringFlag{
			Name:  "bounce_classes, b",
			Value: "",
			Usage: "Optional comma-delimited list of bounce classification codes to search.",
		},
		cli.StringFlag{
			Name:  "campaign_ids, i",
			Value: "",
			Usage: "Optional comma-delimited list of campaign ID's to search. Example: \"Example Campaign Name\"",
		},
		cli.StringFlag{
			Name:  "events, e",
			Value: "",
			Usage: "Optional comma-delimited list of event types to search. Defaults to all event types.",
		},
		cli.StringFlag{
			Name:  "friendly_froms",
			Value: "",
			Usage: "Optional comma-delimited list of friendly_froms to search",
		},
		cli.StringFlag{
			Name:  "from, f",
			Value: "",
			Usage: "Optional Datetime in format of YYYY-MM-DDTHH:MM. Example: 2016-02-10T08:00. Default: One hour ago",
		},
		cli.StringFlag{
			Name:  "message_ids",
			Value: "",
			Usage: "Optional Comma-delimited list of message ID's to search. Example: 0e0d94b7-9085-4e3c-ab30-e3f2cd9c273e.",
		},
		cli.StringFlag{
			Name:  "page",
			Value: "",
			Usage: "Optional results page number to return. Used with per_page for paging through result. Example: 25. Default: 1",
		},
		cli.StringFlag{
			Name:  "per_page",
			Value: "",
			Usage: "Optional number of results to return per page. Must be between 1 and 10,000 (inclusive). Example: 100. Default: 1000.",
		},
		cli.StringFlag{
			Name:  "reason",
			Value: "",
			Usage: "Optional bounce/failure/rejection reason that will be matched using a wildcard (e.g., %%reason%%). Example: bounce.",
		},
		cli.StringFlag{
			Name:  "recipients",
			Value: "",
			Usage: "Optional Comma-delimited list of recipients to search. Example: recipient@example.com",
		},
		cli.StringFlag{
			Name:  "template_ids",
			Value: "",
			Usage: "Optional Comma-delimited list of template ID's to search. Example: templ-1234.",
		},
		cli.StringFlag{
			Name:  "timezone",
			Value: "",
			Usage: "Optional Standard timezone identification string. Example: America/New_York. Default: UTC",
		},
		cli.StringFlag{
			Name:  "to",
			Value: "",
			Usage: "Optional Datetime in format of YYYY-MM-DDTHH:MM. Example: 2016-02-10T00:00. Default: now.",
		},
		cli.StringFlag{
			Name:  "transmission_ids",
			Value: "",
			Usage: "Optional Comma-delimited list of transmission ID's to search (i.e. id generated during creation of a transmission). Example: 65832150921904138.",
		},
		cli.StringFlag{
			Name:  "subaccounts",
			Value: "",
			Usage: "Optional Comma-delimited list of subaccount ID's to search. Example: 101. Passing '0' will show data for only the master account",
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

		for i, val := range VALID_PARAMETERS {

			if c.String(VALID_PARAMETERS[i]) != "" {
				parameters[val] = c.String(val)
			}
		}

		e, err := client.MessageEvents(parameters)
		//e, err := client.SearchMessageEvents(nil)
		if err != nil {
			log.Fatalf("Error: %s\n For additional information try using `--verbose true`\n", err)
			return
		} else {

			for index, element := range e.Events {
				log.Printf("%d\t %s%s", index, element, "\n")
				//log.Printf("%d\t %v\n", index, element)
			}

			log.Printf("\t-------------------\n")
			log.Printf("\tResult Count: %d\n", e.TotalCount)
		}
	}
	app.Run(os.Args)

}
