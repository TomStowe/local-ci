package main

import (
	"flag"
	"fmt"
	"local-ci/src/executor"
	"local-ci/src/models"
	"local-ci/src/parsers"
	"os"
	"strings"
)

func main() {
	ciType := flag.String("type", "gitlab", "Specify the CI type: gitlab or github")
	file := flag.String("file", "", "Path to the CI configuration file")
	skipStages := flag.String("skip", "", "Comma-separated list of stages to skip")
	flag.Parse()

	if *file == "" {
		fmt.Println("Please specify a CI configuration file using the -file flag.")
		os.Exit(1)
	}

	// Parse the steps to skip
	stagesToSkip := []string{}
	if *skipStages != "" {
		for _, stageName := range strings.Split(*skipStages, ",") {
			stagesToSkip = append(stagesToSkip, stageName)
		}
	}

	var pipeline models.Pipeline
	var err error

	switch *ciType {
	case "gitlab":
		pipeline, err = parsers.ParseGitLabCI(*file)
	case "github":
		pipeline, err = parsers.ParseGitHubCI(*file)
	default:
		fmt.Println("Unsupported CI type. Please use 'gitlab' or 'github'.")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error parsing CI configuration: %v\n", err)
		os.Exit(1)
	}

	if err := executor.RunPipeline(pipeline, stagesToSkip); err != nil {
		fmt.Printf("Pipeline failed: %v\n", err)
		os.Exit(1)
	}
}
