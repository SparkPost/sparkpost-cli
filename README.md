

## Getting Started

* `go get github.com/codegangsta/cli`
* `go get github.com/SparkPost/gosparkpost`


*NOTE:* this project currently depends on a custom version of [gosparkpost](https://github.com/SparkPost/gosparkpost) which if found in [this](https://github.com/yepher/gosparkpost/tree/feature/enhance_message_events) fork.



## Usage Examples

* `sparkpost-cli --from "2014-07-20T09:00"`



`sparkpost-cli help`

```
   --baseurl, -u "https://api.sparkpost.com"	Optional baseUrl for SparkPost. [$SPARKPOST_BASEURL]
   --apikey, -k 				Required SparkPost API key [$SPARKPOST_API_KEY]
   --bounce_classes, -b 			Optional comma-delimited list of bounce classification codes to search.
   --campaign_ids, -i 				Optional comma-delimited list of campaign ID's to search. Example: "Example Campaign Name"
   --events, -e 				Optional comma-delimited list of event types to search. Defaults to all event types.
   --friendly_froms 				Optional comma-delimited list of friendly_froms to search
   --from, -f 					Optional Datetime in format of YYYY-MM-DDTHH:MM. Example: 2016-02-10T08:00. Default: One hour ago
   --message_ids 				Optional Comma-delimited list of message ID's to search. Example: 0e0d94b7-9085-4e3c-ab30-e3f2cd9c273e.
   --page 					Optional results page number to return. Used with per_page for paging through result. Example: 25. Default: 1
   --per_page 					Optional number of results to return per page. Must be between 1 and 10,000 (inclusive). Example: 100. Default: 1000.
   --reason 					Optional bounce/failure/rejection reason that will be matched using a wildcard (e.g., %%reason%%). Example: bounce.
   --recipients 				Optional Comma-delimited list of recipients to search. Example: recipient@example.com
   --template_ids 				Optional Comma-delimited list of template ID's to search. Example: templ-1234.
   --timezone 					Optional Standard timezone identification string. Example: America/New_York. Default: UTC
   --to 					Optional Datetime in format of YYYY-MM-DDTHH:MM. Example: 2014-07-20T09:00. Default: now.
   --transmission_ids 				Optional Comma-delimited list of transmission ID's to search (i.e. id generated during creation of a transmission). Example: 65832150921904138.
   --help, -h					show help
   --version, -v				print the version


```