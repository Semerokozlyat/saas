package chrome

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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

	cmd := exec.Command(chromeExecPath,
		"--headless",
		"--disable-gpu",
		"--disable-software-rasterizer",
		"--disable-dev-shm-usage",
		"--no-sandbox",
		fmt.Sprintf("--screenshot=%s/new_screenshot1.png", screensDestination),
		"--hide-scrollbars",
		"https://www.stopgame.ru")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command finished with error: %v", err)
	}
	return fmt.Sprintf("Command output is: %s", cmd.Stdout), nil
}
