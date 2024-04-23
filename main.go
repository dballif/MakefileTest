package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	TestJsonFile string
	Makefile     string
	Generation   bool
)

func init() {
	flag.StringVar(&TestJsonFile, "f", "test.json", "JSON file containign test data")
	flag.StringVar(&Makefile, "m", "Makefile", "Makefile to be tested")
	flag.BoolVar(&Generation, "g", false, "Automatic JSON Generation")
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

	if Generation {
		fmt.Println("Starting Automatic JSON Generation")
		targets := parseMakefileTargets(Makefile)

		// Wrap the targets in the TestTarget structure to allow proper marshalling
		var testTargets TestTargets
		testTargets.TestTargets = append(testTargets.TestTargets, targets...)

		// Create JSON based on the arrays we have collected
		jsonVar, _ := json.MarshalIndent(testTargets, "", "   ")
		fmt.Println(string(jsonVar))

		// Save it to a file
		err := os.WriteFile(TestJsonFile, jsonVar, 0644)
		if err != nil {
			fmt.Println("JSON File writing error")
		}

		return
	}

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

func parseMakefileTargets(makefile string) []TargetConfig {
	// Setup arrays to catch targets & phonys in
	finalTargetArray := []TargetConfig{}
	targetArray := []string{}
	phonyArray := []string{}

	// Parse the Makefile Section
	// Open the Makefile
	file, err := os.Open(makefile)
	if err != nil {
		fmt.Println("Opening Makefile Error: ", err)
		return nil
	}

	// Make sure to close the file when done
	defer file.Close()

	// Create regex
	regexPattern := regexp.MustCompile(`^\s*([^\s#]+)\s*:`)

	// Create scanner to read Makefile line by line
	scanner := bufio.NewScanner(file)

	// Scan the Makefile
	for scanner.Scan() {
		line := scanner.Text()
		// Does it match the regex
		if regexPattern.MatchString(line) {

			// If it is not a .PHONY, it means it should correspond ot a file
			if !strings.Contains(line, ".PHONY:") {
				// Remove the colon
				// Get the target
				target := strings.Split(line, ":")[0]
				targetArray = append(targetArray, target)
			} else {
				target := strings.Split(line, ":")[1]
				// Split on Spaces
				splitTarget := strings.Split(target, " ")
				// Skip the first space which is part of Makefile Standard
				for i := 1; i < len(splitTarget); i++ {
					phonyArray = append(phonyArray, splitTarget[i])
				}
			}
		}
	}

	// Iterate over target array, creating target structures
	for i := 0; i < len(targetArray); i++ {
		// Create a new target Structure
		var newTargetConfig TargetConfig

		// Fill the Structure
		newTargetConfig.Name = targetArray[i] + " Test"
		newTargetConfig.Target = targetArray[i]
		newTargetConfig.IgnoreFailure = true

		// FIXME: Add section that looks for rm or some equivalent
		newTargetConfig.FilesDeleted = ""

		// Check if it is a .PHONY, if it is not, it should correspond to a file
		isPhony := false
		for j := 0; j < len(phonyArray); j++ {
			if phonyArray[j] == targetArray[i] {
				isPhony = true
			}
		}
		if !isPhony {
			newTargetConfig.FilesCreated = targetArray[i]
		} else {
			newTargetConfig.FilesCreated = ""
		}

		// Add the new target to the array
		finalTargetArray = append(finalTargetArray, newTargetConfig)
	}
	fmt.Println(targetArray)
	fmt.Println(phonyArray)
	return finalTargetArray
}

func parseMakefileRm(makefile string) []string {
	// Setup arrays to catch targets & phonys in
	rmArray := []string{}

	// Parse the Makefile Section
	// Open the Makefile
	file, err := os.Open(makefile)
	if err != nil {
		fmt.Println("Opening Makefile Error: ", err)
		return nil
	}

	// Make sure to close the file when done
	defer file.Close()

	// Create regexes
	regexTargetPattern := regexp.MustCompile(`^\s*([^\s#]+)\s*:`)
	regexRmPattern := regexp.MustCompile(`^(@?\s*rm\s+)`)

	// Create scanner to read Makefile line by line
	scanner := bufio.NewScanner(file)

	// Create variable for the target that is being parsed
	var target string

	// Scan the Makefile
	for scanner.Scan() {
		line := scanner.Text()

		//Find target we are working with
		targetMatch := regexTargetPattern.FindStringSubmatch(line)
		if len(targetMatch) > 0 {
			target = targetMatch[1]
		}

		// Does it match the rm regex
		matches := regexRmPattern.FindStringSubmatch(line)
		if len(matches) > 0 {
			rmArray = append(rmArray, target)
		}
	}

	return rmArray
}
