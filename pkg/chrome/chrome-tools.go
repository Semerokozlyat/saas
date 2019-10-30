package chrome

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const ChromeBinEnvName string = "CHROME_BIN"
const ScreensDestEnvName string = "SCREENS_DEST"


// Makes screenshot and returns file name, image as byte array
func MakeScreenshot(websiteURL string) (string, []byte, error) {
	chromeExecPath, found := os.LookupEnv(ChromeBinEnvName)
	if !found {
		return "", nil, fmt.Errorf("chrome executable path is not found")
	}

	fileDestinationPath, found := os.LookupEnv(ScreensDestEnvName)
	if !found {
		return "", nil, fmt.Errorf("screenshot destination path variable is not found")
	}

	fileName := strings.Join([]string{
		strconv.Itoa(int(time.Now().Unix())),
		".png",
	}, "")

	cmd := exec.Command(chromeExecPath,
		"--headless",
		"--disable-gpu",
		"--disable-dev-shm-usage",
		"--no-sandbox",
		"--hide-scrollbars",
		"--run-all-compositor-stages-before-draw",
		fmt.Sprintf("--screenshot=%s/%s", fileDestinationPath, fileName),
		fmt.Sprintf("%s", websiteURL),
		)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", nil, fmt.Errorf("command finished with error: %v", err)
	}
	log.Printf("Command output is: %s", cmd.Stdout)

	fileData, err := ioutil.ReadFile(fileDestinationPath+"/"+fileName)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read screenshot file: %v", err)
	}

	log.Printf("read %d bytes of data", len(fileData))

	return fileName, fileData, nil
}
