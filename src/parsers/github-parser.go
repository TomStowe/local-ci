package parsers

import (
	"fmt"
	"local-ci/src/models"
	"os"

	"gopkg.in/yaml.v2"
)

func ParseGitHubCI(filename string) (models.Pipeline, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return models.Pipeline{}, fmt.Errorf("failed to read GitHub CI file: %w", err)
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return models.Pipeline{}, fmt.Errorf("failed to parse GitHub CI YAML: %w", err)
	}

	var pipeline models.Pipeline
	if jobs, ok := config["jobs"].(map[interface{}]interface{}); ok {
		for jobName, jobValue := range jobs {
			jobNameStr, ok := jobName.(string)
			if !ok {
				continue
			}
			jobMap, ok := jobValue.(map[interface{}]interface{})
			if !ok {
				continue
			}
			var stage models.Stage
			stage.Name = jobNameStr
			if steps, ok := jobMap["steps"].([]interface{}); ok {
				for _, stepValue := range steps {
					stepMap, ok := stepValue.(map[interface{}]interface{})
					if !ok {
						continue
					}
					stepName, _ := stepMap["name"].(string)
					stepCommand, _ := stepMap["run"].(string)
					stage.Steps = append(stage.Steps, models.Step{Name: stepName, Command: stepCommand})
				}
			}
			pipeline.Stages = append(pipeline.Stages, stage)
		}
	}
	return pipeline, nil
}
