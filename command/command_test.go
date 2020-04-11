package command

import (
	"testing"
	"time"
)

func TestExecute(t *testing.T) {
	var dateCmd DateCommand

	outputDateCmd, _ := dateCmd.Execute("/date")

	if outputDateCmd != time.Now().Format("Monday, 2006/01/02") {
		t.Errorf("DateCommand.Execute() was incorrect, got: %s, want: %s.", outputDateCmd, time.Now().Format("Monday, 2006/01/02"))
	}
}
