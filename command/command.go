package command

import (
	"fmt"
	"strings"
	"time"
)

type DateCommand struct{}

func (q DateCommand) Execute(s string) (string, error) {
	if strings.Index(s, "/date") == 0 {
		return time.Now().Format("Monday, 2006/01/02"), nil
	}

	return "", fmt.Errorf("this is not me")
}
