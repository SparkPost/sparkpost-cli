## Importing SendGrid Suppression List

SparkPost CLI can import your SendGrid suppression List(s). However, it has to be in specific format.
Read on to learn how to export suppressions from SendGrid and format them so that they're importable by SparkPost.

**TL;DR** The CLI can import CSV with the following fields:
 - `email` - first column
 - `created` - second column


### Import Command
Regardless of suppression type and preparation we need, we will need to run the following command as the final step to import the suppressions to SparkPost.

```
sp-suppression-list-cli --command sendgrid -f PATH_TO_SENDGRID_EXPORT.csv
```
Note: Replace `PATH_TO_SENDGRID_EXPORT.csv` with your CSV file. Check **Environment** section in main [readme](Usage.md) for setting up API Key and API Path.

----------

### Unsubscribes

Sendgrid's unsubscribe suppression list can be imported to SparkPost without any modification.

- [Export from SendGrid](https://sendgrid.com/docs/User_Guide/Suppressions/advanced_suppression_manager.html#-Export-an-Unsubscribe-Group-List).
- Once the CSV file is downloaded, run the following command
- Run [Import command](#import-command).


### Bounces
We need to make some modification to the bounce exports prior to be imported by SparkPost CLI.

- [Export bounces from SendGrid](https://sendgrid.com/docs/User_Guide/Suppressions/bounces.html#-Download-Bounces-as-CSV).
- Open the downloaded CSV with a spreadsheet tool like MS Excel.
- Remove **status** and **reason** column and save file.
- Run [Import command](#import-command).

### Invalid Emails
We need to make some modification to the invalid emails prior to be imported by SparkPost CLI.

- [Export invalid emails from SendGrid](https://sendgrid.com/docs/User_Guide/Suppressions/invalid_emails.html#-Download-Invalid-Emails-as-CSV).
- Open the downloaded CSV with a spreadsheet tool like MS Excel.
- Remove **reason** column and save file.
- Run [Import command](#import-command).

### Spams
We need to make some modification to the spam reports prior to be imported by SparkPost CLI.

- [Export spam reports from SendGrid](https://sendgrid.com/docs/User_Guide/Suppressions/spam_reports.html#-Download-Spam-Reports-as-CSV).
 - Open the downloaded CSV with a spreadsheet tool like MS Excel.
 - Remove **ip** column and save file.
 - Run [Import command](#import-command).

### Others
- Export the list from SendGrid in CSV.
- Open the downloaded CSV with a spreadsheet tool like MS Excel.
- Remove any columns except **email** and **created**. Also make sure **email** is first column and **created** is second. Save file.
- Run [Import command](#import-command).
