package chrome

import (
	"fmt"
	"os/exec"
)

func MakeScreenshot() (string, error) {
	out, err := exec.Command("hostname").Output()
	if err != nil {
		return "", fmt.Errorf("command finished with error: %v", err)
	}
	return fmt.Sprintf("Command output is: %s", out), nil
}
