<a href="https://www.sparkpost.com"><img src="https://www.sparkpost.com/sites/default/files/attachments/SparkPost_Logo_2-Color_Gray-Orange_RGB.svg" width="200px"/></a>

[Sign up](https://app.sparkpost.com/sign-up?src=Dev-Website&sfdcid=70160000000pqBb) for a SparkPost account and visit our [Developer Hub](https://developers.sparkpost.com) for even more content.

# SparkPost CLI

[![Slack Status](http://slack.sparkpost.com/badge.svg)](http://slack.sparkpost.com)

This is the official SparkPost Command Line Interface (CLI) for the [SparkPost API](https://www.sparkpost.com/api).

## Contributing to sparkpost-cli

Transparency is one of our core values, and we encourage developers to contribute and become part of the SparkPost developer community.

The following is a set of guidelines for contributing to sparkpost-cli, which is hosted in the [SparkPost Organization](https://github.com/sparkpost) on GitHub. These are just guidelines, not rules, use your best judgment and feel free to propose changes to this document in a pull request.

## Submitting Issues

* Before logging an issue, please [search existing issues](https://github.com/SparkPost/sparkpost-cli/issues?q=is%3Aissue+is%3Aopen) first.

* You can create an issue [here](https://github.com/SparkPost/sparkpost-cli/issues/new).  Please include the version number and the Git hash with as much detail as possible in your report.

## Local Development

1. Fork this repo
1. Clone your fork
1. Write some code!
1. Please follow the pull request submission steps in the next section

## Contribution Steps

To contribute to sparkpost-cli:

1. Create a new branch named after the issue you’ll be fixing (include the issue number as the branch name, example: Issue in GH is #8 then the branch name should be _feature/ISSUE-8_))
1. Write corresponding tests and code (only what is needed to satisfy the issue and please test)
    * Include tests in the 'test' directory in an appropriate test file
    * Write code to satisfy the tests
1. Ensure any automated tests pass
1. Submit a new Pull Request applying your feature/fix branch to the `master` branch

## Building

* `go get github.com/urfave/cli`
* `go get github.com/SparkPost/gosparkpost`
* change to the cli tool you want to build
	* `go build`


### Releasing

When a new release is created the script in `tools/bin/package.sh` will be run to create target binaries.
