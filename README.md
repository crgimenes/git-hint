# Git Hint - Your Personal Git Commit Assistant

Git Hint is a command-line tool that helps you to create descriptive and meaningful commit messages for your Git repository. With the power of OpenAI's GPT-3, Git Hint provides you with useful hints and suggestions on what to write in your commit message.

## Getting Started

### Prerequisites

Git Hint requires Go 1.16 or higher and a valid OpenAI API key.

### Installation

1. Using Go:

```bash
go install crg.eti.br/go/git-hint
```

Make sure that your Go bin directory is in your PATH.

2. Set your OpenAI API key:

```bash
export OPENAI_API_KEY="your_api_key"
```

Recommended: Add the above line to your shell's startup script (e.g. ~/.bashrc).

## Usage

Git Hint is a command-line tool that can be used in any Git repository. To use Git Hint, simply run the \`git hint\` command in your repository.

### Options

Git Hint supports the following command-line options:

- `-h` or `--help`: Displays the help message.

- `-openai_api_key`: Sets the OpenAI API key. If not specified, Git Hint will use the value of the `OPENAI_API_KEY` environment variable.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://github.com/your_username/git-hint/blob/main/LICENSE)

## Acknowledgments

Git Hint is built with the help of [OpenAI GPT-3 API](https://openai.com/).
