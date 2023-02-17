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

	fmt.Printf("dir %q\n", dir)

	////////////////////////////////////////////
	gitRepo, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		log.Fatal(err)
	}

	// get base name
	gitRepoBase := filepath.Base(string(gitRepo))

	fmt.Printf("gitRepo %q\n", gitRepoBase)

	/////////////////////////////////////////////

	// get git branch
	gitBranch, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("gitBranch %q\n", string(gitBranch))

	///////////////////////////////////////////////

	// exec git get git status
	gitStatus, err := exec.Command("git", "status").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("gitStatus %q\n", string(gitStatus))

	///////////////////////////////////////////////
	gitDiff, err := exec.Command("git", "diff").Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("gitDiff %q\n", string(gitDiff))

	diff := ""
	if len(gitDiff) > 0 {
		p := `The git diff command returned the following:\n%v\n\m`
		diff = string(gitDiff)
		diff = fmt.Sprintf(p, diff)
	}

	prompt := `The current directory is "%v" and the current git branch is "%v" of the repository "%v".
	The status of the repository:
	%v
	%v
	Write the git commands including names and comments to create a branch 
	and commits as if you were the programmer typing them in the terminal.
	`

	prompt = fmt.Sprintf(
		prompt,
		dir,
		string(gitBranch),
		gitRepoBase,
		string(gitStatus),
		diff,
	)

	fmt.Printf("prompt %s\n", prompt)

	//buf := strings.Builder{}
	err = client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			prompt,
		},
		MaxTokens:   gpt3.IntPtr(1000),
		Temperature: gpt3.Float32Ptr(0.6),
	}, func(resp *gpt3.CompletionResponse) {
		//buf.WriteString(resp.Choices[0].Text)
		fmt.Print(resp.Choices[0].Text)
	})
	if err != nil {
		log.Printf("GPT-3 error: %v", err)
		return
	}
}
