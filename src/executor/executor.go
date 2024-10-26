package executor

import (
	"fmt"
	"local-ci/src/logger"
	"local-ci/src/models"
	string_utils "local-ci/src/utils"
	"os"
	"os/exec"
	"time"
)

// RunPipeline executes the stages and steps of a given pipeline.
// It handles printing the progress using the logger functions.
func RunPipeline(pipeline models.Pipeline, stagesToSkip []string) error {
	for i, stage := range pipeline.Stages {
		// Skip a stage if needed
		if string_utils.In(stage.Name, stagesToSkip) {
			continue
		}

		for j, step := range stage.Steps {
			// Print the current state of the pipeline
			logger.PrintPipeline(pipeline, i, j, false, stagesToSkip)

			// Execute the step command
			fmt.Printf("Running: %s\n", step.Command)
			cmd := exec.Command("sh", "-c", step.Command)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// Run the command and check for errors
			if err := cmd.Run(); err != nil {
				// If there's an error, print the pipeline state and error
				logger.PrintPipeline(pipeline, i, j, true, stagesToSkip)
				logger.PrintError(stage.Name, step.Name, err) // Print the error using logger
				return fmt.Errorf("error in stage '%s', step '%s': %v", stage.Name, step.Name, err)
			}

			// Small delay to simulate progress
			time.Sleep(500 * time.Millisecond)
		}
	}

	logger.PrintPipeline(pipeline, len(pipeline.Stages), 1, false, stagesToSkip)
	return nil
}
