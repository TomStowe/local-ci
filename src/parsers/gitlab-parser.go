package parsers

import (
	"fmt"
	"local-ci/src/models"
	"os"

	"gopkg.in/yaml.v2"
)

// ParseGitLabCI parses a GitLab CI configuration file and returns a Pipeline struct
func ParseGitLabCI(filename string) (models.Pipeline, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return models.Pipeline{}, fmt.Errorf("failed to read GitLab CI file: %w", err)
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return models.Pipeline{}, fmt.Errorf("failed to parse GitLab CI YAML: %w", err)
	}

	var pipeline models.Pipeline
	stagesMap := make(map[string]models.Stage)

	// Iterate over jobs in the GitLab CI configuration
	for jobName, jobValue := range config {
		jobMap, ok := jobValue.(map[interface{}]interface{})
		if !ok {
			continue
		}

		// Get the stage name, defaulting to "default" if not specified
		stageName := "default"
		if stage, ok := jobMap["stage"].(string); ok {
			stageName = stage
		}

		// Create a new stage if it doesn't already exist
		stage, exists := stagesMap[stageName]
		if !exists {
			stage = models.Stage{Name: stageName}
		}

		// Get the script commands
		if script, ok := jobMap["script"].([]interface{}); ok {
			for _, step := range script {
				command, ok := step.(string)
				if ok {
					stage.Steps = append(stage.Steps, models.Step{
						Name:    jobName,
						Command: command,
					})
				}
			}
		}

		pipeline.Stages = append(pipeline.Stages, stage)
	}

	return pipeline, nil
}
