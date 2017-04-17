package main

import (
	"log"
	"os"
	"time"

	"github.com/codegangsta/cli"

	sp "github.com/SparkPost/gosparkpost"
)

func main() {

	validParameters := []string{
		"bounce_classes", "campaign_ids", "events", "friendly_froms", "from",
		"message_ids", "page", "per_page", "reason", "recipients", "template_ids",
		"timezone", "to", "transmission_ids", "subaccounts",
	}

	app := cli.NewApp()
	app.Version = "0.0.2"
	app.Name = "sparkpost-message-event-cli"
	app.Usage = "SparkPost Message Event CLI \n\n\tSee also https://developers.sparkpost.com/api/message-events.html#message-events-message-events-get"
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
		cli.StringFlag{
			Name:  "pause",
			Value: "0",
			Usage: "Seconds to pause before fetching next page of results. Used to guard against rate limit errors.",
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
			Usage: "Optional Comma-delimited list of subaccount ID's to search. Example: 101",
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

		for i, val := range validParameters {

			if c.String(validParameters[i]) != "" {
				parameters[val] = c.String(val)
			}
		}

		eventPage := &sp.EventsPage{}
		eventPage.Params = parameters

		r, err := client.MessageEventsSearch(eventPage)
		totalCount := eventPage.TotalCount

		if err != nil {
			log.Fatalf("Error: %s\n For additional information try using `--verbose true`\n", err)
			return
		}

		sleepTimeout := time.Duration(c.Int64("pause")) * time.Second
		for {
			if eventPage == nil {
				if isVerbose {
					log.Printf("Event page nil")
				}
				break
			}

			if eventPage.Errors != nil {
				log.Fatalf("Error: %v\n For additional information try using `--verbose true`\n", eventPage.Errors)
				break
			}

			if len(eventPage.Events) == 0 {
				if isVerbose {
					log.Printf("Dump: %v", r)
					log.Printf("No more events")
				}
				break
			}

			printEvents(eventPage)

			if c.String("page") != "" {
				break
			}

			if isVerbose {
				log.Printf("NextPage(): %s", eventPage.NextPage)
			}
			if sleepTimeout != 0 {
				if isVerbose {
					log.Printf("Sleep: %d seconds", c.Int64("pause"))
				}
				time.Sleep(sleepTimeout)
			}
			eventPage, r, err = eventPage.Next()
			if err != nil {
				log.Fatalf("Error: %s\n For additional information try using `--verbose true`\n", err)
				break
			}

		}

		log.Printf("\t-------------------\n")
		log.Printf("\tResult Count: %d\n", totalCount)

	}
	app.Run(os.Args)

}

func printEvents(eventPage *sp.EventsPage) {
	for index, event := range eventPage.Events {
		log.Printf("%d\t %s%s", index, event, "\n")
	}
}
