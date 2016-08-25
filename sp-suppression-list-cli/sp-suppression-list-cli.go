package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	sp "github.com/SparkPost/gosparkpost"
	"github.com/codegangsta/cli"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Column mapping for Mandrill Blacklist
const (
	MANDRILL_EMAIL_COL       = 0
	MANDRILL_REASON_COL      = 1
	MANDRILL_DETAIL_COL      = 2
	MANDRILL_CREATED_COL     = 3
	MANDRILL_EXPIRES_AT_COL  = 4
	MANDRILL_LAST_EVENT_COL  = 5
	MANDRILL_EXPIRES_AT2_COL = 6
	MANDRILL_SUBACCOUNT_COL  = 7
)

// Column mapping for SendGrid Blacklist
const (
	SENDGRID_EMAIL_COL = 0
	SENDGRID_CREATED   = 1
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	VALID_PARAMETERS := []string{
		"from", "to", "types", "limit",
	}

	app := cli.NewApp()

	app.Version = "0.0.1"
	app.Name = "suppression-sparkpost-cli"
	app.Usage = "SparkPost suppression list cli"
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
			Usage: "Mandrill blocklist CSV. See https://mandrill.zendesk.com/hc/en-us/articles/205582997",
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
			Name:  "batchsize",
			Value: "100000",
			Usage: "Optional the max suppression records to upload per request. Default value is 100000",
		},
		cli.StringFlag{
			Name:  "delay",
			Value: "0",
			Usage: "Optional how long to wait in seconds between suppression list requests",
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

		batchSize, _ := strconv.ParseInt(c.String("batchsize"), 10, 0)
		batchDelay, _ := strconv.ParseInt(c.String("delay"), 10, 0)

		var client sp.Client
		err := client.Init(cfg)
		if err != nil {
			log.Fatalf("SparkPost client init failed: %s\n", err)
			return
		}

		parameters := make(map[string]string)

		for i, val := range VALID_PARAMETERS {

			if c.String(VALID_PARAMETERS[i]) != "" {
				parameters[val] = c.String(val)
			}
		}

		switch c.String("command") {
		case "list":
			e, err := client.SuppressionList()

			if err != nil {
				log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
				return
			} else {
				csvEntryPrinter(e, true)

			}
		case "retrieve":
			recpipient := c.String("recipient")
			if recpipient == "" {
				log.Fatalf("ERROR: The `retrieve` command requires a recipient.")
				return
			}

			e, err := client.SuppressionRetrieve(recpipient)

			if err != nil {
				log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
				return
			} else {
				csvEntryPrinter(e, false)

			}
		case "search":
			parameters := make(map[string]string)

			for i, val := range VALID_PARAMETERS {

				if c.String(VALID_PARAMETERS[i]) != "" {
					parameters[val] = c.String(val)
				}
			}

			e, err := client.SuppressionSearch(parameters)

			if err != nil {
				log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
				return
			} else {
				csvEntryPrinter(e, true)

			}
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
			} else {
				fmt.Println("OK")

			}
		case "mandrill":
			fmt.Printf("Processing: %s\n", c.String("file"))
			file := c.String("file")
			if file == "" {
				log.Fatalf("ERROR: The `mandrill` command requires a CSV file.")
				return
			}

			f, err := os.Open(file)
			check(err)

			var entries = []sp.SuppressionEntry{}

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

				if record[MANDRILL_EMAIL_COL] == "email" {
					// Skip over header row
					continue
				}

				if record[MANDRILL_REASON_COL] != "hard-bounce" &&
					record[MANDRILL_REASON_COL] != "unsub" &&
					record[MANDRILL_REASON_COL] != "spam" &&
					record[MANDRILL_REASON_COL] != "custom" {
					// Ignore soft-bounce
					continue
				}

				if strings.Count(record[MANDRILL_EMAIL_COL], "@") != 1 {
					fmt.Printf("WARN: Ignoring '%s'. It is not a valid email address.\n", record[MANDRILL_EMAIL_COL])
					continue
				}

				entry := sp.SuppressionEntry{}

				if record[MANDRILL_EMAIL_COL] == "" {
					// Must have email as it is suppression list primary key
					continue
				}

				entry.Email = record[MANDRILL_EMAIL_COL]
				entry.Transactional = false
				entry.NonTransactional = true
				entry.Description = fmt.Sprintf("MBL: %s", record[MANDRILL_DETAIL_COL])

				entries = append(entries, entry)

				if len(entries) > int(batchSize) {
					fmt.Printf("Uploading batch %d\n", batchCount)
					err = client.SuppressionInsertOrUpdate(entries)

					if err != nil {
						log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
						return
					}
					entries = []sp.SuppressionEntry{}
					batchCount++
					if batchDelay > 0 {
						fmt.Printf("Sleeping")
						time.Sleep( time.Duration(batchDelay) * time.Second)
						fmt.Printf("Awake")
					}
				}
			}

			if len(entries) > 0 {
				fmt.Printf("Uploading batch %d\n", batchCount)
				err = client.SuppressionInsertOrUpdate(entries)

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

			var entries = []sp.SuppressionEntry{}

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

				if record[SENDGRID_EMAIL_COL] == "email" {
					// Skip over header row
					continue
				}

				entry := sp.SuppressionEntry{}

				if record[SENDGRID_EMAIL_COL] == "" {
					// Must have email as it is suppression list primary key
					continue
				}

				// SendGrid suppression lists are very dirty and tend to have invalid data. Some examples of invalid addresses are:
				// 	#02232014, gmail.com, To, 8/27/2015, name@yahoo.comett@domain.com"
				if strings.Count(record[MANDRILL_EMAIL_COL], "@") != 1 {
					fmt.Printf("WARN: Ignoring '%s'. It is not a valid email address.\n", record[MANDRILL_EMAIL_COL])
					continue
				}

				entry.Email = record[SENDGRID_EMAIL_COL]
				entry.Transactional = false
				entry.NonTransactional = true
				entry.Description = fmt.Sprintf("SBL: imported from SendGrid")

				entries = append(entries, entry)

				if len(entries) > (1024 * 100) {
					fmt.Printf("Uploading batch %d\n", batchCount)
					err = client.SuppressionInsertOrUpdate(entries)

					if err != nil {
						log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
						return
					}
					entries = []sp.SuppressionEntry{}
					batchCount++
				}

			}

			if len(entries) > 0 {
				fmt.Printf("Uploading batch %d\n", batchCount)
				err = client.SuppressionInsertOrUpdate(entries)

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

func csvEntryPrinter(suppressionList *sp.SuppressionListWrapper, summary bool) {
	entries := suppressionList.Results

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
