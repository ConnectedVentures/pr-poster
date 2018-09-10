package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

var (
	pullRequestIDRegexp = regexp.MustCompile(`pull/(\d+)`)
)

type payload struct {
	body string
}

func main() {
	pullRequestURL := os.Getenv("CIRCLE_PULL_REQUEST")
	if pullRequestURL == "" {
		fmt.Println("No CIRCLE_PULL_REQUEST env set, skipping")
		os.Exit(0)
	}

	githubAPIUsername := os.Getenv("GITHUB_API_USERNAME")
	if githubAPIUsername == "" {
		fmt.Println("No GITHUB_API_USERNAME env set and is required")
		os.Exit(1)
	}

	githubAPIToken := os.Getenv("GITHUB_API_TOKEN")
	if githubAPIToken == "" {
		fmt.Println("No GITHUB_API_TOKEN env set and is required")
		os.Exit(1)
	}

	username := os.Getenv("CIRCLE_PROJECT_USERNAME")
	if username == "" {
		fmt.Println("No CIRCLE_PROJECT_USERNAME env set and is required")
		os.Exit(1)
	}

	repositoryName := os.Getenv("CIRCLE_PROJECT_REPONAME")
	if repositoryName == "" {
		fmt.Println("No CIRCLE_PROJECT_REPONAME env set and is required")
		os.Exit(1)
	}

	pullRequestID := getPullRequestID(pullRequestURL)
	if pullRequestID == "" {
		fmt.Println(fmt.Sprintf("Could not extract pull request ID from env: %q", pullRequestURL))
		os.Exit(1)
	}

	var bodyBuffer bytes.Buffer
	if len(os.Args) == 1 {
		_, err := io.Copy(&bodyBuffer, os.Stdin)
		if err != nil {
			fmt.Println(fmt.Sprintf("Could not copy contents stdin: %q", err))
			os.Exit(1)
		}
	} else {
		filename := os.Args[1]
		fileHandler, err := os.Open(filename)
		if err != nil {
			fmt.Println(fmt.Sprintf("Could not open %q: %q", filename, err))
			os.Exit(1)
		}

		_, err = io.Copy(&bodyBuffer, fileHandler)
		if err != nil {
			fmt.Println(fmt.Sprintf("Could not copy contents from %q: %q", filename, err))
			os.Exit(1)
		}
	}

	githubPayload := payload{
		body: bodyBuffer.String(),
	}

	githubPayloadBuffer, err := json.Marshal(githubPayload)
	if err != nil {
		fmt.Println(fmt.Sprintf("json.Marshal() error: %v\n", err))
		os.Exit(1)
	}

	url := getPullRequestPostURL(username, repositoryName, pullRequestID)
	fmt.Printf("Sending request: %v\n", githubPayloadBuffer)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(githubPayloadBuffer))
	if err != nil {
		fmt.Println(fmt.Sprintf("http.NewRequest() error: %v\n", err))
		os.Exit(1)
	}

	request.Header.Add("Accept", "application/vnd.github.v3+json")
	request.SetBasicAuth(githubAPIUsername, githubAPIToken)

	c := &http.Client{}
	response, err := c.Do(request)
	if err != nil {
		fmt.Printf("http.Do() error: %v\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("ioutil.ReadAll() error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Comment not created, error: %v\nbody: %q", err, string(responseBody))
		os.Exit(1)
	}

	fmt.Println("PR comment posted successfully")
	os.Exit(0)
}

func getPullRequestID(pullRequestURL string) string {
	if len(pullRequestIDRegexp.FindStringSubmatch(pullRequestURL)) == 2 {
		return pullRequestIDRegexp.FindStringSubmatch(pullRequestURL)[1]
	}
	return ""
}

func getPullRequestPostURL(username, repositoryName, pullRequestID string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%s/comments", username, repositoryName, pullRequestID)
}
