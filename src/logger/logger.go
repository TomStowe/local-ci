// logger.go
package logger

import (
	"fmt"
	"local-ci/src/models"
	string_utils "local-ci/src/utils"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func formatStatusHeader(stage models.Stage, stagesToSkip []string, stageToPrint, currentStage int, errOccurred bool) (*strings.Builder, bool) {
	var tableBuilder strings.Builder
	skipStage := string_utils.In(stage.Name, stagesToSkip)
	stageName := stage.Name
	colorFunc := color.New(color.FgWhite).SprintFunc()
	if skipStage {
		return &tableBuilder, true
	} else if stageToPrint < currentStage {
		colorFunc = color.New(color.FgGreen).SprintFunc()
		stageName += " (Success)"
		skipStage = true
	} else if stageToPrint > currentStage {
		colorFunc = color.New(color.FgCyan).SprintFunc()
		stageName += " (Waiting)"
		skipStage = true
	} else if errOccurred {
		colorFunc = color.New(color.FgRed).SprintFunc()
		stageName += " (Failed)"
	}

	// Add the stage name header
	headerColumnStart := "╠"
	headerColumnEnd := "╣"
	if skipStage {
		headerColumnStart = "╚"
		headerColumnEnd = "╝"
	}

	stageHeader := fmt.Sprintf(
		"╔%s╗\n║   %-63s   ║\n%s%s%s",
		createThickBorder(totalWidth),
		colorFunc(cases.Title(language.English, cases.Compact).String(stageName)),
		headerColumnStart,
		createThickBorder(totalWidth),
		headerColumnEnd)
	tableBuilder.WriteString(stageHeader + "\n")
	return &tableBuilder, skipStage
}

// PrintPipeline prints the current state of the pipeline in separate vertical tables.
func PrintPipeline(pipeline models.Pipeline, currentStage int, currentStep int, errOccurred bool, stagesToSkip []string) {
	clearScreen()

	// Iterate through each stage and build the table
	for i, stage := range pipeline.Stages {
		tableBuilder, skipTableBody := formatStatusHeader(stage, stagesToSkip, i, currentStage, errOccurred)

		if skipTableBody {
			fmt.Println(tableBuilder.String())
			fmt.Println()
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
			tableBuilder.WriteString(formatStepStatus(step.Name, status, colorFunc) + "\n")
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
