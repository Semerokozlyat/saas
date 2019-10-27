package chrome

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const ChromeBinEnvName string = "CHROME_BIN"

func MakeScreenshot() (string, error) {
	chromeExecPath, found := os.LookupEnv(ChromeBinEnvName)
	if !found {
		log.Printf("Chrome executable path is not found.")
	}

	out, err := exec.Command(
		chromeExecPath,
		"--headless",
		"--disable-gpu",
		"--disable-software-rasterizer",
		"--disable-dev-shm-usage",
		"--no-sandbox",
		"--screenshot=new_screenshot1.png",
		"--disable-gpu",
		"--hide-scrollbars",
		"https://www.drom.ru").Output()
	//out, err := exec.Command("hostname").Output()
	if err != nil {
		return "", fmt.Errorf("command finished with error: %v", err)
	}
	return fmt.Sprintf("Command output is: %s", out), nil
}
