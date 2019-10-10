package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	. "github.com/logrusorgru/aurora"
	"os"
	"strings"
)

func executor(t string) {
	switch t {
	case "list clusters":
		fmt.Println(Magenta("List of clusters"))
	case "bye":
		fallthrough
	case "exit":
		fallthrough
	case "quit":
		fmt.Println(Magenta("Quitting"))
		os.Exit(0)
	case "help":
		fmt.Println("HELP:\nexit\nquit")
	default:
		fmt.Println("Command not found")
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "help", Description: "show help with all commands"},
		{Text: "exit", Description: "quit the application"},
		{Text: "quit", Description: "quit the application"},
		{Text: "bye", Description: "quit the application"},
	}
	blocks := strings.Split(in.TextBeforeCursor(), " ")

	if in.GetWordBeforeCursor() == "" {
		return nil
	} else {
		return prompt.FilterHasPrefix(s, blocks[0], true)
	}
}

func main() {
	p := prompt.New(executor, completer)
	p.Run()
}
