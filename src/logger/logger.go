// logger.go
package logger

import (
	"fmt"
	"local-ci/src/models"
	string_utils "local-ci/src/utils"
	"strings"

	"github.com/fatih/color"
)

const (
	maxCommandLength = 50 // Maximum length for command display
	totalWidth       = 60 // Total width of the table
)

// clearScreen clears the terminal screen.
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// createThickBorder generates a thick border line for the table.
func createThickBorder(length int) string {
	return strings.Repeat("═", length)
}

// formatCommand formats a command string with ellipsis if it's too long.
func formatCommand(command string) string {
	if len(command) > maxCommandLength {
		return command[:maxCommandLength-3] + "..."
	}
	return command
}

// formatStepStatus formats a command with its corresponding status in color.
func formatStepStatus(command string, status string, colorFunc func(a ...interface{}) string) string {
	formattedCommand := formatCommand(command)
	statusFormatted := colorFunc(status)
	// Adjust the command and status to fit within the total width
	return fmt.Sprintf("║ %-50s %s ║", formattedCommand, statusFormatted)
}

// PrintPipeline prints the current state of the pipeline in separate vertical tables.
func PrintPipeline(pipeline models.Pipeline, currentStage int, currentStep int, errOccurred bool, stagesToSkip []string) {
	clearScreen()

	// Iterate through each stage and build the table
	for i, stage := range pipeline.Stages {
		var tableBuilder strings.Builder
		skipStage := string_utils.In(stage.Name, stagesToSkip)

		stageName := stage.Name
		if skipStage {
			stageName += " (Skip)"
		}

		// Add the stage name header
		stageHeader := fmt.Sprintf("╔%s╗\n    %-20s    \n╠%s╣",
			createThickBorder(totalWidth), "Stage: "+stageName, createThickBorder(totalWidth))
		tableBuilder.WriteString(stageHeader + "\n")
		if skipStage {
			continue
		}

		// Add each step row to the table
		for j, step := range stage.Steps {
			status := "Pending"
			colorFunc := color.New(color.FgYellow).SprintFunc()

			if i < currentStage || (i == currentStage && j < currentStep) {
				status = "Success"
				colorFunc = color.New(color.FgGreen).SprintFunc()
			} else if errOccurred && i == currentStage && j == currentStep {
				status = "Failed "
				colorFunc = color.New(color.FgRed).SprintFunc()
			}

			// Add formatted step status to the table
			tableBuilder.WriteString(formatStepStatus(step.Command, status, colorFunc) + "\n")
		}

		// Add the bottom border to the table
		tableBuilder.WriteString(fmt.Sprintf("╚%s╝", createThickBorder(totalWidth)))

		// Print the table for the current stage
		fmt.Println(tableBuilder.String())
		fmt.Println() // Print an extra newline for separation between tables
	}
}

// PrintError displays the error message when a step fails.
func PrintError(stageName string, stepName string, err error) {
	errorColor := color.New(color.FgRed).SprintFunc()
	fmt.Printf("%s: Error in stage '%s', step '%s': %v\n",
		errorColor("ERROR"), stageName, stepName, err)
}
