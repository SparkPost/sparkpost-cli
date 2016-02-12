package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"

	sp "github.com/SparkPost/gosparkpost"
)

func main() {

	VALID_PARAMETERS := []string{
		"bounce_classes", "campaign_ids", "events", "friendly_froms", "from",
		"message_ids", "page", "per_page", "reason", "recipients", "template_ids",
		"timezone", "to", "transmission_ids",
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
			Usage: "Optional Datetime in format of YYYY-MM-DDTHH:MM. Example: 2014-07-20T09:00. Default: now.",
		},
		cli.StringFlag{
			Name:  "transmission_ids",
			Value: "",
			Usage: "Optional Comma-delimited list of transmission ID's to search (i.e. id generated during creation of a transmission). Example: 65832150921904138.",
		},
	}
	app.Action = func(c *cli.Context) {

		if c.String("baseurl") == "" {
			log.Fatalf("Error: SparkPost BaseUrl must be set\n")
			return
		}

		if c.String("apikey") == "" {
			log.Fatalf("Error: SparkPost API key must be set\n")
			return
		}

		//println("SparkPost baseUrl: ", c.String("baseurl"))

		cfg := &sp.Config{
			BaseUrl:    c.String("baseurl"),
			ApiKey:     c.String("apikey"),
			ApiVersion: 1,
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

		e, err := client.SearchMessageEvents(parameters)
		//e, err := client.SearchMessageEvents(nil)
		if err != nil {
			log.Fatalf("Error: %s\n", err)
			return
		} else if e.Errors != nil {
			fmt.Println("ERROR: ", e.Errors)
		} else {

			for index, element := range e.Results {
				fmt.Printf("%d\t %s%s", index, client.EventAsString(element), "\n")
				//fmt.Printf("%d\t %v\n", index, element)
			}
			
			fmt.Printf("\t-------------------\n")
			fmt.Printf("\tResult Count: %d\n", e.TotalCount)
		}
	}
	app.Run(os.Args)

}
