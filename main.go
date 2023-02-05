package main

import (
	"fmt"
	"os"
	"encoding/json"
	"strings"
	"net/http"
)

func main() {
	// Read the API token from an environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY environment variable not set")
		os.Exit(1)
	}

	// Get Arguments & User Input
	// Set Execution Mode according to Argument
	var mode string
	var input string
	if len(os.Args) == 2 {
		mode = "create"
		input = os.Args[1]
	} else if len(os.Args) == 3 {
		switch os.Args[1] {
			case "-d", "d", "describe":
				mode = "describe"
			case "-c", "c", "create":
				mode = "create"
			case "-h", "h", "help":
				mode = "help"
			default:
				mode = "create"
	}
		if (mode != "help") {
			input = os.Args[2]
		} else {
			printHelp()
		}
	}

	var prompt string
	switch mode {
		case "create":
			prompt = fmt.Sprintf(`Prompt is "%q".
	Try to generate a Linux shell command from it. I want you to only reply with the terminal output inside one unique code block, and nothing else. Do not write explanations.`, input)
		case "describe":
			prompt = fmt.Sprintf(`Prompt is "%q".
	The Prompt is a Linux Shell command. Reply with a short Description of what it does.`, input)
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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Call the OpenAI API to generate a response
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error calling OpenAI API:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Parse the API response
	var result struct {
		Choices []struct {
			Text string `json:"text"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding API response:", err)
		os.Exit(1)
	}

	// Print the response
	fmt.Println(strings.TrimSpace(result.Choices[0].Text))
}

func printHelp() {
	fmt.Println("HELP")
	os.Exit(0)
}