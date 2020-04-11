package command

import (
	"testing"
	"time"
)

func TestExecute(t *testing.T) {
	var dateCmd DateCommand
	inputCmd := []struct {
		command  string
		expected string
	}{
		{"/date", time.Now().Format("Monday, 2006/01/02")},
		{"/plop", ""},
	}

	for _, cmd := range inputCmd {
		result, _ := dateCmd.Execute(cmd.command)

		if result != cmd.expected {
			t.Errorf("DateCommand.Execute() was incorrect, got: %s, want: %s.", result, cmd.expected)
		}
	}
}
