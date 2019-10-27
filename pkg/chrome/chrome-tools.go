package chrome

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const ChromeBinEnvName string = "CHROME_BIN"
const ScreensDestEnvName string = "SCREENS_DEST"

func MakeScreenshot() (string, error) {
	chromeExecPath, found := os.LookupEnv(ChromeBinEnvName)
	if !found {
		log.Printf("Chrome executable path is not found.")
	}

	screensDestination, found := os.LookupEnv(ScreensDestEnvName)
	if !found {
		log.Printf("Screens destination path is not found.")
	}

	resultedCommand := strings.Join([]string{chromeExecPath,
		"--headless",
		"--disable-gpu",
		"--disable-software-rasterizer",
		"--disable-dev-shm-usage",
		"--no-sandbox",
		fmt.Sprintf("--screenshot=%s/new_screenshot1.png", screensDestination),
		"--disable-gpu",
		"--hide-scrollbars",
		"https://www.drom.ru"}, " ")

	log.Printf("Going to exec command: %s", resultedCommand)

	out, err := exec.Command(resultedCommand).Output()
	//out, err := exec.Command("hostname").Output()
	if err != nil {
		return "", fmt.Errorf("command finished with error: %v", err)
	}
	return fmt.Sprintf("Command output is: %s", out), nil
}
