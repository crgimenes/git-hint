package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"crg.eti.br/go/config"
	_ "crg.eti.br/go/config/ini"
	"github.com/PullRequestInc/go-gpt3"
)

type Config struct {
	OpenAIAPIKey string `json:"openai_api_key" ini:"openai_api_key" cfg:"openai_api_key" cfgDescription:"OpenAI API key." cfgRequired:"true"`
}

func main() {
	var cfg Config
	config.File = "config.ini"
	err := config.Parse(&cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := gpt3.NewClient(cfg.OpenAIAPIKey)
	ctx := context.Background()

	////////////////////////////////////////////
	// get current directory

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dirBase := filepath.Base(dir)

	////////////////////////////////////////////
	gitRepo, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		log.Fatal(err)
	}

	// get base name
	gitRepoBase := filepath.Base(string(gitRepo))

	if gitRepoBase == "" {
		log.Fatal("Not a git repository")
	}

	/////////////////////////////////////////////

	// get git branch
	gitBranch, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		log.Fatal(err)
	}

	///////////////////////////////////////////////

	// exec git get git status
	gitStatus, err := exec.Command("git", "status").Output()
	if err != nil {
		log.Fatal(err)
	}

	///////////////////////////////////////////////
	gitDiff, err := exec.Command("git", "diff").Output()
	if err != nil {
		log.Fatal(err)
	}

	///////////////////////////////////////////////

	diff := ""
	if len(gitDiff) > 0 {
		p := `The git diff command returned the following:\n%v\n\n`
		diff = string(gitDiff)
		diff = fmt.Sprintf(p, diff)
	}

	suggestNewBranch := ""
	if string(gitBranch) == "master" || string(gitBranch) == "main" {
		suggestNewBranch = "suggest creating a new branch for your changes."
	}

	prompt := `The current directory is "%v" and the current branch is "%v" of the repository "%v".
The status of the repository:
%v
%v
%v
Write as if you were typing the necessary git commands in the terminal to commit the changes to an appropriately named branch creating commits with insightful but concise descriptions in the present tense.
`

	prompt = fmt.Sprintf(
		prompt,
		dirBase,
		string(gitBranch),
		gitRepoBase,
		string(gitStatus),
		diff,
		suggestNewBranch,
	)

	//fmt.Printf("prompt %s\n", prompt)

	//buf := strings.Builder{}
	//err = client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
	err = client.CompletionStreamWithEngine(ctx, "gpt-3.5-turbo", gpt3.CompletionRequest{
		Prompt: []string{
			prompt,
		},
		MaxTokens:   gpt3.IntPtr(1000),
		Temperature: gpt3.Float32Ptr(1.3),
	}, func(resp *gpt3.CompletionResponse) {
		//buf.WriteString(resp.Choices[0].Text)
		fmt.Print(resp.Choices[0].Text)
	})
	if err != nil {
		log.Printf("GPT-3 error: %v", err)
		return
	}

	fmt.Println()
}
