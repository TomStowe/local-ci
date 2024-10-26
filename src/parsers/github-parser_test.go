package parsers

import (
	"os"
	"testing"

	"local-ci/src/models"
)

// TestParseGitHubCI tests the ParseGitHubCI function for various scenarios.
func TestParseGitHubCI(t *testing.T) {
	// Test case 1: Valid GitHub CI YAML
	validYAML := `
jobs:
  build:
    steps:
      - name: Checkout code
        run: echo "Checkout code"
      - name: Build
        run: echo "Building..."
  test:
    steps:
      - name: Run tests
        run: echo "Running tests"
  deploy:
    steps:
      - name: Deploy application
        run: echo "Deploying..."
`
	tempFile := writeTempYAML(t, validYAML)
	defer os.Remove(tempFile)

	pipeline, err := ParseGitHubCI(tempFile)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedStages := []models.Stage{
		{Name: "build", Steps: []models.Step{
			{Name: "Checkout code", Command: "echo \"Checkout code\""},
			{Name: "Build", Command: "echo \"Building...\""},
		}},
		{Name: "test", Steps: []models.Step{
			{Name: "Run tests", Command: "echo \"Running tests\""},
		}},
		{Name: "deploy", Steps: []models.Step{
			{Name: "Deploy application", Command: "echo \"Deploying...\""},
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
jobs:
  build:
    steps: 
      - name: Checkout code
        run: echo "Checkout code"
      - name: Build
        run: echo "Building..."
  test
    steps:
      - name: Run tests
        run: echo "Running tests"
`
	tempFile = writeTempYAML(t, invalidYAML)
	defer os.Remove(tempFile)

	_, err = ParseGitHubCI(tempFile)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}

	// Test case 3: File read error
	_, err = ParseGitHubCI("non_existent_file.yml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
