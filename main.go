package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var ApiKey string = ""

func main() {
	// Read the API token from an environment variable
	ApiKey = os.Getenv("OPENAI_API_KEY")
	if ApiKey == "" {
		fmt.Println("OPENAI_API_KEY environment variable not set")
		os.Exit(1)
	}

	// Get Arguments & User Input
	// Set Execution Mode according to Argument
	var mode string
	var input string
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "-c", "c", "create":
			mode = "create"
			input = os.Args[1]
		case "-h", "h", "help":
			mode = "help"
		default:
			mode = "help"
		}
	} else if len(os.Args) == 3 {
		switch os.Args[1] {
		case "-d", "d", "describe":
			mode = "describe"
		case "-c", "c", "create":
			mode = "create"
		case "-h", "h", "help":
			mode = "help"
		default:
			mode = "help"
		}
	} else {
		mode = "help"
	}

	if mode != "help" {
		input = os.Args[2]
	} else {
		printHelp()
	}

	result := runPrompt(input, mode)
	fmt.Println(result)

	// Offer to run or describe command
	if mode == "create" {
		fmt.Println("")
		fmt.Println("1) Get a Description")
		fmt.Println("2) Run the Command")
		fmt.Println("3) Exit")
		fmt.Println("WARNING: Don't blindly run commands you do not understand.")
		fmt.Print("\n1 / 2 / 3: ")

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input := scanner.Text()

			switch input {
			case "1":
				resultDescription := runPrompt(result, "describe")
				fmt.Println("\n" + resultDescription)
				// Offer to run Command after Describing it
				fmt.Println("Do you want to run this command? (y/n)")
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					input := scanner.Text()
					if input == "y" {
						runCommand(result)
					}
				}
			case "2":
				err := runCommand(result)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			case "3":
				os.Exit(0)
			default:
				os.Exit(0)
			}
		}

	}

	os.Exit(0)
}

func printHelp() {
	fmt.Printf(`
Helferlein - powered by ChatGPT

Usage: helferlein <arg> "<prompt>"
-c, c, create:
	Describe what you want to achieve, get a Linux Command in return
-d, d, describe:
	Enter a Linux command, get a Description of what it does
-h, h, help:
	Display this Help Section

Before you run Helferlein, set your API-Token by running "export OPENAI_API_KEY=<TOKEN>"
Create a Token at https://platform.openai.com/account/api-keys
`)
	os.Exit(0)
}

func runPrompt(input string, mode string) string {
	// Sanitise User Input
	input = strings.ReplaceAll(input, "\"", "")
	input = strings.TrimSpace(input)
	input = html.EscapeString(input)

	var prompt string
	switch mode {
		case "create":
			prompt = fmt.Sprintf(`Prompt is %q. Try to generate a Linux shell command from it. I want you to only reply with the terminal output inside one unique code block, and nothing else. Do not write explanations.`, input)
		case "describe":
			prompt = fmt.Sprintf(`Prompt is %q. The Prompt is a Linux Shell command. Reply with a short Description of what it does.`, input)
	}

	// Build the API request
	url := fmt.Sprintf("https://api.openai.com/v1/engines/text-davinci-003/completions")
	req, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf(`{
		"prompt": %q,
		"temperature": 0,
		"max_tokens": 100
	}`, prompt)))
	if err != nil {
		fmt.Println("Error building API request:", err)
		os.Exit(1)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ApiKey))

	// Call the OpenAI API to generate a response
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error calling OpenAI API:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Parse the API response
	type errorInfo struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   *string `json:"param,omitempty"`
		Code    *string `json:"code,omitempty"`
	}

	var result struct {
		Choices []struct {
			Text string `json:"text"`
		} `json:"choices"`
		Error *errorInfo `json:"error,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding API response:", err)
		os.Exit(1)
	}

	if result.Error != nil {
		fmt.Println("Error in API response: ", result.Error.Message)
		os.Exit(1)
	}

	// Print the response
	var formattedResult = strings.ReplaceAll(strings.TrimSpace(result.Choices[0].Text), "`", "")
	return formattedResult
}

func runCommand(command string) error {

	// Run the Shell Command
	cmd := exec.Command("sh", "-c", command)

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	// Print the output.
	fmt.Println("\nOutput:")
	fmt.Println(string(output))
	return nil
}
