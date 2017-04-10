package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/codegangsta/cli"
)

// Column mapping for Mandrill Blacklist
const (
	MandrillEmailCol      = 0
	MandrillReasonCol     = 1
	MandrillDetailCol     = 2
	MandrillCreatedCol    = 3
	MandrillExpiresAtCol  = 4
	MandrillLastEventCol  = 5
	MandrillExpiresAt2Col = 6
	MandrillSubAccountCol = 7
)

// Column mapping for SendGrid Blacklist
const (
	SendgridEmailCol = 0
	SendgridCreated  = 1
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	ValidParameters := []string{
		"to", "from", "domain", "cursor", "limit", "per_page", "page", "sources", "types", "description",
	}

	app := cli.NewApp()

	app.Version = "0.0.2"
	app.Name = "suppression-sparkpost-cli"
	app.Usage = "SparkPost suppression list cli\n\n\tSee https://developers.sparkpost.com/api/suppression-list.html"
	app.Flags = []cli.Flag{
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
			Name:  "verbose",
			Value: "false",
			Usage: "Dumps additional information to console",
		},
		cli.StringFlag{
			Name:  "file, f",
			Value: "",
			Usage: "Compatible blacklist CSV file. See README.md for more info.",
		},
		cli.StringFlag{
			Name:  "command",
			Value: "list",
			Usage: "Optional one of list, retrieve, search, delete, mandrill, sendgrid",
		},
		cli.StringFlag{
			Name:  "recipient",
			Value: "",
			Usage: "Recipient email address. Example rcpt_1@example.com",
		},

		// Search Parameters
		cli.StringFlag{
			Name:  "from",
			Value: "",
			Usage: "Optional datetime the entries were last updated, in the format of YYYY-MM-DDTHH:mm:ssZ (2015-04-10T00:00:00)",
		},
		cli.StringFlag{
			Name:  "to",
			Value: "",
			Usage: "Optional datetime the entries were last updated, in the format YYYY-MM-DDTHH:mm:ssZ (2015-04-10T00:00:00)",
		},
		cli.StringFlag{
			Name:  "types",
			Value: "",
			Usage: "Optional types of entries to include in the search, i.e. entries with \"transactional\" and/or \"non_transactional\" keys set to true",
		},
		cli.StringFlag{
			Name:  "limit",
			Value: "",
			Usage: "Optional maximum number of results to return. Must be between 1 and 100000. Default value is 100000",
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
			Name:  "cursor",
			Value: "",
			Usage: "Optional the results cursor location to return, to start paging with cursor, use the value of ‘initial’. When cursor is provided the page parameter is ignored. ( Note: SparkPost only). Example initial",
		},
		cli.StringFlag{
			Name:  "domain",
			Value: "",
			Usage: "Domain of entries to include in the search. ( Note: SparkPost only). Example yahoo.com",
		},
		cli.StringFlag{
			Name:  "sources",
			Value: "",
			Usage: "Types of entries to include in the search, i.e. entries that are transactional or non_transactional",
		},
		cli.StringFlag{
			Name:  "description",
			Value: "",
			Usage: "Description of the entries to include in the search, i.e descriptions that include the text submitted. ( Note: SparkPost only)",
		},
	}
	app.Action = func(c *cli.Context) {

		if c.String("apikey") == "" {
			log.Fatalf("Error: SparkPost API key must be set\n")
			return
		}

		isVerbose := false
		if c.String("verbose") == "true" {
			isVerbose = true
		}

		cfg := &sp.Config{
			BaseUrl:    c.String("baseurl"),
			ApiKey:     c.String("apikey"),
			ApiVersion: 1,
			Verbose:    isVerbose,
		}

		var client sp.Client
		err := client.Init(cfg)
		if err != nil {
			log.Fatalf("SparkPost client init failed: %s\n", err)
			return
		}

		parameters := make(map[string]string)

		for i, val := range ValidParameters {
			if c.String(ValidParameters[i]) != "" {
				parameters[val] = c.String(val)
			}
		}

		switch c.String("command") {
		case "list", "search":
			var err error
			suppressionPage := &sp.SuppressionPage{}

			parameters := make(map[string]string)
			parameters["cursor"] = "initial"

			for i, val := range ValidParameters {

				if c.String(ValidParameters[i]) != "" {
					parameters[val] = c.String(val)
				}
			}

			suppressionPage.Params = parameters
			_, err = client.SuppressionSearch(suppressionPage)

			if err != nil {
				log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
				return
			}

			for {

				if suppressionPage.Errors != nil {
					log.Fatalf("Error: %v\n For additional information try using `--verbose true`\n", suppressionPage.Errors)
					break
				}

				csvEntryPrinter(suppressionPage, true)

				// If user requested a specific page don't page through rest of results
				if c.String("page") != "" {
					return
				}

				if suppressionPage.NextPage == "" {
					return
				}

				if isVerbose {
					log.Printf("NextPage(): %s", suppressionPage.NextPage)
				}
				suppressionPage, _, err = suppressionPage.Next()
				if err != nil {
					log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
					return
				}
			}
		case "retrieve":
			recpipient := c.String("recipient")
			if recpipient == "" {
				log.Fatalf("ERROR: The `retrieve` command requires a recipient.")
				return
			}

			suppressionPage := &sp.SuppressionPage{}
			_, err := client.SuppressionRetrieve(recpipient, suppressionPage)

			if err != nil {
				log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
				return
			}
			csvEntryPrinter(suppressionPage, false)

		case "delete":
			recpipient := c.String("recipient")
			if recpipient == "" {
				log.Fatalf("ERROR: The `delete` command requires a recipient.")
				return
			}

			_, err := client.SuppressionDelete(recpipient)

			if err != nil {
				log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
				return
			}
			fmt.Println("OK")

		case "mandrill":
			fmt.Printf("Processing: %s\n", c.String("file"))
			file := c.String("file")
			if file == "" {
				log.Fatalf("ERROR: The `mandrill` command requires a CSV file.")
				return
			}

			f, err := os.Open(file)
			check(err)

			var entries = []sp.WritableSuppressionEntry{}

			batchCount := 1

			blackListRow := csv.NewReader(bufio.NewReader(f))
			blackListRow.FieldsPerRecord = 8

			for {
				record, err := blackListRow.Read()
				if err == io.EOF {
					break
				}

				if err != nil {
					log.Fatalf("ERROR: Failed to process '%s':\n\t%s", file, err)

					return
				}

				if record[MandrillEmailCol] == "email" {
					// Skip over header row
					continue
				}

				if record[MandrillReasonCol] != "hard-bounce" {
					// Ignore soft-bounce
					continue
				}

				if strings.Count(record[MandrillEmailCol], "@") != 1 {
					fmt.Printf("WARN: Ignoring '%s'. It is not a valid email address.\n", record[MandrillEmailCol])
					continue
				}

				entry := sp.WritableSuppressionEntry{}

				if record[MandrillEmailCol] == "" {
					// Must have email as it is suppression list primary key
					continue
				}

				entry.Recipient = record[MandrillEmailCol]
				entry.Type = "non_transactional"
				entry.Description = fmt.Sprintf("MBL: %s", record[MandrillDetailCol])

				entries = append(entries, entry)

				if len(entries) > (1024 * 100) {
					fmt.Printf("Uploading batch %d\n", batchCount)
					_, err := client.SuppressionUpsert(entries)

					if err != nil {
						log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
						return
					}
					entries = []sp.WritableSuppressionEntry{}
					batchCount++
				}
			}

			if len(entries) > 0 {
				fmt.Printf("Uploading batch %d\n", batchCount)
				_, err := client.SuppressionUpsert(entries)

				if err != nil {
					log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
					return
				}
			}
			fmt.Println("DONE")

		case "sendgrid":
			file := c.String("file")
			if file == "" {
				log.Fatalf("ERROR: The `sendgrid` command requires a CSV file.")
				return
			}

			f, err := os.Open(file)
			check(err)

			var entries = []sp.WritableSuppressionEntry{}

			batchCount := 1

			blackListRow := csv.NewReader(bufio.NewReader(f))
			blackListRow.FieldsPerRecord = 2

			for {
				record, err := blackListRow.Read()
				if err == io.EOF {
					break
				}

				if err != nil {
					log.Fatalf("ERROR: Failed to process '%s':\n\t%s", file, err)

					return
				}

				if record[SendgridEmailCol] == "email" {
					// Skip over header row
					continue
				}

				entry := sp.WritableSuppressionEntry{}

				if record[SendgridEmailCol] == "" {
					// Must have email as it is suppression list primary key
					continue
				}

				// SendGrid suppression lists are very dirty and tend to have invalid data. Some examples of invalid addresses are:
				// 	#02232014, gmail.com, To, 8/27/2015, name@yahoo.comett@domain.com"
				if strings.Count(record[SendgridEmailCol], "@") != 1 {
					fmt.Printf("WARN: Ignoring '%s'. It is not a valid email address.\n", record[SendgridEmailCol])
					continue
				}

				entry.Recipient = record[SendgridEmailCol]
				entry.Type = "non_transactional"
				entry.Description = fmt.Sprintf("SBL: imported from SendGrid")

				entries = append(entries, entry)

				if len(entries) > (1024 * 100) {
					fmt.Printf("Uploading batch %d\n", batchCount)
					_, err := client.SuppressionUpsert(entries)

					if err != nil {
						log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
						return
					}
					entries = []sp.WritableSuppressionEntry{}
					batchCount++
				}

			}

			if len(entries) > 0 {
				fmt.Printf("Uploading batch %d\n", batchCount)
				_, err := client.SuppressionUpsert(entries)

				if err != nil {
					log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
					return
				}
			}
			fmt.Println("DONE")

		default:
			fmt.Printf("\n\nERROR: Unknown Commnad[%s]\n\n", c.String("command"))

			return
		}

	}
	app.Run(os.Args)

}

func csvEntryPrinter(suppressionPage *sp.SuppressionPage, summary bool) {
	entries := suppressionPage.Results

	if summary {
		fmt.Printf("Recipient, Transactional, NonTransactional, Source, Updated, Created\n")
	} else {
		fmt.Printf("Recipient, Transactional, NonTransactional, Source, Updated, Created, Description\n")
	}

	for i := range entries {
		entry := entries[i]
		if summary {
			fmt.Printf("%s, %t, %t, %s, %s, %s\n", entry.Recipient, entry.Transactional, entry.NonTransactional, entry.Source, entry.Updated, entry.Created)
		} else {
			fmt.Printf("%s, %t, %t, %s,%s, %s, %s\n", entry.Recipient, entry.Transactional, entry.NonTransactional, entry.Source, entry.Updated, entry.Created, sanatize(entry.Description))
		}
	}
}

func sanatize(str string) string {

	return stripchars(str, ",\n\r")
}

func stripchars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}
