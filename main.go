package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

var (
	TestJsonFile string
	Makefile     string
)

func init() {
	flag.StringVar(&TestJsonFile, "f", "test.json", "JSON file containign test data")
	flag.StringVar(&Makefile, "m", "Makefile", "Makefile to be tested")
	flag.Parse()
}

type TargetConfig struct {
	Name          string `json: "name"`
	Target        string `json: "targetToRun"`
	FilesCreated  string `json: "filesCreated"`
	FilesDeleted  string `json: "filesDeleted"`
	IgnoreFailure bool   `json: searchForFailureInOutput`
}

type TestTargets struct {
	TestTargets []TargetConfig `json: "testTargets`
}

func main() {
	// Prepare variables for counting
	var failCount int = 0
	var passCount int = 0

	// Prepare text Coloring
	colorReset := "\033[0m"
	failColor := "\033[31m"
	passColor := "\033[32m"

	// Parse the JSON file to find targets info
	targetTestInfo := parseJson(TestJsonFile)

	for _, target := range targetTestInfo.TestTargets {
		// Create a variable for each target
		targetPass := true
		fmt.Println("--------------------------------------------------")
		fmt.Println("Running: " + target.Name)
		// Run the target and find any failures
		pass := runTarget(Makefile, target.Target)

		// Check if failures are allowed
		if !target.IgnoreFailure {
			if !pass {
				fmt.Println(target.Name + ": Failure found in output")
				targetPass = false
			}
		}

		// Check for any output files
		if target.FilesCreated != "" {
			splitString := strings.Split(target.FilesCreated, ",")
			for _, fileCreated := range splitString {
				if _, err := os.Stat(fileCreated); os.IsNotExist(err) {
					fmt.Println(target.Target + " missing file: " + fileCreated)
					targetPass = false
				}
			}
		}

		// Check for any deleted files
		if len(target.FilesDeleted) != 0 {
			splitString := strings.Split(target.FilesDeleted, ",")
			for _, fileDeleted := range splitString {
				if _, err := os.Stat(fileDeleted); err == nil {
					fmt.Println(target.Target + " still contains file: " + fileDeleted)
					targetPass = false
				}
			}
		}

		// Check for overall failure
		if !targetPass {
			fmt.Println("Target: " + target.Target + " --> " + string(failColor) + "FAIL" + string(colorReset))
			failCount++
		} else {
			fmt.Println("Target: " + target.Target + " -->  " + string(passColor) + "PASS" + string(colorReset))
			passCount++
		}

	}
	fmt.Println("--------------------------------------------------")
	fmt.Println("Total Tests: " + fmt.Sprint(passCount+failCount))
	fmt.Println(string(passColor) + "# Passed: " + fmt.Sprint(passCount) + string(colorReset))
	if failCount != 0 {
		fmt.Println(string(failColor) + "# Failed: " + fmt.Sprint(failCount) + string(colorReset))
	}
	fmt.Println("--------------------------------------------------")
}

// Function to read the Makefile and search for targets
func readMakefile(makefile string) {

}

// Function to parse the JSON tests
func parseJson(jsonFile string) TestTargets {
	var testTargets TestTargets

	// Check if file even exists
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		fmt.Println("File does not exist")
	}

	// Open the file
	file, err := os.Open(jsonFile)
	if err != nil {
		fmt.Println("Error opening JSON file")
	}

	// Read the file
	jsonData, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("JSON reading error")
	}

	// Unmarshal JSON
	err = json.Unmarshal(jsonData, &testTargets)
	if err != nil {
		fmt.Println("JSON Unmarshalling error")
	}

	return testTargets
}

// Function to run the Makefile target - return PASS/FAIL (TRUE/FALSE) bool
func runTarget(makefile string, target string) bool {
	// Create a cmd that will run the target from the specified Makefile
	fmt.Println("targetCmd: " + target)
	// Turn it into a command
	// FIXME: Currently, this is capturing all output, probably should just capture stderr
	output, err := exec.Command("make", target, "-f"+makefile).CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}

	// Check if the output contains FAIL or fail
	// FIXME: This should just check if stderr was populated instead of looking for "fail"
	containsFailure := strings.Contains(string(output), "FAIL") || strings.Contains(string(output), "fail")
	//fmt.Println("Output: " + string(output))

	// If string contains failure, then return false
	return !containsFailure
}
