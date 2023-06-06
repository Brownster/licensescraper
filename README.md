LicenseScraper

LicenseScraper is a tool that scrapes license data from XML files, specifically UsageHistory.xml, and provides metrics for the usage and status of licenses. The tool is implemented in Go and uses the maas library for scheduling and executing scraping tasks.
Features

    Scheduled scraping: The tool performs scraping tasks every 12 hours, with a timeout of 5 seconds for each task​1​.
    Flexible data source: The tool can be pointed to any directory containing UsageHistory.xml files​1​.
    Comprehensive metrics: The tool provides a wide range of metrics related to license usage and status, including the peak used licenses, total available licenses, feature name, feature display name, feature value, feature expiration date, and the remaining time until the license expires​1​.

Usage

LicenseScraper requires the path to the data folder containing UsageHistory.xml files as an argument. The argument should be passed with the flag license.datapath:

css

<executable> --license.datapath <path to data folder>

Please replace <executable> with the actual name of the compiled executable and <path to data folder> with the path to the data folder containing UsageHistory.xml files.
Building from Source

To build LicenseScraper from source, you will need a working Go development environment.

bash

go get github.com/Brownster/licensescraper
cd $GOPATH/src/github.com/Brownster/licensescraper
go build

This will create an executable in the current directory that you can use as described in the Usage section.
