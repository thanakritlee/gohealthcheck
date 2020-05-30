package main

import (
	"fmt"
	"gohealthcheck/healthchecker"
	"gohealthcheck/reader/csvreader"
	yamlConfig "gohealthcheck/utility/config/yaml"
	"log"
	"net/http"
	"os"

	"golang.org/x/xerrors"
)

func main() {
	config, err := yamlConfig.NewConfig("appsettings.yaml")
	if err != nil {
		log.Fatal(err)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		err := xerrors.Errorf("main: %w", csvreader.ErrorMissingCSVFile)
		log.Fatal(err)
		return
	}

	input := healthchecker.Input{
		InputType:   healthchecker.TypeCSV,
		CSVFileName: args[0],
	}

	headers := http.Header{}
	headers.Add("Authorization", fmt.Sprintf("Bearer %s", config.GetConfig("token")))

	report := []healthchecker.Report{
		{
			ReportType: healthchecker.TypeHTTP,
			URL:        config.GetConfig("reportendpoint"),
			Headers:    headers,
		},
		{
			ReportType: healthchecker.TypeConsole,
		},
	}

	checker, err := healthchecker.NewHealthChecker(&input, &report)
	if err != nil {
		log.Fatal(err)
	}

	err = checker.CheckHealth()
	if err != nil {
		log.Fatal(err)
	}

	err = checker.ReportResult()
	if err != nil {
		log.Fatal(err)
	}
}
