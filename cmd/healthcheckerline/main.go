package main

import (
	"fmt"
	yamlConfig "gohealthcheck/pkg/config/yaml"
	"gohealthcheck/pkg/healthchecker"
	"gohealthcheck/pkg/reader/csvreader"
	"log"
	"net/http"
	"os"

	"golang.org/x/xerrors"
)

func main() {

	args := os.Args[1:]
	if len(args) != 2 {
		err := xerrors.Errorf("main: %w", csvreader.ErrorMissingCSVFile)
		log.Fatal(err)
		return
	}

	config, err := yamlConfig.NewConfig(args[0])
	if err != nil {
		log.Fatal(err)
	}

	input := healthchecker.Input{
		InputType:   healthchecker.TypeCSV,
		CSVFileName: args[1],
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
