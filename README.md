# Go HealthCheck

A console program implemented in Go for checking the health status of websites given a CSV file input as argument.

## Prerequisites
- Go 1.14.3

## Setup

```sh
go build
```

The above commands will build the binary from the Go code.

## Configuration
Included in the current directory exist a file called `appsettings.yaml`. This file is use for specifying the configurations that are used in the program.

Currently there exist 3 configuration settings that you must specify:
- token - this is the API token received from Line Login.
- reportendpoint - the API endpoint for sending the health check result.
- timeout - the request time out value in seconds.

## Run

```sh
./gohealthcheck input.csv
```

Assuming that your input file which contain a list of web sites you want to check is called `input.csv`. The above command will run the program binary with the input file.