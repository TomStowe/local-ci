package logger

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/fatih/color"
)

func TestFormatCommand(t *testing.T) {
	longCommand := "echo This is a long command that exceeds the maximum length"
	shortCommand := "echo Short command"

	if got := formatCommand(longCommand); got != "echo This is a long command that exceeds the ma..." {
		t.Errorf("formatCommand() = %v, want %v", got, "echo This is a long command that exceeds the ma...")
	}

	if got := formatCommand(shortCommand); got != shortCommand {
		t.Errorf("formatCommand() = %v, want %v", got, shortCommand)
	}
}

func TestFormatStepStatus(t *testing.T) {
	buf := new(bytes.Buffer)
	color.Output = buf // Capture the output

	command := "echo Hello World"
	status := "Success"
	colorFunc := color.New(color.FgGreen).SprintFunc()

	got := formatStepStatus(command, status, colorFunc)
	expected := fmt.Sprintf("║ %-50s %s ║", command, colorFunc(status))

	if got != expected {
		t.Errorf("formatStepStatus() = %v, want %v", got, expected)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
