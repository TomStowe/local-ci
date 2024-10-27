package parsers

import (
	"os"
	"testing"

	"github.com/TomStowe/local-ci/src/models"
)

// Helper function to write a temporary YAML file for testing
func writeTempYAML(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "*.yml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	return tmpFile.Name()
}

// TestParseGitLabCI tests the ParseGitLabCI function for various scenarios.
func TestParseGitLabCI(t *testing.T) {
	// Test case 1: Valid GitLab CI YAML
	validYAML := `
steps:
  - build
  - test
  - deploy


build:
  stage: build
  script:
    - echo "Building..."
test:
  stage: test
  script:
    - echo "Running tests"
deploy:
  stage: deploy
  script:
    - echo "Deploying application"
`
	tempFile := writeTempYAML(t, validYAML)
	defer os.Remove(tempFile)

	pipeline, err := ParseGitLabCI(tempFile)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedStages := []models.Stage{
		{Name: "build", Steps: []models.Step{
			{Name: "echo \"Building...\"", Command: "echo \"Building...\""},
		}},
		{Name: "test", Steps: []models.Step{
			{Name: "echo \"Running tests\"", Command: "echo \"Running tests\""},
		}},
		{Name: "deploy", Steps: []models.Step{
			{Name: "echo \"Deploying application\"", Command: "echo \"Deploying application\""},
		}},
	}

	if len(pipeline.Stages) != len(expectedStages) {
		t.Fatalf("Expected %d stages, got %d", len(expectedStages), len(pipeline.Stages))
	}

	for i, stage := range expectedStages {
		if pipeline.Stages[i].Name != stage.Name {
			t.Errorf("Expected stage name %s, got %s", stage.Name, pipeline.Stages[i].Name)
		}
		if len(pipeline.Stages[i].Steps) != len(stage.Steps) {
			t.Fatalf("Stage %s: Expected %d steps, got %d", stage.Name, len(stage.Steps), len(pipeline.Stages[i].Steps))
		}
		for j, step := range stage.Steps {
			if pipeline.Stages[i].Steps[j].Name != step.Name {
				t.Errorf("Stage %s, Step %d: Expected step name %s, got %s", stage.Name, j, step.Name, pipeline.Stages[i].Steps[j].Name)
			}
			if pipeline.Stages[i].Steps[j].Command != step.Command {
				t.Errorf("Stage %s, Step %d: Expected command %s, got %s", stage.Name, j, step.Command, pipeline.Stages[i].Steps[j].Command)
			}
		}
	}

	// Test case 2: Invalid YAML
	invalidYAML := `
build:
  stage: build
  script: 
    - echo "Building..."
test
  stage: test
  script:
    - echo "Running tests"
`
	tempFile = writeTempYAML(t, invalidYAML)
	defer os.Remove(tempFile)

	_, err = ParseGitLabCI(tempFile)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}

	// Test case 3: File read error
	_, err = ParseGitLabCI("non_existent_file.yml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
