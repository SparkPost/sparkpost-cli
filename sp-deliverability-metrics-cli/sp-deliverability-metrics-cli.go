package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"

	sp "github.com/SparkPost/gosparkpost"
)

func main() {

	validParameters := []string{
		"from", "to", "domains", "campaigns", "templates", "nodes", "bindings",
		"binding_groups", "protocols", "metrics", "timezone", "limit", "order_by", "subaccounts",
	}

	app := cli.NewApp()

	app.Version = "0.0.2"
	app.Name = "sparkpost-message-event-cli"
	app.Usage = "SparkPost Message Event CLI\n\n\tSee https://developers.sparkpost.com/api/metrics.html"
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

		// Metrics Parameters
		cli.StringFlag{
			Name:  "command",
			Value: "domain",
			Usage: "Optional one of domain, binding, binding-group, campaign, template, watched-domain, time-series, bounce-reason, bounce-reason/domain, bounce-classification, rejection-reason, rejection-reason/domain, delay-reason, delay-reason/domain, link-name, attempt",
		},
		cli.StringFlag{
			Name:  "from, f",
			Value: "",
			Usage: "Required Datetime in format of YYYY-MM-DDTHH:MM. Example: 2016-02-10T08:00. Default: One hour ago",
		},
		cli.StringFlag{
			Name:  "to",
			Value: "",
			Usage: "Optional Datetime in format of YYYY-MM-DDTHH:MM. Example: 2016-02-10T00:00. Default: now.",
		},
		cli.StringFlag{
			Name:  "domains, d",
			Value: "",
			Usage: "Optional Comma-delimited list of domains to include Example: gmail.com,yahoo.com,hotmail.com.",
		},
		cli.StringFlag{
			Name:  "campaigns, c",
			Value: "",
			Usage: "Optional Comma-delimited list of campaigns to include. Example: Black Friday",
		},
		cli.StringFlag{
			Name:  "metrics, m",
			Value: "count_injected,count_bounce,count_rejected,count_delivered,count_delivered_first,count_delivered_subsequent,total_delivery_time_first,total_delivery_time_subsequent,total_msg_volume,count_policy_rejection,count_generation_rejection,count_generation_failed,count_inband_bounce,count_outofband_bounce,count_soft_bounce,count_hard_bounce,count_block_bounce,count_admin_bounce,count_undetermined_bounce,count_delayed,count_delayed_first,count_rendered,count_unique_rendered,count_unique_confirmed_opened,count_clicked,count_unique_clicked,count_targeted,count_sent,count_accepted,count_spam_complaint",
			Usage: "Required Comma-delimited list of metrics for filtering",
		},
		cli.StringFlag{
			Name:  "templates",
			Value: "",
			Usage: "Optioanl comma-delimited list of template IDs to include Example: summer-sale",
		},
		cli.StringFlag{
			Name:  "nodes",
			Value: "",
			Usage: "Optional comma-delimited list of nodes to include ( Note: SparkPost Elite only ) Example: Email-MSys-1,Email-MSys-2,Email-MSys-3",
		},
		cli.StringFlag{
			Name:  "bindings",
			Value: "",
			Usage: "Optional comma-delimited list of bindings to include (Note: SparkPost Elite only) Example: Confirmation",
		},
		cli.StringFlag{
			Name:  "binding_groups",
			Value: "",
			Usage: "Optional comma-delimited list of binding groups to include (Note: SparkPost Elite only) Example: Transaction",
		},
		cli.StringFlag{
			Name:  "protocols",
			Value: "",
			Usage: "Optional comma-delimited list of protocols for filtering (Note: SparkPost Elite only) Example: smtp",
		},
		cli.StringFlag{
			Name:  "timezone",
			Value: "",
			Usage: "Standard timezone identification string, defaults to UTC Example: America/New_York.",
		},
		cli.StringFlag{
			Name:  "limit",
			Value: "",
			Usage: "Optional maximum number of results to return Example: 5",
		},
		cli.StringFlag{
			Name:  "order_by",
			Value: "",
			Usage: "Optional metric by which to order results Example: count_injected",
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

		m := &sp.Metrics{}
		m.Params = parameters

		e, err := client.QueryMetrics(m)

		if err != nil {
			log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
			return
		} else if e.Errors != nil {
			log.Fatalf("ERROR: %q.\n\nFor additional information try using `--verbose true`\n\n\n", e.Errors)
			return
		} else {

			metrics := c.String("metrics")
			log.Printf(metrics)
			fields := strings.Split(metrics, ",")

			// TODO: add an HTML output
			csvHeaderPrinter(fields)

			//log.Printf("DUMP: %q\n", m.Results)
			for _, element := range m.Results {
				csvEntryPrinter(fields, c.String("command"), &element)
			}
		}
	}
	app.Run(os.Args)

}

func csvEntryPrinter(fields []string, command string, metricItem *sp.MetricItem) {
	row := ""

	switch command {
	case "domain":
		row = fmt.Sprintf("%s%s, ", row, metricItem.Domain)
	case "campaign":
		row = fmt.Sprintf("%s%s, ", row, metricItem.CampaignId)
	case "template":
		row = fmt.Sprintf("%s%s, ", row, metricItem.TemplateId)
	case "time-series":
		row = fmt.Sprintf("%s%s, ", row, metricItem.TimeStamp)
	case "watched-domain":
		row = fmt.Sprintf("%s%s, ", row, metricItem.WatchedDomain)
	case "binding":
		row = fmt.Sprintf("%s%s, ", row, metricItem.Binding)
	case "binding-group":
		row = fmt.Sprintf("%s%s, ", row, metricItem.BindingGroup)
	default:
		row = fmt.Sprintf("%sUnknown Commnad[%s], ", row, command)
	}

	for i := range fields {
		switch fields[i] {
		case "count_injected":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountInjected)
		case "count_bounce":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountBounce)
		case "count_rejected":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountRejected)
		case "count_delivered":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountDelivered)
		case "count_delivered_first":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountDeliveredFirst)
		case "count_delivered_subsequent":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountDeliveredSubsequent)
		case "total_delivery_time_first":
			row = fmt.Sprintf("%s%d, ", row, metricItem.TotalDeliveryTimeFirst)
		case "total_delivery_time_subsequent":
			row = fmt.Sprintf("%s%d, ", row, metricItem.TotalDeliveryTimeSubsequent)
		case "total_msg_volume":
			row = fmt.Sprintf("%s%d, ", row, metricItem.TotalMsgVolume)
		case "count_policy_rejection":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountPolicyRejection)
		case "count_generation_rejection":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountGenerationRejection)
		case "count_generation_failed":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountGenerationFailed)
		case "count_inband_bounce":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountInbandBounce)
		case "count_outofband_bounce":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountOutofbandBounce)
		case "count_soft_bounce":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountSoftBounce)
		case "count_hard_bounce":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountHardBounce)
		case "count_block_bounce":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountBlockBounce)
		case "count_admin_bounce":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountAdminBounce)
		case "count_undetermined_bounce":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountUndeterminedBounce)
		case "count_delayed":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountDelayed)
		case "count_delayed_first":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountDelayedFirst)
		case "count_rendered":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountRendered)
		case "count_unique_rendered":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountUniqueRendered)
		case "count_unique_confirmed_opened":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountUniqueConfirmedOpened)
		case "count_clicked":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountClicked)
		case "count_unique_clicked":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountUniqueClicked)
		case "count_targeted":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountTargeted)
		case "count_sent":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountSent)
		case "count_accepted":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountAccepted)
		case "count_spam_complaint":
			row = fmt.Sprintf("%s%d, ", row, metricItem.CountSpamComplaint)
		default:
			row = fmt.Sprintf("unknown field: Invalid Field")

		}

	}

	fmt.Println(row)
}

func csvHeaderPrinter(fields []string) {
	row := "domain, "
	for i := range fields {
		row = fmt.Sprintf("%s%s, ", row, fields[i])
	}

	fmt.Println(row)
}
